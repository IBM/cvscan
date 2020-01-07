package main

import (
	"errors"

	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"

	"github.ibm.com/certauto/cvscan/internal/scan"
)

type scanCmd struct {
	opts            metav1.ListOptions
	extra           bool
	clusterWideOnly bool
	configOverrides clientcmd.ConfigOverrides
	kubeconfig      string
}

func newScanCmd() *cobra.Command {
	s := &scanCmd{}

	cmd := &cobra.Command{
		Use:   "cvscan output_path",
		Short: "Take a snapshot of live kubernetes resources",
		RunE: func(cmd *cobra.Command, args []string) error {
			return s.run(args)
		},
	}

	flags := cmd.Flags()
	flags.StringVarP(&s.opts.LabelSelector, "selector", "l", "", "Selector (label query) to filter on")
	flags.StringVar(&s.opts.FieldSelector, "field-selector", "", "Selector (field query) to filter on")
	flags.StringVar(&s.kubeconfig, "kubeconfig", "", "Path to the kubeconfig file to use for CLI requests")
	flags.BoolVar(&s.clusterWideOnly, "cluster-wide-only", false, "ignore all namespace-scoped resources")

	flags.BoolVar(&s.extra, "extra", false, "Flag to enable collecting extra resources")
	flags.MarkHidden("extra")

	clientcmd.BindOverrideFlags(&s.configOverrides, flags, clientcmd.RecommendedConfigOverrideFlags(""))

	return cmd
}

func (sCmd *scanCmd) run(args []string) error {
	if len(args) == 0 {
		return errors.New("output_path argument is required")
	}
	outputPath := args[0]

	config, err := sCmd.getClusterConfig()
	if err != nil {
		return err
	}

	s, err := scan.New(config, sCmd.clusterWideOnly)
	if err != nil {
		return err
	}

	if sCmd.extra {
		err := s.ListAll("", sCmd.opts, outputPath)
		if err != nil {
			return err
		}

		err = s.ListAll(sCmd.configOverrides.Context.Namespace, metav1.ListOptions{}, outputPath)
		if err != nil {
			return err
		}
	}

	return s.ListAll(sCmd.configOverrides.Context.Namespace, sCmd.opts, outputPath)
}

func (sCmd *scanCmd) getClusterConfig() (*rest.Config, error) {
	// First, try in-cluster
	config, err := rest.InClusterConfig()
	if err == nil {
		return config, nil
	} else if err == rest.ErrNotInCluster {
		// Next, try out-of-cluster
		loadingRules := clientcmd.NewDefaultClientConfigLoadingRules()
		loadingRules.ExplicitPath = sCmd.kubeconfig
		kubeConfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(loadingRules, &sCmd.configOverrides)
		return kubeConfig.ClientConfig()
	} else {
		return nil, err
	}
}
