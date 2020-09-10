/*
Copyright 2020 DevSpace Technologies Inc.

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

package tenancy

import (
	"context"
	"fmt"

	configv1alpha1 "github.com/kiosk-sh/kiosk/pkg/apis/config/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/apis/meta/internalversion"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apiserver/pkg/registry/generic"
	"k8s.io/apiserver/pkg/registry/rest"
	"sigs.k8s.io/apiserver-builder-alpha/pkg/builders"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type NewRESTFunc func(cachedClient client.Client, uncachedClient client.Client, scheme *runtime.Scheme) rest.Storage

var (
	TenancyAccountStorage = builders.NewApiResourceWithStorage( // Resource status endpoint
		InternalAccount,
		func() runtime.Object { return &Account{} },     // Register versioned resource
		func() runtime.Object { return &AccountList{} }, // Register versioned resource list
		NewAccountREST,
	)
	NewAccountREST = func(getter generic.RESTOptionsGetter) rest.Storage {
		return NewAccountRESTFunc(CachedClient, UncachedClient, Scheme)
	}
	NewAccountRESTFunc  NewRESTFunc
	TenancySpaceStorage = builders.NewApiResourceWithStorage( // Resource status endpoint
		InternalSpace,
		func() runtime.Object { return &Space{} },     // Register versioned resource
		func() runtime.Object { return &SpaceList{} }, // Register versioned resource list
		NewSpaceREST,
	)
	NewSpaceREST = func(getter generic.RESTOptionsGetter) rest.Storage {
		return NewSpaceRESTFunc(CachedClient, UncachedClient, Scheme)
	}
	NewSpaceRESTFunc NewRESTFunc
	InternalAccount  = builders.NewInternalResource(
		"accounts",
		"Account",
		func() runtime.Object { return &Account{} },
		func() runtime.Object { return &AccountList{} },
	)
	InternalAccountStatus = builders.NewInternalResourceStatus(
		"accounts",
		"AccountStatus",
		func() runtime.Object { return &Account{} },
		func() runtime.Object { return &AccountList{} },
	)
	InternalSpace = builders.NewInternalResource(
		"spaces",
		"Space",
		func() runtime.Object { return &Space{} },
		func() runtime.Object { return &SpaceList{} },
	)
	InternalSpaceStatus = builders.NewInternalResourceStatus(
		"spaces",
		"SpaceStatus",
		func() runtime.Object { return &Space{} },
		func() runtime.Object { return &SpaceList{} },
	)
	// Registered resources and subresources
	ApiVersion = builders.NewApiGroup("tenancy.kiosk.sh").WithKinds(
		InternalAccount,
		InternalAccountStatus,
		InternalSpace,
		InternalSpaceStatus,
	)

	// Required by code generated by go2idl
	AddToScheme = (&runtime.SchemeBuilder{
		ApiVersion.SchemeBuilder.AddToScheme,
		RegisterDefaults,
	}).AddToScheme
	SchemeBuilder      = ApiVersion.SchemeBuilder
	localSchemeBuilder = &SchemeBuilder
	SchemeGroupVersion = ApiVersion.GroupVersion
)

// Required by code generated by go2idl
// Kind takes an unqualified kind and returns a Group qualified GroupKind
func Kind(kind string) schema.GroupKind {
	return SchemeGroupVersion.WithKind(kind).GroupKind()
}

// Required by code generated by go2idl
// Resource takes an unqualified resource and returns a Group qualified GroupResource
func Resource(resource string) schema.GroupResource {
	return SchemeGroupVersion.WithResource(resource).GroupResource()
}

type FinalizerName string
type NamespacePhase string

// +genclient
// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type Account struct {
	metav1.TypeMeta
	metav1.ObjectMeta
	Spec   AccountSpec
	Status AccountStatus
}

type AccountSpec struct {
	configv1alpha1.AccountSpec
}

type AccountStatus struct {
	configv1alpha1.AccountStatus
}

// +genclient
// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type Space struct {
	metav1.TypeMeta
	metav1.ObjectMeta
	Spec   SpaceSpec
	Status SpaceStatus
}

type SpaceSpec struct {
	Account    string
	Finalizers []corev1.FinalizerName
}

type SpaceStatus struct {
	Phase corev1.NamespacePhase
}

//
// Account Functions and Structs
//
// +k8s:deepcopy-gen=false
type AccountStrategy struct {
	builders.DefaultStorageStrategy
}

// +k8s:deepcopy-gen=false
type AccountStatusStrategy struct {
	builders.DefaultStatusStorageStrategy
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type AccountList struct {
	metav1.TypeMeta
	metav1.ListMeta
	Items []Account
}

func (Account) NewStatus() interface{} {
	return AccountStatus{}
}

func (pc *Account) GetStatus() interface{} {
	return pc.Status
}

func (pc *Account) SetStatus(s interface{}) {
	pc.Status = s.(AccountStatus)
}

func (pc *Account) GetSpec() interface{} {
	return pc.Spec
}

func (pc *Account) SetSpec(s interface{}) {
	pc.Spec = s.(AccountSpec)
}

func (pc *Account) GetObjectMeta() *metav1.ObjectMeta {
	return &pc.ObjectMeta
}

func (pc *Account) SetGeneration(generation int64) {
	pc.ObjectMeta.Generation = generation
}

func (pc Account) GetGeneration() int64 {
	return pc.ObjectMeta.Generation
}

// Registry is an interface for things that know how to store Account.
// +k8s:deepcopy-gen=false
type AccountRegistry interface {
	ListAccounts(ctx context.Context, options *internalversion.ListOptions) (*AccountList, error)
	GetAccount(ctx context.Context, id string, options *metav1.GetOptions) (*Account, error)
	CreateAccount(ctx context.Context, id *Account) (*Account, error)
	UpdateAccount(ctx context.Context, id *Account) (*Account, error)
	DeleteAccount(ctx context.Context, id string) (bool, error)
}

// NewRegistry returns a new Registry interface for the given Storage. Any mismatched types will panic.
func NewAccountRegistry(sp builders.StandardStorageProvider) AccountRegistry {
	return &storageAccount{sp}
}

// Implement Registry
// storage puts strong typing around storage calls
// +k8s:deepcopy-gen=false
type storageAccount struct {
	builders.StandardStorageProvider
}

func (s *storageAccount) ListAccounts(ctx context.Context, options *internalversion.ListOptions) (*AccountList, error) {
	if options != nil && options.FieldSelector != nil && !options.FieldSelector.Empty() {
		return nil, fmt.Errorf("field selector not supported yet")
	}
	st := s.GetStandardStorage()
	obj, err := st.List(ctx, options)
	if err != nil {
		return nil, err
	}
	return obj.(*AccountList), err
}

func (s *storageAccount) GetAccount(ctx context.Context, id string, options *metav1.GetOptions) (*Account, error) {
	st := s.GetStandardStorage()
	obj, err := st.Get(ctx, id, options)
	if err != nil {
		return nil, err
	}
	return obj.(*Account), nil
}

func (s *storageAccount) CreateAccount(ctx context.Context, object *Account) (*Account, error) {
	st := s.GetStandardStorage()
	obj, err := st.Create(ctx, object, nil, &metav1.CreateOptions{})
	if err != nil {
		return nil, err
	}
	return obj.(*Account), nil
}

func (s *storageAccount) UpdateAccount(ctx context.Context, object *Account) (*Account, error) {
	st := s.GetStandardStorage()
	obj, _, err := st.Update(ctx, object.Name, rest.DefaultUpdatedObjectInfo(object), nil, nil, false, &metav1.UpdateOptions{})
	if err != nil {
		return nil, err
	}
	return obj.(*Account), nil
}

func (s *storageAccount) DeleteAccount(ctx context.Context, id string) (bool, error) {
	st := s.GetStandardStorage()
	_, sync, err := st.Delete(ctx, id, nil, &metav1.DeleteOptions{})
	return sync, err
}

//
// Space Functions and Structs
//
// +k8s:deepcopy-gen=false
type SpaceStrategy struct {
	builders.DefaultStorageStrategy
}

// +k8s:deepcopy-gen=false
type SpaceStatusStrategy struct {
	builders.DefaultStatusStorageStrategy
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type SpaceList struct {
	metav1.TypeMeta
	metav1.ListMeta
	Items []Space
}

func (Space) NewStatus() interface{} {
	return SpaceStatus{}
}

func (pc *Space) GetStatus() interface{} {
	return pc.Status
}

func (pc *Space) SetStatus(s interface{}) {
	pc.Status = s.(SpaceStatus)
}

func (pc *Space) GetSpec() interface{} {
	return pc.Spec
}

func (pc *Space) SetSpec(s interface{}) {
	pc.Spec = s.(SpaceSpec)
}

func (pc *Space) GetObjectMeta() *metav1.ObjectMeta {
	return &pc.ObjectMeta
}

func (pc *Space) SetGeneration(generation int64) {
	pc.ObjectMeta.Generation = generation
}

func (pc Space) GetGeneration() int64 {
	return pc.ObjectMeta.Generation
}

// Registry is an interface for things that know how to store Space.
// +k8s:deepcopy-gen=false
type SpaceRegistry interface {
	ListSpaces(ctx context.Context, options *internalversion.ListOptions) (*SpaceList, error)
	GetSpace(ctx context.Context, id string, options *metav1.GetOptions) (*Space, error)
	CreateSpace(ctx context.Context, id *Space) (*Space, error)
	UpdateSpace(ctx context.Context, id *Space) (*Space, error)
	DeleteSpace(ctx context.Context, id string) (bool, error)
}

// NewRegistry returns a new Registry interface for the given Storage. Any mismatched types will panic.
func NewSpaceRegistry(sp builders.StandardStorageProvider) SpaceRegistry {
	return &storageSpace{sp}
}

// Implement Registry
// storage puts strong typing around storage calls
// +k8s:deepcopy-gen=false
type storageSpace struct {
	builders.StandardStorageProvider
}

func (s *storageSpace) ListSpaces(ctx context.Context, options *internalversion.ListOptions) (*SpaceList, error) {
	if options != nil && options.FieldSelector != nil && !options.FieldSelector.Empty() {
		return nil, fmt.Errorf("field selector not supported yet")
	}
	st := s.GetStandardStorage()
	obj, err := st.List(ctx, options)
	if err != nil {
		return nil, err
	}
	return obj.(*SpaceList), err
}

func (s *storageSpace) GetSpace(ctx context.Context, id string, options *metav1.GetOptions) (*Space, error) {
	st := s.GetStandardStorage()
	obj, err := st.Get(ctx, id, options)
	if err != nil {
		return nil, err
	}
	return obj.(*Space), nil
}

func (s *storageSpace) CreateSpace(ctx context.Context, object *Space) (*Space, error) {
	st := s.GetStandardStorage()
	obj, err := st.Create(ctx, object, nil, &metav1.CreateOptions{})
	if err != nil {
		return nil, err
	}
	return obj.(*Space), nil
}

func (s *storageSpace) UpdateSpace(ctx context.Context, object *Space) (*Space, error) {
	st := s.GetStandardStorage()
	obj, _, err := st.Update(ctx, object.Name, rest.DefaultUpdatedObjectInfo(object), nil, nil, false, &metav1.UpdateOptions{})
	if err != nil {
		return nil, err
	}
	return obj.(*Space), nil
}

func (s *storageSpace) DeleteSpace(ctx context.Context, id string) (bool, error) {
	st := s.GetStandardStorage()
	_, sync, err := st.Delete(ctx, id, nil, &metav1.DeleteOptions{})
	return sync, err
}
