// +build ignore

package main

import (
	"encoding/json"
	"fmt"
	"os"
	"syscall"

	"github.com/minio/minl/seccomp/seccomp"
)

// PR_SET_NO_NEW_PRIVS isn't exposed in Golang so we define it ourselves copying the value the kernel
const PR_SET_NO_NEW_PRIVS = 0x26

func prctl(option int, arg2, arg3, arg4, arg5 uintptr) (err error) {
	_, _, e1 := syscall.Syscall6(syscall.SYS_PRCTL, uintptr(option), arg2, arg3, arg4, arg5, 0)
	if e1 != 0 {
		err = e1
	}
	return
}

func main() {
	fmt.Println("Validate if seccomp enabled", seccomp.IsEnabled())
	f, err := os.Open("sample.json")
	if err != nil {
		fmt.Println("Unable to open sample.json", err)
		return
	}
	d := json.NewDecoder(f)
	scomp := &seccomp.Seccomp{}
	d.Decode(scomp)
	if err = prctl(PR_SET_NO_NEW_PRIVS, 1, 0, 0, 0); err != nil {
		//		fmt.Println("Unable to set privileges", err)
		return
	}
	seccomp.InitSeccomp(scomp)
}
