package cmd

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/golang/protobuf/proto"
	"github.com/hfam/kne/topo"
	"github.com/spf13/cobra"
	"k8s.io/client-go/util/homedir"
)

var (
	kubecfg  string
	topofile string

	rootCmd = &cobra.Command{
		Use:   "kne_cli",
		Short: "Kubernetes Network Emulation CLI",
		Long: `Kubernetes Network Emulation CLI.  Works with meshnet to create 
layer 2 topology used by containers to layout networks in a k8s
environment.`,
		SilenceUsage: true,
	}
)

// ExecuteContext executes the root command.
func ExecuteContext(ctx context.Context) error {
	return rootCmd.ExecuteContext(ctx)
}

func init() {
	defaultKubeCfg := ""
	if home := homedir.HomeDir(); home != "" {
		defaultKubeCfg = filepath.Join(home, ".kube", "config")
	}
	rootCmd.SetOut(os.Stdout)
	rootCmd.PersistentFlags().StringVar(&kubecfg, "kubecfg", defaultKubeCfg, "kubeconfig file")
	rootCmd.AddCommand(createCmd)
	//rootCmd.AddCommand(destroyCmd)
	rootCmd.AddCommand(showCmd)
	//rootCmd.AddCommand(graphCmd)

}

var (
	createCmd = &cobra.Command{
		Use:       "create <topology file>",
		Short:     "Create Topology",
		PreRunE:   validateTopology,
		RunE:      createFn,
		ValidArgs: []string{"topology"},
	}
	showCmd = &cobra.Command{
		Use:       "show <topology file>",
		Short:     "Show Topology",
		PreRunE:   validateTopology,
		RunE:      showFn,
		ValidArgs: []string{"topology"},
	}
)

func validateTopology(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("%s: topology must be provided", cmd.Use)
	}
	return nil
}

func createFn(cmd *cobra.Command, args []string) error {
	topopb, err := topo.Load(args[0])
	if err != nil {
		return fmt.Errorf("%s: %w", cmd.Use, err)
	}
	t, err := topo.New(kubecfg, topopb)
	if err != nil {
		return fmt.Errorf("%s: %w", cmd.Use, err)
	}
	out := cmd.OutOrStdout()
	fmt.Fprintf(out, "Topology:\n%s\n", proto.MarshalTextString(topopb))
	return t.Create(cmd.Context())
}

func showFn(cmd *cobra.Command, args []string) error {
	topopb, err := topo.Load(args[0])
	if err != nil {
		return fmt.Errorf("%s: %w", cmd.Use, err)
	}
	t, err := topo.New(kubecfg, topopb)
	if err != nil {
		return fmt.Errorf("%s: %w", cmd.Use, err)
	}
	return t.Topology(cmd.Context())
}
