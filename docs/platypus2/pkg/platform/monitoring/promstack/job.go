// Copyright (c) 2023 Volvo Car Corporation
// SPDX-License-Identifier: Apache-2.0

// Code generated by lingon. EDIT AS MUCH AS YOU LIKE.

package promstack

import (
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var AdmissionCreateJOBS = &batchv1.Job{
	ObjectMeta: metav1.ObjectMeta{
		Annotations: map[string]string{
			"helm.sh/hook":               "pre-install,pre-upgrade",
			"helm.sh/hook-delete-policy": "before-hook-creation,hook-succeeded",
		},
		Labels: map[string]string{
			"app":                          "kube-prometheus-stack-admission-create",
			"app.kubernetes.io/instance":   "kube-prometheus-stack",
			"app.kubernetes.io/managed-by": "Helm",
			"app.kubernetes.io/part-of":    "kube-prometheus-stack",
			"app.kubernetes.io/version":    "45.27.2",
			"chart":                        "kube-prometheus-stack-45.27.2",
			"heritage":                     "Helm",
			"release":                      "kube-prometheus-stack",
		},
		Name:      "kube-prometheus-stack-admission-create",
		Namespace: namespace,
	},
	Spec: batchv1.JobSpec{
		Template: corev1.PodTemplateSpec{
			ObjectMeta: metav1.ObjectMeta{
				Labels: map[string]string{
					"app":                          "kube-prometheus-stack-admission-create",
					"app.kubernetes.io/instance":   "kube-prometheus-stack",
					"app.kubernetes.io/managed-by": "Helm",
					"app.kubernetes.io/part-of":    "kube-prometheus-stack",
					"app.kubernetes.io/version":    "45.27.2",
					"chart":                        "kube-prometheus-stack-45.27.2",
					"heritage":                     "Helm",
					"release":                      "kube-prometheus-stack",
				},
				Name: "kube-prometheus-stack-admission-create",
			},
			Spec: corev1.PodSpec{
				Containers: []corev1.Container{
					{
						Args: []string{
							"create",
							"--host=kube-prometheus-stack-operator,kube-prometheus-stack-operator.monitoring.svc",
							"--namespace=monitoring",
							"--secret-name=kube-prometheus-stack-admission",
						},
						Image:           "registry.k8s.io/ingress-nginx/kube-webhook-certgen:v20221220-controller-v1.5.1-58-g787ea74b6",
						ImagePullPolicy: corev1.PullPolicy("IfNotPresent"),
						Name:            "create",
					},
				},
				RestartPolicy: corev1.RestartPolicy("OnFailure"),
				SecurityContext: &corev1.PodSecurityContext{
					RunAsGroup:   P(int64(2000)),
					RunAsNonRoot: P(true),
					RunAsUser:    P(int64(2000)),
				},
				ServiceAccountName: "kube-prometheus-stack-admission",
			},
		},
	},
	TypeMeta: metav1.TypeMeta{
		APIVersion: "batch/v1",
		Kind:       "Job",
	},
}

var AdmissionPatchJOBS = &batchv1.Job{
	ObjectMeta: metav1.ObjectMeta{
		Annotations: map[string]string{
			"helm.sh/hook":               "post-install,post-upgrade",
			"helm.sh/hook-delete-policy": "before-hook-creation,hook-succeeded",
		},
		Labels: map[string]string{
			"app":                          "kube-prometheus-stack-admission-patch",
			"app.kubernetes.io/instance":   "kube-prometheus-stack",
			"app.kubernetes.io/managed-by": "Helm",
			"app.kubernetes.io/part-of":    "kube-prometheus-stack",
			"app.kubernetes.io/version":    "45.27.2",
			"chart":                        "kube-prometheus-stack-45.27.2",
			"heritage":                     "Helm",
			"release":                      "kube-prometheus-stack",
		},
		Name:      "kube-prometheus-stack-admission-patch",
		Namespace: namespace,
	},
	Spec: batchv1.JobSpec{
		Template: corev1.PodTemplateSpec{
			ObjectMeta: metav1.ObjectMeta{
				Labels: map[string]string{
					"app":                          "kube-prometheus-stack-admission-patch",
					"app.kubernetes.io/instance":   "kube-prometheus-stack",
					"app.kubernetes.io/managed-by": "Helm",
					"app.kubernetes.io/part-of":    "kube-prometheus-stack",
					"app.kubernetes.io/version":    "45.27.2",
					"chart":                        "kube-prometheus-stack-45.27.2",
					"heritage":                     "Helm",
					"release":                      "kube-prometheus-stack",
				},
				Name: "kube-prometheus-stack-admission-patch",
			},
			Spec: corev1.PodSpec{
				Containers: []corev1.Container{
					{
						Args: []string{
							"patch",
							"--webhook-name=kube-prometheus-stack-admission",
							"--namespace=monitoring",
							"--secret-name=kube-prometheus-stack-admission",
							"--patch-failure-policy=",
						},
						Image:           "registry.k8s.io/ingress-nginx/kube-webhook-certgen:v20221220-controller-v1.5.1-58-g787ea74b6",
						ImagePullPolicy: corev1.PullPolicy("IfNotPresent"),
						Name:            "patch",
					},
				},
				RestartPolicy: corev1.RestartPolicy("OnFailure"),
				SecurityContext: &corev1.PodSecurityContext{
					RunAsGroup:   P(int64(2000)),
					RunAsNonRoot: P(true),
					RunAsUser:    P(int64(2000)),
				},
				ServiceAccountName: "kube-prometheus-stack-admission",
			},
		},
	},
	TypeMeta: metav1.TypeMeta{
		APIVersion: "batch/v1",
		Kind:       "Job",
	},
}