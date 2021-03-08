package v1beta1

import (
	"context"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"

	topologyv1 "github.com/hfam/kne/api/types/v1beta1"
)

// TopologyInterface provides access to the Topology CRD.
type TopologyInterface interface {
	List(ctx context.Context, opts metav1.ListOptions) (*topologyv1.TopologyList, error)
	Get(ctx context.Context, name string, options metav1.GetOptions) (*topologyv1.Topology, error)
	Create(ctx context.Context, topology *topologyv1.Topology) (*topologyv1.Topology, error)
	Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error)
}

type ClientInterface interface {
	Topology(namespace string) TopologyInterface
}

// TopologyClient is a client for the topology crds.
type TopologyClient struct {
	restClient rest.Interface
}

// NewForConfig returns a new TopologyClient based on c.
func NewForConfig(c *rest.Config) (*TopologyClient, error) {
	config := *c
	config.ContentConfig.GroupVersion = &schema.GroupVersion{Group: topologyv1.GroupName, Version: topologyv1.GroupVersion}
	config.APIPath = "/apis"
	config.NegotiatedSerializer = scheme.Codecs.WithoutConversion()
	config.UserAgent = rest.DefaultKubernetesUserAgent()

	client, err := rest.RESTClientFor(&config)
	if err != nil {
		return nil, err
	}

	return &TopologyClient{restClient: client}, nil
}

func (t *TopologyClient) Topology(namespace string) TopologyInterface {
	return &topologyClient{
		restClient: t.restClient,
		ns:         namespace,
	}
}

type topologyClient struct {
	restClient rest.Interface
	ns         string
}

func (t *topologyClient) List(ctx context.Context, opts metav1.ListOptions) (*topologyv1.TopologyList, error) {
	result := topologyv1.TopologyList{}
	err := t.restClient.
		Get().
		Namespace(t.ns).
		Resource("topologies").
		VersionedParams(&opts, scheme.ParameterCodec).
		Do(ctx).
		Into(&result)

	return &result, err
}

func (t *topologyClient) Get(ctx context.Context, name string, opts metav1.GetOptions) (*topologyv1.Topology, error) {
	result := topologyv1.Topology{}
	err := t.restClient.
		Get().
		Namespace(t.ns).
		Resource("topology").
		Name(name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Do(ctx).
		Into(&result)

	return &result, err
}

func (t *topologyClient) Create(ctx context.Context, topology *topologyv1.Topology) (*topologyv1.Topology, error) {
	result := topologyv1.Topology{}
	err := t.restClient.
		Post().
		Namespace(t.ns).
		Resource("topology").
		Body(topology).
		Do(ctx).
		Into(&result)

	return &result, err
}

func (t *topologyClient) Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error) {
	opts.Watch = true
	return t.restClient.
		Get().
		Namespace(t.ns).
		Resource("topology").
		VersionedParams(&opts, scheme.ParameterCodec).
		Watch(ctx)
}