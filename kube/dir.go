package kube

import (
	"path/filepath"

	"github.com/chrissfurenes/kcc/cmd"
)

func (k *Kube) ConfigDir() string {
	return filepath.Join(KubePath(), "configs", k.CurrentFolder)
}

func KubePath() string {
	Dirpath := filepath.Join(cmd.UserHomeDir(), ".kube")
	err := cmd.TestFolder(Dirpath)
	if err != nil {
		panic(err)
	}
	return filepath.Join(cmd.UserHomeDir(), ".kube")
}
func KubeConfigPath() string {
	return filepath.Join(KubePath(), "config")
}
func KubeConfigsPath() string {
	Path := KubePath()
	DirPath := filepath.Join(Path, "configs")
	err := cmd.TestFolder(DirPath)
	if err != nil {
		panic(err)
	}
	return DirPath
}
