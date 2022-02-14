package go_net

import (
	"fmt"
	"syscall"
	"unsafe"
)

type IEnumUnknown struct {
	vtbl *IEnumUnknownVtbl
}

type IEnumUnknownVtbl struct {
	QueryInterface uintptr
	AddRef         uintptr
	Release        uintptr
	Next uintptr
	Skip uintptr
	Reset uintptr
	Clone uintptr
}

func (obj *IEnumUnknown) AddRef() uintptr {
	ret, _, _ := syscall.Syscall(
		obj.vtbl.AddRef,
		1,
		uintptr(unsafe.Pointer(obj)),
		0,
		0)
	return ret
}

func (obj *IEnumUnknown) Release() uintptr {
	ret, _, _ := syscall.Syscall(
		obj.vtbl.Release,
		1,
		uintptr(unsafe.Pointer(obj)),
		0,
		0)
	return ret
}

func (obj *IEnumUnknown) Next(celt uint32, pEnumRuntime unsafe.Pointer, pceltFetched *uint32) (hresult int, err error) {
	hr, _, err := syscall.Syscall6(
		obj.vtbl.Next,
		4,
		uintptr(unsafe.Pointer(obj)),
		uintptr(celt),
		uintptr(pEnumRuntime),
		uintptr(unsafe.Pointer(pceltFetched)),
		0,
		0,
	)
	if err != syscall.Errno(0) {
		err = fmt.Errorf("there was an error calling the IEnumUnknown::Next method:\r\n%s", err)
		return
	}
	if hr != 0x00 && hr != 0x01 {
		err = fmt.Errorf("the IEnumUnknown::Next method method returned a non-zero HRESULT: 0x%x", hr)
		return
	}
	err = nil
	hresult = int(hr)
	return
}