package acceptance_test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/coreos/etcd/clientv3"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestAcceptance(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Acceptance Suite")
}

type Config struct {
	EtcdEndpoint string `json:"etcd_endpoint"`
}

var (
	config Config
	client *clientv3.Client
)

var _ = BeforeSuite(func() {
	var err error
	config, err = ReadConfig()
	Expect(err).NotTo(HaveOccurred())

	client, err = clientv3.New(clientv3.Config{
		Endpoints: []string{
			config.EtcdEndpoint,
		},
		DialTimeout: 5 * time.Second,
	})
	Expect(err).NotTo(HaveOccurred())
})

func ReadConfig() (Config, error) {
	configFile := os.Getenv("CONFIG_FILE")
	if configFile == "" {
		return Config{}, fmt.Errorf("CONFIG_FILE must be set to run the acceptance test suite.")
	}

	configContents, err := ioutil.ReadFile(configFile)

	var config Config
	err = json.Unmarshal(configContents, &config)
	if err != nil {
		return Config{}, fmt.Errorf("failed to unmarshal configuration file: %s", err)
	}

	return config, nil
}
