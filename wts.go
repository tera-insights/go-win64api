// +build windows,amd64

package winapi

import (
	"syscall"
	"unsafe"

	so "github.com/tera-insights/go-win64api/shared"
)

// Windows API functions
var (
	modWtsapi32                 = syscall.NewLazyDLL("wtsapi32.dll")
	procWTSFreeMemory           = modWtsapi32.NewProc("WTSFreeMemory")
	procWTSFreeMemoryExW        = modWtsapi32.NewProc("WTSFreeMemoryExW")
	procWTSEnumerateSessionsW   = modWtsapi32.NewProc("WTSEnumerateSessionsW")
	procWTSEnumerateSessionsExW = modWtsapi32.NewProc("WTSEnumerateSessionsExW")
	procWTSLogoffSession        = modWtsapi32.NewProc("WTSLogoffSession")
)

// WTS_CURRENT_SERVER_HANDLE is a handle that can be used to perform operations on
// the server running the application.
const WTS_CURRENT_SERVER_HANDLE uintptr = 0

type wtsSessionInfoW struct {
	SessionId       uint32
	pWinStationName *uint16
	State           int32
}

type wtsSessionInfo1W struct {
	ExecEnvId    uint32
	State        int32
	SessionId    uint32
	pSessionName *uint16
	pHostName    *uint16
	pUserName    *uint16
	pDomainName  *uint16
	pFarmName    *uint16
}

// Maps nil pointers to the empty string
func utf16toStringOrEmpty(p *uint16) string {
	if p == nil {
		return ""
	}
	return UTF16toString(p)
}

func wtsFreeMemory(pointer uintptr) {
	procWTSFreeMemory.Call(pointer)
}

func wtsFreeMemoryEx(classType so.WTSTypeClass, pointer uintptr, count uint64) error {
	ret, _, err := procWTSFreeMemoryExW.Call(
		uintptr(classType),
		pointer,
		uintptr(count),
	)
	if ret == 0 {
		return err
	}
	return nil
}

// WTSEnumerateSessions enumerates Remote Desktop Services sessions on the specified
// server.
//
// See: https://docs.microsoft.com/en-us/windows/win32/api/wtsapi32/nf-wtsapi32-wtsenumeratesessionsw
func WTSEnumerateSessions(serverHandle uintptr) ([]so.WTSSessionInfo, error) {
	var sessionInfoPointer uintptr
	var count uint32
	var sizeTest wtsSessionInfoW
	sessionSize := unsafe.Sizeof(sizeTest)

	ret, _, err := procWTSEnumerateSessionsW.Call(
		serverHandle,       // hServer
		uintptr(uint32(0)), // Reserved, must be 0
		uintptr(uint32(1)), // Version, must be 1
		uintptr(unsafe.Pointer(&sessionInfoPointer)), // ppSessionInfo
		uintptr(unsafe.Pointer(&count)),              // pCount
	)
	// ret is a boolean representing success. 0 = false, anything else = true
	if ret == 0 {
		// err already contains the result of GetLastError
		return nil, err
	}

	if sessionInfoPointer != 0 {
		defer wtsFreeMemory(sessionInfoPointer)
	}

	retVal := make([]so.WTSSessionInfo, 0, count)
	for i := uint32(0); i < count; i++ {
		data := (*wtsSessionInfoW)(unsafe.Pointer(sessionInfoPointer + (uintptr(i) * sessionSize)))

		sessionInfo := so.WTSSessionInfo{
			SessionID:      data.SessionId,
			WinStationName: UTF16toString(data.pWinStationName),
			State:          so.WTSConnectState(data.State),
		}
		retVal = append(retVal, sessionInfo)
	}

	return retVal, nil
}

// WTSEnumerateSessionsEx enumerates sessions on a specified Remote Desktop Session Host
// (RD Session Host) server or Remote Desktop Virtualization Host (RD Virtualization Host)
// server.
//
// See: https://docs.microsoft.com/en-us/windows/win32/api/wtsapi32/nf-wtsapi32-wtsenumeratesessionsexw
func WTSEnumerateSessionsEx(serverHandle uintptr) ([]so.WTSSessionInfo1, error) {
	var sessionInfoPointer uintptr
	var count uint32
	var sizeTest wtsSessionInfo1W
	sessionSize := unsafe.Sizeof(sizeTest)
	pLevel := uint32(1)

	ret, _, err := procWTSEnumerateSessionsExW.Call(
		serverHandle,                                 // hServer
		uintptr(unsafe.Pointer(&pLevel)),             // pLevel, must be set to 1
		uintptr(uint32(0)),                           // Filter, must be 0
		uintptr(unsafe.Pointer(&sessionInfoPointer)), // ppSessionInfo
		uintptr(unsafe.Pointer(&count)),              // pCount
	)
	// ret is a boolean representing success. 0 = false, anything else = true
	if ret == 0 {
		// err already contains the result of GetLastError
		return nil, err
	}
	if sessionInfoPointer != 0 {
		defer wtsFreeMemoryEx(so.WTSTypeSessionInfoLevel1, sessionInfoPointer, uint64(count))
	}

	retVal := make([]so.WTSSessionInfo1, 0, count)
	for i := uint32(0); i < count; i++ {
		curPtr := unsafe.Pointer(sessionInfoPointer + (uintptr(i) * sessionSize))
		data := (*wtsSessionInfo1W)(curPtr)

		sessionInfo := so.WTSSessionInfo1{
			ExecEnvID:   data.ExecEnvId,
			State:       so.WTSConnectState(data.State),
			SessionID:   data.SessionId,
			SessionName: UTF16toString(data.pSessionName),
			HostName:    utf16toStringOrEmpty(data.pHostName),
			UserName:    utf16toStringOrEmpty(data.pUserName),
			DomainName:  utf16toStringOrEmpty(data.pDomainName),
			FarmName:    utf16toStringOrEmpty(data.pFarmName),
		}
		retVal = append(retVal, sessionInfo)
	}

	return retVal, nil
}

// WTSLogoffSession logs of a specified Remote Desktop Services session.
//
// If `wait` is true, then this function will only return once the session have been
// logged off.
func WTSLogoffSession(serverHandle uintptr, sessionID uint32, wait bool) error {
	bWait := int(0)
	if wait {
		bWait = 1
	}
	ret, _, err := procWTSLogoffSession.Call(
		serverHandle,       // hServer,
		uintptr(sessionID), // SessionId
		uintptr(bWait),     // bWait
	)
	// ret is a boolean representing success. 0 = false, anything else = true
	if ret == 0 {
		// err already contains the result of GetLastError
		return err
	}
	return nil
}
