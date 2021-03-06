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

// Config is the per node specific configuration.
type Config struct {
	command      []string
	args         []string
	image        string
	env          []string
	sleep        int
	entryCmd     string
	cfgPath      string
	bootFileName string
	bootFile     string
}

// Node is a node in the cluster.
type Node struct {
	id        int
	namespace string
	pb        *topopb.Node
	kClient   *kubernetes.Clientset
	cfg       Config
}

// NewNode creates a new node for use in the k8s cluster.  Configure will push the node to
// the cluster.
func NewNode(namespace string, pb *topopb.Node, kClient *kubernetes.Clientset) *Node {
	return &Node{
		namespace: namespace,
		pb:        pb,
		kClient:   kClient,
	}
}

// Configure creates the node on the k8s cluster.
func (n *Node) Configure(ctx context.Context) error {
	cm := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name: fmt.Sprintf("%s-config", n.pb.Name),
		},
		Data: map[string]string{
			n.cfg.bootFileName: n.cfg.bootFile,
		},
	}
	sCM, err := n.kClient.CoreV1().ConfigMaps(n.namespace).Create(ctx, cm, metav1.CreateOptions{})
	if err != nil {
		return err
	}
	log.Infof("Server Config Map:\n%v\n", sCM)
	return nil
}

// Delete removes the Node from the cluster.
func (n *Node) Delete(ctx context.Context) error {
	return n.kClient.CoreV1().ConfigMaps(n.namespace).Delete(ctx, fmt.Sprintf("%s-config", n.pb.Name), metav1.DeleteOptions{})
}

// Pod returns the pod definition for the node.
func (n *Node) Pod(ctx context.Context) (*corev1.Pod, error) {
	return n.kClient.CoreV1().Pods(n.namespace).Get(ctx, n.pb.Name, metav1.GetOptions{})
}

type Link struct {
	namespace string
	pb        *topopb.Link
	kClient   *kubernetes.Clientset
}