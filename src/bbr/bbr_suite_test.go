package bbr_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestBbr(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Bbr Suite")
}
