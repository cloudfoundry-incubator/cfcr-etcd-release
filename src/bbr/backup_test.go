package bbr_test

import (
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
	yaml "gopkg.in/yaml.v2"
)

func MustHaveEnv(keyname string) string {
	val := os.Getenv(keyname)
	Expect(val).NotTo(BeEmpty(), "Need "+keyname+" for the test")
	return val
}

type BBRArtifact struct {
	Name      string            `yaml:"name"`
	Checksums map[string]string `yaml:"checksums"`
}

type BBRInstance struct {
	Name      string        `yaml:"name"`
	Artifacts []BBRArtifact `yaml:"artifacts"`
}

type BBRMetadata struct {
	Instances []BBRInstance `yaml:"instances"`
}

var _ = Describe("Backup", func() {
	var (
		bbrDir string
		err    error
	)

	BeforeEach(func() {
		bbrDir, err = ioutil.TempDir("", "")
		Expect(err).ToNot(HaveOccurred())
	})

	It("should not give empty artefact", func() {
		deploymentName := MustHaveEnv("BOSH_DEPLOYMENT")
		command := exec.Command("bbr", "deployment",
			"--target", MustHaveEnv("BOSH_ENVIRONMENT"),
			"--username", MustHaveEnv("BOSH_CLIENT"),
			"--password", MustHaveEnv("BOSH_CLIENT_SECRET"),
			"--deployment", deploymentName,
			"--ca-cert", MustHaveEnv("BOSH_CA_CERT_PATH"),
			"backup",
		)
		command.Dir = bbrDir
		session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
		Expect(err).ToNot(HaveOccurred())
		Eventually(session, "1m").Should(gexec.Exit(0))

		globbedFiles, err := filepath.Glob(bbrDir + "/" + deploymentName + "*/metadata")
		Expect(err).ToNot(HaveOccurred())
		Expect(globbedFiles).To(HaveLen(1))

		rawMetadata, err := ioutil.ReadFile(globbedFiles[0])
		Expect(err).ToNot(HaveOccurred())

		metadata := BBRMetadata{}
		err = yaml.Unmarshal(rawMetadata, &metadata)

		Expect(err).ToNot(HaveOccurred())
		Expect(metadata.Instances[0].Artifacts[0].Checksums).ToNot(HaveLen(0))
	})
})
