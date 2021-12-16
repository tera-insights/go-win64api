// +build windows,amd64

package winapi

var (
	conFreeConsole = modAdvapi32.NewProc("FreeConsole")
)

// FreeConsole detaches the program from the console so that the terminal gets closed without stopping the program
FreeConsole(){
	r1, _, err := conFreeConsole.Call()
	if r1 == 0 {
		return err
	}
	return nil
}