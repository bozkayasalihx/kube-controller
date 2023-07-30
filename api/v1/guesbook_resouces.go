package v1

import (
	"fmt"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (g *Guestbook) Construct() []metav1.Object {
	var obj []metav1.Object

	frontendService := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("%s-frontend", g.Name),
			Namespace: g.Namespace,
			Labels: map[string]string{
				"app.kubernetes.io/name":      "guestbook",
				"app.kubernetes.io/component": "frontend",
			},
		},
		Spec: corev1.ServiceSpec{
			Selector: map[string]string{},
			Ports: []corev1.ServicePort{
				{
					Protocol: corev1.ProtocolTCP,
					Port:     80,
				},
			},
		},
	}

	obj = append(obj, frontendService)

	frontendDeployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("%s-frontend", g.Name),
			Namespace: g.Namespace,
			Labels: map[string]string{
				"app.kubernetes.io/name":      "guestbook",
				"app.kubernetes.io/component": "frontend",
			},
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &g.Spec.FrontendReplicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app.kubernetes.io/name":      "guestbook",
					"app.kubernetes.io/component": "frontend",
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app.kubernetes.io/name":      "guestbook",
						"app.kubernetes.io/component": "frontend",
					},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  g.Name,
							Image: g.Spec.FrontendImage,
							Ports: []corev1.ContainerPort{
								{
									ContainerPort: 80,
								},
							},
							Resources: corev1.ResourceRequirements{
								Requests: corev1.ResourceList{
									corev1.ResourceCPU:    *resource.NewMilliQuantity(int64(100), resource.DecimalSI),
									corev1.ResourceMemory: *resource.NewQuantity(int64(100*1024*1024), resource.BinarySI),
								},
							},
							Env: []corev1.EnvVar{
								{
									Name:  "GET_ENV_FROM",
									Value: "dns",
								},
							},
						},
					},
				},
			},
		},
	}

	obj = append(obj, frontendDeployment)

	return obj
}
