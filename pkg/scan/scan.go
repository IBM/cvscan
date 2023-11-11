package scan

import (
	"errors"
	"fmt"
	"os"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"

	"github.com/IBM/cvscan/internal/scan"
)

type Scanner struct {
	Opts            metav1.ListOptions
	Extra           bool
	ClusterWideOnly bool
	ConfigOverrides clientcmd.ConfigOverrides
	Kubeconfig      string
}

func NewScanner() *Scanner {
	return &Scanner{}
}

func (sCmd *Scanner) Run(args []string) error {
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

	s, err := scan.New(config, sCmd.ClusterWideOnly)
	if err != nil {
		return fmt.Errorf("initialize scanner: %v", err)
	}

	err = s.WriteCaps(sCmd.ConfigOverrides.Context.Namespace, sCmd.Opts, outputPath)
	if err != nil {
		return fmt.Errorf("writing capabilities: %v", err)
	}

	if sCmd.Extra {
		err := s.ListAll("", sCmd.Opts, outputPath)
		if err != nil {
			return fmt.Errorf("listing resources with options: %v", err)
		}

		err = s.ListAll(sCmd.ConfigOverrides.Context.Namespace, metav1.ListOptions{}, outputPath)
		if err != nil {
			return fmt.Errorf("listing resources in namespace: %v", err)
		}
	} else {
		err = s.ListAll(sCmd.ConfigOverrides.Context.Namespace, sCmd.Opts, outputPath)
		if err != nil {
			return fmt.Errorf("listing resources: %v", err)
		}
	}
	return nil
}

func (sCmd *Scanner) getClusterConfig() (*rest.Config, error) {

	//use kubeconfig flag if present
	if sCmd.Kubeconfig != "" {
		loadingRules := clientcmd.NewDefaultClientConfigLoadingRules()
		loadingRules.ExplicitPath = sCmd.Kubeconfig
		kubeConfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(loadingRules, &sCmd.ConfigOverrides)
		return kubeConfig.ClientConfig()
	}
	//if no kubeconfig flag, check if in cluster
	//otherwise load from default config path
	config, err := rest.InClusterConfig()
	if err == nil {
		return config, nil
	} else if err == rest.ErrNotInCluster {
		loadingRules := clientcmd.NewDefaultClientConfigLoadingRules()
		kubeConfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(loadingRules, &sCmd.ConfigOverrides)
		return kubeConfig.ClientConfig()
	} else {
		return nil, err
	}

}
