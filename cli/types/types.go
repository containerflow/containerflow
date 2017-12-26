package types

import (
	"context"
	"fmt"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
)

type Metadata struct {
	Name      string
	Namespace string
	Labels    map[string]string
}

type Task struct {
	Name      string
	Box       string
	Workspace string
	Plugin    string
	Metadata  map[string]string
}

type Stage struct {
	Name      string
	Tasks     []Task
	Workspace string
	cnum      chan int
}

type Spec struct {
	Stages []Stage
}

type Pipeline struct {
	Version  string
	Kind     string
	Metadata Metadata
	Spec     Spec
}

func (pipeline Pipeline) Execute(root string, number int) {
	for _, stage := range pipeline.Spec.Stages {
		fmt.Println("## Stage:" + stage.Name + " Process")
		stage.execute(pipeline.Metadata.Namespace, pipeline.Metadata.Name, root)
		fmt.Println("--> Success")
	}
}

func (stage Stage) execute(namespace, name, root string) {
	var size = len(stage.Tasks)
	stage.cnum = make(chan int, size)
	for _, task := range stage.Tasks {
		go task.execute(namespace, name, stage.Name, root, stage.Workspace, stage.cnum)
	}
	for i := 0; i < size; i++ {
		<-stage.cnum
	}
}

func (task Task) execute(namespace string, name string, stage string, root string, workspace string, cnum chan int) {
	switch task.Plugin {
	case "script":
		executeContainer(namespace, name, stage, root, workspace, task)
	}
	cnum <- 1
}

func executeContainer(namespace string, name string, stage string, root string, workspace string, task Task) bool {
	containerName := fmt.Sprintf("%s-%s-%s-%s", namespace, name, stage, task.Name)
	containerName = strings.Replace(containerName, " ", "_", -1)

	cli, err := client.NewEnvClient()
	if err != nil {
		return false
	}

	_, err = cli.ImagePull(
		context.Background(),
		task.Box,
		types.ImagePullOptions{
			All: true,
		},
	)

	if err != nil {
		fmt.Println(err)
		return false
	}

	fmt.Println("--> Start Container:" + containerName)

	if workspace == "" {
		workspace = "/workspace"
	}

	c, err := cli.ContainerCreate(context.Background(), &container.Config{
		Tty:          true,
		Image:        task.Box,
		AttachStdout: true,
		AttachStderr: true,
		Cmd:          strings.Split(task.Metadata["script"], " "),
		WorkingDir:   workspace,
	}, &container.HostConfig{
		AutoRemove: false,
		Mounts: []mount.Mount{
			mount.Mount{
				Type:   mount.TypeBind,
				Source: root,
				Target: workspace,
			},
		},
	}, &network.NetworkingConfig{}, containerName)

	if err != nil {
		fmt.Println(err)
		return false
	}

	err = cli.ContainerStart(context.Background(), c.ID, types.ContainerStartOptions{})

	if err != nil {
		fmt.Println(err)
		return false
	}

	return true
}
