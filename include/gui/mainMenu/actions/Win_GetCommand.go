//+build windows

package actions

func getCommand(args... string)*exec.Cmd{
	var cmd *exec.Cmd
	cmd = exec.CommandContext(executableForConvertingWindows,args...)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	return cmd
}