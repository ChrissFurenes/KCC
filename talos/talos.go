package talos

import (
	"github.com/chrissfurenes/kcc/cmd"
)

type TalosConfigInformation struct {
	Context  string                        `yaml:"context"`
	Contexts map[string]TalosContextConfig `yaml:"contexts"`
}
type TalosContextConfig struct {
	Target    string   `yaml:"target,omitempty"`
	Endpoints []string `yaml:"endpoints"`
	Nodes     []string `yaml:"nodes,omitempty"`
}

type Talos struct {
	TalosPath     string
	CurrentFolder string
	TalosConfig   TalosConfigInformation
}

func NewTalos() *Talos {
	t := &Talos{
		TalosPath:     cmd.TalosPath(),
		CurrentFolder: "",
	}
	return t
}

func (t *Talos) ImportConfig(path string) error {
	return nil
}

func (t *Talos) Validate(path string) error {
	return nil
}

func (t *Talos) GetServerVersion() string {
	return ""
}

func (t *Talos) GetClientVersion() string {
	return ""
}

func (t *Talos) ApplyConfig(path string) error {
	return nil
}

func (t *Talos) ChangeClientTalosVersion(version string) error { // go big or go home XD
	return nil
}

func (t *Talos) DownloadClientTalosVersion(version string) error {
	return nil
}

func (t *Talos) FindAllClientTalosVersion() string {
	return ""
}

func (t *Talos) FindAllTalosVersionInAllCluster() string {
	return ""
}
