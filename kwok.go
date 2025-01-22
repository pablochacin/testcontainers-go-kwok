package kwok

import (
	"bytes"
	"context"
	"fmt"
	"regexp"

	"github.com/docker/go-connections/nat"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/exec"
	"github.com/testcontainers/testcontainers-go/wait"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

const defaultKubeSecurePort = "6443/tcp"

var serverRe = regexp.MustCompile(`https:\/\/127.0.0.1:[0-9]+`)

// KwokContainer represents the Kwok container type used in the module
type KwokContainer struct {
	testcontainers.Container
}

// Run creates an instance of the Kwok container type
func Run(ctx context.Context, image string, opts ...testcontainers.ContainerCustomizer) (*KwokContainer, error) {
	genericContainerReq := testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        image,
			ExposedPorts: []string{"6443/tcp"},
			WaitingFor: wait.ForLog("Cluster is ready"),
			LifecycleHooks: []testcontainers.ContainerLifecycleHooks{
				{
					PostReadies: []testcontainers.ContainerHook{
						func(ctx context.Context, c testcontainers.Container) error {
							_, _, err := c.Exec(ctx, []string{"kwokctl", "scale", "node", "--replicas", "1"})
							return err
						},
					},
				},
			},
		},
		Started: true,
	}

	for _, opt := range opts {
		if err := opt.Customize(&genericContainerReq); err != nil {
			return nil, fmt.Errorf("customize: %w", err)
		}
	}

	container, err := testcontainers.GenericContainer(ctx, genericContainerReq)
	var c *KwokContainer
	if container != nil {
		c = &KwokContainer{Container: container}
	}

	return c, err
}

// GetKubeConfig returns the modified kubeconfig with server url
func (c *KwokContainer) GetKubeConfig(ctx context.Context) ([]byte, error) {
	hostIP, err := c.Host(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get hostIP: %w", err)
	}

	mappedPort, err := c.MappedPort(ctx, nat.Port(defaultKubeSecurePort))
	if err != nil {
		return nil, fmt.Errorf("failed to get mapped port: %w", err)
	}

	_, output, err := c.Exec(ctx, []string{"kwokctl", "get", "kubeconfig"}, exec.Multiplexed())
	if err != nil {
		return nil, fmt.Errorf("failed to get kubeconfig from container: %w", err)
	}

	kubeConfigYaml := &bytes.Buffer{}
	_, err = kubeConfigYaml.ReadFrom(output)
	if err != nil {
		return nil, fmt.Errorf("failed to read file from container: %w", err)
	}

	if !serverRe.MatchString(kubeConfigYaml.String()) {
		return nil, fmt.Errorf("failed to edit kubeconfig: api server url not found")
	}
	modifiedKubeconfig := serverRe.ReplaceAllString(
		kubeConfigYaml.String(),
		fmt.Sprintf("https://%s:%s", hostIP, mappedPort.Port()),
	)

	return []byte(modifiedKubeconfig), nil
}

// GetClient returns a Kubernetes client to access the cluster
func (c *KwokContainer) GetClient(ctx context.Context) (*kubernetes.Clientset, error) {
	kubeConfigYaml , err := c.GetKubeConfig(ctx)
	if err != nil {
		return nil, err
	}

	restcfg, err := clientcmd.RESTConfigFromKubeConfig(kubeConfigYaml)
	if err != nil {
		return nil, fmt.Errorf("failed to create rest config: %w", err)
	}

	k8s, err := kubernetes.NewForConfig(restcfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create client: %w", err)
	}

	return k8s, nil
}