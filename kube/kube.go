package kube

import (
	"fmt"
	"net"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/chrissfurenes/kcc/cmd"
	"gopkg.in/yaml.v3"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
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
	KubePath        string
	ClusterName     string
	CurrentFolder   string
	CurrentFilePath string
	User            string
	Address         string
	Port            string
	Reachable       bool
	Nodes           string
	Pods            string
	Namespaces      string
	Services        string
	Kubeconfig      *rest.Config
	Clientset       *kubernetes.Clientset
	Status          string
	KubeConfig      KubeConfigInformation
}

func NewKube(CurrentFilePath string) *Kube {
	k, err := LoadKube(CurrentFilePath)
	if err != nil {
		return &Kube{CurrentFilePath: CurrentFilePath, Status: err.Error()}
	}
	return k
}

func LoadKube(CurrentFilePath string) (*Kube, error) {
	k := &Kube{
		KubePath:        cmd.KubePath(),
		CurrentFilePath: CurrentFilePath,
	}
	if err := k.Init(); err != nil {
		return nil, err
	}
	return k, nil
}

func (k *Kube) Init() error {
	var path = filepath.Clean(k.CurrentFilePath)
	err := k.Validate(k.CurrentFilePath)
	if err != nil {
		return err
	}
	file, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	if err := yaml.Unmarshal(file, &k.KubeConfig); err != nil {
		return err
	}
	if len(k.KubeConfig.Clusters) == 0 || len(k.KubeConfig.Contexts) == 0 {
		return fmt.Errorf("invalid kubeconfig %s: missing clusters or contexts", path)
	}
	kubeurl, err := url.Parse(k.KubeConfig.Clusters[0].Cluster.Server)
	if err != nil || kubeurl.Hostname() == "" {
		return fmt.Errorf("invalid API server URL in %s", path)
	}
	k.ClusterName = k.KubeConfig.Clusters[0].Name
	k.Address = kubeurl.Hostname()
	k.Port = kubeurl.Port()
	if k.Port == "" {
		if kubeurl.Scheme == "https" {
			k.Port = "443"
		} else {
			k.Port = "80"
		}
	}
	k.User = k.KubeConfig.Contexts[0].Name

	err = k.InitCluster()
	if err != nil {
		return err

	}
	//k.Ping()
	if k.Reachable {
		k.SetPods()
		k.SetNodes()
		k.SetNamespaces()
	}
	return nil
}

func (k *Kube) ConfigDir() string {
	return filepath.Join(k.KubePath, "configs", k.CurrentFolder)
}

func ImportConfig(from string, to string) error {
	if to == "" || filepath.IsAbs(to) || filepath.Base(to) != to || to == "." || to == ".." {
		return fmt.Errorf("invalid destination filename: %q", to)
	}
	if _, err := os.Stat(from); err != nil {
		return err
	}
	return cmd.Copy(from, filepath.Join(cmd.KubeConfigsPath(), to))
}

func (k *Kube) Ping() bool {
	if k.Address != "" && k.Port != "" {
		address := net.JoinHostPort(k.Address, k.Port)
		conn, err := net.DialTimeout("tcp", address, 2*time.Second)
		if err != nil {
			k.Reachable = false
			return false
		}
		defer func(conn net.Conn) {
			err := conn.Close()
			if err != nil {

			}
		}(conn)
		k.Reachable = true
		return true
	}

	k.Reachable = false
	return false
}

func (k *Kube) Validate(path string) error { // TODO: Needs to validate config
	return nil
}
func (k *Kube) InitCluster() error {
	Kubeconfig, err := clientcmd.BuildConfigFromFlags("", k.CurrentFilePath)
	if err != nil {
		return err
	}
	k.Kubeconfig = Kubeconfig
	Clientset, err := kubernetes.NewForConfig(k.Kubeconfig)
	if err != nil {
		return err
	}
	k.Clientset = Clientset
	return nil
}

func (k *Kube) SetNodes() {
	k.Nodes = strconv.Itoa(k.GetNodes())
}

func (k *Kube) SetPods() {
	k.Pods = strconv.Itoa(k.GetPods())
}

func (k *Kube) SetNamespaces() {
	k.Namespaces = strconv.Itoa(k.GetNamespaces())
}

func (k *Kube) IsCurrentCluster() bool {
	readfile, err := os.Stat(k.CurrentFilePath)
	if err != nil {
		return false
	}
	currfile, err := os.Stat(cmd.KubePath())
	if err != nil {
		return false
	}
	return os.SameFile(readfile, currfile)
}
func ApplyConfig(path string) error { // TODO: fix this shit when [enter]
	return nil
}
