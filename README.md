ContainerFlow
==============

## Example

Install

```
cd cli && go build -o ../bin/cflow
cd ../ && export PATH=$PATH:`pwd`/bin
```

Start Pipeline

```
$ cflow
Start Container FLow At Workspace: /go/src/github.com/containerflow/containerflow/cli
## Stage:stage1 Process
--> Start Container:namespaces-name-stage1-go_build
total 12
drwxr-xr-x 6 root root  204 Dec 26 07:45 .
drwxr-xr-x 3 root root 4096 Dec 26 07:46 ..
drwxr-xr-x 3 root root  102 Dec 26 04:38 .workspace
-rw-r--r-- 1 root root  521 Dec 26 07:42 cflow.yml
-rw-r--r-- 1 root root 1120 Dec 26 07:26 main.go
drwxr-xr-x 3 root root  102 Dec 26 03:16 types

--> Success
## Stage:stage2 Process
--> Start Container:namespaces-name-stage2-go_test

--> Success
```