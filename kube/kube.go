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
	Port            int64
	Reachable       bool
	Nodes           []string
	Pods            int
	Status          string
	KubeConfig      KubeConfigInformation
	Cluster         ClusterData
}

type ClusterData struct {
	User         string
	Address      string
	Port         int64
	Reachable    bool
	Nodes        int
	Pods         int
	TalosVersion string
	Status       string
	Test         string
}

func NewKube(CurrentFilePath string) *Kube {
	k := &Kube{
		KubePath:        cmd.KubePath(),
		CurrentFilePath: CurrentFilePath,
	}
	err := k.Init()
	if err != nil {
		panic(err)
	}
	return k
}

func (k *Kube) Init() error {
	var path = filepath.Clean(k.CurrentFilePath)
	file, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	if err := yaml.Unmarshal(file, &k.KubeConfig); err != nil {
		return err
	}
	k.ClusterName = k.KubeConfig.Clusters[0].Name
	k.Address = k.KubeConfig.Clusters[0].Cluster.Server
	k.User = k.KubeConfig.Contexts[0].Name
	return nil
}

func (k *Kube) ConfigDir() string {
	return filepath.Join(k.KubePath, "configs", k.CurrentFolder)
}

func ImportConfig(from string, to string) error {
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

func (k *Kube) Ping() (bool, error) {

	return true, nil
}

func (k *Kube) Validate(path string) error {
	return nil
}
func (k *Kube) GetName() string {
	return k.ClusterName
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

func (c *ClusterData) Ping() bool {
	address := net.JoinHostPort(c.Address, strconv.FormatInt(c.Port, 5))

	conn, err := net.DialTimeout("tcp", address, 2*time.Second)

	if err != nil {
		c.Reachable = false
		return false
	}
	defer conn.Close()
	c.Reachable = true
	return true
}
func (c *ClusterData) SetServer(server string) error {
	u, err := url.Parse(server)
	if err != nil {
		return err
	}
	if u.Hostname() == "" {
		return fmt.Errorf("invalid server address: %s", server)
	}

	c.Address = u.Hostname()
	c.Port, err = strconv.ParseInt(u.Port(), 10, 64)
	if err != nil {
		return err
	}
	if len(strconv.FormatInt(c.Port, 5)) <= 0 {
		switch u.Scheme {
		case "https":
			c.Port = 443
		case "http":
			c.Port = 80
		default:
			return fmt.Errorf(
				"server has no port: %s",
				server,
			)
		}
	}
	return nil
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
func ApplyConfig(path string) error {
	return nil
}
