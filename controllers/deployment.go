package controllers

import (
	"context"
	"fmt"
	"github.com/jarpsimoes/git-http-server-operator/controllers/utils"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	githttpserver1alpha1 "github.com/jarpsimoes/git-http-server-operator/api/v1alpha1"
)

func labels(v *githttpserver1alpha1.GitHttpServer, tier string) map[string]string {
	// Fetches and sets labels

	return map[string]string{
		"operator": "git-http-server-operator",
		"app":      v.Name,
		"tier":     tier,
	}
}

// ensureDeployment ensures Deployment resource presence in given namespace.
func (r *GitHttpServerReconciler) ensureDeployment(request reconcile.Request,
	instance *githttpserver1alpha1.GitHttpServer,
	dep *appsv1.Deployment,
) (*reconcile.Result, error) {

	// See if deployment already exists and create if it doesn't
	found := &appsv1.Deployment{}
	err := r.Get(context.TODO(), types.NamespacedName{
		Name:      dep.Name,
		Namespace: instance.Namespace,
	}, found)
	if err != nil && errors.IsNotFound(err) {

		// Create the deployment
		err = r.Create(context.TODO(), dep)

		if err != nil {
			// Deployment failed
			return &reconcile.Result{}, err
		} else {
			// Deployment was successful
			return nil, nil
		}
	} else if err != nil {
		// Error that isn't due to the deployment not existing
		return &reconcile.Result{}, err
	}

	return nil, nil
}

// backendDeployment is a code for Creating Deployment
func (r *GitHttpServerReconciler) backendDeployment(v *githttpserver1alpha1.GitHttpServer) *appsv1.Deployment {
	logger := log.Log

	if v.Spec.Image == "" {
		v.Spec.Image = "jarpsimoes/git_http_server:latest"
	}

	if v.Spec.HttpPort == 0 {
		v.Spec.HttpPort = 8081
	}

	var livenessProbe = utils.GetProbe(*v)
	coreToleations := utils.ConvertTolerations(v.Spec.Tolerations)

	labels := labels(v, "backend")
	size := int32(1)
	dep := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("%s-deployment", v.Name),
			Namespace: v.Namespace,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &size,
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
				},
				Spec: corev1.PodSpec{
					Tolerations: coreToleations,
					Containers: []corev1.Container{{
						Image:           v.Spec.Image,
						ImagePullPolicy: corev1.PullAlways,
						Name:            fmt.Sprintf("%s-pod", v.Name),
						Env:             utils.MergeConfigurationWithEnvironmentVariables(*v),
						Ports: []corev1.ContainerPort{{
							ContainerPort: v.Spec.HttpPort,
							Name:          "http",
						}},
						LivenessProbe: &livenessProbe,
						StartupProbe:  &livenessProbe,
					}},
				},
			},
		},
	}

	err := controllerutil.SetControllerReference(v, dep, r.Scheme)
	if err != nil {
		logger.Error(err, "Error on backend deployment")
	}
	return dep
}
