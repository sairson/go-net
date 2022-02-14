package go_net

import (
	"bytes"
	"fmt"
	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"
	"syscall"
	"unsafe"
)

var OleAut32 = syscall.MustLoadDLL("OleAut32.dll")

// SafeArray SafeArray结构类型
type SafeArray struct {
	cDims      uint16
	fFeatures  uint16
	cbElements uint32
	cLocks     uint32
	pvData     uintptr
	rgsabound  [1]SafeArrayBound
}

type SafeArrayBound struct {
	CElements uint32
	LLbound   int32
}

//SafeArrayCreate SafeArrayCreate指针创建函数
func SafeArrayCreate(vt uint16,cDims uint32,rgsabound *SafeArrayBound)(safeArray *SafeArray,err error) {
	// OleAut32.dll
	// SafeArrayCreate函数
	SafeArrayCreateHandle := OleAut32.MustFindProc("SafeArrayCreate")

	ret,_,err := SafeArrayCreateHandle.Call(
		uintptr(vt),
		uintptr(cDims),
		uintptr(unsafe.Pointer(rgsabound)))
	if err != syscall.Errno(0){
		return safeArray,err
	}
	// 没有错误,将err设置为nil
	err = nil
	// 判断ret,为0创建失败，但有返回值
	if ret == 0 {
		err = fmt.Errorf("the OleAut32:SafeArrayCreate function return 0x%x but the SafeArray was not create",ret)
	}
	// 创建成功
	safeArray = (*SafeArray)(unsafe.Pointer(ret))
	return safeArray,err
}


// CreateSafeArray 创建一个SafeArray指针
func CreateSafeArray(rawBytes []byte)(*SafeArray,error) {
	safeArrayBound := SafeArrayBound{
		CElements: uint32(len(rawBytes)),
		LLbound: int32(0),
	}
	safeArray,err := SafeArrayCreate(0x0011,1,&safeArrayBound)
	if err != nil {
		return nil, err
	}
	// 申请内存copy，将数据复制到pvData当中
	RtlCopyMemoryHandle := syscall.MustLoadDLL("ntdll.dll").MustFindProc("RtlCopyMemory")
	_, _, err = RtlCopyMemoryHandle.Call(
		safeArray.pvData,
		uintptr(unsafe.Pointer(&rawBytes[0])),
		uintptr(len(rawBytes)),
	)
	// 判断调用情况，err则调用失败
	if err != syscall.Errno(0) {
		return nil, err
	}
	return safeArray, nil
}

// SysAllocString SysAllocString用于分配一个新的字符串
func SysAllocString(str string) (unsafe.Pointer, error) {

	sysAllocString := OleAut32.MustFindProc("SysAllocString")

	input := utf16Le(str)
	ret, _, err := sysAllocString.Call(
		uintptr(unsafe.Pointer(&input[0])),
	)

	if err != syscall.Errno(0) {
		return nil, err
	}
	return unsafe.Pointer(ret), nil
}

func utf16Le(s string) []byte {
	enc := unicode.UTF16(unicode.LittleEndian, unicode.IgnoreBOM).NewEncoder()
	var buf bytes.Buffer
	t := transform.NewWriter(&buf, enc)
	t.Write([]byte(s))
	return buf.Bytes()
}

// SafeArrayPutElement SafeArrayPutElement用于将数据元素存储在数组中的指定位置
func SafeArrayPutElement(psa *SafeArray, rgIndices int32, pv unsafe.Pointer) error {
	safeArrayPutElement := OleAut32.MustFindProc("SafeArrayPutElement")

	hr, _, err := safeArrayPutElement.Call(
		uintptr(unsafe.Pointer(psa)),
		uintptr(unsafe.Pointer(&rgIndices)),
		uintptr(pv),
	)
	if err != syscall.Errno(0) {
		return err
	}
	if hr != 0x00 {
		return fmt.Errorf("the OleAut32:SafeArrayPutElement call returned a non-zero HRESULT: 0x%x", hr)
	}
	return nil
}

// SafeArrayDestroy 销毁安全指针
func SafeArrayDestroy(psa *SafeArray) error {

	safeArrayDestroy := OleAut32.MustFindProc("SafeArrayDestroy")

	hr, _, err := safeArrayDestroy.Call(
		uintptr(unsafe.Pointer(psa)),
		0,
		0,
	)
	if err != syscall.Errno(0) {
		return fmt.Errorf("the oleaut32!SafeArrayDestroy function call returned an error:\n%s", err)
	}
	if hr != 0x00 {
		return fmt.Errorf("the oleaut32!SafeArrayDestroy function returned a non-zero HRESULT: 0x%x", hr)
	}
	return nil
}