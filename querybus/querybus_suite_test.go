package querybus_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestQuerybus(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Querybus Suite")
}
