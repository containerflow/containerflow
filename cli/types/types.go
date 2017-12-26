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
func (pipeline Pipeline) Execute(root string, number int) (err error) {
	fmt.Println("## Set Build Environment")
	services := make(map[string]string)

	for _, service := range pipeline.Spec.Services {
		fmt.Println("--> Start service " + service.Name)
		cid, err := service.start(pipeline.Metadata.Namespace, pipeline.Metadata.Name)
		if err == nil {
			services[service.Name] = cid
		} else {
			log.Fatalln(err)
		}
	}

	fmt.Println("## Start Build")
	for _, stage := range pipeline.Spec.Stages {
		fmt.Println("--> Stage:" + stage.Name + " Process")
		stage.execute(pipeline.Metadata.Namespace, pipeline.Metadata.Name, root, services)
	}

	fmt.Println("## CleanUp Build Environment")

	for name, cid := range services {
		fmt.Println("--> Stop Service " + name)
		forceRemoveContainer(cid)
	}

	return
}

func (stage Stage) execute(namespace, name, root string, links map[string]string) {
	var size = len(stage.Tasks)
	stage.cnum = make(chan int, size)
	for _, task := range stage.Tasks {
		go task.execute(namespace, name, stage.Name, root, stage.Workspace, stage.cnum, links)
	}
	for i := 0; i < size; i++ {
		<-stage.cnum
	}
}

func (task Task) execute(namespace string, name string, stage string, root string, workspace string, cnum chan int, links map[string]string) {
	executeContainer(namespace, name, stage, root, workspace, task, links)
	cnum <- 1
}

func (service Service) start(namespace, name string) (cid string, err error) {
	containerName := fmt.Sprintf("%s-%s-%s", namespace, name, service.Name)
	containerName = strings.Replace(containerName, " ", "_", -1)

	cli, err := client.NewEnvClient()
	if err != nil {
		return
	}

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

	cid, err = startContainer(&container.Config{
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

func executeContainer(namespace string, name string, stage string, root string, workspace string, task Task, linkMap map[string]string) bool {
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

	if workspace == "" {
		workspace = "/workspace"
	}

	entrypoint := []string{
		"sh",
		"-c",
		strings.Join(task.Commands, "&&"),
	}

	links := make([]string, len(linkMap))

	for name, container := range linkMap {
		links = append(links, fmt.Sprintf("%s:%s", container, name))
	}

	fmt.Println(links)

	cid, err := startContainer(&container.Config{
		Tty:          true,
		Image:        task.Box,
		AttachStdout: true,
		AttachStderr: true,
		WorkingDir:   workspace,
		Entrypoint:   entrypoint,
		Env:          task.Environment,
	}, &container.HostConfig{
		AutoRemove: false,
		// Links:      links,
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
	err = forceRemoveContainer(cid)

	if err != nil {
		log.Fatalln(err)
		return false
	}

	return true
}

func startContainer(config *container.Config, hostConfig *container.HostConfig, containerName string) (containerID string, err error) {
	cli, err := client.NewEnvClient()
	if err != nil {
		log.Fatalln(err)
		return
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

func forceRemoveContainer(cid string) (err error) {
	cli, err := client.NewEnvClient()
	if err != nil {
		log.Fatalln(err)
		return
	}

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
