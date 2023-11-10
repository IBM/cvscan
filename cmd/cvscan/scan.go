package main

import (
	"fmt"

	"github.com/spf13/cobra"
	"k8s.io/client-go/tools/clientcmd"

	"github.com/IBM/cvscan/pkg/scan"
	"github.com/IBM/cvscan/pkg/version"
)

func newScanCmd() *cobra.Command {
	s := scan.NewScanner()
	versionFlag := false

	cmd := &cobra.Command{
		Use:   "cvscan output_path",
		Short: "Take a snapshot of live kubernetes resources",
		RunE: func(cmd *cobra.Command, args []string) error {
			if versionFlag {
				fmt.Println(version.Version())
				return nil
			}

			return s.Run(args)
		},
	}

	flags := cmd.Flags()
	flags.StringVarP(&s.Opts.LabelSelector, "selector", "l", "", "Selector (label query) to filter on")
	flags.StringVar(&s.Opts.FieldSelector, "field-selector", "", "Selector (field query) to filter on")
	flags.StringVar(&s.Kubeconfig, "kubeconfig", "", "Path to the kubeconfig file to use for CLI requests")
	flags.BoolVar(&s.ClusterWideOnly, "cluster-wide-only", false, "ignore all namespace-scoped resources")
	flags.BoolVar(&versionFlag, "version", false, "print version information and exit")

	flags.BoolVar(&s.Extra, "extra", false, "Flag to enable collecting extra resources")
	flags.MarkHidden("extra")

	clientcmd.BindOverrideFlags(&s.ConfigOverrides, flags, clientcmd.RecommendedConfigOverrideFlags(""))

	return cmd
}
