package golangutils

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"
)

func GetFileNameWithoutExtension(p string) string {
	fileName := filepath.Base(p)
	fileExt := filepath.Ext(fileName)
	return strings.TrimSuffix(fileName, fileExt)
}

func Gofmt(f string) {
	cmd := exec.Command("gofmt", "-w", f)
	err := cmd.Run()

	if err != nil {
		fmt.Println(err)
	}
}
