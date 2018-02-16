package acceptance_test

import (
	"context"
	"fmt"
	"time"

	"github.com/satori/go.uuid"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

const (
	ETCD_REQUEST_TIMEOUT = 5 * time.Second
)

var _ = Describe("Cluster Health", func() {
	It("checks if the cluster is operational", func() {
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
})
