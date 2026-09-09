package cmd

import (
	"errors"
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
	Dirpath := filepath.Join(UserHomeDir(), ".kube")
	err := TestFolder(Dirpath)
	if err != nil {
		panic(err)
	}
	return filepath.Join(UserHomeDir(), ".kube")
}
func KubeConfigPath() string {
	return filepath.Join(KubePath(), "config")
}
func KubeConfigsPath() string {
	Path := KubePath()
	DirPath := filepath.Join(Path, "configs")
	err := TestFolder(DirPath)
	if err != nil {
		panic(err)
	}
	return DirPath
}

func TalosPath() string {
	Dirpath := filepath.Join(UserHomeDir(), ".talos")
	err := TestFolder(Dirpath)
	if err != nil {
		panic(err)
	}
	return filepath.Join(UserHomeDir(), ".talos")
}
func TalosConfigsPath() string {
	Path := TalosPath()
	DirPath := filepath.Join(Path, "configs")
	err := TestFolder(DirPath)
	if err != nil {
		panic(err)
	}
	return DirPath
}

func TestFolder(name string) error {
	info, err := os.Stat(name)
	if err == nil {
		if info.IsDir() {
			return nil
		} else if !info.IsDir() {
			return fmt.Errorf("%s exists but is not a directory", name)
		}
	} else if errors.Is(err, os.ErrNotExist) {
		return CreateFolder(name)
	}
	return err
}

func CreateFolder(name string) error {
	return os.Mkdir(name, 0600)
}

func Help() {
	fmt.Println("KCC usage:")
	fmt.Println("version .................... prints version")
	fmt.Println("i/import ..[talos/kube].... import new config")
	fmt.Println("help/h ..................... show help")
}
