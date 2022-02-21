package cmd

import (
	"os"
	"path"
)

func expandPath(p string) string {
	if path.IsAbs(p) {
		return p
	}
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	return path.Join(cwd, p)
}
