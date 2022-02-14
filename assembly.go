package go_net

import (
	"fmt"
	"golang.org/x/sys/windows"
	"syscall"
	"unsafe"
)

type Assembly struct {
	lpVtbl *AssemblyVtbl
}

type AssemblyVtbl struct {
	QueryInterface              uintptr
	AddRef                      uintptr
	Release                     uintptr
	GetTypeInfoCount            uintptr
	GetTypeInfo                 uintptr
	GetIDsOfNames               uintptr
	Invoke                      uintptr
	get_ToString                uintptr
	Equals                      uintptr
	GetHashCode                 uintptr
	GetType                     uintptr
	get_CodeBase                uintptr
	get_EscapedCodeBase         uintptr
	GetName                     uintptr
	GetName_2                   uintptr
	get_FullName                uintptr
	get_EntryPoint              uintptr
	GetType_2                   uintptr
	GetType_3                   uintptr
	GetExportedTypes            uintptr
	GetTypes                    uintptr
	GetManifestResourceStream   uintptr
	GetManifestResourceStream_2 uintptr
	GetFile                     uintptr
	GetFiles                    uintptr
	GetFiles_2                  uintptr
	GetManifestResourceNames    uintptr
	GetManifestResourceInfo     uintptr
	get_Location                uintptr
	get_Evidence                uintptr
	GetCustomAttributes         uintptr
	GetCustomAttributes_2       uintptr
	IsDefined                   uintptr
	GetObjectData               uintptr
	add_ModuleResolve           uintptr
	remove_ModuleResolve        uintptr
	GetType_4                   uintptr
	GetSatelliteAssembly        uintptr
	GetSatelliteAssembly_2      uintptr
	LoadModule                  uintptr
	LoadModule_2                uintptr
	CreateInstance              uintptr
	CreateInstance_2            uintptr
	CreateInstance_3            uintptr
	GetLoadedModules            uintptr
	GetLoadedModules_2          uintptr
	GetModules                  uintptr
	GetModules_2                uintptr
	GetModule                   uintptr
	GetReferencedAssemblies     uintptr
	get_GlobalAssemblyCache     uintptr
}

func NewAssembly(ppv uintptr) *Assembly {
	return (*Assembly)(unsafe.Pointer(ppv))
}

func (obj *Assembly) QueryInterface(riid *windows.GUID, ppvObject *uintptr) uintptr {
	ret, _, _ := syscall.Syscall(
		obj.lpVtbl.QueryInterface,
		3,
		uintptr(unsafe.Pointer(obj)),
		uintptr(unsafe.Pointer(riid)),
		uintptr(unsafe.Pointer(ppvObject)))
	return ret
}

func (obj *Assembly) AddRef() uintptr {
	ret, _, _ := syscall.Syscall(
		obj.lpVtbl.AddRef,
		1,
		uintptr(unsafe.Pointer(obj)),
		0,
		0)
	return ret
}

func (obj *Assembly) Release() uintptr {
	ret, _, _ := syscall.Syscall(
		obj.lpVtbl.Release,
		1,
		uintptr(unsafe.Pointer(obj)),
		0,
		0)
	return ret
}

func (obj *Assembly) GetEntryPoint() (pRetVal *MethodInfo, err error) {
	hr, _, err := syscall.Syscall(
		obj.lpVtbl.get_EntryPoint,
		2,
		uintptr(unsafe.Pointer(obj)),
		uintptr(unsafe.Pointer(&pRetVal)),
		0,
	)
	if err != syscall.Errno(0) {
		err = fmt.Errorf("the Assembly::GetEntryPoint method returned an error:%s", err)
		return
	}
	if hr != 0x0 {
		err = fmt.Errorf("the Assembly::GetEntryPoint method returned a non-zero HRESULT: 0x%x", hr)
		return
	}
	err = nil
	return
}