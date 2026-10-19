package kcc

import (
	"path/filepath"

	"github.com/chrissfurenes/kcc/cmd"
)

func (k *Kcc) KccConfigDir() string {
	return filepath.Join(k.KccPath, "configs", k.CurrentFolder)
}

func KccPath() string {
	Dirpath := filepath.Join(cmd.UserHomeDir(), ".kcc")
	err := cmd.TestFolder(Dirpath)
	if err != nil {
		panic(err)
	}
	return filepath.Join(cmd.UserHomeDir(), ".kcc")
}
func KccConfigPath() string {
	return filepath.Join(KccPath(), "config")
}
func KccConfigsPath() string {
	Path := KccPath()
	DirPath := filepath.Join(Path, "configs")
	err := cmd.TestFolder(DirPath)
	if err != nil {
		panic(err)
	}
	return DirPath
}
