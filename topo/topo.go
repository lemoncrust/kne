package topo

import (
	"context"
	"fmt"
	"io/ioutil"

	log "github.com/sirupsen/logrus"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"

	"github.com/golang/protobuf/proto"
	topopb "github.com/hfam/kne/proto/topo"
)

var (
	meshNetCRD = map[string]string{
		"group":   "networkop.co.uk",
		"version": "v1beta1",
		"plural":  "topologies",
	}
)

// Manager is a topology instance manager for k8s cluster instance.
type Manager struct {
	clientset *kubernetes.Clientset
	tpb       *topopb.Topology
	nodes     map[string]*Node
	links     map[string]*Link
}

// New creates a new topology manager based on the provided kubecfg and topology.
func New(kubecfg string, tpb *topopb.Topology) (*Manager, error) {
	log.Infof("Creating manager for: %s", tpb.Name)
	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", kubecfg)
	if err != nil {
		return nil, err
	}

	// create the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return &Manager{
		clientset: clientset,
		tpb:       tpb,
	}, nil
}

// Create creates an instance of the managed topology.
func (m *Manager) Create(ctx context.Context) error {
	m.GetPods(ctx)
	for _, n := range m.tpb.Nodes {
		log.Infof("Adding Node: %s:%s", n.Name, n.Type)
		m.nodes[n.Name] = &Node{
			namespace: m.tpb.Name,
			pb:        n,
			kClient:   m.clientset,
		}
	}
	for _, l := range m.tpb.Links {
		log.Info("Adding Link: %s:%s %s:%s", l.ANode, l.AInt, l.ZNode, l.ZInt)
		m.links[k] = &Link{
			namespace: m.tpb.Name,
			pb:        l,
			kClient:   m.clientset,
		}
	}
	return nil
}

// GetPods gets all pods in the managed k8s cluster.
func (m *Manager) GetPods(ctx context.Context) error {
	pods, err := m.clientset.CoreV1().Pods("").List(ctx, metav1.ListOptions{})
	if err != nil {
		return err
	}
	for _, p := range pods.Items {
		fmt.Println(p.Namespace, p.Name)
	}
	return nil
}

// Load loads a Topology from fName.
func Load(fName string) (*topopb.Topology, error) {
	b, err := ioutil.ReadFile(fName)
	if err != nil {
		return nil, err
	}
	t := &topopb.Topology{}
	if err := proto.UnmarshalText(string(b), t); err != nil {
		return nil, err
	}
	return t, nil
}

type Node struct {
	namespace string
	pb        *topopb.Node
	kClient   *kubernetes.Clientset
}

func NewNode(namespace string, pb *topopb.Node, kClient *kubernetes.Clientset) *Node {
	return &Node{
		namespace: namespace,
		pb:        pb,
		kClient:   kClient,
	}
}

func (n *Node) Pod(ctx context.Context) (*corev1.Pod, error) {
	return n.kClient.CoreV1().Pods(n.namespace).Get(ctx, n.pb.Name, metav1.GetOptions{})
}

type Link struct {
	namespace string
	pb        *topopb.Link
	kClient   *kubernetes.Clientset
}
