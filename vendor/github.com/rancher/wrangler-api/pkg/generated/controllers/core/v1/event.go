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
	"context"
	"time"

	"github.com/rancher/wrangler/pkg/generic"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/equality"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/watch"
	informers "k8s.io/client-go/informers/core/v1"
	clientset "k8s.io/client-go/kubernetes/typed/core/v1"
	listers "k8s.io/client-go/listers/core/v1"
	"k8s.io/client-go/tools/cache"
)

type EventHandler func(string, *v1.Event) (*v1.Event, error)

type EventController interface {
	generic.ControllerMeta
	EventClient

	OnChange(ctx context.Context, name string, sync EventHandler)
	OnRemove(ctx context.Context, name string, sync EventHandler)
	Enqueue(namespace, name string)
	EnqueueAfter(namespace, name string, duration time.Duration)

	Cache() EventCache
}

type EventClient interface {
	Create(*v1.Event) (*v1.Event, error)
	Update(*v1.Event) (*v1.Event, error)

	Delete(namespace, name string, options *metav1.DeleteOptions) error
	Get(namespace, name string, options metav1.GetOptions) (*v1.Event, error)
	List(namespace string, opts metav1.ListOptions) (*v1.EventList, error)
	Watch(namespace string, opts metav1.ListOptions) (watch.Interface, error)
	Patch(namespace, name string, pt types.PatchType, data []byte, subresources ...string) (result *v1.Event, err error)
}

type EventCache interface {
	Get(namespace, name string) (*v1.Event, error)
	List(namespace string, selector labels.Selector) ([]*v1.Event, error)

	AddIndexer(indexName string, indexer EventIndexer)
	GetByIndex(indexName, key string) ([]*v1.Event, error)
}

type EventIndexer func(obj *v1.Event) ([]string, error)

type eventController struct {
	controllerManager *generic.ControllerManager
	clientGetter      clientset.EventsGetter
	informer          informers.EventInformer
	gvk               schema.GroupVersionKind
}

func NewEventController(gvk schema.GroupVersionKind, controllerManager *generic.ControllerManager, clientGetter clientset.EventsGetter, informer informers.EventInformer) EventController {
	return &eventController{
		controllerManager: controllerManager,
		clientGetter:      clientGetter,
		informer:          informer,
		gvk:               gvk,
	}
}

func FromEventHandlerToHandler(sync EventHandler) generic.Handler {
	return func(key string, obj runtime.Object) (ret runtime.Object, err error) {
		var v *v1.Event
		if obj == nil {
			v, err = sync(key, nil)
		} else {
			v, err = sync(key, obj.(*v1.Event))
		}
		if v == nil {
			return nil, err
		}
		return v, err
	}
}

func (c *eventController) Updater() generic.Updater {
	return func(obj runtime.Object) (runtime.Object, error) {
		newObj, err := c.Update(obj.(*v1.Event))
		if newObj == nil {
			return nil, err
		}
		return newObj, err
	}
}

func UpdateEventDeepCopyOnChange(client EventClient, obj *v1.Event, handler func(obj *v1.Event) (*v1.Event, error)) (*v1.Event, error) {
	if obj == nil {
		return obj, nil
	}

	copyObj := obj.DeepCopy()
	newObj, err := handler(copyObj)
	if newObj != nil {
		copyObj = newObj
	}
	if obj.ResourceVersion == copyObj.ResourceVersion && !equality.Semantic.DeepEqual(obj, copyObj) {
		return client.Update(copyObj)
	}

	return copyObj, err
}

func (c *eventController) AddGenericHandler(ctx context.Context, name string, handler generic.Handler) {
	c.controllerManager.AddHandler(ctx, c.gvk, c.informer.Informer(), name, handler)
}

func (c *eventController) AddGenericRemoveHandler(ctx context.Context, name string, handler generic.Handler) {
	removeHandler := generic.NewRemoveHandler(name, c.Updater(), handler)
	c.controllerManager.AddHandler(ctx, c.gvk, c.informer.Informer(), name, removeHandler)
}

