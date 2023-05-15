// Copyright (c) 2023 Volvo Car Corporation
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"

	"github.com/volvo-cars/lingon/pkg/kube"
	"github.com/volvo-cars/lingon/pkg/terra"
	"github.com/volvo-cars/lingoneks/pkg/infra"
	"github.com/volvo-cars/lingoneks/pkg/platform/awsauth"
	"github.com/volvo-cars/lingoneks/pkg/platform/karpenter"
	"github.com/volvo-cars/lingoneks/pkg/terraclient"
	"golang.org/x/exp/slog"
)

var S = terra.String

func main() {
	var apply bool
	var destroy bool
	var plan bool
	flag.BoolVar(
		&apply,
		"apply",
		false,
		"Apply the terraform changes (default: false)",
	)
	flag.BoolVar(
		&destroy,
		"destroy",
		false,
		"Destroy the terraform resources (default: false)",
	)
	flag.BoolVar(
		&plan,
		"plan",
		false,
		"Plan the terraform changes (default: false)",
	)
	flag.Parse()

	ap := AWSParams{
		BackendS3Key: "terriyaki-tf-experiment",
		Region:       "eu-north-1",
		Profile:      "vcc-cdds-prod-legacy",
	}
	p := runParams{
		Apply:          apply,
		Destroy:        destroy,
		Plan:           plan,
		AWSParams:      ap,
		KubeconfigPath: "kubeconfig",
		ManifestPath:   ".lingon/k8s",
		ClusterParams: ClusterParams{
			Name:    "platypus-2",
			Version: "1.24",
			ID:      1,
		},
		TFLabels: map[string]string{
			infra.TagEnv: "dev",
			"terraform":  "true",
		},
		KLabels: map[string]string{
			infra.TagEnv: "dev",
		},
	}

	if err := run(p); err != nil {
		slog.Error("run", "err", err)
		os.Exit(1)
	}
	slog.Info("done")
}

type runParams struct {
	AWSParams      AWSParams
	KubeconfigPath string
	ManifestPath   string
	ClusterParams  ClusterParams
	TFLabels       map[string]string
	KLabels        map[string]string
	Apply          bool
	Destroy        bool
	Plan           bool
}
type AWSParams struct {
	BackendS3Key string
	Region       string
	Profile      string
}
type ClusterParams struct {
	Name    string
	Version string
	ID      int
}

func StepSep(name string) {
	fmt.Printf("\n\n> %s  \n =====================\n\n", name)
}

