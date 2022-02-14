package go_net

import (
	"fmt"
	"strings"
	"syscall"
	"unsafe"
)

// ICORRuntimeHost  ICORRuntimeHost对象结构
type ICORRuntimeHost struct {
	lpVtbl *ICORRuntimeHostVtbl
}

type ICORRuntimeHostVtbl struct {
	QueryInterface                uintptr
	AddRef                        uintptr
	Release                       uintptr
	CreateLogicalThreadState      uintptr
	DeleteLogicalThreadState      uintptr
	SwitchInLogicalThreadState    uintptr
	SwitchOutLogicalThreadState   uintptr
	LocksHeldByLogicalThreadState uintptr
	MapFile                       uintptr
	GetConfiguration              uintptr
	Start                         uintptr
	Stop                          uintptr
	CreateDomain                  uintptr
	GetDefaultDomain              uintptr
	EnumDomains                   uintptr
	NextDomain                    uintptr
	CloseEnum                     uintptr
	CreateDomainEx                uintptr
	CreateDomainSetup             uintptr
	CreateEvidence                uintptr
	UnloadDomain                  uintptr
	CurrentDomain                 uintptr
}


func (obj *ICORRuntimeHost) AddRef() uintptr {
	ret, _, _ := syscall.Syscall(
		obj.lpVtbl.AddRef,
		1,
		uintptr(unsafe.Pointer(obj)),
		0,
		0)
	return ret
}

func (obj *ICORRuntimeHost) Release() uintptr {
	ret, _, _ := syscall.Syscall(
		obj.lpVtbl.Release,
		1,
		uintptr(unsafe.Pointer(obj)),
		0,
		0)
	return ret
}

func (obj *ICORRuntimeHost)QueryInterface() uintptr {
	ret,_,_ := syscall.Syscall(
		obj.lpVtbl.QueryInterface,
		1,
		uintptr(unsafe.Pointer(obj)),
		0,
		0)
	return ret
}


// Start 运行程序相关runtimes
func (obj *ICORRuntimeHost) Start() uintptr {
	ret, _, _ := syscall.Syscall(
		obj.lpVtbl.Start,
		1,
		uintptr(unsafe.Pointer(obj)),
		0,
		0)
	return ret
}

// GetDefaultDomain 获取程序默认域
func (obj *ICORRuntimeHost) GetDefaultDomain() (IUnknown *IUnknown, err error) {
	hr, _, err := syscall.Syscall(
		obj.lpVtbl.GetDefaultDomain,
		2,
		uintptr(unsafe.Pointer(obj)),
		uintptr(unsafe.Pointer(&IUnknown)),
		0,
	)
	if err != syscall.Errno(0) && !strings.Contains(err.Error(),"specified procedure could not be found") {
		// 会返回一个错误
		return
	}
	if hr != 0x00 {
		err = fmt.Errorf("the ICORRuntimeHost::GetDefaultDomain method method returned a non-zero HRESULT: 0x%x", hr)
		return
	}
	err = nil
	return
}

func GetICORRuntimeHost(runtimeInfo *ICLRRuntimeInfo) (*ICORRuntimeHost, error) {
	var runtimeHost *ICORRuntimeHost
	err := runtimeInfo.GetInterface(CLSID_CorRuntimeHost, IID_ICorRuntimeHost, unsafe.Pointer(&runtimeHost))
	if err != nil {
		return nil, err
	}

	_ = runtimeHost.Start()
	return runtimeHost, err
}