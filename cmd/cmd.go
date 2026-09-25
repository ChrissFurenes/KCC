package cmd

import (
	"errors"
	"fmt"
	"io"
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

func CreateBackup() error {
	fmt.Println("Creating backup to ~/.kcc/backup/")
	err := CreateFolder(filepath.Join(UserHomeDir(), ".kcc"))
	if err != nil && !errors.Is(err, os.ErrExist) {
		return err
	}
	err = CreateFolder(filepath.Join(UserHomeDir(), ".kcc/backup"))
	if err != nil && !errors.Is(err, os.ErrExist) {
		return err
	}
	err = CreateFolder(filepath.Join(UserHomeDir(), ".kcc/backup/kube/"))
	if err != nil && !errors.Is(err, os.ErrExist) {
		return err
	}
	err = CreateFolder(filepath.Join(UserHomeDir(), ".kcc/backup/talos/"))
	if err != nil && !errors.Is(err, os.ErrExist) {
		return err
	}
	destKubeBackup := filepath.Join(UserHomeDir(), ".kcc/backup/kube/")
	destTalosBackup := filepath.Join(UserHomeDir(), ".kcc/backup/talos/")

	err = copyDir(KubeConfigsPath(), destKubeBackup)
	if err != nil {
		return err
	}
	err = copyDir(TalosConfigsPath(), destTalosBackup)
	if err != nil {
		return err
	}
	return nil
}

func Help() {
	fmt.Println("KCC usage:")
	fmt.Println("version .................... prints version")
	fmt.Println("i/import ..[talos/kube].... import new config arg 2 decide witch of talos and kube")
	fmt.Println("backup -b --backup..........backups talos and kube configs to ~/.kcc/backup/")
	fmt.Println("help/h ..................... show help")
}

func Copy(srcFile, dstFile string) error {
	out, err := os.Create(dstFile)
	if err != nil {
		return err
	}

	defer out.Close()

	in, err := os.Open(srcFile)
	if err != nil {
		return err
	}

	defer in.Close()

	_, err = io.Copy(out, in)
	if err != nil {
		return err
	}

	return nil
}

func copyDir(srcDir string, dstDir string) error {
	entries, err := os.ReadDir(srcDir)
	if err != nil {
		return err
	}
	for _, entry := range entries {
		srcPath := filepath.Join(srcDir, entry.Name())
		dstPath := filepath.Join(dstDir, entry.Name())
		fileInfo, err := os.Stat(srcPath)
		if err != nil {
			return err
		}
		switch fileInfo.Mode() & os.ModeType {
		case os.ModeDir:
			if err := CreateIfNotExists(dstPath, 0600); err != nil {
				return err
			}
			if err := copyDir(srcPath, dstPath); err != nil {
				return err
			}
		default:
			if err := Copy(srcPath, dstPath); err != nil {
				return err
			}
		}
	}
	return nil
}

func Exists(filePath string) bool {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return false
	}

	return true
}

func CreateIfNotExists(dir string, perm os.FileMode) error {
	if Exists(dir) {
		return nil
	}

	if err := os.MkdirAll(dir, perm); err != nil {
		return fmt.Errorf("failed to create directory: '%s', error: '%s'", dir, err.Error())
	}

	return nil
}
