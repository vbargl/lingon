// Copyright 2023 Volvo Car Corporation
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"errors"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"runtime/debug"
	"strings"

	"github.com/golingon/lingon/pkg/kube"
	"github.com/golingon/lingon/pkg/kubeutil"
	"k8s.io/apimachinery/pkg/runtime"
	kubescheme "k8s.io/client-go/kubernetes/scheme"

	certmanager "github.com/cert-manager/cert-manager/pkg/api"
	externalsecretsv1beta1 "github.com/external-secrets/external-secrets/apis/externalsecrets/v1beta1"
	traefikv1alpha1 "github.com/traefik/traefik/v3/pkg/provider/kubernetes/crd/traefikio/v1alpha1"
	metallbv1beta2 "go.universe.tf/metallb/api/v1beta2"
	apiextensions "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	apiextensionsbeta "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	// helmv2 "github.com/fluxcd/helm-controller/api/v2beta1"
	// profilev1 "github.com/kubeflow/kubeflow/components/profile-controller/api/v1"
	// profilev1beta1 "github.com/kubeflow/kubeflow/components/profile-controller/api/v1beta1"
	// metacontrolleralpha "github.com/metacontroller/metacontroller/pkg/apis/metacontroller/v1alpha1"
	// otelv1alpha1 "github.com/open-telemetry/opentelemetry-operator/apis/v1alpha1"
	// slothv1alpha1 "github.com/slok/sloth/pkg/kubernetes/api/sloth/v1"
	// tektonpipelinesv1 "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1"
	// tektontriggersv1alpha1 "github.com/tektoncd/triggers/pkg/apis/triggers/v1alpha1"
	// istionetworkingv1alpha3 "istio.io/client-go/pkg/apis/networking/v1alpha3"
	// istionetworkingv1beta1 "istio.io/client-go/pkg/apis/networking/v1beta1"
	// istiosecurityv1beta1 "istio.io/client-go/pkg/apis/security/v1beta1"
	// utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	// knativecachingalpha1 "knative.dev/caching/pkg/apis/caching/v1alpha1"
	// knativeservingv1 "knative.dev/serving/pkg/apis/serving/v1"
	// capiahelm "sigs.k8s.io/cluster-api-addon-provider-helm/api/v1alpha1"
	// gatewayv1alpha2 "sigs.k8s.io/gateway-api/apis/v1alpha2"
	// gatewayv1beta1 "sigs.k8s.io/gateway-api/apis/v1beta1"
	// secretsstorev1 "sigs.k8s.io/secrets-store-csi-driver/apis/v1"
)

const crdMsg = "IF there is an issue with CRDs. Please visit this page to solve it https://github.com/golingon/lingon/tree/main/docs/kubernetes/crd"

func main() {
	var in, out, appName, pkgName string
	var version, verbose, ignoreErr bool

	groupByKind := true
	removeAppName := true
	flag.StringVar(
		&in,
		"in",
		"-",
		"specify the input directory of the yaml manifests, '-' for stdin",
	)
	flag.StringVar(
		&out,
		"out",
		"out",
		"specify the output directory for manifests.",
	)
	flag.StringVar(
		&appName,
		"app",
		"myapp",
		"specify the app name. This will be used as the package name if none is specified.",
	)
	flag.StringVar(
		&pkgName,
		"pkg",
		"",
		"specify the package name. If none is specified the app name will be used. Cannot contain a dash.",
	)
	flag.BoolVar(
		&groupByKind,
		"group",
		true,
		"specify if the output should be grouped by kind (default) or split by name.",
	)
	flag.BoolVar(
		&removeAppName,
		"clean-name",
		true,
		"specify if the app name should be removed from the variable, struct and file name.",
	)
	flag.BoolVar(&version, "version", false, "show version")
	flag.BoolVar(&verbose, "v", false, "show logs")
	flag.BoolVar(
		&ignoreErr,
		"ignore-errors",
		false,
		"ignore errors, useful to generate as much as possible",
	)
	flag.Parse()

	if version {
		printVersion()
		return
	}

	if pkgName == "" {
		pkgName = strings.ReplaceAll(appName, "-", "")
	}

	slog.Info(
		"flags",
		slog.String("in", in),
		slog.String("out", out),
		slog.String("app", appName),
		slog.Bool("group", groupByKind),
		slog.Bool("clean-name", removeAppName),
		slog.Bool("verbose", verbose),
		slog.Bool("ignore-errors", ignoreErr),
	)

	if err := run(
		in,
		out,
		appName,
		pkgName,
		groupByKind,
		removeAppName,
		verbose,
		ignoreErr,
	); err != nil {
		slog.Error(
			"run",
			slog.Any("error", err),
			slog.String("CRD", crdMsg),
		)
		os.Exit(1)
	}

	slog.Info("done")
}

