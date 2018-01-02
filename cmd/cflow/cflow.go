package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	model "github.com/containerflow/containerflow/pkg/types"
	yaml "gopkg.in/yaml.v2"
)

func check(e error) {
	if e != nil {
		log.Fatalf("error: %v", e)
		panic(e)
	}
}

func getCurrentDirectory(config string) string {
	dir, err := filepath.Abs(filepath.Dir(config))
	check(err)
	return strings.Replace(dir, "\\", "/", -1)
}

func main() {

	var configFile = "./cflow.yml"

	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		fmt.Println("cflow.yml not found")
		return
	}

	maxProcs := runtime.NumCPU()
	runtime.GOMAXPROCS(maxProcs)

	root := getCurrentDirectory(configFile)
	fmt.Println("Start Container FLow At Workspace:", root)

	data, err := ioutil.ReadFile(configFile)

	check(err)

	pipeline := model.Pipeline{}
	err = yaml.Unmarshal([]byte(data), &pipeline)
	check(err)

	var buildNumber = 0

	os.Mkdir(".workspace", os.ModePerm)
	os.Mkdir(".workspace/builds", os.ModePerm)

	workPath := fmt.Sprintf(".workspace/builds/%v", buildNumber)
	os.Mkdir(workPath, os.ModePerm)

	pipeline.Execute(root, buildNumber)
}
