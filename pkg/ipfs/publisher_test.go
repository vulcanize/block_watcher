package ipfs_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vulcanize/block_watcher/pkg/ipfs"
	"github.com/vulcanize/block_watcher/test_helpers"
	ipfs_wrapper "github.com/vulcanize/block_watcher/test_helpers/mocks/ipfs"
)

var _ = Describe("IPFS publisher", func() {
	It("calls dag put with the passed data", func() {
		mockDagPutter := ipfs_wrapper.NewMockDagPutter()
		publisher := ipfs.NewIpfsPublisher(mockDagPutter)
		fakeBytes := []byte{1, 2, 3, 4, 5}

		_, err := publisher.DagPut(fakeBytes)

		Expect(err).NotTo(HaveOccurred())
		Expect(mockDagPutter.Called).To(BeTrue())
		Expect(mockDagPutter.PassedBytes).To(Equal(fakeBytes))
	})

	It("returns error if dag put fails", func() {
		mockDagPutter := ipfs_wrapper.NewMockDagPutter()
		mockDagPutter.SetError(test_helpers.FakeError)
		publisher := ipfs.NewIpfsPublisher(mockDagPutter)

		_, err := publisher.DagPut([]byte{1, 2, 3, 4, 5})

		Expect(err).To(HaveOccurred())
		Expect(err).To(MatchError(test_helpers.FakeError))
	})
})
