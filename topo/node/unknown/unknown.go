package unknown

import (
	log "github.com/sirupsen/logrus"

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
	log.Infof("Custom node: %v\n", pb)
	return nil
}

func init() {
	node.Register(topopb.Node_Unknown, New)
}
