package go_net

import (
	"fmt"
	"golang.org/x/sys/windows"
	"syscall"
	"unsafe"
)

// ICLRRuntimeInfo ICIRLRuntimeInfo 实例对象结构
type ICLRRuntimeInfo struct {
	lpVtbl *ICLRRuntimeInfoVtbl
}

type ICLRRuntimeInfoVtbl struct {
	QueryInterface uintptr
	AddRef uintptr
	Release uintptr
	GetVersionString uintptr
	GetRuntimeDirectory uintptr
	IsLoaded uintptr
	LoadErrorString uintptr
	LoadLibrary uintptr
	GetProcAddress uintptr
	GetInterface uintptr
	IsLoadable uintptr
	SetDefaultStartupFlags uintptr
	GetDefaultStartupFlags uintptr
	BindAsLegacyV2Runtime uintptr
	IsStarted uintptr
}

func (obj *ICLRRuntimeInfo)AddRef()uintptr {
	ret,_,_ := syscall.Syscall(
		obj.lpVtbl.AddRef,
		1,
		uintptr(unsafe.Pointer(obj)),
		0,
		0)
	return ret
}

func (obj *ICLRRuntimeInfo)Release()uintptr{
	ret,_,_ := syscall.Syscall(
		obj.lpVtbl.Release,
		1,
		uintptr(unsafe.Pointer(obj)),
		0,
		0)
	return ret
}

func (obj *ICLRRuntimeInfo)QueryInterface() uintptr {
	ret,_,_ := syscall.Syscall(
		obj.lpVtbl.QueryInterface,
		1,
		uintptr(unsafe.Pointer(obj)),
		0,
		0)
	return ret
}

// BindAsLegacyV2Runtime 激活策略并绑定当前runtimes
func (obj *ICLRRuntimeInfo) BindAsLegacyV2Runtime() uintptr {
	ret, _, _ := syscall.Syscall(
		obj.lpVtbl.BindAsLegacyV2Runtime,
		1,
		uintptr(unsafe.Pointer(obj)),
		0,
		0,
	)
	return ret
}

// IsLoadable 判断接口相关的runtimes是否可以加载到当前进程
func (obj *ICLRRuntimeInfo) IsLoadable() (pbLoadable bool, err error) {

	hr, _, err := syscall.Syscall(
		obj.lpVtbl.IsLoadable,
		2,
		uintptr(unsafe.Pointer(obj)),
		uintptr(unsafe.Pointer(&pbLoadable)),
		0)
	if err != syscall.Errno(0) {
		err = fmt.Errorf("the ICLRRuntimeInfo::IsLoadable method returned an error:\r\n%s", err)
		return
	}
	if hr != 0x00 {
		err = fmt.Errorf("the ICLRRuntimeInfo::IsLoadable method  returned a non-zero HRESULT: 0x%x", hr)
		return
	}
	err = nil
	return
}

func GetRuntimeInfo(metahost *ICLRMetaHost, version string) (*ICLRRuntimeInfo, error) {
	pwzVersion, err := syscall.UTF16PtrFromString(version)
	if err != nil {
		return nil, err
	}
	return metahost.GetRuntime(pwzVersion, IID_ICLRRuntimeInfo)
}


// GetInterface 返回ICORRuntimeHost 接口对象
func (obj *ICLRRuntimeInfo) GetInterface(rclsid windows.GUID, riid windows.GUID, ppUnk unsafe.Pointer) error {
	hr, _, err := syscall.Syscall6(
		obj.lpVtbl.GetInterface,
		4,
		uintptr(unsafe.Pointer(obj)),
		uintptr(unsafe.Pointer(&rclsid)),
		uintptr(unsafe.Pointer(&riid)),
		uintptr(ppUnk),
		0,
		0,
	)
	if err != syscall.Errno(0) && err.Error() != "The requested lookup key was not found in any active activation context." {
		return fmt.Errorf("the ICLRRuntimeInfo::GetInterface method returned an error:%s", err)
	}
	if hr != 0x00 {
		return fmt.Errorf("the ICLRRuntimeInfo::GetInterface method returned a non-zero HRESULT: 0x%x", hr)
	}
	return nil
}

func (obj *ICLRRuntimeInfo) GetVersionString() (version string, err error) {
	var pchBuffer uint32
	hr, _, err := syscall.Syscall(
		obj.lpVtbl.GetVersionString,
		3,
		uintptr(unsafe.Pointer(obj)),
		0,
		uintptr(unsafe.Pointer(&pchBuffer)),
	)
	if err != syscall.Errno(0) {
		err = fmt.Errorf("there was an error calling the ICLRRuntimeInfo::GetVersionString method during preallocation:%s", err)
		return
	}
	// 0x8007007a = The data area passed to a system call is too small, expected when passing a nil buffer for preallocation
	if hr != 0x00 && hr != 0x8007007a {
		err = fmt.Errorf("the ICLRRuntimeInfo::GetVersionString method (preallocation) returned a non-zero HRESULT: 0x%x", hr)
		return
	}

	pwzBuffer := make([]uint16, 20)

	hr, _, err = syscall.Syscall(
		obj.lpVtbl.GetVersionString,
		3,
		uintptr(unsafe.Pointer(obj)),
		uintptr(unsafe.Pointer(&pwzBuffer[0])),
		uintptr(unsafe.Pointer(&pchBuffer)),
	)
	if err != syscall.Errno(0) {
		err = fmt.Errorf("there was an error calling the ICLRRuntimeInfo::GetVersionString method:%s", err)
		return
	}
	if hr != 0x00 {
		err = fmt.Errorf("the ICLRRuntimeInfo::GetVersionString method returned a non-zero HRESULT: 0x%x", hr)
		return
	}
	err = nil
	version = syscall.UTF16ToString(pwzBuffer)
	return
}