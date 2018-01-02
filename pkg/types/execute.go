package types

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
)

// Execute Start pipeline process
func (pipeline Pipeline) Execute(root string, number int) (err error) {
	fmt.Println("## Set Build Environment")

	cli, err := client.NewEnvClient()
	if err != nil {
		log.Fatalln(err)
		return
	}

	defer cli.Close()

	services := []string{}
	links := []string{}

	for _, service := range pipeline.Spec.Services {
		fmt.Println("--> Start service " + service.Name)
		cid, err := service.start(cli, pipeline.Metadata.Namespace, pipeline.Metadata.Name)
		if err == nil {
			services = append(services, cid)
			links = append(links, fmt.Sprintf("%s:%s", cid, service.Name))
		} else {
			log.Fatalln(err)
		}
	}

	fmt.Println("## Start Build")
	fmt.Println(links, "links", len(links))

	for _, stage := range pipeline.Spec.Stages {
		fmt.Println("--> Stage:" + stage.Name + " Process")
		stage.execute(cli, pipeline.Metadata.Namespace, pipeline.Metadata.Name, root, links)
	}

	fmt.Println("## CleanUp Build Environment")

	for _, cid := range services {
		forceRemoveContainer(cli, cid)
	}

	return
}

func (stage Stage) execute(cli *client.Client, namespace, name, root string, links []string) {
	var size = len(stage.Tasks)
	stage.cnum = make(chan int, size)
	for _, task := range stage.Tasks {
		go task.execute(cli, namespace, name, stage.Name, root, stage.Workspace, stage.cnum, links)
	}
	for i := 0; i < size; i++ {
		<-stage.cnum
	}
}

func (task Task) execute(cli *client.Client, namespace string, name string, stage string, root string, workspace string, cnum chan int, links []string) {
	executeContainer(cli, namespace, name, stage, root, workspace, task, links)
	cnum <- 1
}

func (service Service) start(cli *client.Client, namespace, name string) (cid string, err error) {
	containerName := fmt.Sprintf("%s-%s-%s", namespace, name, service.Name)
	containerName = strings.Replace(containerName, " ", "_", -1)

	_, err = cli.ImagePull(
		context.Background(),
		service.Box,
		types.ImagePullOptions{
			All: true,
		},
	)

	if err != nil {
		return
	}

	cid, err = startContainer(cli, &container.Config{
		Tty:          true,
		Image:        service.Box,
		AttachStdout: true,
		AttachStderr: true,
		Env:          service.Environment,
	}, &container.HostConfig{
		AutoRemove: false,
	}, containerName)

	if err != nil {
		return
	}

	return
}

func executeContainer(cli *client.Client, namespace string, name string, stage string, root string, workspace string, task Task, links []string) bool {
	containerName := fmt.Sprintf("%s-%s-%s-%s", namespace, name, stage, task.Name)
	containerName = strings.Replace(containerName, " ", "_", -1)

	if workspace == "" {
		workspace = "/workspace"
	}

	entrypoint := []string{
		"sh",
		"-c",
		strings.Join(task.Commands, "&&"),
	}

	cid, err := startContainer(cli, &container.Config{
		Tty:          true,
		Image:        task.Box,
		AttachStdout: true,
		AttachStderr: true,
		WorkingDir:   workspace,
		Entrypoint:   entrypoint,
		Env:          task.Environment,
	}, &container.HostConfig{
		AutoRemove: false,
		Links:      links,
		Mounts: []mount.Mount{
			mount.Mount{
				Type:   mount.TypeBind,
				Source: root,
				Target: workspace,
			},
		},
	}, containerName)

	if err != nil {
		log.Fatalln(err)
		forceRemoveContainer(cli, cid)
		return false
	}

	reader, err := cli.ContainerLogs(context.Background(), cid, types.ContainerLogsOptions{
		ShowStdout: true,
	})

	if err != nil {
		log.Fatalln(err)
		return false
	}

	b, err := ioutil.ReadAll(reader)

	if err != nil {
		log.Fatalln(err)
		return false
	}

	fmt.Println(string(b))
	err = forceRemoveContainer(cli, cid)

	if err != nil {
		log.Fatalln(err)
		return false
	}

	return true
}

func startContainer(cli *client.Client, config *container.Config, hostConfig *container.HostConfig, containerName string) (containerID string, err error) {
	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{
		All: true,
	})

	for _, container := range containers {
		for _, name := range container.Names {
			if name == fmt.Sprintf("/%s", containerName) {
				forceRemoveContainer(cli, container.ID)
			}
		}
	}

	c, err := cli.ContainerCreate(context.Background(), config, hostConfig, &network.NetworkingConfig{}, containerName)
	if err != nil {
		log.Fatalln(err)
		return
	}

	err = cli.ContainerStart(context.Background(), c.ID, types.ContainerStartOptions{})

	if err != nil {
		log.Fatalln(err)
		return
	}

	return c.ID, err
}

func forceRemoveContainer(cli *client.Client, cid string) (err error) {
	err = cli.ContainerRemove(context.Background(), cid, types.ContainerRemoveOptions{
		Force: true,
	})
	return
}

// checkError args[0] force to interrupt programe
func checkError(err error, args ...interface{}) (isError bool) {
	if err != nil {
		isError = true
		log.Fatalln(err)
	}
	return
}
