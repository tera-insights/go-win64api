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
