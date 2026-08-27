package kube

import "strconv"

func ImportConfig(path string) error {
	return nil
}

func Validate(path string) error {
	return nil
}

func GetNodes() int {
	return 0
}
func GetPods() int {
	return 0
}
func GetNamespaces() int {
	return 0
}
func GetControlplanes() int {
	return 0
}
func GetWorkers() int {
	return 0
}
func GetClusterInfo() string {
	nodes := GetNodes()
	pods := GetPods()
	namespaces := GetNamespaces()
	controlplanes := GetControlplanes()
	workers := GetWorkers()
	return strconv.Itoa(nodes) + strconv.Itoa(pods) + strconv.Itoa(namespaces) + strconv.Itoa(controlplanes) + strconv.Itoa(workers)
}

func ApplyConfig(path string) error {
	return nil
}

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
