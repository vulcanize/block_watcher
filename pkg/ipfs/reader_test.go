package ipfs_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vulcanize/eth-block-extractor/pkg/ipfs"
	"github.com/vulcanize/eth-block-extractor/test_helpers"
	ipfs_mocks "github.com/vulcanize/eth-block-extractor/test_helpers/mocks/ipfs"
)

var _ = Describe("Ipld reader", func() {
	It("gets block corresponding to cid", func() {
		getter := ipfs_mocks.NewMockGetter()
		resolver := ipfs_mocks.NewMockResolver()
		reader := ipfs.NewIpldReader(getter, resolver)

		_, err := reader.Read(test_helpers.FakeCid)

		Expect(err).NotTo(HaveOccurred())
		getter.AssertGetCalledWith(test_helpers.FakeCid)
	})

	It("returns error if getting block returns err", func() {
		getter := ipfs_mocks.NewMockGetter()
		getter.SetReturnErr(test_helpers.FakeError)
		resolver := ipfs_mocks.NewMockResolver()
		reader := ipfs.NewIpldReader(getter, resolver)

		_, err := reader.Read(test_helpers.FakeCid)

		Expect(err).To(HaveOccurred())
		Expect(err).To(MatchError(test_helpers.FakeError))
	})

	It("resolves block into IPLD node", func() {
		getter := ipfs_mocks.NewMockGetter()
		getter.SetReturnBlock(test_helpers.FakeIpfsBlock)
		resolver := ipfs_mocks.NewMockResolver()
		reader := ipfs.NewIpldReader(getter, resolver)

		_, err := reader.Read(test_helpers.FakeCid)

		Expect(err).NotTo(HaveOccurred())
		resolver.AssertResolveCalledWith(test_helpers.FakeIpfsBlock)
	})

	It("returns error if resolving block returns err", func() {
		getter := ipfs_mocks.NewMockGetter()
		resolver := ipfs_mocks.NewMockResolver()
		resolver.SetReturnErr(test_helpers.FakeError)
		reader := ipfs.NewIpldReader(getter, resolver)

		_, err := reader.Read(test_helpers.FakeCid)

		Expect(err).To(HaveOccurred())
		Expect(err).To(MatchError(test_helpers.FakeError))
	})
})
