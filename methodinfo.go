package go_net

import (
	"fmt"
	"golang.org/x/sys/windows"
	"syscall"
	"unicode/utf16"
	"unsafe"
)

type MethodInfo struct {
	lpVtbl *MethodInfoVtbl
}

type MethodInfoVtbl struct {
	QueryInterface                 uintptr
	AddRef                         uintptr
	Release                        uintptr
	GetTypeInfoCount               uintptr
	GetTypeInfo                    uintptr
	GetIDsOfNames                  uintptr
	Invoke                         uintptr
	get_ToString                   uintptr
	Equals                         uintptr
	GetHashCode                    uintptr
	GetType                        uintptr
	get_MemberType                 uintptr
	get_name                       uintptr
	get_DeclaringType              uintptr
	get_ReflectedType              uintptr
	GetCustomAttributes            uintptr
	GetCustomAttributes_2          uintptr
	IsDefined                      uintptr
	GetParameters                  uintptr
	GetMethodImplementationFlags   uintptr
	get_MethodHandle               uintptr
	get_Attributes                 uintptr
	get_CallingConvention          uintptr
	Invoke_2                       uintptr
	get_IsPublic                   uintptr
	get_IsPrivate                  uintptr
	get_IsFamily                   uintptr
	get_IsAssembly                 uintptr
	get_IsFamilyAndAssembly        uintptr
	get_IsFamilyOrAssembly         uintptr
	get_IsStatic                   uintptr
	get_IsFinal                    uintptr
	get_IsVirtual                  uintptr
	get_IsHideBySig                uintptr
	get_IsAbstract                 uintptr
	get_IsSpecialName              uintptr
	get_IsConstructor              uintptr
	Invoke_3                       uintptr
	get_returnType                 uintptr
	get_ReturnTypeCustomAttributes uintptr
	GetBaseDefinition              uintptr
}

func NewMethodInfo(ppv uintptr) *MethodInfo {
	return (*MethodInfo)(unsafe.Pointer(ppv))
}

func (obj *MethodInfo) QueryInterface(riid windows.GUID, ppvObject unsafe.Pointer) error {
	hr, _, err := syscall.Syscall(
		obj.lpVtbl.QueryInterface,
		3,
		uintptr(unsafe.Pointer(obj)),
		uintptr(unsafe.Pointer(&riid)), // A reference to the interface identifier (IID) of the interface being queried for.
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
func (obj *MethodInfo) AddRef() uintptr {
	ret, _, _ := syscall.Syscall(
		obj.lpVtbl.AddRef,
		1,
		uintptr(unsafe.Pointer(obj)),
		0,
		0)
	return ret
}

func (obj *MethodInfo) Release() uintptr {
	ret, _, _ := syscall.Syscall(
		obj.lpVtbl.Release,
		1,
		uintptr(unsafe.Pointer(obj)),
		0,
		0)
	return ret
}

func (obj *MethodInfo) GetType(pRetVal *uintptr) uintptr {
	ret, _, _ := syscall.Syscall(
		obj.lpVtbl.GetType,
		2,
		uintptr(unsafe.Pointer(obj)),
		uintptr(unsafe.Pointer(pRetVal)),
		0)
	return ret
}

func (obj *MethodInfo) Invoke_3(variantObj Variant, parameters *SafeArray) (err error) {
	var pRetVal *Variant
	hr, _, err := syscall.Syscall6(
		obj.lpVtbl.Invoke_3,
		4,
		uintptr(unsafe.Pointer(obj)),
		uintptr(unsafe.Pointer(&variantObj)),
		uintptr(unsafe.Pointer(parameters)),
		uintptr(unsafe.Pointer(pRetVal)),
		0,
		0,
	)
	if err != syscall.Errno(0) {
		err = fmt.Errorf("the MethodInfo::Invoke_3 method returned an error:\r\n%s", err)
		return
	}

	// If the HRESULT is a TargetInvocationException, attempt to get the inner error
	// This currentl doesn't work
	if uint32(hr) == 0x80131604 {

		var iSupportErrorInfo *ISupportErrorInfo
		// See if MethodInfo supports the ISupportErrorInfo interface
		err = obj.QueryInterface(IID_ISupportErrorInfo, unsafe.Pointer(&iSupportErrorInfo))
		if err != nil {
			err = fmt.Errorf("the MethodInfo::QueryInterface method returned an error when looking for the ISupportErrorInfo interface:\r\n%s", err)
			return
		}
		err = iSupportErrorInfo.InterfaceSupportsErrorInfo(IID_ICorRuntimeHost)
		if err != nil {
			err = fmt.Errorf("there was an error with the ISupportErrorInfo::InterfaceSupportsErrorInfo method:\r\n%s", err)
			return
		}

		// Get the IErrorInfo object
		iErrorInfo, errG := GetErrorInfo()
		if errG != nil {
			err = fmt.Errorf("there was an error getting the IErrorInfo object:\r\n%s", errG)
			return err
		}

		// Read the IErrorInfo description
		desc, errD := iErrorInfo.GetDescription()
		if errD != nil {
			err = fmt.Errorf("the IErrorInfo::GetDescription method returned an error:\r\n%s", errD)
			return err
		}
		if desc == nil {
			err = fmt.Errorf("the Assembly::Invoke_3 method returned a non-zero HRESULT: 0x%x", hr)
			return
		}
		err = fmt.Errorf("the Assembly::Invoke_3 method returned a non-zero HRESULT: 0x%x with an IErrorInfo description of: %s", hr, *desc)
	}
	if hr != 0x0 {
		err = fmt.Errorf("the Assembly::Invoke_3 method returned a non-zero HRESULT: 0x%x", hr)
		return
	}

	if pRetVal != nil {
		err = fmt.Errorf("the Assembly::Invoke_3 method returned a non-zero pRetVal: %+v", pRetVal)
		return
	}
	err = nil
	return
}
// GetString returns a string version of the method's signature
func (obj *MethodInfo) GetString() (str string, err error) {
	var object *string
	hr, _, err := syscall.Syscall(
		obj.lpVtbl.get_ToString,
		2,
		uintptr(unsafe.Pointer(obj)),
		uintptr(unsafe.Pointer(&object)),
		0,
	)
	if err != syscall.Errno(0) {
		err = fmt.Errorf("the MethodInfo::ToString method returned an error:%s", err)
		return
	}
	if hr != 0x0 {
		err = fmt.Errorf("the Assembly::ToString method returned a non-zero HRESULT: 0x%x", hr)
		return
	}
	err = nil
	str = ReadUnicodeStr(unsafe.Pointer(object))
	return
}

func ReadUnicodeStr(ptr unsafe.Pointer) string {

	var byteVal uint16
	out := make([]uint16, 0)
	for i := 0; ; i++ {
		byteVal = *(*uint16)(unsafe.Pointer(ptr))
		if byteVal == 0x0000 {
			break
		}
		out = append(out, byteVal)
		ptr = unsafe.Pointer(uintptr(ptr) + 2)
	}
	return string(utf16.Decode(out))
}