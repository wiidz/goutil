package systemHelper

import (
	"os/exec"
)

// ExecCmd 执行shell命令
// exec.Command 的第一个参数应该是命令，后续参数是命令的参数。如果 restart 是一个脚本文件，你可能需要使用 bash 或 sh 来执行它。例如：
// cmd := exec.Command("/bin/bash", "-c", "/home/hujiayilu/ros/RoServercode/bin/Release/restart")
func ExecCmd(name string, arg ...string) (res string, err error) {

	// 定义要执行的命令
	cmd := exec.Command(name, arg...)

	// 执行命令并捕获输出和错误
	output, err := cmd.CombinedOutput()
	if err != nil {
		return
	}

	// 打印命令输出
	res = string(output)
	return
}
