# go-net
带有stdout和stderr的go-clr库,程序并没有大的修改主要是程序调用和代码变量，测试，封装<br>
# 使用
```
go get github.com/sairson/go-net
```
example
```
package main

import (
	"fmt"
	"github.com/sairson/go-net"
	"io/ioutil"
)

func main(){
	install, err := go_net.RuntimeInstall("v4")
	if err != nil {
		fmt.Println(err)
		return 
	}
	err = go_net.RedirectStdoutAndStderr()
	if err != nil {
		fmt.Println(err)
		return
	}
	body, err := ioutil.ReadFile("F:\\tools\\go-net\\example\\SharpSQLTools.exe")
	if err != nil {
		fmt.Println(err)
		return
	}
	assembly, err := go_net.LoadAssembly(install, body)
	if err != nil {
		fmt.Println(err)
		return
	}
	stdout,stderr := go_net.InvokeAssembly(assembly,[]string{"clr_pwd"})
	fmt.Println(stdout,stderr)
	_,_ = stderr,stdout

}
```
# 参考和原版
```
https://github.com/ropnop/go-clr
https://github.com/Ne0nd0g/go-clr
https://blog.ropnop.com/hosting-clr-in-golang/
```
