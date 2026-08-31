package kube

import (
	"os"
	"path/filepath"
	"strconv"

	"github.com/chrissfurenes/kcc/cmd"
)

type KubeConfigInformation struct {
	Clusters []struct {
		Name    string `yaml:"name"`
		Cluster struct {
			Server string `yaml:"server"`
		} `yaml:"cluster"`
	} `yaml:"clusters"`

	Contexts []struct {
		Context struct {
			Cluster string `yaml:"cluster"`
			User    string `yaml:"user"`
		} `yaml:"context"`
		Name string `yaml:"name"`
	} `yaml:"contexts"`
}

type Kube struct {
	KubePath      string
	CurrentFolder string
	User          string
	Address       string
	Port          int64
	Reachable     bool
	Nodes         []string
	Pods          int
	Status        string
	KubeConfig    KubeConfigInformation
}

func NewKube(CurrentFolder string) *Kube {
	k := &Kube{
		KubePath:      cmd.KubePath(),
		CurrentFolder: CurrentFolder,
	}
	return k
}

func (k *Kube) ConfigDir() string {
	return filepath.Join(k.KubePath, "configs", k.CurrentFolder)
}

func (k *Kube) ImportConfig(ftype string, from string, to string) error {
	source := from
	if !filepath.IsAbs(source) {
		dir, err := os.Getwd()
		if err != nil {
			return err
		}
		source = filepath.Join(dir, source)
	}
	return nil
}

func (k *Kube) Validate(path string) error {
	return nil
}

func (k *Kube) GetNodes() int {
	return 0
}
func (k *Kube) GetPods() int {
	return 0
}
func (k *Kube) GetNamespaces() int {
	return 0
}
func (k *Kube) GetControlplanes() int {
	return 0
}
func (k *Kube) GetWorkers() int {
	return 0
}
func (k *Kube) GetClusterInfo() string {
	nodes := k.GetNodes()
	pods := k.GetPods()
	namespaces := k.GetNamespaces()
	controlplanes := k.GetControlplanes()
	workers := k.GetWorkers()
	return strconv.Itoa(nodes) + strconv.Itoa(pods) + strconv.Itoa(namespaces) + strconv.Itoa(controlplanes) + strconv.Itoa(workers)
}

func ApplyConfig(path string) error {
	return nil
}
