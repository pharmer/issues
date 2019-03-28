package util

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"time"

	"github.com/appscode/go/homedir"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"sigs.k8s.io/cluster-api/pkg/client/clientset_generated/clientset"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
)

const scriptDirectory = "scripts"

func RunScript(script string, args ...string) error {
	wd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get working directory: %s\n", err)
	}

	return runCommand(path.Join(wd, scriptDirectory, script), args...)
}

func runCommand(cmd string, args ...string) error {
	c := exec.Command(cmd, args...)
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	c.Env = append(c.Env, append(os.Environ())...)
	fmt.Printf("Running command %q\n", cmd)
	return c.Run()
}

func WaitForNodeReady(kc kubernetes.Interface, numNodes int32) error {
	var (
		nodes      []corev1.Node
		numReadyNodes int32
	)

	count := 1
	if err := wait.Poll(5*time.Second, 15*time.Minute, func() (bool, error) {
		fmt.Println("Attempt",count,": Waiting for the Nodes to be Ready . . . .")
		nl, err := kc.CoreV1().Nodes().List(metav1.ListOptions{})
		if err != nil {
			return false, err
		}
		nodes = nl.Items
		numReadyNodes = 0
	Nodes:
		for _, node := range nodes {
			_, ok := node.Labels["node-role.kubernetes.io/master"];if ok{
				continue
			}

			for _, taint := range node.Spec.Taints {
				if taint.Key == "node.cloudprovider.kubernetes.io/uninitialized" {
					continue Nodes
				}
			}

			for _, cond := range node.Status.Conditions {
				if cond.Type == corev1.NodeReady && cond.Status == corev1.ConditionTrue {
					numReadyNodes++
				}
			}
		}
		count++
		return numReadyNodes == numNodes, nil

	}); err != nil {
		return err
	}
	return nil
}

func KubeClient() (kubernetes.Interface, error) {
	kubeconfig := kubeConfigPath()
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return nil, err
	}
	return kubernetes.NewForConfig(config)
}

func ClusterApiClient() (clientset.Interface, error){
	kubeconfig := kubeConfigPath()
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return nil, err
	}
	return clientset.NewForConfig(config)
}

func kubeConfigPath() string {
	return homedir.HomeDir() + "/.kube/config"
}