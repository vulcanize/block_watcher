package rocks_test

import (
	"errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/8thlight/block_watcher/pkg/db/rocks"
	"github.com/8thlight/block_watcher/test_helpers"
)

var _ = Describe("Rocks database", func() {
	It("converts block number to byte array to fetch block hash from the reader", func() {
		reader := test_helpers.NewMockRocksDatabaseReader()
		reader.SetReturnHashBytes([]byte{9, 9, 9, 9, 9})
		reader.SetReturnHeaderBytes([]byte{8, 8, 8, 8, 8})
		rocksDB := rocks.NewRocksDatabase(test_helpers.NewMockDecompressor(), reader)

		_, err := rocksDB.GetBlockHeaderByBlockNumber(0)

		Expect(err).NotTo(HaveOccurred())
		Expect(reader.GetBlockHashCalled).To(BeTrue())
		blockZeroAsByteArray := []byte{0, 0, 0, 0}
		expectedBlockHashKey := append([]byte{1}, blockZeroAsByteArray...)
		Expect(reader.PassedBlockHashKey).To(Equal(expectedBlockHashKey))
	})

	It("returns error if fetching block hash returns error", func() {
		reader := test_helpers.NewMockRocksDatabaseReader()
		fakeErr := errors.New("failed")
		reader.SetGetHashError(fakeErr)
		rocksDB := rocks.NewRocksDatabase(test_helpers.NewMockDecompressor(), reader)

		_, err := rocksDB.GetBlockHeaderByBlockNumber(1234)

		Expect(err).To(HaveOccurred())
		Expect(err).To(MatchError(fakeErr))
	})

	It("returns error if reader returns empty byte array for block hash", func() {
		reader := test_helpers.NewMockRocksDatabaseReader()
		reader.SetReturnHashBytes([]byte{})
		rocksDB := rocks.NewRocksDatabase(test_helpers.NewMockDecompressor(), reader)

		_, err := rocksDB.GetBlockHeaderByBlockNumber(0)

		Expect(err).To(HaveOccurred())
		Expect(err).To(MatchError(rocks.ErrBlockHashNotFound))
	})

	It("pops first byte off of block hash to fetch block header from the reader", func() {
		reader := test_helpers.NewMockRocksDatabaseReader()
		blockHash := []byte{160, 1, 2, 3, 4}
		reader.SetReturnHashBytes(blockHash)
		reader.SetReturnHeaderBytes([]byte{8, 8, 8, 8, 8})
		rocksDB := rocks.NewRocksDatabase(test_helpers.NewMockDecompressor(), reader)

		_, err := rocksDB.GetBlockHeaderByBlockNumber(1234)

		Expect(err).NotTo(HaveOccurred())
		Expect(reader.GetBlockHeaderCalled).To(BeTrue())
		Expect(reader.PassedBlockHeaderKey).To(Equal(blockHash[1:]))
	})

	It("returns error if fetching block header returns error", func() {
		reader := test_helpers.NewMockRocksDatabaseReader()
		reader.SetReturnHashBytes([]byte{9, 9, 9, 9, 9})
		fakeErr := errors.New("failed")
		reader.SetGetHeaderError(fakeErr)
		rocksDB := rocks.NewRocksDatabase(test_helpers.NewMockDecompressor(), reader)

		_, err := rocksDB.GetBlockHeaderByBlockNumber(1234)

		Expect(err).To(HaveOccurred())
		Expect(err).To(MatchError(fakeErr))
	})

	It("returns error if reader returns empty byte array for block header", func() {
		reader := test_helpers.NewMockRocksDatabaseReader()
		reader.SetReturnHashBytes([]byte{9, 9, 9, 9, 9})
		reader.SetReturnHeaderBytes([]byte{})
		rocksDB := rocks.NewRocksDatabase(test_helpers.NewMockDecompressor(), reader)

		_, err := rocksDB.GetBlockHeaderByBlockNumber(0)

		Expect(err).To(HaveOccurred())
		Expect(err).To(MatchError(rocks.ErrBlockHeaderNotFound))
	})

	It("decompresses block header returned by the reader", func() {
		decompressor := test_helpers.NewMockDecompressor()
		reader := test_helpers.NewMockRocksDatabaseReader()
		reader.SetReturnHashBytes([]byte{9, 9, 9, 9, 9})
		bytesToDecode := []byte{1, 2, 3, 4, 5}
		reader.SetReturnHeaderBytes(bytesToDecode)
		rocksDB := rocks.NewRocksDatabase(decompressor, reader)

		_, err := rocksDB.GetBlockHeaderByBlockNumber(1234)

		Expect(err).NotTo(HaveOccurred())
		Expect(decompressor.Called).To(BeTrue())
		Expect(decompressor.PassedBytes).To(Equal(bytesToDecode))
	})

	It("returns error if decompression returns error", func() {
		decompressor := test_helpers.NewMockDecompressor()
		fakeError := errors.New("failed")
		decompressor.SetError(fakeError)
		reader := test_helpers.NewMockRocksDatabaseReader()
		reader.SetReturnHashBytes([]byte{9, 9, 9, 9, 9})
		reader.SetReturnHeaderBytes([]byte{8, 8, 8, 8, 8})
		rocksDB := rocks.NewRocksDatabase(decompressor, reader)

		_, err := rocksDB.GetBlockHeaderByBlockNumber(1234)

		Expect(err).To(HaveOccurred())
		Expect(err).To(MatchError(fakeError))
	})
})
