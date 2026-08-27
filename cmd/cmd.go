package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

func UserHomeDir() string {
	if runtime.GOOS == "windows" {
		home := os.Getenv("HOMEDRIVE") + os.Getenv("HOMEPATH")
		if home == "" {
			home = os.Getenv("USERPROFILE")
		}
		return filepath.Clean(home)
	}
	return filepath.Clean(os.Getenv("HOME"))
}

func KubePath() string {
	return filepath.Join(UserHomeDir(), ".kube")
}

func TalosPath() string {
	return filepath.Join(UserHomeDir(), ".talos")
}

func Help() {
	fmt.Println("KCC usage:")
	fmt.Println("version ........ prints version")
	fmt.Println("i/import ....... import new config")
	fmt.Println("help/h .......... show help")
}
