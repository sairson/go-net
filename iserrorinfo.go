package go_net

import (
	"fmt"
	"golang.org/x/sys/windows"
	"syscall"
	"unsafe"
)

// IErrorInfo IErrorInfo 接口结构
type IErrorInfo struct {
	lpVtbl *IErrorInfoVtbl
}

type IErrorInfoVtbl struct {
	QueryInterface uintptr
	AddRef uintptr
	Release uintptr
	GetDescription uintptr
	GetGUID uintptr
	GetHelpContext uintptr
	GetHelpFile uintptr
	GetSource uintptr
}

func (obj *IErrorInfo) AddRef() uintptr {
	ret, _, _ := syscall.Syscall(
		obj.lpVtbl.AddRef,
		1,
		uintptr(unsafe.Pointer(obj)),
		0,
		0)
	return ret
}

func (obj *IErrorInfo)Release()uintptr{
	ret, _, _ := syscall.Syscall(
		obj.lpVtbl.Release,
		1,
		uintptr(unsafe.Pointer(obj)),
		0,
		0)
	return ret
}

func (obj *IErrorInfo)QueryInterface()uintptr {
	ret, _, _ := syscall.Syscall(
		obj.lpVtbl.QueryInterface,
		1,
		uintptr(unsafe.Pointer(obj)),
		0,
		0)
	return ret
}


func (obj *IErrorInfo) GetDescription() (pbstrDescription *string, err error) {
	hr, _, err := syscall.Syscall(
		obj.lpVtbl.GetDescription,
		2,
		uintptr(unsafe.Pointer(obj)),
		uintptr(unsafe.Pointer(&pbstrDescription)),
		0,
	)

	if err != syscall.Errno(0) {
		err = fmt.Errorf("the IErrorInfo::GetDescription method returned an error:%s", err)
		return
	}
	if hr != 0x00 {
		err = fmt.Errorf("the IErrorInfo::GetDescription method method returned a non-zero HRESULT: 0x%x", hr)
		return
	}
	err = nil
	return
}


func (obj *IErrorInfo) GetGUID() (pGUID *windows.GUID, err error) {
	hr, _, err := syscall.Syscall(
		obj.lpVtbl.GetGUID,
		2,
		uintptr(unsafe.Pointer(obj)),
		uintptr(unsafe.Pointer(pGUID)),
		0,
	)

	if err != syscall.Errno(0) {
		err = fmt.Errorf("the IErrorInfo::GetGUID method returned an error:%s", err)
		return
	}
	if hr != 0x00 {
		err = fmt.Errorf("the IErrorInfo::GetGUID method method returned a non-zero HRESULT: 0x%x", hr)
		return
	}
	err = nil
	return
}

// GetErrorInfo  获取上次调用指针返回的错误
func GetErrorInfo() (pperrinfo *IErrorInfo, err error) {

	procGetErrorInfo := OleAut32.MustFindProc("GetErrorInfo")
	hr, _, err := procGetErrorInfo.Call(0, uintptr(unsafe.Pointer(&pperrinfo)))
	if err != syscall.Errno(0) {
		err = fmt.Errorf("the OleAu32.GetErrorInfo procedure call returned an error:%s", err)
		return
	}
	if hr != 0x00 {
		err = fmt.Errorf("the OleAu32.GetErrorInfo procedure call returned a non-zero HRESULT code: 0x%x", hr)
		return
	}
	err = nil
	return
}
