package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"

	"github.ibm.com/certauto/cvscan/internal/scan"
	"github.ibm.com/certauto/cvscan/pkg/version"
)

type scanCmd struct {
	opts            metav1.ListOptions
	extra           bool
	clusterWideOnly bool
	configOverrides clientcmd.ConfigOverrides
	kubeconfig      string
	version         bool
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
	flags.BoolVar(&s.version, "version", false, "print version information and exit")

	flags.BoolVar(&s.extra, "extra", false, "Flag to enable collecting extra resources")
	flags.MarkHidden("extra")

	clientcmd.BindOverrideFlags(&s.configOverrides, flags, clientcmd.RecommendedConfigOverrideFlags(""))

	return cmd
}

func (sCmd *scanCmd) run(args []string) error {
	if sCmd.version {
		fmt.Println(version.Version())
		return nil
	}

	if len(args) == 0 {
		return errors.New("output_path argument is required")
	}
	outputPath := args[0]

	err := os.MkdirAll(outputPath, os.ModePerm)
	if err != nil {
		return fmt.Errorf("creating output directory: %v", err)
	}

	config, err := sCmd.getClusterConfig()
	if err != nil {
		return fmt.Errorf("getting config: %v", err)
	}

	s, err := scan.New(config, sCmd.clusterWideOnly)
	if err != nil {
		return fmt.Errorf("initialize scanner: %v", err)
	}

	err = s.WriteCaps(sCmd.configOverrides.Context.Namespace, sCmd.opts, outputPath)
	if err != nil {
		return fmt.Errorf("writing capabilities: %v", err)
	}

	if sCmd.extra {
		err := s.ListAll("", sCmd.opts, outputPath)
		if err != nil {
			return fmt.Errorf("listing resources with options: %v", err)
		}

		err = s.ListAll(sCmd.configOverrides.Context.Namespace, metav1.ListOptions{}, outputPath)
		if err != nil {
			return fmt.Errorf("listing resources in namespace: %v", err)
		}
	} else {
		err = s.ListAll(sCmd.configOverrides.Context.Namespace, sCmd.opts, outputPath)
		if err != nil {
			return fmt.Errorf("listing resources: %v", err)
		}
	}
	return nil
}

func (sCmd *scanCmd) getClusterConfig() (*rest.Config, error) {
	// First, try in-cluster
	config, err := rest.InClusterConfig()
	//if no error and kubeconfig parm isn't set return config we found
	//otherwise if we aren't in cluster or they do have a kubeconfig we want to use it
	if err == nil && sCmd.kubeconfig == "" {
		return config, nil

	} else if err == rest.ErrNotInCluster || sCmd.kubeconfig != "" {
		// Next, try out-of-cluster
		loadingRules := clientcmd.NewDefaultClientConfigLoadingRules()
		loadingRules.ExplicitPath = sCmd.kubeconfig
		kubeConfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(loadingRules, &sCmd.configOverrides)
		return kubeConfig.ClientConfig()
	} else {
		return nil, err
	}
}
