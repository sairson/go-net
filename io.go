package go_net

import (
	"bufio"
	"bytes"
	"fmt"
	"golang.org/x/sys/windows"
	"os"
	"sync"
	"time"
)


//origSTDOUT windows句柄的原始输出
var origSTDOUT = windows.Stdout

// origSTDERR windows句柄的错误原始输出
var origSTDERR = windows.Stderr

// rSTDOUT io.Reader的输出
var rSTDOUT *os.File

// wSTDOUT io.Writer的输出
var wSTDOUT *os.File

// rSTDERR io.Reader的stderr输出
var rSTDERR *os.File

// wSTDERR io.Writer的stderr输出
var wSTDERR *os.File


var Mutex = &sync.Mutex{}
var errors = make(chan error)
var Stdout bytes.Buffer
var Stderr bytes.Buffer


func RedirectStdoutAndStderr()(err error) {
	rSTDOUT, wSTDOUT, err = os.Pipe()
	if err != nil {
		err = fmt.Errorf("there was an error calling the os.Pipe() function to create a new STDOUT:%s", err)
		return
	}
	rSTDERR, wSTDERR, err = os.Pipe()
	if err != nil {
		err = fmt.Errorf("there was an error calling the os.Pipe() function to create a new STDERR:%s", err)
		return
	}
	// 设置一个新的文件对于STDOUT
	if err = windows.SetStdHandle(windows.STD_OUTPUT_HANDLE, windows.Handle(wSTDOUT.Fd())); err != nil {
		err = fmt.Errorf("there was an error calling the windows.SetStdHandle function for STDOUT:%s", err)
		return
	}
	//设置一个新的文件对于STDERR
	if err = windows.SetStdHandle(windows.STD_ERROR_HANDLE, windows.Handle(wSTDERR.Fd())); err != nil {
		err = fmt.Errorf("there was an error calling the windows.SetStdHandle function for STDERR:%s", err)
		return
	}
	go BufferStdout()
	go BufferStderr()

	return err
}



func ReadStdoutStderr()(stdout string,stderr string,err error) {
	// 等待一下，防止没有执行完就读取
	time.Sleep(1 * time.Microsecond)

	// 判断错误是否不为空
	if len(errors) > 0 {
		var totalErrors string
		// 所有错误结果
		for e := range errors {
			totalErrors += e.Error()
		}
		err = fmt.Errorf(totalErrors)
		return
	}

	// 读取stdout的buffer
	if Stdout.Len() > 0 {
		stdout = Stdout.String()
		Stdout.Reset()
	}

	// 读取stderr的buffer
	if Stderr.Len() > 0 {
		stderr = Stderr.String()
		Stderr.Reset()
	}
	return stdout, stderr, err
}

func BufferStdout() {
	stdoutReader := bufio.NewReader(rSTDOUT)
	for {
		// Standard STDOUT buffer size is 4k
		buf := make([]byte, 4096)
		line, err := stdoutReader.Read(buf)
		if err != nil {
			errors <- fmt.Errorf("there was an error reading from STDOUT in io.BufferStdout:%s", err)
		}
		if line > 0 {
			// Remove null bytes and add contents to the buffer
			Stdout.Write(bytes.TrimRight(buf, "\x00"))
		}
	}
}

func BufferStderr() {
	stderrReader := bufio.NewReader(rSTDERR)
	for {
		// Standard STDOUT buffer size is 4k
		buf := make([]byte, 4096)
		line, err := stderrReader.Read(buf)
		if err != nil {
			errors <- fmt.Errorf("there was an error reading from STDOUT in io.BufferStdout:%s", err)
		}
		if line > 0 {
			// Remove null bytes and add contents to the buffer
			Stderr.Write(bytes.TrimRight(buf, "\x00"))
		}
	}
}


func CloseStdoutStderr() (err error) {
	err = rSTDOUT.Close()
	if err != nil {
		err = fmt.Errorf("there was an error closing the STDOUT Reader:%s", err)
		return
	}

	err = wSTDOUT.Close()
	if err != nil {
		err = fmt.Errorf("there was an error closing the STDOUT Writer:%s", err)
		return
	}

	err = rSTDERR.Close()
	if err != nil {
		err = fmt.Errorf("there was an error closing the STDERR Reader:%s", err)
		return
	}

	err = wSTDERR.Close()
	if err != nil {
		err = fmt.Errorf("there was an error closing the STDERR Writer:%s", err)
		return
	}
	return nil
}

func RestoreStdoutStderr() error {
	if err := windows.SetStdHandle(windows.STD_OUTPUT_HANDLE, origSTDOUT); err != nil {
		return fmt.Errorf("there was an error calling the windows.SetStdHandle function to restore the original STDOUT handle:%s", err)
	}
	if err := windows.SetStdHandle(windows.STD_ERROR_HANDLE, origSTDERR); err != nil {
		return fmt.Errorf("there was an error calling the windows.SetStdHandle function to restore the original STDERR handle:%s", err)
	}
	return nil
}