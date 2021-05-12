//+build windows

package actions

import (
	"os/exec"
	"syscall"
)

func getCommand(args... string)*exec.Cmd{
	var cmd *exec.Cmd
	cmd = exec.Command(executableForConvertingWindows,args...)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	return cmd
}