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

package argoprojio

import (
	"context"
	"encoding/json"

	kafka "github.com/IBM/sarama"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	logf "sigs.k8s.io/controller-runtime/pkg/log"

	argoprojiov1alpha1 "github.com/shashank-coindcx/coindcx-cr-watcher/api/argoproj.io/v1alpha1"
)

type KafkaMesage struct {
	Name        string `json:"name"`
	Namespace   string `json:"namespace"`
	Destination string `json:"destination"`
	Source      string `json:"source"`
	Status      string `json:"status"`
}

func (r *ApplicationReconciler) sendKafkaMessage(app *argoprojiov1alpha1.Application) error {
	logger := logf.FromContext(context.Background())
	message := KafkaMesage{
		Name:        app.Name,
		Namespace:   app.Namespace,
		Destination: app.Spec.Destination.Name,
		Source:      app.Spec.Source.RepoURL,
		Status:      app.Status.Sync.Status,
	}

	data, _ := json.Marshal(message)

	msg := &kafka.ProducerMessage{
		Topic: "external-secret-events",
		Value: kafka.StringEncoder(data),
	}

	if _, _, err := r.Producer.SendMessage(msg); err != nil {
		return err
	}

	logger.Info("Message sent to Kafka", "message", message)

	return nil
}

// ApplicationReconciler reconciles a Application object
type ApplicationReconciler struct {
	client.Client
	Scheme   *runtime.Scheme
	Producer kafka.SyncProducer
}

// +kubebuilder:rbac:groups=argoproj.io.internal.coindcx.com,resources=applications,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=argoproj.io.internal.coindcx.com,resources=applications/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=argoproj.io.internal.coindcx.com,resources=applications/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Application object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.21.0/pkg/reconcile
func (r *ApplicationReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := logf.FromContext(ctx)

	var app argoprojiov1alpha1.Application
	if err := r.Client.Get(ctx, req.NamespacedName, &app); err != nil {
		logger.Error(err, "unable to fetch Application", "name", req.Name, "namespace", req.Namespace)

		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	logger.Info("Reconcile called for Application",
		"name", app.Name,
		"namespace", app.Namespace, "destination", app.Spec.Destination.Name,
		"source", app.Spec.Source.RepoURL,
		"status", app.Status.Sync.Status)

	if err := r.sendKafkaMessage(&app); err != nil {
		logger.Error(err, "Failed to send message to Kafka", "application", app.Name)
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *ApplicationReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&argoprojiov1alpha1.Application{}).
		Named("argoproj.io-application").
		Complete(r)
}
