package go_net

import (
	"fmt"
	"golang.org/x/sys/windows"
	"syscall"
	"unsafe"
)

// ICLRMetaHost ICLRMetaHost 结构实现
type ICLRMetaHost struct {
	lpVtbl *ICLRMetaHostVtbl
}

type ICLRMetaHostVtbl struct {
	QueryInterface uintptr
	AddRef uintptr
	Release uintptr
	GetRuntime uintptr
	GetVersionFromFile uintptr
	EnumerateInstalledRuntimes uintptr
	EnumerateLoadedRuntimes uintptr
	RequestRuntimeLoadedNotification uintptr
	QueryLegacyV2RuntimeBinding uintptr
	ExitProcess uintptr
}

// 函数后期无用，不做实例化(甚至可以不声明)
/*
func (obj *ICLRMetaHost)EnumerateLoadedRuntimes(){}
func (obj *ICLRMetaHost)ExitProcess(){}
func (obj *ICLRMetaHost)GetVersionFromFile(){}
func (obj *ICLRMetaHost)RequestRuntimeLoadedNotification(){}
func (obj *ICLRMetaHost)QueryLegacyV2RuntimeBinding(){}
 */

func (obj *ICLRMetaHost) GetRuntime(pwzVersion *uint16, riid windows.GUID) (ppRuntime *ICLRRuntimeInfo, err error) {

	hr, _, err := syscall.Syscall6(
		obj.lpVtbl.GetRuntime,
		4,
		uintptr(unsafe.Pointer(obj)),
		uintptr(unsafe.Pointer(pwzVersion)),
		uintptr(unsafe.Pointer(&IID_ICLRRuntimeInfo)),
		uintptr(unsafe.Pointer(&ppRuntime)),
		0,
		0,
	)

	if err != syscall.Errno(0) {
		err = fmt.Errorf("there was an error calling the ICLRMetaHost::GetRuntime method:%s", err)
		return
	}
	if hr != 0x00 {
		err = fmt.Errorf("the ICLRMetaHost::GetRuntime method returned a non-zero HRESULT: 0x%x", hr)
		return
	}
	err = nil
	return
}

func (obj *ICLRMetaHost)AddRef() uintptr {
	ret,_,_ := syscall.Syscall(
		obj.lpVtbl.AddRef,
		1,
		uintptr(unsafe.Pointer(obj)),
		0,
		0)
	return ret
}

func (obj *ICLRMetaHost)Release() uintptr {
	ret,_,_ := syscall.Syscall(
		obj.lpVtbl.Release,
		1,
		uintptr(unsafe.Pointer(obj)),
		0,
		0)
	return ret
}

func (obj *ICLRMetaHost)QueryInterface() uintptr {
	ret,_ ,_ := syscall.Syscall(
		obj.lpVtbl.QueryInterface,
		1,
		uintptr(unsafe.Pointer(obj)),
		0,
		0)
	return ret
}

//https://docs.microsoft.com/zh-cn/dotnet/framework/unmanaged-api/hosting/iclrmetahost-enumerateinstalledruntimes-method

func (obj *ICLRMetaHost)EnumerateInstalledRuntimes()(ppEnumerator *IEnumUnknown, err error){
	hr, _, err := syscall.Syscall(
		obj.lpVtbl.EnumerateInstalledRuntimes,
		2,
		uintptr(unsafe.Pointer(obj)),
		uintptr(unsafe.Pointer(&ppEnumerator)),
		0,
	)
	if err != syscall.Errno(0) {
		err = fmt.Errorf("there was an error calling the ICLRMetaHost::EnumerateInstalledRuntimes method:\r\n%s", err)
		return
	}
	if hr != 0x00 {
		err = fmt.Errorf("the ICLRMetaHost::EnumerateInstalledRuntimes method returned a non-zero HRESULT: 0x%x", hr)
		return
	}
	err = nil
	return
}