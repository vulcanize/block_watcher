package ipfs_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vulcanize/eth-block-extractor/pkg/ipfs"
	"github.com/vulcanize/eth-block-extractor/test_helpers"
	ipfs_mocks "github.com/vulcanize/eth-block-extractor/test_helpers/mocks/ipfs"
)

var _ = Describe("IPLD publisher", func() {
	It("calls dag put with the passed data", func() {
		mockDagPutter := ipfs_mocks.NewMockDagPutter()
		publisher := ipfs.NewIpldPublisher(mockDagPutter)
		fakeBytes := []byte{1, 2, 3, 4, 5}

		_, err := publisher.DagPut(fakeBytes)

		Expect(err).NotTo(HaveOccurred())
		Expect(mockDagPutter.Called).To(BeTrue())
		Expect(mockDagPutter.PassedInterface).To(Equal(fakeBytes))
	})

	It("returns error if dag put fails", func() {
		mockDagPutter := ipfs_mocks.NewMockDagPutter()
		mockDagPutter.SetError(test_helpers.FakeError)
		publisher := ipfs.NewIpldPublisher(mockDagPutter)

		_, err := publisher.DagPut([]byte{1, 2, 3, 4, 5})

		Expect(err).To(HaveOccurred())
		Expect(err).To(MatchError(test_helpers.FakeError))
	})
})
