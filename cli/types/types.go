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

type Metadata struct {
	Name      string
	Namespace string
	Labels    map[string]string
}

type Task struct {
	Name        string
	Box         string
	Workspace   string
	Plugin      string
	Commands    []string
	Environment []string
	Metadata    map[string]interface{}
}

type Stage struct {
	Name      string
	Tasks     []Task
	Workspace string
	cnum      chan int
}

type Service struct {
	Name        string
	Box         string
	Environment []string
}

type Spec struct {
	Stages   []Stage
	Services []Service
}

type Pipeline struct {
	Version  string
	Kind     string
	Metadata Metadata
	Spec     Spec
}

// Execute Start pipeline process
func (pipeline Pipeline) Execute(root string, number int) {
	fmt.Println("## Set Environment")

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
	executeContainer(namespace, name, stage, root, workspace, task)
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

	entrypoint := []string{
		"sh",
		"-c",
		strings.Join(task.Commands, "&&"),
	}

	c, err := cli.ContainerCreate(context.Background(), &container.Config{
		Tty:          true,
		Image:        task.Box,
		AttachStdout: true,
		AttachStderr: true,
		WorkingDir:   workspace,
		Entrypoint:   entrypoint,
		Env:          task.Environment,
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
		log.Fatalln(err)
		return false
	}

	err = cli.ContainerStart(context.Background(), c.ID, types.ContainerStartOptions{})

	if err != nil {
		log.Fatalln(err)
		return false
	}

	reader, err := cli.ContainerLogs(context.Background(), c.ID, types.ContainerLogsOptions{
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

	cli.ContainerRemove(context.Background(), c.ID, types.ContainerRemoveOptions{
		Force: true,
	})

	return true
}

// checkError args[0] force to interrupt programe
func checkError(err error, args ...interface{}) (isError bool) {
	if err != nil {
		isError = true
		log.Fatalln(err)
	}
	return
}
