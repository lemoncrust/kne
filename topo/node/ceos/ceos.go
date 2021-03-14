package ceos

import (
	"fmt"

	topopb "github.com/hfam/kne/proto/topo"
	"github.com/hfam/kne/topo/node"
)

func New(pb *topopb.Node) (node.Interface, error) {
	if err := defaults(pb); err != nil {
		return nil, err
	}
	return &Node{
		pb: pb,
	}, nil
}

type Node struct {
	pb *topopb.Node
}

func (n *Node) Proto() *topopb.Node {
	return n.pb
}

func defaults(pb *topopb.Node) error {
	if _, ok := pb.Services[443]; !ok {
		pb.Services[443] = &topopb.Service{
			Name:    "ssl",
			Inside:  443,
			Outside: node.GetNextPort(),
		}
	}
	for _, v := range pb.Services {
		if v.Outside == 0 {
			v.Outside = node.GetNextPort()
		}
	}
	pb.Labels["type"] = topopb.Node_AristaCEOS.String()
	if pb.Config.Image == "" {
		pb.Config.Image = "ceos:latest"
	}
	pb.Config.Command = []string{"/sbin/init"}
	pb.Config.Env = map[string]string{
		"CEOS":                                "1",
		"EOS_PLATFORM":                        "ceoslab",
		"container":                           "docker",
		"ETBA":                                "1",
		"SKIP_ZEROTOUCH_BARRIER_IN_SYSDBINIT": "1",
		"INTFTYPE":                            "eth",
	}
	pb.Config.EntryCommand = fmt.Sprintf("kubectl exec -it %s -- Cli", pb.Name)
	if pb.Config.ConfigData != nil {
		pb.Config.ConfigPath = "/mnt/flash"
		pb.Config.ConfigFile = "startup-config"
	}
	pb.Constraints["cpu"] = "0.5"
	pb.Constraints["memory"] = "1Gi"
	return nil
}

func init() {
	node.Register(topopb.Node_AristaCEOS, New)
}
