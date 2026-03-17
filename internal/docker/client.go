package docker

import (
	"context"
	"fmt"
	"strings"

	dockerclient "github.com/fsouza/go-dockerclient"
	"github.com/aymantoumi/portlens/internal/types"
)

type Client struct {
	dc *dockerclient.Client
}

func NewClient() (*Client, error) {
	dc, err := dockerclient.NewClientFromEnv()
	if err != nil {
		return nil, fmt.Errorf("docker: connect: %w", err)
	}
	return &Client{dc: dc}, nil
}

func (c *Client) Close() error {
	return nil
}

func (c *Client) ScanPorts(ctx context.Context) ([]*types.PortEntry, error) {
	list, err := c.dc.ListContainers(dockerclient.ListContainersOptions{
		All: false,
	})
	if err != nil {
		return nil, fmt.Errorf("docker: list containers: %w", err)
	}

	var entries []*types.PortEntry
	for _, ctr := range list {
		for _, port := range ctr.Ports {
			if port.PublicPort == 0 {
				continue
			}
			e := &types.PortEntry{
				Port:           int(port.PublicPort),
				Proto:          "tcp",
				BindIP:         port.IP,
				ContainerName:  containerName(ctr.Names),
				ContainerID:    shortID(ctr.ID),
				ImageName:      ctr.Image,
				ComposeProject: ctr.Labels["com.docker.compose.project"],
				ComposeService: ctr.Labels["com.docker.compose.service"],
				ComposeFile:    ctr.Labels["com.docker.compose.project.config_files"],
			}
			if e.ComposeProject != "" {
				e.Kind = types.SourceDockerCompose
			} else {
				e.Kind = types.SourceDocker
			}
			entries = append(entries, e)
		}
	}
	return entries, nil
}

func containerName(names []string) string {
	if len(names) == 0 {
		return ""
	}
	return strings.TrimPrefix(names[0], "/")
}

func shortID(id string) string {
	if len(id) > 12 {
		return id[:12]
	}
	return id
}
