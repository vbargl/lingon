// Code generated by lingon. EDIT AS MUCH AS YOU LIKE.

package metricsserver

import (
	"context"
	"errors"
	"os"
	"os/exec"

	promoperator "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	"github.com/volvo-cars/lingon/pkg/kube"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	apiregistrationv1 "k8s.io/kube-aggregator/pkg/apis/apiregistration/v1"
)

// validate the struct implements the interface
var _ kube.Exporter = (*MetricsServer)(nil)

const namespace = "monitoring"

// MetricsServer contains kubernetes manifests
type MetricsServer struct {
	kube.App

	AuthReaderRB                   *rbacv1.RoleBinding
	Deploy                         *appsv1.Deployment
	SA                             *corev1.ServiceAccount
	SVC                            *corev1.Service
	ServiceMonitor                 *promoperator.ServiceMonitor
	SystemAggregatedReaderCR       *rbacv1.ClusterRole
	SystemAuthDelegatorCRB         *rbacv1.ClusterRoleBinding
	SystemCR                       *rbacv1.ClusterRole
	SystemCRB                      *rbacv1.ClusterRoleBinding
	V1Beta1MetricsK8SIoApiservices *apiregistrationv1.APIService
}

// New creates a new MetricsServer
func New() *MetricsServer {
	return &MetricsServer{
		AuthReaderRB:                   AuthReaderRB,
		Deploy:                         Deploy,
		SA:                             SA,
		SVC:                            SVC,
		ServiceMonitor:                 ServiceMonitor,
		SystemAggregatedReaderCR:       SystemAggregatedReaderCR,
		SystemAuthDelegatorCRB:         SystemAuthDelegatorCRB,
		SystemCR:                       SystemCR,
		SystemCRB:                      SystemCRB,
		V1Beta1MetricsK8SIoApiservices: V1Beta1MetricsK8SIoApiservices,
	}
}

// Apply applies the kubernetes objects to the cluster
func (a *MetricsServer) Apply(ctx context.Context) error {
	return Apply(ctx, a)
}

// Export exports the kubernetes objects to YAML files in the given directory
func (a *MetricsServer) Export(dir string) error {
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
