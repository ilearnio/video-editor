package config

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

var CWD_PATH string

func LoadPaths() {
	isTestEnv := flag.Lookup("test.v") != nil
	if isTestEnv {
		_, filename, _, _ := runtime.Caller(1)
		CWD_PATH = filepath.Dir(filepath.Dir(filename))
	} else {
		cwd, err := os.Getwd()
		if err != nil {
			fmt.Println(err)
		}
		CWD_PATH = cwd
	}
}
