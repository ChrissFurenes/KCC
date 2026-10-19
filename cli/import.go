package cli

import (
	"os"

	"github.com/chrissfurenes/kcc/kube"
	"github.com/chrissfurenes/kcc/talos"
	"gopkg.in/yaml.v3"
)

type kccfile struct {
	kube  kube.KubeConfigInformation   `yaml:"kube"`
	talos talos.TalosConfigInformation `yaml:"talos"`
}

type importexport struct {
	kubeConfigFilePath string
	talosConfigFilePath string
	inputFilePath      string
	outputFilePath     string
	kccFile kccfile
}

func (k *importexport) Import() error { // TODO: make import work
	// import kube config
	kubefile ,err := os.ReadFile(k.kubeConfigFilePath)
	if err != nil {
		return err
	}
	if err := yaml.Unmarshal(kubefile, &k.kccFile.kube); err != nil { return err }

	// import talos config if added
	if len(k.talosConfigFilePath) != 0 {
		talosfile ,err := os.ReadFile(k.talosConfigFilePath)
		if err != nil {
			return err
		}
		if err := yaml.Unmarshal(talosfile, &k.kccFile.talos); err != nil { return err }
	}

	f, err := os.OpenFile(k.outputFilePath, os.O_WRONLY|os.O_CREATE,600)
	if err != nil {
		return err
	}
	enc := yaml.NewEncoder(f)
	enc.SetIndent(2)
	_ = enc.Encode(&k.kccFile)
	_ = enc.Close()
	return nil
}

func (k *kccfile) Export() error { // TODO: make export work
	return nil
}
