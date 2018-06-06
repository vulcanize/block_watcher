package state_computation_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestStateComputation(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "StateComputation Suite")
}
