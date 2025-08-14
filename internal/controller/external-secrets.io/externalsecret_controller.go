/*
Copyright 2025 coindcx.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package externalsecretsio

import (
	"context"
	"github.com/shashank-coindcx/coindcx-cr-watcher/internal/utils"

	externalsecretsv1 "github.com/external-secrets/external-secrets/apis/externalsecrets/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

const (
	kafkaTopic = "external-secret-events"
	cr         = "external-secrets.io/externalsecrets"
)

// ExternalSecretReconciler reconciles a ExternalSecret object
type ExternalSecretReconciler struct {
	client.Client
	Scheme *runtime.Scheme
	Kafka  *utils.Kafka
}

// +kubebuilder:rbac:groups=external-secrets.io,resources=externalsecrets,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=external-secrets.io,resources=externalsecrets/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=external-secrets.io,resources=externalsecrets/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the ExternalSecret object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.21.0/pkg/reconcile
func (r *ExternalSecretReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := logf.FromContext(ctx)

	var es externalsecretsv1.ExternalSecret
	if err := r.Client.Get(ctx, req.NamespacedName, &es); err != nil {
		logger.Error(err, "unable to fetch ExternalSecret", "name", req.Name, "namespace", req.Namespace)

		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	var status string

	for _, cond := range es.Status.Conditions {
		status = string(cond.Message)
	}

	logger.Info("Reconcile called for ExternalSecret",
		"name", es.Name,
		"namespace", es.Namespace,
		"status", status)

	message := utils.KafkaMesage{
		Name:      es.Name,
		CR:        cr,
		Namespace: es.Namespace,
		Status:    status,
	}

	if err := r.Kafka.SendKafkaMessage(kafkaTopic, &message, &logger); err != nil {
		logger.Error(err, "Failed to send message to Kafka", "externalsecret", es.Name)
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *ExternalSecretReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&externalsecretsv1.ExternalSecret{}).
		Named("external-secrets.io-externalsecret").
		Complete(r)
}
