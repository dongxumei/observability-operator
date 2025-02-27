/*
Copyright 2023.

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

package apps

import (
	"context"
	"sync"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/grafana/dskit/services"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/utils/pointer"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	"github.com/udmire/observability-operator/api/v1alpha1"
	"github.com/udmire/observability-operator/pkg/apps/reconcile"
	"github.com/udmire/observability-operator/pkg/apps/specs"
	"github.com/udmire/observability-operator/pkg/operator/manager"
	"github.com/udmire/observability-operator/pkg/templates/provider"
)

// AppsReconciler reconciles a Apps object
type AppsReconciler struct {
	*services.BasicService

	client.Client
	Scheme *runtime.Scheme

	cfg Config

	mgr ctrl.Manager
	cnp manager.ClusterNameProvider

	handler       specs.AppHandler
	appReconciler reconcile.AppReconciler
	logger        log.Logger
}

func New(client client.Client, schema *runtime.Scheme, config Config, tp provider.TemplateProvider, logger log.Logger) *AppsReconciler {
	reconciler := &AppsReconciler{
		Client: client,
		Scheme: schema,

		cfg:           config,
		handler:       specs.New(tp, logger),
		appReconciler: reconcile.New(logger, client),
		logger:        logger,
	}
	reconciler.BasicService = services.NewIdleService(func(serviceContext context.Context) error {
		return reconciler.SetupWithManager(reconciler.mgr)
	}, nil)
	return reconciler
}

func (r *AppsReconciler) SetManager(mgr ctrl.Manager) {
	r.mgr = mgr
}

func (r *AppsReconciler) SetClusterNameProvider(cnp manager.ClusterNameProvider) {
	r.cnp = cnp
}

//+kubebuilder:rbac:groups=udmire.cn,resources=apps,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=udmire.cn,resources=apps/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=udmire.cn,resources=apps/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Apps object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.15.0/pkg/reconcile
func (r *AppsReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	level.Info(r.logger).Log("msg", "reconciling applications")
	defer level.Info(r.logger).Log("msg", "done reconciling applications")

	instance := &v1alpha1.Apps{}
	if err := r.Get(ctx, req.NamespacedName, instance); apierrors.IsNotFound(err) {
		level.Error(r.logger).Log("msg", "detected deleted Apps", "err", err)
		return ctrl.Result{}, nil
	} else if err != nil {
		level.Error(r.logger).Log("msg", "unable to get Apps", "err", err)
		return ctrl.Result{}, nil
	}

	r.normalizeApps(instance)

	finalizerName := "apps.udmire.cn/finalizer"
	if instance.ObjectMeta.DeletionTimestamp.IsZero() {
		if !controllerutil.ContainsFinalizer(instance, finalizerName) {
			controllerutil.AddFinalizer(instance, finalizerName)
			if err := r.Update(ctx, instance); err != nil {
				return ctrl.Result{}, err
			}
		}
	} else {
		// The object is being deleted
		if controllerutil.ContainsFinalizer(instance, finalizerName) {
			for _, apploy := range instance.Spec.Apployments {
				selector := r.handler.Selector(apploy)
				if err := r.appReconciler.CleanClusterLayerResources(instance.UID, selector); err != nil {
					return ctrl.Result{}, err
				}
			}

			// remove our finalizer from the list and update it.
			controllerutil.RemoveFinalizer(instance, finalizerName)
			if err := r.Update(ctx, instance); err != nil {
				return ctrl.Result{}, err
			}
		}

		// Stop reconciliation as the item is being deleted
		return ctrl.Result{}, nil
	}

	owner := metav1.OwnerReference{
		APIVersion:         instance.APIVersion,
		BlockOwnerDeletion: pointer.Bool(true),
		Controller:         pointer.Bool(true),
		Kind:               instance.Kind,
		Name:               instance.Name,
		UID:                instance.UID,
	}

	var wg sync.WaitGroup
	semaphore := make(chan struct{}, r.cfg.Concurrency)
	for _, apploy := range instance.Spec.Apployments {
		wg.Add(1)
		go func(app v1alpha1.AppSpec) {
			defer wg.Done()
			semaphore <- struct{}{}

			manifest, err := r.handler.Handle(app)
			if err != nil {
				level.Error(r.logger).Log("msg", "failed to generate manifests", "instance", instance.Name, "applicaion", app.Name, "err", err)
				<-semaphore
				return
			}

			err = r.appReconciler.Reconcile(owner, "application", app.Name, manifest)
			if err != nil {
				level.Error(r.logger).Log("msg", "failed to apply manifests", "instance", instance.Name, "application", app.Name, "err", err)
			}

			<-semaphore
		}(apploy)
	}
	wg.Wait()

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *AppsReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&v1alpha1.Apps{}).
		Complete(r)
}

func (r *AppsReconciler) normalizeApps(instance *v1alpha1.Apps) {
	for name, app := range instance.Spec.Apployments {
		app.Name = name
		app.Namespace = instance.Namespace
	}
}
