package winapi

import (
	"fmt"
	"log"
	"os/exec"
)

var required_loadable_dlls []string = []string{
	"userenv.dll",
}

func init() {
	// cmd := exec.Command("cmd", "/C", "for /f %s in ('dir /b *.dll') do regsvr32 /s %s")
	// cmd.Dir = "C:\\Windows\\system32"
	// if out, err := cmd.CombinedOutput(); err != nil {
	// 	log.Printf("[Warning] failed to reload DLL %v with error: %s and output: %s", nil, err.Error(), out)
	// } else {
	// 	log.Printf("[Info] successfully reloaded system DLL %v", nil)
	// }

	for _, dll := range required_loadable_dlls {
		cmd := exec.Command("cmd", "/C", fmt.Sprintf("regsvr32 /s %s", dll))
		cmd.Dir = "C:\\Windows\\System32"
		if out, err := cmd.CombinedOutput(); err != nil {
			log.Printf("[Warning] failed to register DLL %s with error: %s and output: %s", dll, err.Error(), out)
		} else {
			log.Printf("[Info] successfully registered system DLL %s", dll)
		}
	}
}
