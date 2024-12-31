package systemHelper

import (
	"os/exec"
)

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
