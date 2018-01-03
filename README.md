ContainerFlow
==============

> The Experiment Project

## Required

* cli: containerFlow command interface
    * docker
* server: containerFlow API
    * mongodb
* website: react base web ui

## How To Develope

Start Denpendency

```
docker-compose up -d
```

## Example

Install

```
cd cmd/cflow
go install

export $PATH=$GOPATH/bin:$PATH
```

Start Pipeline

```
$ cflow
Start Container FLow At Workspace: /Users/zhengyunlong/Workspace/go/src/github.com/containerflow/containerflow
## Set Build Environment
--> Start service database
--> Start service cache
## Start Build
[825c5715f8a48127bd35c0bcc47eda53473cd7973c46a7ceb06b4d824601d42e:database 647de7d799a567075b24bf795693a3c0c9d481596ef0f7e2e98ee1fbfac3fbdd:cache] links 2
--> Stage:stage1 Process
...
--> Stage:stage2 Process
...
## CleanUp Build Environment
```

