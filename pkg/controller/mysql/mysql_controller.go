/*

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

package mysql

import (
	"context"
	"log"

	myappsv1 "github.com/droot/findemo/pkg/apis/myapps/v1"
	storagev1beta1 "github.com/droot/findemo/pkg/apis/storage/v1beta1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

/**
* USER ACTION REQUIRED: This is a scaffold file intended for the user to modify with their own Controller
* business logic.  Delete these comments after modifying this file.*
 */

// Add creates a new MySQL Controller and adds it to the Manager with default RBAC. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
// USER ACTION REQUIRED: update cmd/manager/main.go to call this myapps.Add(mgr) to install this Controller
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileMySQL{Client: mgr.GetClient(), scheme: mgr.GetScheme()}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("mysql-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	log.Printf("registering controller: ")
	ctrl = c

	// Watch for changes to MySQL
	err = c.Watch(&source.Kind{Type: &myappsv1.MySQL{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	// TODO(user): Modify this to be the types you create
	// Uncomment watch a Deployment created by MySQL - change this for objects you create
	// err = c.Watch(&source.Kind{Type: &appsv1.Deployment{}}, &handler.EnqueueRequestForOwner{
	// 	IsController: true,
	// 	OwnerType:    &myappsv1.MySQL{},
	// })
	if err != nil {
		return err
	}

	return nil
}

var _ reconcile.Reconciler = &ReconcileMySQL{}
var ctrl controller.Controller

// ReconcileMySQL reconciles a MySQL object
type ReconcileMySQL struct {
	client.Client
	scheme       *runtime.Scheme
	watchStorage bool
}

// Reconcile reads that state of the cluster for a MySQL object and makes changes based on the state read
// and what is in the MySQL.Spec
// TODO(user): Modify this Reconcile function to implement your Controller logic.  The scaffolding writes
// a Deployment as an example
// Automatically generate RBAC rules to allow the Controller to read and write Deployments
// +kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=myapps.findemo.com,resources=mysqls,verbs=get;list;watch;create;update;patch;delete
func (r *ReconcileMySQL) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	// Fetch the MySQL instance
	instance := &myappsv1.MySQL{}
	err := r.Get(context.TODO(), request.NamespacedName, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			// Object not found, return.  Created objects are automatically garbage collected.
			// For additional cleanup logic use finalizers.
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return reconcile.Result{}, err
	}

	if !r.watchStorage {
		log.Printf("watching deployments dynamically")
		// TODO(user): Modify this to be the types you create
		// Uncomment watch a Deployment created by MySQL - change this for objects you create
		// err = ctrl.Watch(&source.Kind{Type: &storagev1beta1.Disk{}}, &handler.EnqueueRequestForOwner{
		// 	IsController: true,
		// 	OwnerType:    &myappsv1.MySQL{}})
		u := &unstructured.Unstructured{}
		u.SetGroupVersionKind(schema.GroupVersionKind{
			Group:   "storage.findemo.com",
			Version: "v1beta1",
			Kind:    "Disk",
		})
		err = ctrl.Watch(&source.Kind{Type: u}, &handler.EnqueueRequestForOwner{
			IsController: true,
			OwnerType:    &myappsv1.MySQL{}})
		if err != nil {
			log.Printf("failed to watch deployment: %v", err)
			return reconcile.Result{}, err
		}
		log.Printf("enabled watching deployments")
		r.watchStorage = true
	}

	//
	// install pre-delete hooks (finalizers) if not present and object is not
	// being deleted
	//

	// myFinalizerName := "finalizers.myapps.com"
	//
	// if instance.ObjectMeta.DeletionTimestamp.IsZero() {
	// 	// if does not contain the finalizer, then add it and update it.
	// 	if !containsString(instance.ObjectMeta.Finalizers, myFinalizerName) {
	// 		instance.ObjectMeta.Finalizers = append(instance.ObjectMeta.Finalizers, myFinalizerName)
	// 		if err := r.Update(context.Background(), instance); err != nil {
	// 			return reconcile.Result{Requeue: true}, nil
	// 		}
	// 	}
	// } else {
	// 	if containsString(instance.ObjectMeta.Finalizers, myFinalizerName) {
	// 		log.Printf("deleting the external dependencies")
	// 		// delete the external dependency here
	// 		// if fail to delete the external dependency here, return with error
	// 		// so that it can be retried
	// 		// also write your deletion logic to be idempotent
	//
	// 		instance.ObjectMeta.Finalizers = removeString(instance.ObjectMeta.Finalizers, myFinalizerName)
	// 		if err := r.Update(context.Background(), instance); err != nil {
	// 			return reconcile.Result{Requeue: true}, nil
	// 		}
	// 	}
	// }

	// TODO(user): Change this to be the object type created by your controller
	// Define the desired Deployment object
	deploy := &storagev1beta1.Disk{
		ObjectMeta: metav1.ObjectMeta{
			Name:      instance.Name + "-disk",
			Namespace: instance.Namespace,
		},
		Spec: storagev1beta1.DiskSpec{},
	}
	if err := controllerutil.SetControllerReference(instance, deploy, r.scheme); err != nil {
		return reconcile.Result{}, err
	}

	// TODO(user): Change this for the object type created by your controller
	// Check if the Deployment already exists
	found := &storagev1beta1.Disk{}
	err = r.Get(context.TODO(), types.NamespacedName{Name: deploy.Name, Namespace: deploy.Namespace}, found)
	if err != nil && errors.IsNotFound(err) {
		log.Printf("Creating Disk %s/%s\n", deploy.Namespace, deploy.Name)
		err = r.Create(context.TODO(), deploy)
		if err != nil {
			return reconcile.Result{}, err
		}
	} else if err != nil {
		return reconcile.Result{}, err
	}

	// TODO(user): Change this for the object type created by your controller
	// Update the found object and write the result back if there are any changes
	// if !reflect.DeepEqual(deploy.Spec, found.Spec) {
	// 	found.Spec = deploy.Spec
	// 	log.Printf("Updating Deployment %s/%s\n", deploy.Namespace, deploy.Name)
	// 	err = r.Update(context.TODO(), found)
	// 	if err != nil {
	// 		return reconcile.Result{}, err
	// 	}
	// }
	return reconcile.Result{}, nil
}

func containsString(slice []string, s string) bool {
	for _, item := range slice {
		if item == s {
			return true
		}
	}
	return false
}

func removeString(slice []string, s string) (result []string) {
	for _, item := range slice {
		if item == s {
			continue
		}
		result = append(result, item)
	}
	return
}
