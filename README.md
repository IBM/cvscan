[![Build Status](https://app.travis-ci.com/IBM/cvscan.svg?branch=master)](https://app.travis-ci.com/IBM/cvscan)

## cvscan

`cvscan` is a utility that takes snapshots of resources running in a Kubernetes
cluster.

## Build
1. Prerequisites: [`go` 1.19+](https://golang.org/dl/),
   [`make`](https://www.gnu.org/software/make/)
1. `make build` will build and put a `cvscan` executable in the current directory or `make
   install` to build the executable and put it in `$GOBIN`.

## Usage

    $ cvscan help
    Take a snapshot of live kubernetes resources

    Usage:
    cvscan output_path [flags]

    Flags:
        --as string                      Username to impersonate for the operation
        --as-group stringArray           Group to impersonate for the operation, this flag can be repeated to specify multiple groups.
        --certificate-authority string   Path to a cert file for the certificate authority
        --client-certificate string      Path to a client certificate file for TLS
        --client-key string              Path to a client key file for TLS
        --cluster string                 The name of the kubeconfig cluster to use
        --cluster-wide-only              ignore all namespace-scoped resources
        --context string                 The name of the kubeconfig context to use
        --field-selector string          Selector (field query) to filter on
    -h, --help                           help for cvscan
        --insecure-skip-tls-verify       If true, the server's certificate will not be checked for validity. This will make your HTTPS connections insecure
        --kubeconfig string              Path to the kubeconfig file to use for CLI requests
    -n, --namespace string               If present, the namespace scope for this CLI request
        --password string                Password for basic authentication to the API server
        --request-timeout string         The length of time to wait before giving up on a single server request. Non-zero values should contain a corresponding time unit (e.g. 1s, 2m, 3h). A value of zero means don't timeout requests. (default "0")
    -l, --selector string                Selector (label query) to filter on
        --server string                  The address and port of the Kubernetes API server
        --token string                   Bearer token for authentication to the API server
        --user string                    The name of the kubeconfig user to use
        --username string                Username for basic authentication to the API server

The only required argument to `cvscan` is an output directory where resource YAML defintions will be put. Output files are named as `scanned-{kind}-{namespace}-{name}.yaml`. A `caps.json` file is also generated in the output directory which contains other cluster information. Files are overwritten if they already exist.

By default, `cvscan` will collect every resource in the cluster. Results can be filtered by field or label selectors with the `--field-selector` and `-l/--selector` flags, which accept the same syntax as `kubectl`.

A namespace may also be specified with `-n/--namespace`. When a namespace is specified, non-namespaced resources like Nodes and PersistentVolumes are not captured. Cluster-scoped resources can be captured in isolation with the `--cluster-wide-only` flag.

Filtering options are strictly additive, so specifying both a label selector and a namespace will only collect resources that match the selector AND are in the namespace. A logical OR may be implemented by calling `cvscan` multiple times with the same output directory.

Some resources are filtered automatically as they have been classified as noise. This behavior is not configurable.
