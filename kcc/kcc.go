package kcc

import (
	"github.com/chrissfurenes/kcc/kube"
	"github.com/chrissfurenes/kcc/talos"
)

type Item struct {
	Name        string
	Path        string
	FileName    string
	File        []byte
	IsDir       bool
	IsConfig    bool
	IsActive    bool
	IsTalos     bool
	IsBack      bool
	Config      kube.KubeConfigInformation
	ClusterData kube.ClusterData
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
}
