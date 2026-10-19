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
	FileType    string
	IsDir       bool
	IsConfig    bool
	IsActive    bool
	IsTalos     bool
	IsLocked    bool
	IsBack      bool
	KccFile     Kccfile
	Kube        kube.Kube   // TODO: May remove... hmmm
	Talos       talos.Talos // TODO: May remove... hmmm
}
type Kcc struct {
	Name          string
	KccPath       string
	CurrentFolder string
	FilePath      string
}

type Kccfile struct {
	FileType string                       `yaml:"file-type"`
	Kube     kube.KubeConfigInformation   `yaml:"kube"`
	Talos    talos.TalosConfigInformation `yaml:"talos"`
}

//spinner := []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}

func NewItem(file string) *Item {
	k := NewKubeConfig(file)

	return &Item{
		Name:     k.ClusterName,
		File:     nil,
		IsDir:    false,
		IsConfig: true,
		IsTalos:  TalosConfigExists(), // TODO: need to be fixed, but later me problem XD
		//Config:   k.KubeConfig,
	}

}
func NewKccConfig(file string) *Kcc {
	return NewKcc(file)
}
func NewKcc(file string) *Kcc {
	k := &Kcc{
		FilePath: file,
	}
	return k
}
func NewKubeConfig(kubefile string) *kube.Kube { // NEW. TODO: FIX this
	return kube.NewKube(kubefile)
}
func NewTalosConfig(talosfile string) *talos.Talos { // NEW. TODO: FIX this
	return talos.NewTalos() // talos.NewTalos(talosfile)
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
			"FileType: " + i.FileType + "\n\n" +
			"User:.. " + i.Kube.User + "\n" +
			"IP:.... " + i.Kube.Address + "\n" +
			"Port:.. " + i.Kube.Port + "\n" +
			"Ping:.. " + i.PingText() + "\n" + // TODO: needs to be fixed
			"Path:.. " + i.Path + "\n" // to debugging
		if i.Kube.Reachable { // TODO: Change to run when get info from cluster (nodes, pods ....)
			clusterinfo := "\n" +
				"Nodes:..... " + i.Kube.Nodes + "\n" +
				"Pods:...... " + i.Kube.Pods + "\n" +
				"Namespace:. " + i.Kube.Namespaces + "\n"
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
	currfile, err := os.Stat(kube.KubeConfigPath())
	if err != nil {
		return false
	}
	return os.SameFile(readfile, currfile)
}

func (i *Item) PingText() string {
	return cmd.StatusText(i.Kube.Reachable)

}

///////////////
