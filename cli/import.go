package cli

import (
	"os"
	"path/filepath"

	"github.com/chrissfurenes/kcc/kcc"
	"gopkg.in/yaml.v3"
)

type Importexport struct {
	KubeConfigFilePath  string
	TalosConfigFilePath string
	InputFilePath       string
	OutputFilePath      string
	kccFile             kcc.Kccfile
}

func (k *Importexport) Import() error { // TODO: make import work
	// import kube config
	kubefile, err := os.ReadFile(k.KubeConfigFilePath)
	if err != nil {
		return err
	}
	if err := yaml.Unmarshal(kubefile, &k.kccFile.Kube); err != nil {
		return err
	}

	// import talos config if added
	if len(k.TalosConfigFilePath) != 0 {
		talosfile, err := os.ReadFile(k.TalosConfigFilePath)
		if err != nil {
			return err
		}
		if err := yaml.Unmarshal(talosfile, &k.kccFile.Talos); err != nil {
			return err
		}
	}
	k.kccFile.FileType = "KCC"

	f, err := os.OpenFile(filepath.Join(kcc.KccConfigsPath(), k.OutputFilePath, k.kccFile.Kube.Clusters[0].Name)+".kcc", os.O_RDWR|os.O_CREATE, 640)
	if err != nil {
		return err
	}

	enc := yaml.NewEncoder(f)
	enc.SetIndent(2)

	_ = enc.Encode(k.kccFile)
	_ = enc.Close()
	return nil
}

func (k *Importexport) Export() error { // TODO: make export work
	innfile, err := os.ReadFile(k.InputFilePath)
	if err != nil {
		return err
	}
	if err := yaml.Unmarshal(innfile, &k.kccFile); err != nil {
		return err
	}
	kubefile, err := os.OpenFile(filepath.Join(k.OutputFilePath, k.kccFile.Kube.Contexts[0].Context.Cluster)+".kubeconfig", os.O_RDWR|os.O_CREATE, 640)
	if err != nil {
		return err
	}
	enc := yaml.NewEncoder(kubefile)
	enc.SetIndent(2)
	_ = enc.Encode(k.kccFile.Kube)
	_ = enc.Close()
	if len(k.kccFile.Talos.Context) != 0 {
		talosfile, err := os.OpenFile(filepath.Join(k.OutputFilePath, k.kccFile.Kube.Contexts[0].Context.Cluster)+".talosconfig", os.O_RDWR|os.O_CREATE, 640)
		if err != nil {
			return err
		}
		encT := yaml.NewEncoder(talosfile)
		encT.SetIndent(2)
		_ = encT.Encode(k.kccFile.Talos)
		_ = encT.Close()
	}
	return nil
}
