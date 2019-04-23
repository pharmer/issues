package e2e

import (
	"github.com/appscode/go/crypto/rand"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/pharmer/pharmer/test/e2e/util"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var (
	replicas = "1"
)

var _ = Describe(*provider, func() {
	BeforeEach(func() {
		KUBERNETES_VERSION = *kv2
		cluster = "pharmer-test-" + rand.Characters(6)
	})

	var createCluster = func(kv string) {
		By("Running create_cluster.sh")
		err := RunScript("/create_cluster.sh", kv, cluster, *provider, *zone, *nodes, replicas, *masters)
		Expect(err).NotTo(HaveOccurred())

		By("Getting KubeConfig")
		kc, err := KubeClient()
		Expect(err).NotTo(HaveOccurred())

		By("Waiting for the Nodes to be Ready")
		err = WaitForNodeReady(kc, 1)
		Expect(err).NotTo(HaveOccurred())
	}

	var updateReplica = func(n int32) error {
		By("Getting Cluster API client")
		caClient, err := ClusterApiClient()
		Expect(err).NotTo(HaveOccurred())

		By("Getting MachineSet")
		machineSets, err := caClient.ClusterV1alpha1().MachineSets(metav1.NamespaceDefault).List(metav1.ListOptions{})
		Expect(err).NotTo(HaveOccurred())

		By("Updating Machines")
		for _, machineSet := range machineSets.Items {
			machineSet.Spec.Replicas = &n
			_, err = caClient.ClusterV1alpha1().MachineSets(metav1.NamespaceDefault).Update(&machineSet)
			Expect(err).NotTo(HaveOccurred())
		}

		By("Getting Kubernetes Client")
		kc, err := KubeClient()
		Expect(err).NotTo(HaveOccurred())

		By("Waiting for " + string(n) + " Nodes to become ready")

		err = WaitForNodeReady(kc, n)
		return err
	}

	var upgradeCluster = func() error {
		By("Running upgrade_cluster.sh")
		err := RunScript("/upgrade_cluster.sh", *kv2, cluster)
		return err
	}

	var deleteCluster = func() error {
		By("Running delete_cluster.sh")
		err := RunScript("/delete_cluster.sh", cluster)
		return err
	}

	JustBeforeEach(func() {
		createCluster(KUBERNETES_VERSION)
	})

	AfterEach(func() {
		Skip("skipping deleting cluster")
		By("Deleting Cluster")
		err = deleteCluster()
		Expect(err).NotTo(HaveOccurred())
	})

	Context("when kubernetes version is "+*kv2, func() {
		It("should create a cluster", func() {
			Expect(err).NotTo(HaveOccurred())
		})

		It("should increase the number of MachineSet replica from 1 to 2", func() {
			err = updateReplica(2)
			Expect(err).NotTo(HaveOccurred())
		})
	})

	Context("when kubernetes version is "+*kv1, func() {
		BeforeEach(func() {
			KUBERNETES_VERSION = *kv1
		})

		It("should create a cluster", func() {
			Expect(err).NotTo(HaveOccurred())
		})

		It("should upgrade the cluster from "+*kv1+" to "+*kv2, func() {
			err = upgradeCluster()
			Expect(err).NotTo(HaveOccurred())
		})
	})
})
