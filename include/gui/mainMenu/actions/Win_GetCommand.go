//+build windows

package actions

import (
	"os/exec"
	"syscall"
	"fmt"
)

func getCommand(args... string)*exec.Cmd{
	var cmd *exec.Cmd
	cmd = exec.Command(executableForConvertingWindows,args...)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	fmt.Println("getCommand",executableForConvertingWindows,args...)
	return cmd
}