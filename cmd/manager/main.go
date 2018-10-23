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

package main

import (
	"log"

	"github.com/droot/findemo/pkg/apis"
	"github.com/droot/findemo/pkg/controller"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/runtime/signals"
)

func main() {
	// Get a config to talk to the apiserver
	cfg, err := config.GetConfig()
	if err != nil {
		log.Fatal(err)
	}

	// Create a new Cmd to provide shared dependencies and start components
	mgr, err := manager.New(cfg, manager.Options{ /* Scheme: runtime.NewScheme() */ })
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Registering Components.")
	// gvk := schema.GroupVersionKind{
	// 	Group:   "extensions",
	// 	Version: "v1beta1",
	// 	Kind:    "Deployment",
	// }
	// mgr.GetScheme().AddKnownTypeWithName(gvk, &unstructured.Unstructured{})
	// metav1.AddToGroupVersion(mgr.GetScheme(), schema.GroupVersion{
	// 	Group:   gvk.Group,
	// 	Version: gvk.Version,
	// })
	// gvkList := schema.GroupVersionKind{
	// 	Group:   "extensions",
	// 	Version: "v1beta1",
	// 	Kind:    "DeploymentList",
	// }
	// mgr.GetScheme().AddKnownTypeWithName(gvkList, &unstructured.Unstructured{})
	// metav1.AddToGroupVersion(mgr.GetScheme(), schema.GroupVersion{
	// 	Group:   gvk.Group,
	// 	Version: gvk.Version,
	// })

	// Setup Scheme for all resources
	if err := apis.AddToScheme(mgr.GetScheme()); err != nil {
		log.Fatal(err)
	}

	// Setup all Controllers
	if err := controller.AddToManager(mgr); err != nil {
		log.Fatal(err)
	}

	log.Printf("Starting the Cmd.")

	// Start the Cmd
	log.Fatal(mgr.Start(signals.SetupSignalHandler()))
}