func run(p runParams) error {
	slog.Info("run", "params", p)
	ctx := context.Background()
	uniqueName := p.ClusterParams.Name

	tf := terraclient.NewClient(
		terraclient.WithDefaultPlan(p.Plan),
		terraclient.WithDefaultApply(p.Apply),
	)
	StepSep("vpc")

	// VPC

	vpcName := uniqueName + "-vpc"
	vpcOpts := infra.Opts{
		Name: uniqueName,
		AZs: [3]string{
			"eu-north-1a", "eu-north-1b", "eu-north-1c",
		},
		CIDR: "10.0.0.0/16",
		PublicSubnetCIDRs: [3]string{
			"10.0.101.0/24", "10.0.102.0/24", "10.0.103.0/24",
		},
		PrivateSubnetCIDRs: [3]string{
			"10.0.1.0/24", "10.0.2.0/24", "10.0.3.0/24",
		},
	}
	slog.Info(
		"vpc stack",
		slog.String("name", vpcName),
		slog.Any("opts", vpcOpts),
	)

	vpc := vpcStack{
		AWSStackConfig: newAWSStackConfig(vpcName, p),
		AWSVPC:         *infra.NewAWSVPC(vpcOpts),
	}

	if err := tf.Run(ctx, &vpc); err != nil {
		return fmt.Errorf("tfrun: handling vpc: %w", err)
	}
	if !vpc.IsStateComplete() {
		slog.Info("VPC state not in sync, finishing here. Is it Applied ?")
		return finishAndDestroy(ctx, p, tf)
	}

	vpcState := vpc.AWSVPC.VPC.StateMust()
	privateSubnetIDs := [3]string{}
	for i, subnet := range vpc.AWSVPC.PrivateSubnets {
		privateSubnetIDs[i] = subnet.StateMust().Id
	}

	// EKS
	StepSep("eks")

	vpcID := vpcState.Id
	eksName := uniqueName + "-eks"
	eksOpts := infra.ClusterOpts{
		Name:             p.ClusterParams.Name,
		Version:          p.ClusterParams.Version,
		VPCID:            vpcID,
		PrivateSubnetIDs: privateSubnetIDs,
	}
	slog.Info(
		"eks stack",
		slog.String("name", eksName),
		slog.Any("opts", eksOpts),
	)
	eks := eksStack{
		AWSStackConfig: newAWSStackConfig(eksName, p),
		Cluster:        *infra.NewCluster(eksOpts),
	}
	if err := tf.Run(ctx, &eks); err != nil {
		return fmt.Errorf("tfrun: handling cluster: %w", err)
	}
	if !eks.IsStateComplete() {
		slog.Info("EKS cluster state not in sync, finishing here")
		return finishAndDestroy(ctx, p, tf)
	}

	eksState := eks.EKSCluster.StateMust()
	oidcState := eks.IAMOIDCProvider.StateMust()

	// KARPENTER INFRA
	StepSep("karpenter infra")

	karpenterName := uniqueName + "-karpenter"
	karinfraOpts := karpenter.InfraOpts{
		Name:             eksState.Name + "-karpenter",
		ClusterName:      eksState.Name,
		ClusterARN:       eksState.Arn,
		PrivateSubnetIDs: privateSubnetIDs,
		OIDCProviderArn:  oidcState.Arn,
		OIDCProviderURL:  oidcState.Url,
	}

	slog.Info("karpenter infra", slog.Any("opts", karinfraOpts))

	ks := karpenterStack{
		AWSStackConfig: newAWSStackConfig(karpenterName, p),
		Infra:          karpenter.NewInfra(karinfraOpts),
	}
	if err := tf.Run(ctx, &ks); err != nil {
		return fmt.Errorf("terraforming karpenter: %w", err)
	}
	if !ks.IsStateComplete() {
		slog.Info(
			"stack state not in sync",
			slog.String("stack", ks.StackName()),
		)
		return finishAndDestroy(ctx, p, tf)
	}

	// CSI EBS INFRA

	StepSep("csi ebs infra")

	csiEbsOpts := infra.CSIOpts{
		ClusterName:     eksState.Name,
		OIDCProviderArn: oidcState.Arn,
		OIDCProviderURL: oidcState.Url,
	}
	cs := csiEbsStack{
		AWSStackConfig: newAWSStackConfig(uniqueName+"-csi-ebs", p),
		CSI:            *infra.NewCSIEBS(csiEbsOpts),
	}
	if err := tf.Run(ctx, &cs); err != nil {
		return fmt.Errorf("terraforming csi-ebs: %w", err)
	}
	if !cs.IsStateComplete() {
		slog.Info(
			"stack state not in sync",
			slog.String("stack", cs.StackName()),
		)
		return finishAndDestroy(ctx, p, tf)
	}

	// KUBECONFIG

	StepSep("kubeconfig")

	slog.Info(
		"getting kubeconfig from aws",
		slog.String("profile", p.AWSParams.Profile),
		slog.String("cluster", p.ClusterParams.Name),
		slog.String("region", p.AWSParams.Region),
		slog.String("kubeconfig", p.KubeconfigPath),
	)

	if err := kubeconfigFromAWSCmd(
		ctx,
		p.AWSParams.Profile,
		p.ClusterParams.Name,
		p.AWSParams.Region,
		p.KubeconfigPath,
	); err != nil {
		return fmt.Errorf("kubeconfig from aws: %w", err)
	}

	kctlOpts := kubectlOpts{
		CtxName:        p.ClusterParams.Name,
		KubeConfigPath: p.KubeconfigPath,
		ManifestPath:   p.ManifestPath,
	}

	// KARPENTER KUBERNETES

	StepSep("k8s karpenter")

	clusterName := eks.EKSCluster.StateMust().Name
	karOpts := karpenter.Opts{
		ClusterName:            eks.EKSCluster.StateMust().Name,
		ClusterEndpoint:        eks.EKSCluster.StateMust().Endpoint,
		IAMRoleArn:             ks.Controller.Role.StateMust().Arn,
		DefaultInstanceProfile: ks.InstanceProfile.InstanceProfile.StateMust().Name,
		InterruptQueue:         ks.SimpleQueue.StateMust().Name,
	}

	slog.Info("karpenter k8s", slog.Any("opts", karOpts))

	kap := karpenter.New(karOpts)

	if err := kubeExportApply(
		ctx,
		kap,
		"karpenter",
		kctlOpts,
		"--server-side=true",
		"apply", "-f", "-",
	); err != nil {
		return err
	}

	// Wait for Karpenter to start before applying CRDs
	// otherwise the webhooks fail.
	// Could take a while for the Fargate nodes to become available.
	// Usually it happens within 2 minutes, but just to be sure...
	timeout := "5m"
	objID := fmt.Sprintf(
		"%s/%s",
		kap.Deploy.TypeMeta.GetObjectKind().GroupVersionKind().GroupKind().String(),
		kap.Deploy.ObjectMeta.Name,
	)

	slog.Info(
		"waiting for karpenter deployment...",
		slog.String("object ID", objID),
		slog.String("timeout", timeout),
	)

	if err := kubectl(
		ctx,
		os.Stdin,
		os.Stdout,
		os.Stderr,
		"--kubeconfig",
		p.KubeconfigPath,
		"--context",
		p.ClusterParams.Name,
		"wait",
		"--namespace", kap.Deploy.Namespace,
		objID,
		"--for=condition=available",
		"--timeout="+timeout,
	); err != nil {
		return fmt.Errorf("waiting for karpenter deployment: %w", err)
	}

	// KARPENTER PROVISIONERS

	StepSep("k8s karpenter provisioners")

	karProvOpts := karpenter.ProvisionersOpts{
		ClusterName:       clusterName,
		AvailabilityZones: vpcOpts.AZs,
	}

	slog.Info("karpenter provisioners", slog.Any("opts", karProvOpts))
	kapProvisioners := karpenter.NewProvisioners(karProvOpts)

	if err := kubeExportApply(
		ctx,
		kapProvisioners,
		"karpenter-provisioners",
		kctlOpts,
		"--server-side=true",
		"apply", "-f", "-",
	); err != nil {
		return err
	}

	// FARGATE AWS AUTH CONFIGMAP

	StepSep("k8s fragate aws auth")

	kmNodeRoleARN := ks.InstanceProfile.IAMRole.StateMust().Arn
	kmFargateRoleARN := ks.FargateProfile.IAMRole.StateMust().Arn
	// Apply the aws-auth configmap
	awsAuth, err := awsauth.NewConfigMap(
		&awsauth.Data{
			MapRoles: karpenter.AWSAuthMapRoles(
				kmNodeRoleARN,
				kmFargateRoleARN,
			),
		},
	)
	if err != nil {
		return fmt.Errorf("creating aws-auth configmap: %w", err)
	}

	if err := kubeExportApply(
		ctx,
		awsAuth,
		"aws-auth",
		kctlOpts,
		"--server-side=true",
		"apply", "-f", "-",
		"--force-conflicts", // Required to become owner
	); err != nil {
		return err
	}

	// This needs to come last, in case state is in sync but destroy flag was
	// passed
	if p.Destroy {
		return finishAndDestroy(ctx, p, tf)
	}

	fmt.Printf("\nTerriyaki Summary:\n")
	for _, mod := range tf.Stacks() {
		diff := "no plan"
		if plan := mod.Plan(); plan != nil {
			diff = fmt.Sprintf(
				"add: %d, destroy: %d",
				len(plan.AddResources), len(plan.DestroyResources),
			)
		}
		fmt.Printf(
			"%s: resources: %s\n",
			mod.StackName(),
			diff,
		)
	}

	StepSep("end")
	return nil
}

