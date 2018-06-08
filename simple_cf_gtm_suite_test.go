package gtm_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestSimpleCfGtm(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "SimpleCfGtm Suite")
}
