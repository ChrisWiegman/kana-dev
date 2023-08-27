package docker

/**
 * Docker code examples currently from https://willschenk.com/articles/2021/controlling_docker_in_golang/
 **/

import (
	"context"
	"fmt"
	"os/exec"
	"runtime"
	"time"

	"github.com/ChrisWiegman/kana-cli/internal/console"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

var execCommand = exec.Command

var maxRetries = 12
var sleepDuration = 5

// DockerClient is an interface the must be implemented to provide Docker services through this package.
type DockerClient struct {
	apiClient       APIClient
	imageUpdateData ViperClient
	checkedImages   []string
}

func NewDockerClient(consoleOutput *console.Console, appDirectory string) (dockerClient *DockerClient, err error) {
	dockerClient = new(DockerClient)

	dockerClient.apiClient, err = client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return nil, err
	}

	err = ensureDockerIsAvailable(consoleOutput, dockerClient.apiClient)
	if err != nil {
		return nil, err
	}

	dockerClient.imageUpdateData, _ = dockerClient.loadImageUpdateData(appDirectory)

	return dockerClient, nil
}

func ensureDockerIsAvailable(consoleOutput *console.Console, apiClient APIClient) error {
	_, err := apiClient.ContainerList(context.Background(), types.ContainerListOptions{})
	if err != nil {
		if runtime.GOOS == "darwin" { //nolint:goconst
			consoleOutput.Println("Docker doesn't appear to be running. Trying to start Docker.")
			err = execCommand("open", "-a", "Docker").Run()
			if err != nil {
				return fmt.Errorf("error: unable to start Docker for Mac")
			}

			retries := 0

			for retries <= maxRetries {
				retries++

				if retries == maxRetries {
					consoleOutput.Println("Restarting Docker is taking too long. We seem to have hit an error")
					return fmt.Errorf("error: unable to start Docker for Mac")
				}

				time.Sleep(time.Duration(sleepDuration) * time.Second)

				_, err = apiClient.ContainerList(context.Background(), types.ContainerListOptions{})
				if err != nil {
					return err
				}
			}
		}
	}

	return err
}
