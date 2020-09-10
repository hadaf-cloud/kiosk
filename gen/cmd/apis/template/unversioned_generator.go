package template

import (
	"io"
	"text/template"

	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/gengo/generator"
	"k8s.io/gengo/namer"
)

type unversionedGenerator struct {
	generator.DefaultGen
	apigroup *APIGroup
}

var _ generator.Generator = &unversionedGenerator{}

func CreateUnversionedGenerator(apigroup *APIGroup, filename string) generator.Generator {
	return &unversionedGenerator{
		generator.DefaultGen{OptionalName: filename},
		apigroup,
	}
}

func (d *unversionedGenerator) Imports(c *generator.Context) []string {
	imports := sets.NewString(
		"fmt",
		"context",
		"k8s.io/apimachinery/pkg/apis/meta/internalversion",
		"k8s.io/apimachinery/pkg/runtime",
		"k8s.io/apimachinery/pkg/runtime/schema",
		"k8s.io/apiserver/pkg/registry/generic",
		"k8s.io/apiserver/pkg/registry/rest",
		"sigs.k8s.io/apiserver-builder-alpha/pkg/builders",
		"sigs.k8s.io/controller-runtime/pkg/client")

	// Get imports for all fields
	for _, s := range d.apigroup.Structs {
		for _, f := range s.Fields {
			if len(f.UnversionedImport) > 0 {
				imports.Insert(f.UnversionedImport)
			}
		}
	}

	return imports.List()
}

func (d *unversionedGenerator) Finalize(context *generator.Context, w io.Writer) error {
	temp := template.
		Must(template.New("unversioned-wiring-template").Funcs(map[string]interface{}{
			"public": namer.IC,
		}).Parse(UnversionedAPITemplate))

	err := temp.Execute(w, d.apigroup)
	if err != nil {
		return err
	}
	return err
}

