package talos

import (
	"path/filepath"

	"github.com/chrissfurenes/kcc/cmd"
)

func (t *Talos) ConfigDir() string {
	return filepath.Join(t.TalosPath, "configs", t.CurrentFolder)
}

func TalosPath() string {
	Dirpath := filepath.Join(cmd.UserHomeDir(), ".talos")
	err := cmd.TestFolder(Dirpath)
	if err != nil {
		panic(err)
	}
	return filepath.Join(cmd.UserHomeDir(), ".talos")
}
func TalosConfigPath() string {
	return filepath.Join(TalosPath(), "config")
}
func TalosConfigsPath() string {
	Path := TalosPath()
	DirPath := filepath.Join(Path, "configs")
	err := cmd.TestFolder(DirPath)
	if err != nil {
		panic(err)
	}
	return DirPath
}
