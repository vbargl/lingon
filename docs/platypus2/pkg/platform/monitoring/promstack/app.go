// Copyright (c) 2023 Volvo Car Corporation
// SPDX-License-Identifier: Apache-2.0

// Code generated by lingon. EDIT AS MUCH AS YOU LIKE.

package promstack

import (
	"context"
	"errors"
	"os"
	"os/exec"

	prometheusoperatorv1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	"github.com/volvo-cars/lingon/pkg/kube"
	admissionregistrationv1 "k8s.io/api/admissionregistration/v1"
	appsv1 "k8s.io/api/apps/v1"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
)

// validate the struct implements the interface
var _ kube.Exporter = (*KubePrometheusStack)(nil)

const namespace = "monitoring"

// KubePrometheusStack contains kubernetes manifests
type KubePrometheusStack struct {
	kube.App

	AdmissionCR                                     *rbacv1.ClusterRole
	AdmissionCRB                                    *rbacv1.ClusterRoleBinding
	AdmissionCreateJOBS                             *batchv1.Job
	AdmissionMutatingwebhookconfigurations          *admissionregistrationv1.MutatingWebhookConfiguration
	AdmissionPatchJOBS                              *batchv1.Job
	AdmissionRB                                     *rbacv1.RoleBinding
	AdmissionRole                                   *rbacv1.Role
	AdmissionSA                                     *corev1.ServiceAccount
	AdmissionValidatingwebhookconfigurations        *admissionregistrationv1.ValidatingWebhookConfiguration
	AlertmanagerAlertmanager                        *prometheusoperatorv1.Alertmanager
	AlertmanagerAlertmanagerSecrets                 *corev1.Secret
	AlertmanagerOverviewCM                          *corev1.ConfigMap
	AlertmanagerRulesPrometheusRule                 *prometheusoperatorv1.PrometheusRule
	AlertmanagerSA                                  *corev1.ServiceAccount
	AlertmanagerSVC                                 *corev1.Service
	AlertmanagerServiceMonitor                      *prometheusoperatorv1.ServiceMonitor
	ApiserverCM                                     *corev1.ConfigMap
	ApiserverServiceMonitor                         *prometheusoperatorv1.ServiceMonitor
	ClusterTotalCM                                  *corev1.ConfigMap
	ConfigReloadersPrometheusRule                   *prometheusoperatorv1.PrometheusRule
	ControllerManagerCM                             *corev1.ConfigMap
	CorednsSVC                                      *corev1.Service
	CorednsServiceMonitor                           *prometheusoperatorv1.ServiceMonitor
	EtcdCM                                          *corev1.ConfigMap
	EtcdPrometheusRule                              *prometheusoperatorv1.PrometheusRule
	GeneralRulesPrometheusRule                      *prometheusoperatorv1.PrometheusRule
	GrafanaCM                                       *corev1.ConfigMap
	GrafanaCR                                       *rbacv1.ClusterRole
	GrafanaCRB                                      *rbacv1.ClusterRoleBinding
	GrafanaConfigDashboardsCM                       *corev1.ConfigMap
	GrafanaDatasourceCM                             *corev1.ConfigMap
	GrafanaDeploy                                   *appsv1.Deployment
	GrafanaOverviewCM                               *corev1.ConfigMap
	GrafanaRB                                       *rbacv1.RoleBinding
	GrafanaRole                                     *rbacv1.Role
	GrafanaSA                                       *corev1.ServiceAccount
	GrafanaSVC                                      *corev1.Service
	GrafanaSecrets                                  *corev1.Secret
	GrafanaServiceMonitor                           *prometheusoperatorv1.ServiceMonitor
	GrafanaTestCM                                   *corev1.ConfigMap
	GrafanaTestPO                                   *corev1.Pod
	GrafanaTestSA                                   *corev1.ServiceAccount
	K8SCorednsCM                                    *corev1.ConfigMap
	K8SResourcesClusterCM                           *corev1.ConfigMap
	K8SResourcesNamespaceCM                         *corev1.ConfigMap
	K8SResourcesNodeCM                              *corev1.ConfigMap
	K8SResourcesPodCM                               *corev1.ConfigMap
	K8SResourcesWorkloadCM                          *corev1.ConfigMap
	K8SResourcesWorkloadsNamespaceCM                *corev1.ConfigMap
	K8SRulesPrometheusRule                          *prometheusoperatorv1.PrometheusRule
	KubeApiserverAvailabilityRulesPrometheusRule    *prometheusoperatorv1.PrometheusRule
	KubeApiserverBurnrateRulesPrometheusRule        *prometheusoperatorv1.PrometheusRule
	KubeApiserverHistogramRulesPrometheusRule       *prometheusoperatorv1.PrometheusRule
	KubeApiserverSlosPrometheusRule                 *prometheusoperatorv1.PrometheusRule
	KubeControllerManagerSVC                        *corev1.Service
	KubeControllerManagerServiceMonitor             *prometheusoperatorv1.ServiceMonitor
	KubeEtcdSVC                                     *corev1.Service
	KubeEtcdServiceMonitor                          *prometheusoperatorv1.ServiceMonitor
	KubePrometheusGeneralRulesPrometheusRule        *prometheusoperatorv1.PrometheusRule
	KubePrometheusNodeRecordingRulesPrometheusRule  *prometheusoperatorv1.PrometheusRule
	KubeProxySVC                                    *corev1.Service
	KubeProxyServiceMonitor                         *prometheusoperatorv1.ServiceMonitor
	KubeSchedulerRulesPrometheusRule                *prometheusoperatorv1.PrometheusRule
	KubeSchedulerSVC                                *corev1.Service
	KubeSchedulerServiceMonitor                     *prometheusoperatorv1.ServiceMonitor
	KubeStateMetricsCR                              *rbacv1.ClusterRole
	KubeStateMetricsCRB                             *rbacv1.ClusterRoleBinding
	KubeStateMetricsDeploy                          *appsv1.Deployment
	KubeStateMetricsPrometheusRule                  *prometheusoperatorv1.PrometheusRule
	KubeStateMetricsSA                              *corev1.ServiceAccount
	KubeStateMetricsSVC                             *corev1.Service
	KubeStateMetricsServiceMonitor                  *prometheusoperatorv1.ServiceMonitor
	KubeletCM                                       *corev1.ConfigMap
	KubeletRulesPrometheusRule                      *prometheusoperatorv1.PrometheusRule
	KubeletServiceMonitor                           *prometheusoperatorv1.ServiceMonitor
	KubernetesAppsPrometheusRule                    *prometheusoperatorv1.PrometheusRule
	KubernetesResourcesPrometheusRule               *prometheusoperatorv1.PrometheusRule
	KubernetesStoragePrometheusRule                 *prometheusoperatorv1.PrometheusRule
	KubernetesSystemApiserverPrometheusRule         *prometheusoperatorv1.PrometheusRule
	KubernetesSystemControllerManagerPrometheusRule *prometheusoperatorv1.PrometheusRule
	KubernetesSystemKubeProxyPrometheusRule         *prometheusoperatorv1.PrometheusRule
	KubernetesSystemKubeletPrometheusRule           *prometheusoperatorv1.PrometheusRule
	KubernetesSystemPrometheusRule                  *prometheusoperatorv1.PrometheusRule
	KubernetesSystemSchedulerPrometheusRule         *prometheusoperatorv1.PrometheusRule
	NamespaceByPodCM                                *corev1.ConfigMap
	NamespaceByWorkloadCM                           *corev1.ConfigMap
	NodeClusterRsrcUseCM                            *corev1.ConfigMap
	NodeExporterPrometheusRule                      *prometheusoperatorv1.PrometheusRule
	NodeExporterRulesPrometheusRule                 *prometheusoperatorv1.PrometheusRule
	NodeNetworkPrometheusRule                       *prometheusoperatorv1.PrometheusRule
	NodeRsrcUseCM                                   *corev1.ConfigMap
	NodeRulesPrometheusRule                         *prometheusoperatorv1.PrometheusRule
	NodesCM                                         *corev1.ConfigMap
	NodesDarwinCM                                   *corev1.ConfigMap
	OperatorCR                                      *rbacv1.ClusterRole
	OperatorCRB                                     *rbacv1.ClusterRoleBinding
	OperatorDeploy                                  *appsv1.Deployment
	OperatorSA                                      *corev1.ServiceAccount
	OperatorSVC                                     *corev1.Service
	OperatorServiceMonitor                          *prometheusoperatorv1.ServiceMonitor
	PersistentvolumesusageCM                        *corev1.ConfigMap
	PodTotalCM                                      *corev1.ConfigMap
	PrometheusCM                                    *corev1.ConfigMap
	PrometheusCR                                    *rbacv1.ClusterRole
	PrometheusCRB                                   *rbacv1.ClusterRoleBinding
	PrometheusNodeExporterDS                        *appsv1.DaemonSet
	PrometheusNodeExporterSA                        *corev1.ServiceAccount
	PrometheusNodeExporterSVC                       *corev1.Service
	PrometheusNodeExporterServiceMonitor            *prometheusoperatorv1.ServiceMonitor
	PrometheusOperatorPrometheusRule                *prometheusoperatorv1.PrometheusRule
	PrometheusPrometheus                            *prometheusoperatorv1.Prometheus
	PrometheusPrometheusRule                        *prometheusoperatorv1.PrometheusRule
	PrometheusSA                                    *corev1.ServiceAccount
	PrometheusSVC                                   *corev1.Service
	PrometheusServiceMonitor                        *prometheusoperatorv1.ServiceMonitor
	ProxyCM                                         *corev1.ConfigMap
	SchedulerCM                                     *corev1.ConfigMap
	WorkloadTotalCM                                 *corev1.ConfigMap
}

