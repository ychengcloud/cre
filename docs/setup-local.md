# 安装

## 本地安装

此命令会将 `cre` 安装在 `/usr/local/bin`， 也可通过修改命令行参数安装在其他目录。

请打开终端/控制台窗口，输入如下命令：

```bash
curl -sfL https://raw.githubusercontent.com/ychengcloud/cre/main/scripts/install.sh | sh -s -- -b /usr/local/bin
```

在安装 `cre` 后，你应能在 `PATH` 中找到它。确认是否安装成功：

```bash
cre -v
```

## 命令行用法

- init

初始化命令， 会在当前目录生成项目目录和样例配置文件

- generate

生成命令，根据指定模板生成相应项目

```bash
cre -h

Usage:
  cre [flags]
  cre [command]

Available Commands:
  generate    generate go code for the database schema
  help        Help about any command
  init        initialize framework for project

Flags:
  -h, --help      help for cre
  -v, --version   version for cre
```