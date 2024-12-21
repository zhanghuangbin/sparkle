package meta

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/thoas/go-funk"
	"os"
	"strings"
)
import "os/exec"

type OSType int

const (
	WIN OSType = iota
	LINUX
	MAC
)

func OfOSType(t string) (OSType, error) {
	switch t {
	case "windows":
		return WIN, nil
	case "darwin":
		return MAC, nil
	case "linux":
	case "plan9":
	case "solaris":
	case "posix":
	case "freebsd":
	case "openbsd":
		return LINUX, nil
	}
	return -1, errors.New("unsupported os type")
}

type Env interface {
	Apply(alias Alias) error
}

type WinEnv struct {
	global bool
}

// insure interface implementation
var _ Env = WinEnv{}

func New(osType OSType, global bool) (Env, error) {
	switch osType {
	case WIN:
		return WinEnv{global: global}, nil
	default:
		return nil, errors.New("unsupported os type")
	}
}

func (env WinEnv) Apply(alias Alias) error {
	oldVal, err := getPersistEnvVar(alias.Key)
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

	// 更新当前shell的环境变量
	// FIXME 无法更新当前shell的环境变量
	if _, _, err := doExec("cmd", "/c", "set", fmt.Sprintf("%s=%s", alias.Key, newVal)); err != nil {
		return err
	}
	fmt.Fprintf(os.Stdout, "set shell variable:%s successfully, newVal: \n\t%s \noldVal: \n\t%s\n", alias.Key, newVal, oldVal)

	return nil
}

func getPersistEnvVar(key string) (val string, err error) {
	if key == "" {
		return "", errors.New("key is empty")
	}

	formatKey := "%" + key + "%"
	// FIXME 获取到的新的shell的值，存在bug
	_, val, err = doExec("cmd.exe ", "/c", "echo", formatKey)
	if err != nil {
		return "", err
	}
	if val == formatKey+"\r\n" {
		return "", nil
	}

	if strings.LastIndex(val, "\r\n") != -1 {
		val = val[:strings.LastIndex(val, "\r\n")]
	}
	return
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

	var out bytes.Buffer
	command.Stdout = &out
	command.Stderr = &out
	if err := command.Run(); err != nil {
		if exiterr, ok := err.(*exec.ExitError); ok {
			return exiterr.ExitCode(), out.String(), err
		} else {
			return -1, out.String(), err
		}
	}
	return 0, out.String(), nil
}
