//+build linux

package actions

import "os/exec"

func getCommand(args... string)*exec.Cmd{
	return exec.Command(executableForConvertingLinux, args...)
}