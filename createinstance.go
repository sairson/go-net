package go_net

import (
	"fmt"
	"golang.org/x/sys/windows"
	"strings"
	"syscall"
	"unsafe"
)

func CLRCreateInstance(clsid, riid windows.GUID) (ppInterface *ICLRMetaHost, err error) {

	if clsid != CLSID_CLRMetaHost {
		err = fmt.Errorf("the input Class ID (CLSID) is not supported: %s", clsid)
		return
	}
	modMSCoree := syscall.MustLoadDLL("mscoree.dll")
	procCLRCreateInstance := modMSCoree.MustFindProc("CLRCreateInstance")

	// For some reason this procedure call returns "The specified procedure could not be found." even though it works
	hr, _, err := procCLRCreateInstance.Call(
		uintptr(unsafe.Pointer(&clsid)),
		uintptr(unsafe.Pointer(&riid)),
		uintptr(unsafe.Pointer(&ppInterface)),
	)
	if err != nil && !strings.Contains(err.Error(),"The specified procedure could not be found"){
		return nil, err
	}
	if hr != 0x00 {
		err = fmt.Errorf("the mscoree ! CLRCreateInstance function returned a non-zero HRESULT: 0x%x", hr)
		return
	}
	err = nil
	return ppInterface,err
}