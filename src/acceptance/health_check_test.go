package acceptance_test

import (
	"context"
	"fmt"
	"net/url"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/pkg/transport"
	uuid "github.com/satori/go.uuid"
)

const (
	ETCD_REQUEST_TIMEOUT = 5 * time.Second
)

var _ = Describe("The cluster", func() {
	It("is operational", func() {
		clientContext, cancel := context.WithTimeout(context.Background(), ETCD_REQUEST_TIMEOUT)
		defer cancel()

		guid := uuid.NewV4()
		key := fmt.Sprintf("some-key-%s", guid)
		value := fmt.Sprintf("some-value-%s", guid)

		By("writing a key", func() {
			_, err := client.Put(clientContext, key, value)
			Expect(err).NotTo(HaveOccurred())
		})

		By("reading back the key", func() {
			getResponse, err := client.Get(clientContext, key)
			Expect(err).NotTo(HaveOccurred())

			Expect(getResponse.Kvs).To(HaveLen(1))
			Expect(string(getResponse.Kvs[0].Value)).To(Equal(value))
		})

		_, err := client.Delete(clientContext, key)
		Expect(err).NotTo(HaveOccurred())
	})

	It("doesn't accept http connections", func() {
		insecureUrl, err := url.Parse(config.EtcdEndpoint)
		Expect(err).NotTo(HaveOccurred())

		insecureUrl.Scheme = "http"

		client, err = clientv3.New(clientv3.Config{
			Endpoints: []string{
				insecureUrl.String(),
			},
			DialTimeout: 5 * time.Second,
		})

		Expect(err).NotTo(HaveOccurred())

		clientContext, cancel := context.WithTimeout(context.Background(), ETCD_REQUEST_TIMEOUT)
		defer cancel()

		guid := uuid.NewV4()
		key := fmt.Sprintf("some-key-%s", guid)
		value := fmt.Sprintf("some-value-%s", guid)

		By("writing a key", func() {
			_, err := client.Put(clientContext, key, value)
			Expect(err.Error()).To(ContainSubstring("transport is closing"))
		})
	})

	It("doesn't accept https connections without signed client certs", func() {
		tlsInfo := transport.TLSInfo{
			CertFile:      config.EtcdClientSelfSignedCertPath,
			KeyFile:       config.EtcdClientSelfSignedPrivateKeyPath,
			TrustedCAFile: config.EtcdClientCAPath,
		}
		tlsConfig, err := tlsInfo.ClientConfig()
		Expect(err).NotTo(HaveOccurred())

		client, err = clientv3.New(clientv3.Config{
			Endpoints: []string{
				config.EtcdEndpoint,
			},
			DialTimeout: 1 * time.Second,
			TLS:         tlsConfig,
		})

		Expect(err.Error()).To(ContainSubstring("context deadline exceeded"))
	})
})