var UnversionedAPITemplate = `
type NewRESTFunc func(cachedClient client.Client, uncachedClient client.Client, scheme *runtime.Scheme) rest.Storage

var (
	{{ range $api := .UnversionedResources -}}
	{{ if $api.REST -}}
		{{$api.Group|public}}{{$api.Kind}}Storage = builders.NewApiResourceWithStorage( // Resource status endpoint
			Internal{{ $api.Kind }},
			func() runtime.Object { return &{{ $api.Kind }}{} },     // Register versioned resource
			func() runtime.Object { return &{{ $api.Kind }}List{} }, // Register versioned resource list
			New{{ $api.REST }},
		)
		New{{ $api.REST }} = func(getter generic.RESTOptionsGetter) rest.Storage {
			return New{{ $api.REST }}Func(CachedClient, UncachedClient, Scheme)
		}
		New{{ $api.REST }}Func NewRESTFunc
	{{ else -}}
		{{$api.Group|public}}{{$api.Kind}}Storage = builders.NewApiResource( // Resource status endpoint
			Internal{{ $api.Kind }},
			func() runtime.Object { return &{{ $api.Kind }}{} },     // Register versioned resource
			func() runtime.Object { return &{{ $api.Kind }}List{} }, // Register versioned resource list
			&{{ $api.Strategy }}{builders.StorageStrategySingleton},
		)
	{{ end -}}
	{{ end -}}
	{{ range $api := .UnversionedResources -}}
	{{- if $api.ShortName -}}
	Internal{{ $api.Kind }} = builders.NewInternalResourceWithShortcuts(
	{{ else -}}
	Internal{{ $api.Kind }} = builders.NewInternalResource(
	{{ end -}}
		"{{ $api.Resource }}",
        "{{ $api.Kind }}",
		func() runtime.Object { return &{{ $api.Kind }}{} },
		func() runtime.Object { return &{{ $api.Kind }}List{} },
	{{ if $api.ShortName -}}
		[]string{"{{ $api.ShortName }}"},
		[]string{"aggregation"}, // TBD
	{{ end -}}
	)
	Internal{{ $api.Kind }}Status = builders.NewInternalResourceStatus(
		"{{ $api.Resource }}",
        "{{ $api.Kind }}Status",
		func() runtime.Object { return &{{ $api.Kind }}{} },
		func() runtime.Object { return &{{ $api.Kind }}List{} },
	)
	{{ range $subresource := .Subresources -}}
	Internal{{$subresource.Kind}}REST = builders.NewInternalSubresource(
		"{{$subresource.Resource}}", "{{$subresource.Request}}", "{{$subresource.Path}}",
		func() runtime.Object { return &{{$subresource.Request}}{} },
	)
	New{{$subresource.Kind}}REST = func(getter generic.RESTOptionsGetter) rest.Storage {
		return New{{$subresource.Kind}}RESTFunc(CachedClient, UncachedClient, Scheme)
	}
	New{{$subresource.Kind}}RESTFunc NewRESTFunc
	{{ end -}}
	{{ end -}}

	// Registered resources and subresources
	ApiVersion = builders.NewApiGroup("{{.Group}}.{{.Domain}}").WithKinds(
		{{ range $api := .UnversionedResources -}}
		Internal{{$api.Kind}},
		Internal{{$api.Kind}}Status,
		{{ range $subresource := $api.Subresources -}}
		Internal{{$subresource.Kind}}REST,
		{{ end -}}
		{{ end -}}
	)

	// Required by code generated by go2idl
	AddToScheme = (&runtime.SchemeBuilder{
		ApiVersion.SchemeBuilder.AddToScheme, 
		RegisterDefaults, 
	}).AddToScheme
	SchemeBuilder = ApiVersion.SchemeBuilder
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

{{ range $a := .Aliases -}}
type {{ $a.Name }} {{ $a.UnderlyingTypeName }}
{{ end -}}

{{ range $s := .Structs -}}
{{ if $s.GenUnversioned -}}
{{ if $s.GenClient }}// +genclient{{end}}
{{ if $s.GenClient }}// +genclient{{ if $s.NonNamespaced }}:nonNamespaced{{end}}{{end}}
{{ if $s.GenDeepCopy }}// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object{{end}}

type {{ $s.Name }} struct {
{{ range $f := $s.Fields -}}
    {{ $f.Name }} {{ $f.UnversionedType }}
{{ end -}}
}
{{ end -}}
{{ end -}}

{{ range $api := .UnversionedResources -}}
//
// {{.Kind}} Functions and Structs
//
// +k8s:deepcopy-gen=false
type {{.Strategy}} struct {
	builders.DefaultStorageStrategy
}

// +k8s:deepcopy-gen=false
type {{.StatusStrategy}} struct {
	builders.DefaultStatusStorageStrategy
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type {{$api.Kind}}List struct {
	metav1.TypeMeta
	metav1.ListMeta
	Items []{{$api.Kind}}
}

{{ range $subresource := $api.Subresources -}}
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type {{$subresource.Request}}List struct {
	metav1.TypeMeta
	metav1.ListMeta
	Items []{{$subresource.Request}}
}
{{ end -}}

func ({{$api.Kind}}) NewStatus() interface{} {
	return {{$api.Kind}}Status{}
}

func (pc *{{$api.Kind}}) GetStatus() interface{} {
	return pc.Status
}

func (pc *{{$api.Kind}}) SetStatus(s interface{}) {
	pc.Status = s.({{$api.Kind}}Status)
}

func (pc *{{$api.Kind}}) GetSpec() interface{} {
	return pc.Spec
}

func (pc *{{$api.Kind}}) SetSpec(s interface{}) {
	pc.Spec = s.({{$api.Kind}}Spec)
}

func (pc *{{$api.Kind}}) GetObjectMeta() *metav1.ObjectMeta {
	return &pc.ObjectMeta
}

func (pc *{{$api.Kind}}) SetGeneration(generation int64) {
	pc.ObjectMeta.Generation = generation
}

func (pc {{$api.Kind}}) GetGeneration() int64 {
	return pc.ObjectMeta.Generation
}

// Registry is an interface for things that know how to store {{.Kind}}.
// +k8s:deepcopy-gen=false
type {{.Kind}}Registry interface {
	List{{.Kind}}s(ctx context.Context, options *internalversion.ListOptions) (*{{.Kind}}List, error)
	Get{{.Kind}}(ctx context.Context, id string, options *metav1.GetOptions) (*{{.Kind}}, error)
	Create{{.Kind}}(ctx context.Context, id *{{.Kind}}) (*{{.Kind}}, error)
	Update{{.Kind}}(ctx context.Context, id *{{.Kind}}) (*{{.Kind}}, error)
	Delete{{.Kind}}(ctx context.Context, id string) (bool, error)
}

// NewRegistry returns a new Registry interface for the given Storage. Any mismatched types will panic.
func New{{.Kind}}Registry(sp builders.StandardStorageProvider) {{.Kind}}Registry {
	return &storage{{.Kind}}{sp}
}

// Implement Registry
// storage puts strong typing around storage calls
// +k8s:deepcopy-gen=false
type storage{{.Kind}} struct {
	builders.StandardStorageProvider
}

func (s *storage{{.Kind}}) List{{.Kind}}s(ctx context.Context, options *internalversion.ListOptions) (*{{.Kind}}List, error) {
	if options != nil && options.FieldSelector != nil && !options.FieldSelector.Empty() {
		return nil, fmt.Errorf("field selector not supported yet")
	}
	st := s.GetStandardStorage()
	obj, err := st.List(ctx, options)
	if err != nil {
		return nil, err
	}
	return obj.(*{{.Kind}}List), err
}

func (s *storage{{.Kind}}) Get{{.Kind}}(ctx context.Context, id string, options *metav1.GetOptions) (*{{.Kind}}, error) {
	st := s.GetStandardStorage()
	obj, err := st.Get(ctx, id, options)
	if err != nil {
		return nil, err
	}
	return obj.(*{{.Kind}}), nil
}

func (s *storage{{.Kind}}) Create{{.Kind}}(ctx context.Context, object *{{.Kind}}) (*{{.Kind}}, error) {
	st := s.GetStandardStorage()
	obj, err := st.Create(ctx, object, nil, &metav1.CreateOptions{})
	if err != nil {
		return nil, err
	}
	return obj.(*{{.Kind}}), nil
}

func (s *storage{{.Kind}}) Update{{.Kind}}(ctx context.Context, object *{{.Kind}}) (*{{.Kind}}, error) {
	st := s.GetStandardStorage()
	obj, _, err := st.Update(ctx, object.Name, rest.DefaultUpdatedObjectInfo(object), nil, nil, false, &metav1.UpdateOptions{})
	if err != nil {
		return nil, err
	}
	return obj.(*{{.Kind}}), nil
}

func (s *storage{{.Kind}}) Delete{{.Kind}}(ctx context.Context, id string) (bool, error) {
	st := s.GetStandardStorage()
	_, sync, err := st.Delete(ctx, id, nil, &metav1.DeleteOptions{})
	return sync, err
}

{{ end -}}
`