// New creates a new KubePrometheusStack
func New() *KubePrometheusStack {
	return &KubePrometheusStack{
		AdmissionCR:                              AdmissionCR,
		AdmissionCRB:                             AdmissionCRB,
		AdmissionCreateJOBS:                      AdmissionCreateJOBS,
		AdmissionMutatingwebhookconfigurations:   AdmissionMutatingwebhookconfigurations,
		AdmissionPatchJOBS:                       AdmissionPatchJOBS,
		AdmissionRB:                              AdmissionRB,
		AdmissionRole:                            AdmissionRole,
		AdmissionSA:                              AdmissionSA,
		AdmissionValidatingwebhookconfigurations: AdmissionValidatingwebhookconfigurations,
		AlertmanagerAlertmanager:                 AlertmanagerAlertmanager,
		AlertmanagerAlertmanagerSecrets:          AlertmanagerAlertmanagerSecrets,
		AlertmanagerOverviewCM:                   AlertmanagerOverviewCM,
		AlertmanagerRulesPrometheusRule:          AlertmanagerRulesPrometheusRule,
		AlertmanagerSA:                           AlertmanagerSA,
		AlertmanagerSVC:                          AlertmanagerSVC,
		AlertmanagerServiceMonitor:               AlertmanagerServiceMonitor,
		ApiserverCM:                              ApiserverCM,
		ApiserverServiceMonitor:                  ApiserverServiceMonitor,
		ClusterTotalCM:                           ClusterTotalCM,
		ConfigReloadersPrometheusRule:            ConfigReloadersPrometheusRule,
		ControllerManagerCM:                      ControllerManagerCM,
		CorednsSVC:                               CorednsSVC,
		CorednsServiceMonitor:                    CorednsServiceMonitor,
		EtcdCM:                                   EtcdCM,
		EtcdPrometheusRule:                       EtcdPrometheusRule,
		GeneralRulesPrometheusRule:               GeneralRulesPrometheusRule,
		GrafanaCM:                                GrafanaCM,
		GrafanaCR:                                GrafanaCR,
		GrafanaCRB:                               GrafanaCRB,
		GrafanaConfigDashboardsCM:                GrafanaConfigDashboardsCM,
		GrafanaDatasourceCM:                      GrafanaDatasourceCM,
		GrafanaDeploy:                            GrafanaDeploy,
		GrafanaOverviewCM:                        GrafanaOverviewCM,
		GrafanaRB:                                GrafanaRB,
		GrafanaRole:                              GrafanaRole,
		GrafanaSA:                                GrafanaSA,
		GrafanaSVC:                               GrafanaSVC,
		GrafanaSecrets:                           GrafanaSecrets,
		GrafanaServiceMonitor:                    GrafanaServiceMonitor,
		GrafanaTestCM:                            GrafanaTestCM,
		GrafanaTestPO:                            GrafanaTestPO,
		GrafanaTestSA:                            GrafanaTestSA,
		K8SCorednsCM:                             K8SCorednsCM,
		K8SResourcesClusterCM:                    K8SResourcesClusterCM,
		K8SResourcesNamespaceCM:                  K8SResourcesNamespaceCM,
		K8SResourcesNodeCM:                       K8SResourcesNodeCM,
		K8SResourcesPodCM:                        K8SResourcesPodCM,
		K8SResourcesWorkloadCM:                   K8SResourcesWorkloadCM,
		K8SResourcesWorkloadsNamespaceCM:         K8SResourcesWorkloadsNamespaceCM,
		K8SRulesPrometheusRule:                   K8SRulesPrometheusRule,
		KubeApiserverAvailabilityRulesPrometheusRule:   KubeApiserverAvailabilityRulesPrometheusRule,
		KubeApiserverBurnrateRulesPrometheusRule:       KubeApiserverBurnrateRulesPrometheusRule,
		KubeApiserverHistogramRulesPrometheusRule:      KubeApiserverHistogramRulesPrometheusRule,
		KubeApiserverSlosPrometheusRule:                KubeApiserverSlosPrometheusRule,
		KubeControllerManagerSVC:                       KubeControllerManagerSVC,
		KubeControllerManagerServiceMonitor:            KubeControllerManagerServiceMonitor,
		KubeEtcdSVC:                                    KubeEtcdSVC,
		KubeEtcdServiceMonitor:                         KubeEtcdServiceMonitor,
		KubePrometheusGeneralRulesPrometheusRule:       KubePrometheusGeneralRulesPrometheusRule,
		KubePrometheusNodeRecordingRulesPrometheusRule: KubePrometheusNodeRecordingRulesPrometheusRule,
		KubeProxySVC:                                    KubeProxySVC,
		KubeProxyServiceMonitor:                         KubeProxyServiceMonitor,
		KubeSchedulerRulesPrometheusRule:                KubeSchedulerRulesPrometheusRule,
		KubeSchedulerSVC:                                KubeSchedulerSVC,
		KubeSchedulerServiceMonitor:                     KubeSchedulerServiceMonitor,
		KubeStateMetricsCR:                              KubeStateMetricsCR,
		KubeStateMetricsCRB:                             KubeStateMetricsCRB,
		KubeStateMetricsDeploy:                          KubeStateMetricsDeploy,
		KubeStateMetricsPrometheusRule:                  KubeStateMetricsPrometheusRule,
		KubeStateMetricsSA:                              KubeStateMetricsSA,
		KubeStateMetricsSVC:                             KubeStateMetricsSVC,
		KubeStateMetricsServiceMonitor:                  KubeStateMetricsServiceMonitor,
		KubeletCM:                                       KubeletCM,
		KubeletRulesPrometheusRule:                      KubeletRulesPrometheusRule,
		KubeletServiceMonitor:                           KubeletServiceMonitor,
		KubernetesAppsPrometheusRule:                    KubernetesAppsPrometheusRule,
		KubernetesResourcesPrometheusRule:               KubernetesResourcesPrometheusRule,
		KubernetesStoragePrometheusRule:                 KubernetesStoragePrometheusRule,
		KubernetesSystemApiserverPrometheusRule:         KubernetesSystemApiserverPrometheusRule,
		KubernetesSystemControllerManagerPrometheusRule: KubernetesSystemControllerManagerPrometheusRule,
		KubernetesSystemKubeProxyPrometheusRule:         KubernetesSystemKubeProxyPrometheusRule,
		KubernetesSystemKubeletPrometheusRule:           KubernetesSystemKubeletPrometheusRule,
		KubernetesSystemPrometheusRule:                  KubernetesSystemPrometheusRule,
		KubernetesSystemSchedulerPrometheusRule:         KubernetesSystemSchedulerPrometheusRule,
		NamespaceByPodCM:                                NamespaceByPodCM,
		NamespaceByWorkloadCM:                           NamespaceByWorkloadCM,
		NodeClusterRsrcUseCM:                            NodeClusterRsrcUseCM,
		NodeExporterPrometheusRule:                      NodeExporterPrometheusRule,
		NodeExporterRulesPrometheusRule:                 NodeExporterRulesPrometheusRule,
		NodeNetworkPrometheusRule:                       NodeNetworkPrometheusRule,
		NodeRsrcUseCM:                                   NodeRsrcUseCM,
		NodeRulesPrometheusRule:                         NodeRulesPrometheusRule,
		NodesCM:                                         NodesCM,
		NodesDarwinCM:                                   NodesDarwinCM,
		OperatorCR:                                      OperatorCR,
		OperatorCRB:                                     OperatorCRB,
		OperatorDeploy:                                  OperatorDeploy,
		OperatorSA:                                      OperatorSA,
		OperatorSVC:                                     OperatorSVC,
		OperatorServiceMonitor:                          OperatorServiceMonitor,
		PersistentvolumesusageCM:                        PersistentvolumesusageCM,
		PodTotalCM:                                      PodTotalCM,
		PrometheusCM:                                    PrometheusCM,
		PrometheusCR:                                    PrometheusCR,
		PrometheusCRB:                                   PrometheusCRB,
		PrometheusNodeExporterDS:                        PrometheusNodeExporterDS,
		PrometheusNodeExporterSA:                        PrometheusNodeExporterSA,
		PrometheusNodeExporterSVC:                       PrometheusNodeExporterSVC,
		PrometheusNodeExporterServiceMonitor:            PrometheusNodeExporterServiceMonitor,
		PrometheusOperatorPrometheusRule:                PrometheusOperatorPrometheusRule,
		PrometheusPrometheus:                            PrometheusPrometheus,
		PrometheusPrometheusRule:                        PrometheusPrometheusRule,
		PrometheusSA:                                    PrometheusSA,
		PrometheusSVC:                                   PrometheusSVC,
		PrometheusServiceMonitor:                        PrometheusServiceMonitor,
		ProxyCM:                                         ProxyCM,
		SchedulerCM:                                     SchedulerCM,
		WorkloadTotalCM:                                 WorkloadTotalCM,
	}
}

