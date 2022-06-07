# 快速上手

本教程将构建一个基于 `Golang` 的 `Todo` 后端服务，提供标准的 `Graphql API` 以供前端应用调用。在我们开始之前，请确保您的机器上满足了以下前提条件。

## 前提条件

- [Make](setup-make.md)
- [Golang 1.16+](https://golang.org/doc/install)
- [Mysql 8.0+](https://dev.mysql.com/doc/refman/8.0/en/installing.html)
- [cre 0.1.10+](setup-local.md) 
- [Wire](https://github.com/google/wire)
- [gowatch](https://github.com/silenceper/gowatch) (可选)

## 数据字典

设计应用的数据字典 `todo.sql` 如下：

```sql
-- Create a database
CREATE DATABASE `todo` DEFAULT CHARACTER SET = `utf8mb4`;

USE `todo`;
DROP TABLE IF EXISTS `todos`;
CREATE TABLE `todos` (
  `id` char(36) NOT NULL,
  `title` varchar(32) NOT NULL DEFAULT '' COMMENT '标题',
  `completed` tinyint(1) NOT NULL DEFAULT '0' COMMENT '是否完成',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8 COMMENT = 'Todo';

```

导入 `Mysql`

```bash
mysql -h <host> -u <user> < mysql.sql
```
## 创建项目

以 Linux 下 shell 命令为例

## 生成代码

- 下载 `grpc` 模板

```bash
mkdir -p $HOME/.cre/contrib
git clone https://github.com/ychengcloud/contrib $HOME/.cre/contrib

```

# 创建项目目录

```bash
mkdir todo

cd todo

# 复制模板文件到项目目录
cp -r $HOME/.cre/contrib/grpc/skeleton/* .

## 数据字典
数据字典文件，命名为 schema.sql 置于 database 目录下

```

新建 `cre.yaml` 内容如下：

```yaml

# 项目名称
project: todo
package: "github.com/ychengcloud/todo"
# 数据库配置
dsn: "mysql://<user>:<password>@tcp(127.0.0.1:3306)/todo?charset=utf8mb4"
# 模板根路径
root: "$HOME/.cre/contrib/grpc/templates"
# 模板生成根路径
genRoot: "./"

# Golang 的默认变量标识符与模板中的变量标识符相同时，需要修改成不同的
#delim:
#  left: "@@"
#  right: "@@"

# NameFormat 目标路径
# Path 模板路径名，以 root 为相对路径
# genPath 模板生成路径, 以 genRoot 为相对路径
# format 文件名，支持模板变量
# mode 生成模式， single: 单文件， multi: 多文件
templates:
  - path: "proto/v1/api.proto.tmpl"
    genPath: "proto/todo/v1"
    format: "api.proto"
    mode: "single"

  - path: "go.mod.tmpl"
    genPath: "./"
    format: "go.mod"

  - path: "buf.gen.yaml.tmpl"
    genPath: "./"
    format: "buf.gen.yaml"

  - path: "cmd/main.go.tmpl"
    genPath: "cmd"
    format: "main.go"

  - path: "cmd/injector.go.tmpl"
    genPath: "cmd"
    format: "injector.go"

  - path: "models/model.go.tmpl"
    genPath: "models"
    format: "{{.Table}}.go"
    mode: "multi"

  - path: "repositories/repository.go.tmpl"
    genPath: "repositories"
    format: "repository.go"

  - path: "repositories/repository_test.go.tmpl"
    genPath: "repositories"
    format: "repository_test.go"

  - path: "repositories/gorm/gorm.go.tmpl"
    genPath: "repositories/gorm"
    format: "gorm.go"
  - path: "repositories/gorm/gorm_test.go.tmpl"
    genPath: "repositories/gorm"
    format: "gorm_test.go"
  - path: "repositories/gorm/repository.go.tmpl"
    genPath: "repositories/gorm"
    format: "{{.Table}}.go"
    mode: "multi"
  - path: "repositories/gorm/repository_test.go.tmpl"
    genPath: "repositories/gorm"
    format: "{{.Table}}_test.go"
    mode: "multi"

  - path: "server/server.go.tmpl"
    genPath: "server"
    format: "server.go"

  - path: "services/base.go.tmpl"
    genPath: "services"
    format: "service.go"

  - path: "services/base_test.go.tmpl"
    genPath: "services"
    format: "service_test.go"

  - path: "services/service.go.tmpl"
    genPath: "services"
    format: "{{.Table}}.go"
    mode: "multi"

  - path: "services/service_test.go.tmpl"
    genPath: "services"
    format: "{{.Table}}_test.go"
    mode: "multi"

  - path: "test/data.go.tmpl"
    genPath: "test"
    format: "data.go"
  
  - path: "test/e2e/gorm/gorm_test.go.tmpl"
    genPath: "test/e2e/gorm"
    format: "gorm_test.go"
    
# 数据表配置
tables:
  - name: "todos"
    fields:
    - name: "id"
      required: true
      filterable: true
      operations: ["Eq", "In"]
    - name: "title"
      required: true
      filterable: true
      operations: ["Eq", "In"]
  
```


- 生成代码

```bash

make gen
make install
make mock
make proto
go mod tidy
go test ./...

```

## 运行前配置

config.default.yaml 为全部可配置项, 默认新建 config.yaml，只需要添加需要覆盖的配置项即可。

配置模板内容如下：

```yaml
app:
  name: 
  # 运行模式 1. debug 2. release， 默认 release
  mode: debug
  # 是否能访问Api文档, 默认 false
  doc: true
  # 绑定 IP
  host: 127.0.0.1
  # 绑定 Port
  port: 7779


db:
  dialect: mysql
  mysql:
    user: root
    password: ""
    host: "127.0.0.1"
    port: 3306
    name: project
    charset: utf8mb4
    debug: true
logger:
  filename: /tmp/.log
  maxSize: 500
  maxBackups: 3
  maxAge: 3
  level: "debug"
  stdout: false

probes:
  # 是否开启 Kubernetes probes, 默认 false
  enable: false
  readinessPath: /ready
  livenessPath: /live
  port: 8080
prometheus:
  # 是否开启 Prometheus, 默认 false
  enable: false
  path: /metrics
  port: 8003
  checkIntervalSeconds: 10
pprof:
  # 是否能访问Golang Pprof, 默认 false
  enable: false
  port: 6060
tracing: 
  # 是否开启 opentracing, 默认 false
  enable: false
  jaeger:
    serviceName: admin
    logSpans: false
    reporter:
      localAgentHostPort: "jaeger-agent:6831"
    sampler:
      type: const
      param: 1
jwt:
  # dd if=/dev/urandom bs=1 count=32 2>/dev/null | base64 -w 0 | rev | cut -b 2- | rev
  # signingKey: GRuHhzxQm7z0H7jFBHxd0x2UEjvJHgt+286nnJCOHYw
  contextKey: users
  hydraKeysUri: http://localhost:4445/keys/hydra.openid.id-token
  tokenType: bearer
  signingKey: YOUCHENG
  issuer: newx.io
  claimKey: claim
  signingMethod: HS512
  # seconds
  expired: 1000000
oauth:
  endpoint:
    authURL: "http://localhost:4444/oauth2/auth"
    tokenURL: "http://localhost:4444/oauth2/token"
  config:
    redirectURL: "http://localhost:3001/auth/callback"
    clientID: "myclient5"
    clientSecret: "mysecret5"

```

## 运行

```bash
make run
```

## 体验

浏览器打开 [Graphql Playground](http://localhost:7779/api/playground)

#### 创建 Todo

```graphql
# Create
mutation ($input: TodoInput!) {
  todoCreate (input: $input) {
    todo {
      id
      title
      completed
    }
  }
}


# Variables
{
  "input": {
    "title": "My Todo 1",
    "completed": 0
	}
}
```

#### 查询 Todo
```graphql
# Query
query {
  todosOffsetBased {
    edges {
      node {
        id
        title
        completed
      }
      
    }
    totalCount
    pageInfo {
      hasNextPage
      hasPreviousPage
    }
  }
}
```

#### 更新 Todo
```graphql
# Update
mutation($ids: [ID]!, $input: TodoInput!) {
  todoUpdate(ids: $ids, input: $input) {
    count
  }
}

# Variables
{
  "ids": ["595cc895-ca31-4c6f-9897-d376b1ec8bb2"],
  "input": {
    "title": "My Todo 1",
    "completed": 1
	}
}

```

#### 删除 Todo
```graphql
# Delete
mutation {
  todoDelete(ids: ["595cc895-ca31-4c6f-9897-d376b1ec8bb2"]) {
    count
  }
}
```