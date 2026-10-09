package kube

import (
	"context"
	"fmt"
	_ "strconv"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

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
func (k *Kube) GetPods() int {
	pods, err := k.Clientset.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return 0
	}
	return len(pods.Items)
}
func (k *Kube) GetNamespaces() int {
	namespace, err := k.Clientset.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return 0
	}
	return len(namespace.Items)
}
func (k *Kube) GetServices() int {
	services, err := k.Clientset.CoreV1().Services("").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return 0
	}
	return len(services.Items)
}
func (k *Kube) GetClusterInfo() string {
	nodes := k.GetNodes()
	pods := k.GetPods()
	namespaces := k.GetNamespaces()
	controlplanes := k.GetControlplanes()
	workers := k.GetWorkers()
	return fmt.Sprintf("Nodes: %d, Pods: %d, Namespaces: %d, Control planes: %d, Workers: %d", nodes, pods, namespaces, controlplanes, workers)
}
func (k *Kube) GetControlplanes() int {
	controlplane, err := k.Clientset.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{}) // TODO: need testing
	if err != nil {
		return 0
	}
	total := 0
	for _, node := range controlplane.Items {
		if _, ok := node.Labels["node-role.kubernetes.io/control-plane"]; ok {
			total++
			continue
		}
		if _, ok := node.Labels["node-role.kubernetes.io/master"]; ok {
			total++
		}
	}
	return total
}

func (k *Kube) GetWorkers() int { // TODO: add if posable
	return 0
}
