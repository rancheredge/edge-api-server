/*
Copyright The Kubernetes Authors.

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

// Code generated by main. DO NOT EDIT.

package v1

import (
	"github.com/rancher/lasso/pkg/controller"
	"github.com/rancher/wrangler/pkg/schemes"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

func init() {
	v1.AddToScheme(schemes.All)
}

type Interface interface {
	ConfigMap() ConfigMapController
	Endpoints() EndpointsController
	Event() EventController
	Namespace() NamespaceController
	Node() NodeController
	PersistentVolumeClaim() PersistentVolumeClaimController
	Pod() PodController
	Secret() SecretController
	Service() ServiceController
	ServiceAccount() ServiceAccountController
}

func New(controllerFactory controller.SharedControllerFactory) Interface {
	return &version{
		controllerFactory: controllerFactory,
	}
}

type version struct {
	controllerFactory controller.SharedControllerFactory
}

func (c *version) ConfigMap() ConfigMapController {
	return NewConfigMapController(schema.GroupVersionKind{Group: "", Version: "v1", Kind: "ConfigMap"}, "configmaps", c.controllerFactory)
}
func (c *version) Endpoints() EndpointsController {
	return NewEndpointsController(schema.GroupVersionKind{Group: "", Version: "v1", Kind: "Endpoints"}, "endpoints", c.controllerFactory)
}
func (c *version) Event() EventController {
	return NewEventController(schema.GroupVersionKind{Group: "", Version: "v1", Kind: "Event"}, "events", c.controllerFactory)
}
func (c *version) Namespace() NamespaceController {
	return NewNamespaceController(schema.GroupVersionKind{Group: "", Version: "v1", Kind: "Namespace"}, "namespaces", c.controllerFactory)
}
func (c *version) Node() NodeController {
	return NewNodeController(schema.GroupVersionKind{Group: "", Version: "v1", Kind: "Node"}, "nodes", c.controllerFactory)
}
func (c *version) PersistentVolumeClaim() PersistentVolumeClaimController {
	return NewPersistentVolumeClaimController(schema.GroupVersionKind{Group: "", Version: "v1", Kind: "PersistentVolumeClaim"}, "persistentvolumeclaims", c.controllerFactory)
}
func (c *version) Pod() PodController {
	return NewPodController(schema.GroupVersionKind{Group: "", Version: "v1", Kind: "Pod"}, "pods", c.controllerFactory)
}
func (c *version) Secret() SecretController {
	return NewSecretController(schema.GroupVersionKind{Group: "", Version: "v1", Kind: "Secret"}, "secrets", c.controllerFactory)
}
func (c *version) Service() ServiceController {
	return NewServiceController(schema.GroupVersionKind{Group: "", Version: "v1", Kind: "Service"}, "services", c.controllerFactory)
}
func (c *version) ServiceAccount() ServiceAccountController {
	return NewServiceAccountController(schema.GroupVersionKind{Group: "", Version: "v1", Kind: "ServiceAccount"}, "serviceaccounts", c.controllerFactory)
}
