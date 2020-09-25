package util

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func GetCurrentPath() (string) {
	file, err := exec.LookPath(os.Args[0])
	if err != nil {
		return ""
	}
	path, err := filepath.Abs(file)
	if err != nil {
		return ""
	}

	path = strings.Replace(path, "\\", "/", -1)
	i := strings.LastIndex(path, "/")
	if i < 0 {
		return ""
	}
	return string(path[0 : i+1])
}
