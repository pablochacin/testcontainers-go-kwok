package kwok

import (
	"context"
	"fmt"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)


// Container represents the Kwok container type used in the module
type Container struct {
	testcontainers.Container
}

// Run creates an instance of the Kwok container type
func Run(ctx context.Context, image string, opts ...testcontainers.ContainerCustomizer) (*Container, error) {
	genericContainerReq := testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        image,
			ExposedPorts: []string{"6443/tcp"},
			WaitingFor: wait.ForLog("Cluster is ready"),
		},
		Started: true,
	}

	for _, opt := range opts {
		if err := opt.Customize(&genericContainerReq); err != nil {
			return nil, fmt.Errorf("customize: %w", err)
		}
	}

	container, err := testcontainers.GenericContainer(ctx, genericContainerReq)
	var c *Container
	if container != nil {
		c = &Container{Container: container}
	}

	if err != nil {
		return c, fmt.Errorf("creating container: %w", err)
	}

	// scale cluster to 1 node
	// TODO: make number of nodes configurable with opts
	_, _, err = c.Exec(ctx, []string{"kwokctl", "scale", "node", "--replicas", "1"})
	if err != nil {
		return c, fmt.Errorf("creating node: %w", err)
	}
	return c, nil
}

