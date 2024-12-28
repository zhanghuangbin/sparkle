//go:build windows

package meta

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/atotto/clipboard"
	"github.com/thoas/go-funk"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"io"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

type winEnv struct {
	global bool
}

func instance(global bool) winEnv {
	return winEnv{global: global}
}

func (env winEnv) Apply(alias Alias) error {
	oldVal, err := getPersistEnvVar(env.global, alias.Key)
	if err != nil {
		return err
	}
	newVal := alias.Value
	if alias.Type == APPEND {
		newVal = appendOldEnvVal(oldVal, newVal)
	}

	args := make([]string, 0, 3)
	args = append(args, alias.Key, newVal)
	if env.global {
		args = append(args, "/M")
	}

	if _, _, err := doExec("setx", args...); err != nil {
		return err
	}

	globalDesc := ""
	if env.global {
		globalDesc = "global"
	}
	fmt.Fprintf(os.Stdout, "set %s env:%s successfully, newVal: \n\t%s \noldVal: \n\t%s\n", globalDesc, alias.Key, newVal, oldVal)

	setEnvCmd := fmt.Sprintf("set %s=\"%s\"", alias.Key, alias.Value)
	fmt.Fprintf(os.Stdout, fmt.Sprintf(`
----------------------------------------------------------
由于系统限制，无法在修改当前shell的环境变量，请复制以下命令到当前shell中执行。
`))

	if err := clipboard.WriteAll(setEnvCmd); err == nil {
		fmt.Fprintf(os.Stdout, `
      命令已复制到剪贴板，可直接粘贴到当前shell中执行！
`)
	}

	fmt.Fprintf(os.Stdout, "\n\n\t%s\n\n", setEnvCmd)

	fmt.Fprintf(os.Stdout, `

----------------------------------------------------------
`)

	return nil
}

func getPersistEnvVar(global bool, key string) (val string, err error) {
	if key == "" {
		return "", errors.New(fmt.Sprintf("key %s is empty", key))
	}

	args := []string{"query"}
	if global {
		args = append(args, "HKEY_LOCAL_MACHINE\\SYSTEM\\CurrentControlSet\\Control\\Session Manager\\Environment")
	} else {
		args = append(args, "HKEY_CURRENT_USER\\Environment")
	}
	args = append(args, "/v", key)

	_, val, err = doExec("reg", args...)
	if err != nil {
		if strings.Contains(val, "系统找不到指定的注册表项或值") {
			return "", nil
		}
		return "", err
	}

	lines := strings.Split(val, "\r\n")
	targetLine := ""
	for _, line := range lines {
		line = strings.TrimLeft(line, " ")
		if strings.HasPrefix(line, key) {
			targetLine = line
			break
		}
	}
	if targetLine == "" {
		return "", nil
	} else {
		tripleParts := strings.Split(targetLine, "    ")
		return tripleParts[2], nil
	}
}

func appendOldEnvVal(oldVal string, newVal string) string {
	if oldVal == "" {
		return newVal
	}
	pairs := strings.Split(oldVal, ";")
	idx := funk.IndexOf(pairs, newVal)
	if idx == -1 {
		return newVal + ";" + oldVal
	} else {
		r := make([]string, 0, len(pairs))
		r[0] = newVal
		r = append(r, pairs[idx:]...)
		r = append(r, pairs[idx+1:]...)
		return strings.Join(r, ";")
	}
}

func doExec(cmd string, args ...string) (int, string, error) {
	command := exec.Command(cmd, args...)

	log.Printf("exec: %s %s\n", cmd, strings.Join(args, " "))

	var out bytes.Buffer
	command.Stdout = &out
	command.Stderr = &out

	if err := command.Run(); err != nil {
		result, err1 := normalExecResult(out)
		if err1 != nil {
			return -1, "", err1
		}
		if exiterr, ok := err.(*exec.ExitError); ok {
			return exiterr.ExitCode(), result, errors.New(err.Error() + "," + result)
		} else {
			return -1, result, err
		}
	}

	result, err := normalExecResult(out)
	if err != nil {
		return -1, "", err
	}
	return 0, result, nil
}

func normalExecResult(out bytes.Buffer) (string, error) {
	if runtime.GOOS == "windows" {
		reader := transform.NewReader(bytes.NewBuffer(out.Bytes()), simplifiedchinese.GBK.NewDecoder())
		output, err := io.ReadAll(reader)
		if err != nil {
			return out.String(), err
		}

		return string(output), nil
	} else {
		return out.String(), nil
	}
}
