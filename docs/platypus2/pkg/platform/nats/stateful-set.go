// Copyright (c) 2023 Volvo Car Corporation
// SPDX-License-Identifier: Apache-2.0

// Code generated by lingon. EDIT AS MUCH AS YOU LIKE.

package nats

import (
	ku "github.com/volvo-cars/lingon/pkg/kubeutil"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var probe = &corev1.Probe{
	FailureThreshold:    int32(3),
	InitialDelaySeconds: int32(10),
	PeriodSeconds:       int32(30),
	ProbeHandler:        ku.ProbeHTTP("/", PortProbe),
	SuccessThreshold:    int32(1),
	TimeoutSeconds:      int32(5),
}

var startupProbe = &corev1.Probe{
	FailureThreshold:    int32(90),
	InitialDelaySeconds: int32(10),
	PeriodSeconds:       int32(10),
	ProbeHandler: ku.ProbeHTTP(
		ku.PathProbes,
		PortProbe,
	),
	SuccessThreshold: int32(1),
	TimeoutSeconds:   int32(5),
}

var STS = &appsv1.StatefulSet{
	TypeMeta: ku.TypeStatefulSetV1,
	ObjectMeta: metav1.ObjectMeta{
		Labels:    BaseLabels(),
		Name:      appName,
		Namespace: namespace,
	},
	Spec: appsv1.StatefulSetSpec{
		PodManagementPolicy: appsv1.ParallelPodManagement,
		Replicas:            P(int32(replicas)),
		Selector: &metav1.LabelSelector{
			MatchLabels: matchLabels,
		},
		ServiceName: appName,
		Template: corev1.PodTemplateSpec{
			ObjectMeta: metav1.ObjectMeta{
				Annotations: ku.AnnotationPrometheus(
					ku.PathMetrics,
					PortMetrics,
				),
				// TODO: checksum/config
				// Annotations: map[string]string{
				// 	"checksum/config":      "7f38e14ac96d519ba71b7154f94978c89119ea52eb682d730406fcfa3eeee880",
				// 	"prometheus.io/path":   ku.PathMetrics,
				// 	"prometheus.io/port":   fmt.Sprintf("%d", PortMetrics),
				// 	"prometheus.io/scrape": "true",
				// },
				Labels: matchLabels,
			},
			Spec: corev1.PodSpec{
				Affinity: &corev1.Affinity{
					PodAntiAffinity: ku.AntiAffinityHostnameByLabel(
						"app",
						appName,
					),
				},
				Containers: []corev1.Container{
					{
						Name:            appName,
						Image:           ImgNats,
						Command:         cmd[ImgNats],
						ImagePullPolicy: corev1.PullIfNotPresent,

						Env: []corev1.EnvVar{
							ku.EnvVarDownAPI("POD_NAME", "metadata.name"),
							ku.EnvVarDownAPI(
								"POD_NAMESPACE",
								"metadata.namespace",
							),
							{
								Name:  "SERVER_NAME",
								Value: "$(POD_NAME)",
							},
							{
								Name:  "CLUSTER_ADVERTISE",
								Value: "$(POD_NAME).nats.$(POD_NAMESPACE).svc.cluster.local",
							},
						},

						Lifecycle: &corev1.Lifecycle{
							PreStop: &corev1.LifecycleHandler{
								Exec: &corev1.ExecAction{
									Command: []string{
										"nats-server",
										"-sl=ldm=/var/run/nats/nats.pid",
									},
								},
							},
						},

						LivenessProbe:  probe,
						ReadinessProbe: probe,
						StartupProbe:   startupProbe,

						Ports: []corev1.ContainerPort{
							{
								ContainerPort: ports[PortNameClient].Port,
								Name:          ports[PortNameClient].Name,
							},
							{
								ContainerPort: ports[PortNameCluster].Port,
								Name:          ports[PortNameCluster].Name,
							},
							{
								ContainerPort: ports[PortNameMonitor].Port,
								Name:          ports[PortNameMonitor].Name,
							},
						},
						Resources: ku.Resources("2", "4Gi", "2", "4Gi"),
						VolumeMounts: []corev1.VolumeMount{
							cm.VolumeMount,
							{
								MountPath: "/var/run/nats",
								Name:      "pid",
							},
						},
					}, {
						Command:         cmd[ImgConfigReloader],
						Image:           ImgConfigReloader,
						ImagePullPolicy: corev1.PullIfNotPresent,
						Name:            "reloader",
						VolumeMounts: []corev1.VolumeMount{
							cm.VolumeMount,
							{
								MountPath: "/var/run/nats",
								Name:      "pid",
							},
						},
					}, {
						Name:            "promexporter",
						Image:           ImgPromExporter,
						Args:            cmd[ImgPromExporter],
						ImagePullPolicy: corev1.PullIfNotPresent,
						Ports: []corev1.ContainerPort{
							{
								ContainerPort: PortMetrics,
								Name:          PortNameMetrics,
							},
						},
					},
				},
				DNSPolicy:                     corev1.DNSClusterFirst,
				ServiceAccountName:            SA.Name,
				ShareProcessNamespace:         P(true),
				TerminationGracePeriodSeconds: P(int64(60)),
				Volumes: []corev1.Volume{
					cm.VolumeAndMount().Volume(),
					{
						Name:         "pid",
						VolumeSource: corev1.VolumeSource{},
					},
				},
			},
		},
	},
}
