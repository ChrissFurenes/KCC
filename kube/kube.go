package kube

import (
	"context"
	"net"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/chrissfurenes/kcc/cmd"
	"gopkg.in/yaml.v3"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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
	kubeurl, _ := url.Parse(k.KubeConfig.Clusters[0].Cluster.Server)
	k.ClusterName = k.KubeConfig.Clusters[0].Name
	k.Address = kubeurl.Hostname()
	k.Port = kubeurl.Port()
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

func (k *Kube) Ping() {
	address := net.JoinHostPort(k.Address, k.Port)
	conn, err := net.DialTimeout("tcp", address, 2*time.Second)
	if err != nil {
		k.Reachable = false
	}
	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {

		}
	}(conn)
	k.Reachable = true
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
func (k *Kube) GetName() string {
	return k.ClusterName
}

func (k *Kube) GetNodes() int {
	nodes, err := k.Clientset.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return 0
	}
	return len(nodes.Items)
}
func (k *Kube) SetNodes() {
	k.Nodes = strconv.Itoa(k.GetNodes())
}

func (k *Kube) GetPods() int {
	pods, err := k.Clientset.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return 0
	}
	return len(pods.Items)
}
func (k *Kube) SetPods() {
	k.Pods = strconv.Itoa(k.GetPods())
}
func (k *Kube) GetNamespaces() int {
	namespace, err := k.Clientset.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return 0
	}
	return len(namespace.Items)
}
func (k *Kube) SetNamespaces() {
	k.Namespaces = strconv.Itoa(k.GetNamespaces())
}
func (k *Kube) GetControlplanes() int {
	controlplane, err := k.Clientset.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{}) // TODO: need testing
	if err != nil {
		return 0
	}
	return len(controlplane.Kind)
}

func (k *Kube) GetServices() int {
	services, err := k.Clientset.CoreV1().Services("").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return 0
	}
	return len(services.Items)
}

func (k *Kube) GetWorkers() int { // TODO: add if posable
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
