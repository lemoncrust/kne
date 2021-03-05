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
		Use:   "topo",
		Short: "Kubernetes Network Emulation CLI",
		Long: `Kubernetes Network Emulation CLI.  Works with meshnet to create 
		layer 2 topology used by containers to layout networks in a k8s environment.`,
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
	//rootCmd.AddCommand(showCmd)
	//rootCmd.AddCommand(graphCmd)

}

var (
	createCmd = &cobra.Command{
		Use:   "create topology",
		Short: "Create Topology",
		PreRun: func(cmd *cobra.Command, args []string) {
			fmt.Printf("Inside subCmd PreRun with args: %v\n", args)
		},
		RunE: createFn,
	}
)

func createFn(cmd *cobra.Command, args []string) error {
	fmt.Printf("Inside Create with args: %v\n", args)
	topopb, err := topo.Load(args[0])
	if err != nil {
		return err
	}
	t, err := topo.New(kubecfg, topopb)
	if err != nil {
		return err
	}
	out := cmd.OutOrStdout()
	fmt.Fprintf(out, "Topology:\n%s\n", proto.MarshalTextString(topopb))
	return t.Create(cmd.Context())
}
