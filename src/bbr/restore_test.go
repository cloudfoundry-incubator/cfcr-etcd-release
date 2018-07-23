package bbr_test

import (
	"context"
	"io/ioutil"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/pkg/transport"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
)

var _ = Describe("Restore", func() {
	var (
		bbrDir         string
		bbrArgs        []string
		etcd           *clientv3.Client
		deploymentName string
		err            error
	)

	BeforeEach(func() {
		bbrDir, err = ioutil.TempDir("", "")
		Expect(err).ToNot(HaveOccurred())

		deploymentName = MustHaveEnv("BOSH_DEPLOYMENT")

		bbrArgs = []string{"deployment",
			"--target", MustHaveEnv("BOSH_ENVIRONMENT"),
			"--username", MustHaveEnv("BOSH_CLIENT"),
			"--password", MustHaveEnv("BOSH_CLIENT_SECRET"),
			"--deployment", deploymentName,
			"--ca-cert", MustHaveEnv("BOSH_CA_CERT_PATH")}

		etcdEndpoint := MustHaveEnv("ETCD_ENDPOINT")

		tlsInfo := transport.TLSInfo{
			CertFile:           MustHaveEnv("ETCD_CLIENT_CERT"),
			KeyFile:            MustHaveEnv("ETCD_KEY_FILE"),
			InsecureSkipVerify: true,
		}
		tlsConfig, err := tlsInfo.ClientConfig()
		Expect(err).NotTo(HaveOccurred())

		etcd, err = clientv3.New(clientv3.Config{
			Endpoints:   []string{etcdEndpoint},
			DialTimeout: 5 * time.Second,
			TLS:         tlsConfig,
		})
		Expect(err).NotTo(HaveOccurred())

		_, err = etcd.Put(context.TODO(), "key", "orig-value")
		Expect(err).ToNot(HaveOccurred())
	})

	It("should restore backed up artifacts", func() {
		backupCmd := exec.Command("bbr", append(bbrArgs, "backup")...)
		backupCmd.Dir = bbrDir
		session, err := gexec.Start(backupCmd, GinkgoWriter, GinkgoWriter)
		Expect(err).ToNot(HaveOccurred())
		Eventually(session, "1m").Should(gexec.Exit(0))

		_, err = etcd.Put(context.TODO(), "key", "fake-value")
		Expect(err).ToNot(HaveOccurred())

		val, err := etcd.Get(context.TODO(), "key")
		Expect(err).NotTo(HaveOccurred())
		Expect(string(val.Kvs[0].Value)).To(Equal("fake-value"))

		globbedFiles, err := filepath.Glob(bbrDir + "/" + deploymentName + "*")
		Expect(err).ToNot(HaveOccurred())
		Expect(globbedFiles).To(HaveLen(1))

		restoreCmd := exec.Command("bbr", append(bbrArgs, "restore", "--artifact-path", globbedFiles[0])...)
		session, err = gexec.Start(restoreCmd, GinkgoWriter, GinkgoWriter)
		Expect(err).ToNot(HaveOccurred())
		Eventually(session, "1m").Should(gexec.Exit(0))

		val, err = etcd.Get(context.TODO(), "key")
		Expect(err).NotTo(HaveOccurred())
		Expect(string(val.Kvs[0].Value)).To(Equal("orig-value"))
	})
})
