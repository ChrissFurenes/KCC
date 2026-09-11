package kcc

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/chrissfurenes/kcc/cmd"
	"github.com/chrissfurenes/kcc/kube"
	"github.com/chrissfurenes/kcc/talos"
)

type Item struct {
	Name        string
	Path        string
	FileName    string
	DisplayName string
	InfoText    string
	File        []byte
	IsDir       bool
	IsConfig    bool
	IsActive    bool
	IsTalos     bool
	IsLocked    bool
	IsBack      bool

	Config      kube.KubeConfigInformation // hmmm
	ClusterData kube.ClusterData           // hmmm
}

func NewItem(file string) *Item {
	k := NewKubeConfig(file)

	return &Item{
		Name:     k.ClusterName,
		Path:     k.KubePath,
		FileName: k.KubePath,
		File:     nil,
		IsDir:    false,
		IsConfig: true,
		IsTalos:  TalosConfigExists(), // need to be fixed, but later me problem XD
		Config:   k.KubeConfig,
	}

}

func NewKubeConfig(kubefile string) *kube.Kube {
	k := NewKubeConfig(kubefile)
	return k
}
func NewTalosConfig(talosfile string) *talos.Talos {
	t := NewTalosConfig(talosfile)
	return t
}
func TalosConfigExists() bool {
	return false
} // Later work

func (i *Item) GetDisplayName() string { // Done
	if i.IsBack {
		return " << Back to folder: " + i.Name
	}
	if i.IsDir {
		return "📁 " + i.Name
	}
	if !i.IsConfig {
		return "⚠ " + i.Name
	}
	name := "☸  " + i.Name
	i.IsActive = i.IsCurrentItem()
	if i.IsActive {
		name += " - " + GreenText("ACTIVE")
	}
	return name
}

func (i *Item) GetInfoText() string {
	if i.IsConfig {
		i.InfoText = "Name:.. " + i.Name + "\n\n" +
			"User:.. " + i.ClusterData.User + "\n" +
			"IP:.... " + i.ClusterData.Address + "\n" +
			"Port:.. " + strconv.FormatInt(i.ClusterData.Port, 10) + "\n" +
			"Ping:.. " + strings.ToUpper(strconv.FormatBool(false)) + "\n" +
			"Path:.. " + i.Path + "\n" // to debugging
	}
	return i.InfoText
}
func Backup(version string) { // need to add file check so it don`t always create a backup folder
	now := time.Now()
	str1 := now.Format("02-01-2006_15-04-05")
	err := os.Mkdir(filepath.Join(cmd.KubePath(), "backup__V_"+version+"__date_"+str1), 0777)
	if err != nil {
		panic(err)
	}

}

func (i *Item) Apply() {

}

func (i *Item) IsCurrentItem() bool {
	readfile, err := os.Stat(i.Path)
	if err != nil {
		return false
	}
	currfile, err := os.Stat(cmd.KubeConfigPath())
	if err != nil {
		return false
	}
	return os.SameFile(readfile, currfile)
}

func RedText(text string) string {
	return "[red]" + text + "[::-]"
}
func GreenText(text string) string {
	return "[green]" + text + "[::-]"
}
func YellowText(text string) string {
	return "[yellow]" + text + "[::-]"
}
func statusColorIcon(ok bool) (color, icon string) {
	if ok {
		return "[green]", " 🟢"
	}
	return "[red]", "🔴"
}
