# iocc_proxy
# 背景

相信有一部分人喜爱 GO 的初衷大概是：**跨平台静态编译**，如果在没用通过 CGO 引用其他库的话，一般编译出来的可执行二进制文件部署起来非常方便，但人们在实际中发现，使用 Go 语言开发的后端 WEB 程序存在 HTML 模版、图片、JS、CSS、JSON 等静态资源，部署时需要把这些静态资源与二进制程序一起上传到服务器部署，在现今遍地花容器的今天，为了简化部署流程，能不能更进一步的把这些静态文件与二进制程序一起打包起来部署，那样岂不是更方便？对运维来说打包一体化带来的好处就是运维复杂度的降低，对技术团队来说打包一体化还可以保证程序完整性可控性。因此，GO 社区发起了一个期望 Go 编译器支持嵌入静态文件的提案[issue#35950](https://github.com/golang/go/issues/35950)。现在，这个功能将随着 1.16 版本一起发布，目前最新的版本是 Go 1.16 RC1 预览版。

# embed 嵌入

```
└── cmd 测试目录
    ├── assets 静态资源目录
    │   ├── .idea.txt
    │   ├── golang.txt
    │   └── hello.txt
    └── main.go  测试go源文件
```

### 字符串、字节切片、文件嵌入

```
package main

import (
	"embed"
	_ "embed"
	"fmt"
)

//go:embed指令用来嵌入，必须紧跟着嵌入后的变量名
//只支持嵌入为string, byte slice和embed.FS三种类型，这三种类型的别名(alias)和命名类型(如type S string)都不可以

//以字符串形式嵌入 assets/hello.txt
//go:embed assets/hello.txt
var s string

//文件的内容嵌入为slice of byte，也就是一个字节数组
//go:embed assets/hello.txt
var b []byte

//嵌入为一个文件系统 新的文件系统FS
//go:embed assets/hello.txt
//go:embed assets/golang.txt
var f embed.FS

func main() {
	fmt.Println("embed string.", s)
	fmt.Println("embed byte.", string(b))

	data, _ := f.ReadFile("assets/hello.txt")
	fmt.Println("embed fs.", string(data))

	data, _ = f.ReadFile("assets/golang.txt")
	fmt.Println("embed fs.", string(data))
}

```

编译运行后输出：

```
embed string. hello golang!
embed byte. hello golang!
embed fs. hello golang!
embed fs. hello!

```

从上面的代码可以看出，embed 支持嵌入为 string,byte slice 和 embed.FS 这三种类型，另外也不允许从这些类型中派生。

### 嵌入文件

对于 FS 类型的嵌入，也可以支持一个变量嵌入多个文件。

```
//go:embed assets/hello.txt
//go:embed assets/golang.txt
var f embed.FS
```

当然也支持，两个变量嵌入一个文件。虽然两个变量嵌入了同一个文件，但它们在编译的时候会独立分配，彼此之间并不会互相影响。

### 嵌入文件夹

FS 类型的嵌入支持文件夹.分隔符采用正斜杠/,即使是 windows 系统也采用这个模式。

```
//go:embed assets
var f embed.FS

func main() {
	data, _ := f.ReadFile("assets/hello.txt")
	fmt.Println(string(data))
}
```

### 嵌入匹配

go:embed 指令中可以只写文件夹名，此文件夹中除了.和_开头的文件和文件夹都会被嵌入，并且子文件夹也会被递归的嵌入，形成一个此文件夹的文件系统。

如果想嵌入.和_开头的文件和文件夹， 比如.hello.txt 文件，那么就需要使用*，比如 go:embed assets/*。

*不具有递归性，所以子文件夹下的.和_不会被嵌入，除非你在专门使用子文件夹的*进行嵌入:



```
├── assets
│   ├── .idea.txt
│   ├── golang.txt
│   └── hello.txt
└── main.go

package main

import (
	"embed"
	_ "embed"
	"fmt"
)

//go:embed assets/*
var f embed.FS

func main() {
	data, _ := f.ReadFile("assets/.idea.txt")
	fmt.Println(string(data))
}
```

### FS 文件系统

embed.FS 实现了 io/fs.FS 接口，它可以打开一个文件，返回 fs.File:

```
package main

import (
	"embed"
	_ "embed"
	"fmt"
)

//go:embed assets
var f embed.FS

func main() {
	dirEntries, _ := f.ReadDir("assets")
	for _, de := range dirEntries {
		fmt.Println(de.Name(), de.IsDir())
	}
}
```

它还提供了 ReadFileh 和 ReadDir 功能，遍历一个文件下的文件和文件夹信息：

```
package main

import (
	"embed"
	_ "embed"
	"fmt"
)

//go:embed assets
var f embed.FS

func main() {
	dirEntries, _ := f.ReadDir("assets")
	for _, de := range dirEntries {
		fmt.Println(de.Name(), de.IsDir())
	}
}
```

因为它实现了 io/fs.FS 接口，所以可以返回它的子文件夹作为新的文件系统：

```
package main

import (
	"embed"
	_ "embed"
	"fmt"
	"io/fs"
	"io/ioutil"
)

//go:embed assets
var f embed.FS

func main() {
	as, _ := fs.Sub(f, "assets")
	hi, _ := as.Open("hello.txt")
	data, _ := ioutil.ReadAll(hi)
	fmt.Println(string(data))
}
```

### 总结：

- 对于单个的文件，支持嵌入为字符串和 byte slice
- 对于多个文件和文件夹，支持嵌入为新的文件系统 FS
- go:embed 指令用来嵌入，必须紧跟着嵌入后的变量名
- 只支持嵌入为 string, byte slice 和 embed.FS 三种类型，类型派生也不可以。