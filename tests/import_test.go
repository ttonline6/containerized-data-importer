package tests_test

import (
	"fmt"
	"strings"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"

	"kubevirt.io/containerized-data-importer/pkg/common"
	"kubevirt.io/containerized-data-importer/pkg/controller"
	"kubevirt.io/containerized-data-importer/tests"
	"kubevirt.io/containerized-data-importer/tests/framework"
	"kubevirt.io/containerized-data-importer/tests/utils"
)

const (
	testSuiteName                    = "Importer Test Suite"
	namespacePrefix                  = "importer"
	assertionPollInterval            = 2 * time.Second
	controllerSkipPVCCompleteTimeout = 60 * time.Second
	invalidEndpoint                  = "http://gopats.com/who-is-the-goat.iso"
	waitForCDIImportTimeout          = 30 * time.Second
)

var _ = Describe(testSuiteName, func() {
	f := framework.NewFrameworkOrDie(namespacePrefix)

	It("Should not perform CDI operations on PVC without annotations", func() {
		pvc, err := f.CreatePVCFromDefinition(utils.NewPVCDefinition("no-import", "1G", nil, nil))
		By("Verifying PVC with no annotation remains empty")
		Eventually(func() bool {
			log, err := tests.RunKubectlCommand(f, "logs", f.ControllerPod.Name, "-n", f.CdiInstallNs)
			Expect(err).NotTo(HaveOccurred())
			return strings.Contains(log, "pvc annotation \""+controller.AnnEndpoint+"\" not found, skipping pvc \""+f.Namespace.Name+"/no-import\"")
		}, controllerSkipPVCCompleteTimeout, assertionPollInterval).Should(BeTrue())
		Expect(err).ToNot(HaveOccurred())
		// Wait a while to see if CDI puts anything in the PVC.
		Expect(framework.VerifyPVCIsEmpty(f, pvc)).To(BeTrue())
		// Not deleting PVC as it will be removed with the NS removal.
	})

	It("Import pod status should be Fail on unavailable endpoint", func() {
		pvc, err := f.CreatePVCFromDefinition(utils.NewPVCDefinition(
			"no-import",
			"1G",
			map[string]string{controller.AnnEndpoint: invalidEndpoint},
			nil))
		Expect(err).ToNot(HaveOccurred())

		By("Verify the pod status is Failed on the target PVC")
		status, phaseAnnotation, err := utils.WaitForPVCAnnotation(f.K8sClient, f.Namespace.Name, pvc, controller.AnnPodPhase)
		Expect(phaseAnnotation).To(BeTrue())
		Expect(status).Should(BeEquivalentTo(v1.PodPending))
	})
})

var _ = Describe(testSuiteName+"-prometheus", func() {
	f := framework.NewFrameworkOrDie(namespacePrefix)

	BeforeEach(func() {
		_, err := f.CreatePrometheusServiceInNs(f.Namespace.Name)
		Expect(err).NotTo(HaveOccurred(), "Error creating prometheus service")
	})

	It("Import pod should have prometheus stats available while importing", func() {
		c := f.K8sClient
		ns := f.Namespace.Name
		httpEp := fmt.Sprintf("http://%s:%d", utils.FileHostName+"."+utils.FileHostNs, utils.HTTPNoAuthPort)
		pvcAnn := map[string]string{
			controller.AnnEndpoint: httpEp,
			controller.AnnSecret:   "",
		}

		By(fmt.Sprintf("Creating PVC with endpoint annotation %q", httpEp+"/tinyCore.iso"))
		_, err := utils.CreatePVCFromDefinition(c, ns, utils.NewPVCDefinition("import-e2e", "20M", pvcAnn, nil))
		Expect(err).NotTo(HaveOccurred(), "Error creating PVC")

		_, err = utils.FindPodByPrefix(c, ns, common.ImporterPodName, common.CDILabelSelector)
		//importer, err := utils.FindPodByPrefix(c, ns, common.ImporterPodName, common.CDILabelSelector)
		Expect(err).NotTo(HaveOccurred(), fmt.Sprintf("Unable to get importer pod %q", ns+"/"+common.ImporterPodName))

		var endpoint *v1.Endpoints
		var pod v1.Pod
		l, err := labels.Parse("prometheus.kubevirt.io")
		Expect(err).ToNot(HaveOccurred())
		Eventually(func() int {
			endpoint, err = c.CoreV1().Endpoints(ns).Get("kubevirt-prometheus-metrics", metav1.GetOptions{})
			Expect(err).NotTo(HaveOccurred())
			podList, err := c.CoreV1().Pods(ns).List(metav1.ListOptions{LabelSelector: l.String()})
			pod = podList.Items[0]
			Expect(err).ToNot(HaveOccurred())
			return len(endpoint.Subsets)
		}, 60, 1).Should(Equal(1))

		By("checking if the endpoint contains the metrics port and only one matching subset")
		Expect(endpoint.Subsets[0].Ports).To(HaveLen(1))
		Expect(endpoint.Subsets[0].Ports[0].Name).To(Equal("metrics"))
		Expect(endpoint.Subsets[0].Ports[0].Port).To(Equal(int32(443)))
	})
})
