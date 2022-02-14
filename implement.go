package go_net

import (
	"fmt"
	"strings"
	"unsafe"
)

// 方法实现
func getInstalledRuntimes(metahost *ICLRMetaHost)([]string,error) {
	var runtimes []string
	// EnumerateInstalledRuntimes 修改返回方式
	enumICLRRuntimeInfo, err := metahost.EnumerateInstalledRuntimes()
	if err != nil {
		return runtimes,err
	}
	var hr int
	for hr != 0x01 {
		var runtimeInfo *ICLRRuntimeInfo
		var fetched = uint32(0)
		hr, err = enumICLRRuntimeInfo.Next(1,unsafe.Pointer(&runtimeInfo),&fetched)
		if err != nil {
			return runtimes, fmt.Errorf("InstalledRuntimes Next Error:%s", err)
		}
		if hr == 0x01 {
			break
		}
		// 释放
		runtimeInfo.Release()
		version, err := runtimeInfo.GetVersionString()
		if err != nil {
			return runtimes,err
		}
		runtimes = append(runtimes,version)
	}
	if len(runtimes) == 0 {
		return runtimes,fmt.Errorf("could not find any installed runtimes")
	}
	return runtimes,err
}

func PrepareParameters(params []string) (*SafeArray, error) {
	sab := SafeArrayBound{
		CElements: uint32(len(params)),
		LLbound:   0,
	}
	listStrSafeArrayPtr, err := SafeArrayCreate(VT_BSTR, 1, &sab) // VT_BSTR
	if err != nil {
		return nil, err
	}
	for i, p := range params {
		bstr, err := SysAllocString(p)
		if err != nil {
			return nil, err
		}
		SafeArrayPutElement(listStrSafeArrayPtr, int32(i), bstr)
	}

	paramVariant := Variant{
		VT:  VT_BSTR | VT_ARRAY, // VT_BSTR | VT_ARRAY
		Val: uintptr(unsafe.Pointer(listStrSafeArrayPtr)),
	}

	sab2 := SafeArrayBound{
		CElements: uint32(1),
		LLbound:   0,
	}
	paramsSafeArrayPtr, err := SafeArrayCreate(VT_VARIANT, 1, &sab2) // VT_VARIANT
	if err != nil {
		return nil, err
	}
	err = SafeArrayPutElement(paramsSafeArrayPtr, int32(0), unsafe.Pointer(&paramVariant))
	if err != nil {
		return nil, err
	}
	return paramsSafeArrayPtr, nil
}


func RuntimeInstall(runtimes string) (runtimeHost *ICORRuntimeHost, err error) {
	if runtimes == "" {
		// 默认v4版本运行时
		runtimes = "v4"
	}
	// 创建MetaHost实例对象
	metahost, err := CLRCreateInstance(CLSID_CLRMetaHost, IID_ICLRMetaHost)
	if err != nil {
		return runtimeHost,fmt.Errorf("install .NET runtimes has an error:%s",err)
	}
	times, err := getInstalledRuntimes(metahost)
	if err != nil {
		return nil,err
	}
	var latestRuntime string
	for _,r := range times {
		// 如果r与输入的版本相同最后一个版本
		if strings.Contains(r,runtimes) {
			latestRuntime = r
			break
		}else{
			latestRuntime = r
		}
	}
	runtimeInfo, err := GetRuntimeInfo(metahost,latestRuntime)
	if err != nil {
		return nil,err
	}
	isLoadable,err := runtimeInfo.IsLoadable()
	if err != nil {
		return nil, err
	}
	if !isLoadable {
		err = fmt.Errorf("%s is not loadable for some reason", latestRuntime)
	}
	return GetICORRuntimeHost(runtimeInfo)
}

// LoadAssembly 返回程序入口点
func LoadAssembly(runtimeHost *ICORRuntimeHost, rawBytes []byte) (methodInfo *MethodInfo, err error) {
	appDomain, err := GetAppDomain(runtimeHost)
	if err != nil {
		return
	}
	safeArrayPtr, err := CreateSafeArray(rawBytes)
	if err != nil {
		return
	}

	assembly, err := appDomain.Load_3(safeArrayPtr)
	if err != nil {
		return
	}
	return assembly.GetEntryPoint()
}


func InvokeAssembly(methodInfo *MethodInfo, params []string)(stdout string,stderr string) {
	var paramSafeArray *SafeArray
	methodSignature, err := methodInfo.GetString()
	if err != nil {
		stderr = err.Error()
		return stdout,stderr
	}
	if !strings.Contains(methodSignature, "Void Main()") {
		if paramSafeArray, err = PrepareParameters(params);err != nil {
			stderr = err.Error()
			return stdout,stderr
		}
	}
	nullVariant := Variant{
		VT:  1,
		Val: uintptr(0),
	}
	defer func(psa *SafeArray) {
		err := SafeArrayDestroy(psa)
		if err != nil {
			return
		}
	}(paramSafeArray)

	Mutex.Lock()
	defer Mutex.Unlock()
	// 调用程序
	err = methodInfo.Invoke_3(nullVariant, paramSafeArray)
	if err != nil {
		stderr = err.Error()
	}

	if wSTDOUT != nil {
		var e string
		stdout, e, err = ReadStdoutStderr()
		stderr += e
		if err != nil {
			stderr += err.Error()
		}
	}

	return stdout,stderr

}