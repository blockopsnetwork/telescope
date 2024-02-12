package process

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/blockopsnetwork/telescope/internal/logger"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/spf13/viper"
)

const contextTimeout = 10 * time.Second

var (
	ctx                 = context.Background()
	containerName       = "telescope"
	telescopeConfigFile = fmt.Sprintf("%s/.telescope/agent.yaml", os.Getenv("HOME"))
	telescopeImage      = "grafana/agent:v0.37.2"
	vol                 = fmt.Sprintf("%v:%v", telescopeConfigFile, "/etc/agent-config/agent.yaml")
)
type Docker struct {
	client *client.Client
}

// NewDockerClient create a new docker process that manages a container
// also checks if docker daemon is running.
func NewDockerClient() *Docker {
	cl, err := client.NewClientWithOpts(client.FromEnv)
	if err == nil {
		ctx, cancel := context.WithTimeout(ctx, contextTimeout)
		defer cancel()

		_, err = cl.Ping(ctx)

		if err != nil {
			logger.Log.Fatalf("💥 Docker is not installed or running: %v", err)
		}
	}

	if err != nil {
		log.Fatalf("💥 failed to create Docker client: %v", err)
	}

	logger.Log.Info("✅ Docker daemon is running")

	return &Docker{client: cl}
}

// PullImage: pull docker image from remote
func (cl *Docker) PullImage() {
	logger.Log.Info("📥 pulling image")
	logger.Log.Info("⏳ this may take a while...")
	ctx := ctx

	cl.client.NegotiateAPIVersion(ctx)

	out, err := cl.client.ImagePull(ctx, telescopeImage, types.ImagePullOptions{})
	if err != nil {
		logger.Log.Fatalf("💥 image pull failed: %v", err)
	}

	_, err = io.ReadAll(out)
	if err != nil {
		logger.Log.Fatalf("💥 failed to read image pull output: %v", err)
	}

	logger.Log.Info("✅ image has been pulled successfully.")

}

// StartDockerContainer: starts docker container
func (cl *Docker) StartDockerContainer() {
	logger.Log.Info("⛽ starting process...")
	hostConfig := &container.HostConfig{
		Binds:       []string{vol},
		NetworkMode: "host",
		PidMode:     "host",
		CapAdd:      []string{"SYS_TIME"},
		RestartPolicy: container.RestartPolicy{
			Name: "unless-stopped",
		},
	}

	config := &container.Config{
		Hostname: containerName,
		Image:    telescopeImage,
		Cmd: []string{
				"--config.file=/etc/agent-config/agent.yaml",
				"-config.expand-env",
			},
		Env:      []string{
			"REMOTE_WRITE_URL=https://thanos-receiver.blockops.network/api/v1/receive",
			fmt.Sprintf("%s=%s", "PROJECT_ID", viper.GetString("config.projectId")),
			fmt.Sprintf("%s=%s", "PROJECT_NAME", viper.GetString("config.projectName")),
			fmt.Sprintf("%s=%s", "NETWORK", viper.GetString("config.network")),
			fmt.Sprintf("%s=%s", "TELESCOPE_USERNAME", viper.GetString("config.telescope_username")),
			fmt.Sprintf("%s=%s", "TELESCOPE_PASSWORD", viper.GetString("config.telescope_password")),
		},
	}

	resp, err := cl.client.ContainerCreate(ctx, config, hostConfig, nil, nil, containerName)

	if err != nil {
		logger.Log.Fatalf("💥 failed to create Docker container: %v", err)
	}

	err = cl.client.ContainerStart(ctx, resp.ID, container.StartOptions{})

	if err != nil {
		logger.Log.Fatalf("💥 failed to start Docker container: %v", err)
	}
}

func (cl *Docker) ContainerLogs() {
	reader, err := cl.client.ContainerLogs(ctx, containerName, container.LogsOptions{ShowStdout: true, ShowStderr: true, Follow: true})
	if err != nil {
		logger.Log.Fatalf("💥 failed to read container logs: %v", err)
	}

	_, err = io.Copy(os.Stdout, reader)
	if err != nil && err != io.EOF {
		logger.Log.Fatalf("💥 failed to copy container logs: %v", err)
	}
}
