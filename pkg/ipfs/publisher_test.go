package ipfs_test

import (
	"errors"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vulcanize/block_watcher/pkg/ipfs"
	"github.com/vulcanize/block_watcher/test_helpers"
)

var _ = Describe("IPFS publisher", func() {
	It("calls dag put with the passed data", func() {
		mockDagPutter := test_helpers.NewMockDagPutter()
		publisher := ipfs.NewIpfsPublisher(mockDagPutter)
		fakeBytes := []byte{1, 2, 3, 4, 5}

		_, err := publisher.DagPut(fakeBytes)

		Expect(err).NotTo(HaveOccurred())
		Expect(mockDagPutter.Called).To(BeTrue())
		Expect(mockDagPutter.PassedBytes).To(Equal(fakeBytes))
	})

	It("returns error if dag put fails", func() {
		mockDagPutter := test_helpers.NewMockDagPutter()
		fakeError := errors.New("failed")
		mockDagPutter.SetError(fakeError)
		publisher := ipfs.NewIpfsPublisher(mockDagPutter)

		_, err := publisher.DagPut([]byte{1, 2, 3, 4, 5})

		Expect(err).To(HaveOccurred())
		Expect(err).To(MatchError(fakeError))
	})
})
