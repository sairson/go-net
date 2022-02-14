package go_net

import (
	"fmt"
	"golang.org/x/sys/windows"
	"syscall"
	"unsafe"
)


// https://docs.microsoft.com/en-us/windows/win32/api/oaidl/nn-oaidl-isupporterrorinfo

type ISupportErrorInfo struct {
	lpVtbl *ISupportErrorInfoVtbl
}

type ISupportErrorInfoVtbl struct {
	QueryInterface uintptr
	AddRef uintptr
	Release uintptr
	InterfaceSupportsErrorInfo uintptr
}



func (obj *ISupportErrorInfo) QueryInterface(riid windows.GUID, ppvObject unsafe.Pointer) error {

	hr, _, err := syscall.Syscall(
		obj.lpVtbl.QueryInterface,
		3,
		uintptr(unsafe.Pointer(obj)),
		uintptr(unsafe.Pointer(&riid)),
		uintptr(ppvObject),
	)
	if err != syscall.Errno(0) {
		return fmt.Errorf("the IUknown::QueryInterface method returned an error:\r\n%s", err)
	}
	if hr != 0x0 {
		return fmt.Errorf("the IUknown::QueryInterface method method returned a non-zero HRESULT: 0x%x", hr)
	}
	return nil
}

func (obj *ISupportErrorInfo) AddRef() (count uint32, err error) {
	ret, _, err := syscall.Syscall(
		obj.lpVtbl.AddRef,
		1,
		uintptr(unsafe.Pointer(obj)),
		0,
		0,
	)
	if err != syscall.Errno(0) {
		return 0, fmt.Errorf("the IUnknown::AddRef method returned an error:\r\n%s", err)
	}
	err = nil

	count = *(*uint32)(unsafe.Pointer(ret))
	return
}


func (obj *ISupportErrorInfo) Release() (count uint32, err error) {

	ret, _, err := syscall.Syscall(
		obj.lpVtbl.Release,
		1,
		uintptr(unsafe.Pointer(obj)),
		0,
		0,
	)
	if err != syscall.Errno(0) {
		return 0, fmt.Errorf("the IUnknown::Release method returned an error:\r\n%s", err)
	}
	err = nil
	count = *(*uint32)(unsafe.Pointer(ret))
	return
}

func (obj *ISupportErrorInfo) InterfaceSupportsErrorInfo(riid windows.GUID) error {

	hr, _, err := syscall.Syscall(
		obj.lpVtbl.InterfaceSupportsErrorInfo,
		2,
		uintptr(unsafe.Pointer(obj)),
		uintptr(unsafe.Pointer(&riid)),
		0,
	)
	if err != syscall.Errno(0) {
		return fmt.Errorf("the ISupportErrorInfo::InterfaceSupportsErrorInfo method returned an error:\r\n%s", err)
	}
	if hr != 0x0 {
		return fmt.Errorf("the ISupportErrorInfo::InterfaceSupportsErrorInfo method method returned a non-zero HRESULT: 0x%x", hr)
	}
	return nil
}