// Apply applies the kubernetes objects to the cluster
func (a *KubePrometheusStack) Apply(ctx context.Context) error {
	return Apply(ctx, a)
}

// Export exports the kubernetes objects to YAML files in the given directory
func (a *KubePrometheusStack) Export(dir string) error {
	return kube.Export(a, kube.WithExportOutputDirectory(dir))
}

// Apply applies the kubernetes objects contained in Exporter to the cluster
func Apply(ctx context.Context, km kube.Exporter) error {
	cmd := exec.CommandContext(ctx, "kubectl", "apply", "-f", "-")
	cmd.Env = os.Environ()        // inherit environment in case we need to use kubectl from a container
	stdin, err := cmd.StdinPipe() // pipe to pass data to kubectl
	if err != nil {
		return err
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	go func() {
		defer func() {
			err = errors.Join(err, stdin.Close())
		}()
		if errEW := kube.Export(
			km,
			kube.WithExportWriter(stdin),
			kube.WithExportAsSingleFile("stdin"),
		); errEW != nil {
			err = errors.Join(err, errEW)
		}
	}()

	if errS := cmd.Start(); errS != nil {
		return errors.Join(err, errS)
	}

	// waits for the command to exit and waits for any copying
	// to stdin or copying from stdout or stderr to complete
	return errors.Join(err, cmd.Wait())
}

// P converts T to *T, useful for basic types
func P[T any](t T) *T {
	return &t
}