func kubectl(
	ctx context.Context,
	stdin io.Reader,
	stdout io.Writer,
	stderr io.Writer,
	args ...string,
) error {
	cmd := exec.CommandContext(ctx, "kubectl", args...)
	cmd.Env = os.Environ() // inherit environment in case we need to use kubectl from a container

	cmd.Stdin = stdin
	cmd.Stdout = stdout
	cmd.Stderr = stderr

	return cmd.Run()
}

type kubectlOpts struct {
	CtxName        string
	KubeConfigPath string
	ManifestPath   string
}

func kubeExportApply(
	ctx context.Context,
	ka kube.Exporter,
	name string,
	p kubectlOpts,
	args ...string,
) error {
	var buf bytes.Buffer
	if err := kube.Export(
		ka,
		kube.WithExportWriter(&buf),
		kube.WithExportAsSingleFile("stdin"),
	); err != nil {
		return fmt.Errorf("exporting %s: %w", name, err)
	}
	if err := kube.Export(
		ka,
		kube.WithExportOutputDirectory(p.ManifestPath),
		kube.WithExportAsSingleFile("%s.yaml"),
	); err != nil {
		return fmt.Errorf("exporting %s: %w", name, err)
	}
	if err := kubectl(
		ctx,
		&buf,
		os.Stdout,
		os.Stderr,
		append(
			append(
				[]string{},
				"--kubeconfig", p.KubeConfigPath,
				"--context", p.CtxName,
			), args...,
		)...,
	); err != nil {
		return fmt.Errorf("applying %s: %w", name, err)
	}
	return nil
}