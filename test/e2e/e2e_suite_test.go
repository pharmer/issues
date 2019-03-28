package e2e

import (
	"flag"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/pharmer/pharmer/test/e2e/util"
)

func TestE2e(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "E2e Suite")
}

var kv1 = flag.String("current-version", "v12.0.5", "Lowest Kubernetes version")
var kv2 = flag.String("desired-version", "v13.0.4", "Highest Kubernetes version")
var provider = flag.String("provider", "", "Provider name")
var zone = flag.String("zone", "", "Zones name")
var nodes = flag.String("nodes", "", "Node type")
var file = flag.String("from-file", "", "File path for GoogleCloud credential")
var cluster string

var (
	KUBERNETES_VERSION string
	err error
)

var providers = make(map[string]string, 10)
func init(){
	providers = map[string]string{
		"linode": "Linode",
		"aks" : "Azure",
		"azure" : "Azure",
		"aws" : "AWS",
		"eks" : "AWS",
		"digitalocean" : "DigitalOcean",
		"gce" : "GoogleCloud",
		"gke" : "GoogleCloud",
		"packet" : "Packet",
		"vultr" : "Vultr",
	}
}

var createCredential = func() {
	By("Creating "+*provider+" Credential")
	err = RunScript("/create_credential.sh", providers[*provider], *provider, *file)
	Expect(err).NotTo(HaveOccurred())
}

var deleteCredential = func() {
	By("Deleting "+*provider+" Credential")
	err = RunScript("/delete_credential.sh", providers[*provider], *provider)
	Expect(err).NotTo(HaveOccurred())
}

var installDeps = func() {
	By("Installing Dependencies")
	err = RunScript("deps.sh")
	Expect(err).NotTo(HaveOccurred())
}

var _ = BeforeSuite(func() {
	installDeps()
	createCredential()
})

var _ = AfterSuite(func() {
	deleteCredential()
})
