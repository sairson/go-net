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