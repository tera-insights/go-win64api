package shared

// WTSConnectState is an enumeration that specifies the connection state of a
// Remote Desktop Services session.
//
// See: https://docs.microsoft.com/en-us/windows/win32/api/wtsapi32/ne-wtsapi32-wts_connectstate_class
type WTSConnectState int

// The valid values for the WTSConnectState enumeration.
const (
	WTSActive WTSConnectState = iota
	WTSConnected
	WTSConnectQuery
	WTSShadow
	WTSDisconnected
	WTSIdle
	WTSListen
	WTSReset
	WTSDown
	WTSInit
)

// WTSTypeClass is an enumeration that specifies the type of a structure that a Remote
// Desktop Services function has returned in a buffer.
//
// See: https://docs.microsoft.com/en-us/windows/win32/api/wtsapi32/ne-wtsapi32-wts_type_class
type WTSTypeClass int

// The valid values for the WTSTypeClass enumeration
const (
	WTSTypeProcessInfoLevel0 WTSTypeClass = iota
	WTSTypeProcessInfoLevel1
	WTSTypeSessionInfoLevel1
)

// A WTSSessionInfo contains information about a specific Remote Desktop Services session.
type WTSSessionInfo struct {
	SessionID      uint32
	WinStationName string
	State          WTSConnectState
}

// A WTSSessionInfo1 contains extended information about a specific Remote Desktop Services
// session.
//
// The HostName, UserName, DomainName, and FarmName entries may be empty, in which case
// NULL was returned by the underlying API.
//
// See: https://docs.microsoft.com/en-us/windows/win32/api/wtsapi32/ns-wtsapi32-wts_session_info_1w
type WTSSessionInfo1 struct {
	ExecEnvID   uint32
	State       WTSConnectState
	SessionID   uint32
	SessionName string
	HostName    string
	UserName    string
	DomainName  string
	FarmName    string
}