func (c *eventController) OnChange(ctx context.Context, name string, sync EventHandler) {
	c.AddGenericHandler(ctx, name, FromEventHandlerToHandler(sync))
}

func (c *eventController) OnRemove(ctx context.Context, name string, sync EventHandler) {
	removeHandler := generic.NewRemoveHandler(name, c.Updater(), FromEventHandlerToHandler(sync))
	c.AddGenericHandler(ctx, name, removeHandler)
}

func (c *eventController) Enqueue(namespace, name string) {
	c.controllerManager.Enqueue(c.gvk, c.informer.Informer(), namespace, name)
}

func (c *eventController) EnqueueAfter(namespace, name string, duration time.Duration) {
	c.controllerManager.EnqueueAfter(c.gvk, c.informer.Informer(), namespace, name, duration)
}

func (c *eventController) Informer() cache.SharedIndexInformer {
	return c.informer.Informer()
}

func (c *eventController) GroupVersionKind() schema.GroupVersionKind {
	return c.gvk
}

func (c *eventController) Cache() EventCache {
	return &eventCache{
		lister:  c.informer.Lister(),
		indexer: c.informer.Informer().GetIndexer(),
	}
}

func (c *eventController) Create(obj *v1.Event) (*v1.Event, error) {
	return c.clientGetter.Events(obj.Namespace).Create(context.TODO(), obj, metav1.CreateOptions{})
}

func (c *eventController) Update(obj *v1.Event) (*v1.Event, error) {
	return c.clientGetter.Events(obj.Namespace).Update(context.TODO(), obj, metav1.UpdateOptions{})
}

func (c *eventController) Delete(namespace, name string, options *metav1.DeleteOptions) error {
	if options == nil {
		options = &metav1.DeleteOptions{}
	}
	return c.clientGetter.Events(namespace).Delete(context.TODO(), name, *options)
}

func (c *eventController) Get(namespace, name string, options metav1.GetOptions) (*v1.Event, error) {
	return c.clientGetter.Events(namespace).Get(context.TODO(), name, options)
}

func (c *eventController) List(namespace string, opts metav1.ListOptions) (*v1.EventList, error) {
	return c.clientGetter.Events(namespace).List(context.TODO(), opts)
}

func (c *eventController) Watch(namespace string, opts metav1.ListOptions) (watch.Interface, error) {
	return c.clientGetter.Events(namespace).Watch(context.TODO(), opts)
}

func (c *eventController) Patch(namespace, name string, pt types.PatchType, data []byte, subresources ...string) (result *v1.Event, err error) {
	return c.clientGetter.Events(namespace).Patch(context.TODO(), name, pt, data, metav1.PatchOptions{}, subresources...)
}

type eventCache struct {
	lister  listers.EventLister
	indexer cache.Indexer
}

func (c *eventCache) Get(namespace, name string) (*v1.Event, error) {
	return c.lister.Events(namespace).Get(name)
}

func (c *eventCache) List(namespace string, selector labels.Selector) ([]*v1.Event, error) {
	return c.lister.Events(namespace).List(selector)
}

func (c *eventCache) AddIndexer(indexName string, indexer EventIndexer) {
	utilruntime.Must(c.indexer.AddIndexers(map[string]cache.IndexFunc{
		indexName: func(obj interface{}) (strings []string, e error) {
			return indexer(obj.(*v1.Event))
		},
	}))
}

func (c *eventCache) GetByIndex(indexName, key string) (result []*v1.Event, err error) {
	objs, err := c.indexer.ByIndex(indexName, key)
	if err != nil {
		return nil, err
	}
	result = make([]*v1.Event, 0, len(objs))
	for _, obj := range objs {
		result = append(result, obj.(*v1.Event))
	}
	return result, nil
}
