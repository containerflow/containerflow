package types

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
