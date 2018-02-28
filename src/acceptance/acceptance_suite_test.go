package acceptance_test

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/pkg/transport"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestAcceptance(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Acceptance Suite")
}

type Config struct {
	EtcdEndpoints                      []string `json:"etcd_endpoints"`
	EtcdClientCAPath                   string   `json:"etcd_client_ca_path"`
	EtcdClientCertPath                 string   `json:"etcd_client_cert_path"`
	EtcdClientPrivateKeyPath           string   `json:"etcd_client_private_key_path"`
	EtcdClientSelfSignedCertPath       string   `json:"etcd_client_selfsigned_cert_path"`
	EtcdClientSelfSignedPrivateKeyPath string   `json:"etcd_client_selfsigned_private_key_path"`
}

var (
	config Config
	client *clientv3.Client
)

var _ = BeforeSuite(func() {
	var err error

	configFile := os.Getenv("CONFIG_FILE")
	Expect(configFile).NotTo(BeEmpty(), "CONFIG_FILE must be set to run the acceptance test suite.")

	configContents, err := ioutil.ReadFile(configFile)
	Expect(err).NotTo(HaveOccurred())

	Expect(json.Unmarshal(configContents, &config)).To(Succeed())
	tlsInfo := transport.TLSInfo{
		CertFile:      config.EtcdClientCertPath,
		KeyFile:       config.EtcdClientPrivateKeyPath,
		TrustedCAFile: config.EtcdClientCAPath,
	}
	tlsConfig, err := tlsInfo.ClientConfig()
	Expect(err).NotTo(HaveOccurred())

	client, err = clientv3.New(clientv3.Config{
		Endpoints:   config.EtcdEndpoints,
		DialTimeout: 5 * time.Second,
		TLS:         tlsConfig,
	})
	Expect(err).NotTo(HaveOccurred())
})
