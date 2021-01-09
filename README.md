# tools

useful tools.

## tbc

安装：`make tbc_install`. (请先安装 make 工具)

使用：

```sh
$ tbc --help
Usage:
  tbc [flags]
  tbc [command]

Available Commands:
  callgraph   get the callgraph of a function or method
  help        Help about any command
  impl        find implement by given interface in specify path
  interface   gen struct interface
  mock        gen interface mock struct
  replace     replace import path

Flags:
      --depth int          specify depth
      --from string        specify from path with replace
      --func string        specify func or method name
  -h, --help               help for tbc
      --ignore string      specify ignore package
      --interface string   specify interface
  -o, --output string      specify output file
  -p, --path string        specify import path
  -r, --recursive          recursively process dir from current
      --to string          specify to path with replace

Use "tbc [command] --help" for more information about a command.

```

### 生成结构体接口

```sh
gen struct interface, like: 
			type M struct {
				// ...
			}
			func (m *M) String() string {
				return "m.name"
			}
			got: 
			type IM interface {
				String() string
			}
```

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

```sh
gen interface mock struct, like: type I interface { String() string }, 
			gen mock: 
				type Mock struct { StringFunc func() string } 
				var _ I = &Mock{}
				func (mock *Mock) String() string {
					return mock.StringFunc()
				}
			after that, you can use like below:
				var mock = &Mock{
					// init the func like the normal field
					StringFunc: func() string {
						return "jd"								
					},	
				}
				fmt.Println(mock.String())
```

### 找接口实现

```sh
find implement by given interface in specify path, like: 
			'tbc impl --interface=io.Writer'
			will get some structs like
			type MyWriter struct {}
			func (w *MyWriter) Write(data []byte) (n int, err error)
```

### 生成函数/方法调用图

```sh
        path: specify package path
        func：specify function or method, if method, use 'StructName.MethodName', like：A.GetByName
        ignore: ignore package, std means standart packages, others use themself package path
        depth: call depth, use it if you want to skip deep call info

like:
        tbc callgraph --path=xxx.xxx.xxx/a/b --func=[main|normal_func|struct_method] --ignore=std;xxx.xxx.xxx/e/f --depth=2

Usage:
  tbc callgraph [flags]

Flags:
  -h, --help   help for callgraph

Global Flags:
      --depth int          specify depth
      --from string        specify from path with replace
      --func string        specify func or method name
      --ignore string      specify ignore package
      --interface string   specify interface
  -o, --output string      specify output file
  -p, --path string        specify import path
  -r, --recursive          recursively process dir from current
      --to string          specify to path with replace
```

如有以下代码：

```go
func A() {
	B()

	E()
}
func B() {
	C()
}
func C() {
	d := D{}
	d.D()
}
type D struct{}
func (d D) D() {

}
func E() {
	// more call
}
```

生成：

`A -> B -> C -> d.D`

`A -> E`

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
