# sparkle

一些window的工具，用于简化操作。也作为golang的练手项目。

sparkle 是一个快速修改切换系统环境变量的命令行工具。它的实现思路如下：

1. 注册一系列环境变量的信息到一个持久化文件中，并为这些变量信息取一个别名
2. 修改环境变量时，通过这些别名，从而减少修改变量时需要输入的信息，从而快速的切换环境变量
3. 同时，sparkle还支持对环境变量进行分组，从而可以快速切换到某个分组下的环境变量（TODO)
4. 此外，sparkle在修改系统环境变量的同时，会更新当前bash的环境变量，从而无需重启bash（TODO)


## 获取

build:

```shell
go get https://github.com/zhanghuangbin/sparkle
```

## usage

sparkle使用[cobra](https://github.com/spf13/cobra)来解析参数。所以，sparkle也遵循cobra的命令规范，您可以从cobra的文档中，了解更多命令的用法。

```text
sparkle [command]

Available Commands:
alias       别名信息命令
env         修改环境变量，并更新shell的环境变量
help        Help about any command

Flags:
-h, --help           help for sparkle
-s, --store string   store file used by aliasCommand(default is $HOME/.sparkle.yaml)

Use "sparkle [command] --help" for more information about a command.
```

## more

该项目仅仅是本人闲暇之余学习golang写的玩具，不接受issue,也不接受pr。