func defaultSerializer() runtime.Decoder {
	// ADD MORE CRDS HERE
	var (
		errs []error
	)

	errs = append(errs,
		certmanager.AddToScheme(kubescheme.Scheme))
	errs = append(errs,
		externalsecretsv1beta1.AddToScheme(kubescheme.Scheme))
	errs = append(errs,
		apiextensions.AddToScheme(kubescheme.Scheme))
	errs = append(errs,
		apiextensionsv1.AddToScheme(kubescheme.Scheme))
	errs = append(errs,
		apiextensionsbeta.AddToScheme(kubescheme.Scheme))
	errs = append(errs,
		traefikv1alpha1.AddToScheme(kubescheme.Scheme))
	errs = append(errs,
		metallbv1beta2.AddToScheme(kubescheme.Scheme))

	if err := errors.Join(errs...); err != nil {
		slog.Error("add to scheme", "error", err)
		os.Exit(1)
	}

	return kubescheme.Codecs.UniversalDeserializer()
}

func run(
	in, out, appName, pkgName string,
	groupByKind, removeAppName, verbose, ignoreErr bool,
) error {
	opts := []kube.ImportOption{
		kube.WithImportAppName(appName),
		kube.WithImportPackageName(pkgName),
		kube.WithImportOutputDirectory(out),
		kube.WithImportSerializer(defaultSerializer()),
	}
	opts = append(opts, kube.WithImportGroupByKind(groupByKind))
	opts = append(opts, kube.WithImportRemoveAppName(removeAppName))
	opts = append(opts, kube.WithImportVerbose(verbose))
	opts = append(opts, kube.WithImportIgnoreErrors(ignoreErr))

	// stdin
	if in == "-" {
		opts = append(opts, kube.WithImportReadStdIn())
		if err := kube.Import(opts...); err != nil {
			return fmt.Errorf("import: %w", err)
		}
		return nil
	}

	// single file
	fi, err := os.Stat(in)
	if err != nil {
		return err
	}
	if !fi.IsDir() {
		opts = append(opts, kube.WithImportManifestFiles([]string{in}))
		if err := kube.Import(opts...); err != nil {
			return fmt.Errorf("import: %w", err)
		}
		return nil
	}

	// directory
	files, err := kubeutil.ListYAMLFiles(in)
	if err != nil {
		slog.Error("list yaml files", "error", err)
	}

	fmt.Printf("files:\n- %s\n", strings.Join(files, "\n- "))
	opts = append(opts, kube.WithImportManifestFiles(files))
	if err := kube.Import(opts...); err != nil {
		return fmt.Errorf("import: %w", err)
	}
	return nil
}

var (
	ver    = "dev"
	commit = "none"
	date   = "unknown"
)

func printVersion() {
	bi, ok := debug.ReadBuildInfo()
	if !ok {
		_, _ = fmt.Fprintln(os.Stderr, "error reading build-info")
		os.Exit(1)
	}
	fmt.Printf("Build:\n%s\n", bi)
	fmt.Printf("Version: %s\n", ver)
	fmt.Printf("Commit: %s\n", commit)
	fmt.Printf("Date: %s\n", date)
}
