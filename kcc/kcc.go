package kcc

import (
	"os"

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
	Status      string
	StatusText  string
	File        []byte
	IsDir       bool
	IsConfig    bool
	IsActive    bool
	IsTalos     bool
	IsLocked    bool
	IsBack      bool

	Kube  kube.Kube
	Talos talos.Talos
}

//spinner := []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}

func NewItem(file string) *Item {
	k := NewKubeConfig(file)

	return &Item{
		Name:     k.ClusterName,
		Path:     k.KubePath,
		FileName: k.KubePath,
		File:     nil,
		IsDir:    false,
		IsConfig: true,
		IsTalos:  TalosConfigExists(), // TODO: need to be fixed, but later me problem XD
		//Config:   k.KubeConfig,
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

func (i *Item) GetDisplayName() string { // DONE
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
		name += " - " + cmd.GreenText("ACTIVE")
	}
	return name
}

func (i *Item) GetInfoText() string {
	if i.IsConfig {
		i.InfoText = "Name:.. " + i.Name + "\n\n" +
			"User:.. " + i.Kube.User + "\n" +
			"IP:.... " + i.Kube.Address + "\n" +
			"Port:.. " + i.Kube.Port + "\n" +
			"Ping:.. " + i.PingText() + "\n" + // TODO: needs to be fixed
			"Path:.. " + i.Path + "\n" // to debugging
		if i.Kube.Reachable { // TODO: Change to run when get info from cluster (nodes, pods ....)
			clusterinfo := "\n" +
				"Nodes:....." + i.Kube.Nodes + "\n" +
				"Pods:......" + i.Kube.Pods + "\n" +
				"Namespace:." + i.Kube.Namespaces + "\n"
			i.InfoText = i.InfoText + clusterinfo
		} else {
			clusterinfo := "\n" +
				"Status: " + i.Kube.Status
			i.InfoText = i.InfoText + clusterinfo
		}
	}
	return i.InfoText
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

func (i *Item) PingText() string {
	if i.Kube.Reachable {
		return cmd.StatusText(true, "TRUE")
	}
	return cmd.StatusText(false, "FALSE")
}
