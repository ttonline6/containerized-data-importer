/*
Copyright 2019 Red Hat, Inc.

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

// Code generated by client-gen. DO NOT EDIT.

package internalversion

import (
	operators "github.com/operator-framework/operator-lifecycle-manager/pkg/api/apis/operators"
	scheme "github.com/operator-framework/operator-lifecycle-manager/pkg/api/client/clientset/internalversion/scheme"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"
)

// SubscriptionsGetter has a method to return a SubscriptionInterface.
// A group's client should implement this interface.
type SubscriptionsGetter interface {
	Subscriptions(namespace string) SubscriptionInterface
}

// SubscriptionInterface has methods to work with Subscription resources.
type SubscriptionInterface interface {
	Create(*operators.Subscription) (*operators.Subscription, error)
	Update(*operators.Subscription) (*operators.Subscription, error)
	UpdateStatus(*operators.Subscription) (*operators.Subscription, error)
	Delete(name string, options *v1.DeleteOptions) error
	DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error
	Get(name string, options v1.GetOptions) (*operators.Subscription, error)
	List(opts v1.ListOptions) (*operators.SubscriptionList, error)
	Watch(opts v1.ListOptions) (watch.Interface, error)
	Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *operators.Subscription, err error)
	SubscriptionExpansion
}

// subscriptions implements SubscriptionInterface
type subscriptions struct {
	client rest.Interface
	ns     string
}

// newSubscriptions returns a Subscriptions
func newSubscriptions(c *OperatorsClient, namespace string) *subscriptions {
	return &subscriptions{
		client: c.RESTClient(),
		ns:     namespace,
	}
}

// Get takes name of the subscription, and returns the corresponding subscription object, and an error if there is any.
func (c *subscriptions) Get(name string, options v1.GetOptions) (result *operators.Subscription, err error) {
	result = &operators.Subscription{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("subscriptions").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of Subscriptions that match those selectors.
func (c *subscriptions) List(opts v1.ListOptions) (result *operators.SubscriptionList, err error) {
	result = &operators.SubscriptionList{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("subscriptions").
		VersionedParams(&opts, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested subscriptions.
func (c *subscriptions) Watch(opts v1.ListOptions) (watch.Interface, error) {
	opts.Watch = true
	return c.client.Get().
		Namespace(c.ns).
		Resource("subscriptions").
		VersionedParams(&opts, scheme.ParameterCodec).
		Watch()
}

// Create takes the representation of a subscription and creates it.  Returns the server's representation of the subscription, and an error, if there is any.
func (c *subscriptions) Create(subscription *operators.Subscription) (result *operators.Subscription, err error) {
	result = &operators.Subscription{}
	err = c.client.Post().
		Namespace(c.ns).
		Resource("subscriptions").
		Body(subscription).
		Do().
		Into(result)
	return
}

// Update takes the representation of a subscription and updates it. Returns the server's representation of the subscription, and an error, if there is any.
func (c *subscriptions) Update(subscription *operators.Subscription) (result *operators.Subscription, err error) {
	result = &operators.Subscription{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("subscriptions").
		Name(subscription.Name).
		Body(subscription).
		Do().
		Into(result)
	return
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().

func (c *subscriptions) UpdateStatus(subscription *operators.Subscription) (result *operators.Subscription, err error) {
	result = &operators.Subscription{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("subscriptions").
		Name(subscription.Name).
		SubResource("status").
		Body(subscription).
		Do().
		Into(result)
	return
}

// Delete takes name of the subscription and deletes it. Returns an error if one occurs.
func (c *subscriptions) Delete(name string, options *v1.DeleteOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("subscriptions").
		Name(name).
		Body(options).
		Do().
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *subscriptions) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("subscriptions").
		VersionedParams(&listOptions, scheme.ParameterCodec).
		Body(options).
		Do().
		Error()
}

// Patch applies the patch and returns the patched subscription.
func (c *subscriptions) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *operators.Subscription, err error) {
	result = &operators.Subscription{}
	err = c.client.Patch(pt).
		Namespace(c.ns).
		Resource("subscriptions").
		SubResource(subresources...).
		Name(name).
		Body(data).
		Do().
		Into(result)
	return
}
