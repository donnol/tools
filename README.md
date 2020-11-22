# tools

useful tools.

## tbc

安装：`make tbc_install`. (请先安装 make 工具)

使用：

```sh
$ tbc --help
a tool named to be continued

Usage:
  tbc [flags]
  tbc [command]

Available Commands:
  help        Help about any command
  interface   gen struct interface

Flags:
  -h, --help          help for tbc
  -p, --path string   specify import path
  -r, --recursive     recursively process dir from current

Use "tbc [command] --help" for more information about a command.
```

### 生成结构体接口

```go
type M struct {}

func (m M) String() string {
    return "i am m"
}

func (m M) innerMethod() {

}
```

生成接口

```go
type IM interface {
    String() string
}
```

### 替换模块路径

替换源码里的包导入路径，如：

```go
// From
import (
    "github.com/xxx/tools"
)
// To
import (
    "github.com/yyy/tools"
)
```

### 生成mock结构体

## inject

依赖注入

## apitest

api接口测试及文档生成

## dbdoc

数据库文档生成

## cache

练手缓存

## reflectx

反射方法封装

## worker

goroutine工作控制
