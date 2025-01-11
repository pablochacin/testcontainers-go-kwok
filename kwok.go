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


// WithNodes is a container customizer that sets the number of nodes in the cluster
func WithNodes(nodes int) testcontainers.CustomizeRequestOption {
	return func(req *testcontainers.GenericContainerRequest) error {
		req.LifecycleHooks = append(req.LifecycleHooks, testcontainers.ContainerLifecycleHooks{
			PostReadies: []testcontainers.ContainerHook{
				func(ctx context.Context, c testcontainers.Container) error {
					_, _, err := c.Exec(ctx, []string{"kwokctl", "scale", "node", "--replicas", fmt.Sprintf("%d", nodes)})
					return err
				},
			},
		})

		return nil
	}
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
	
	return c, err
}

