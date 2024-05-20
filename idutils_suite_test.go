package idutils

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"testing"
)

func TestIdutils(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Idutils Suite")
}
