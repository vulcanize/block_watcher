package rocks

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/ethereum/go-ethereum/rlp"
)

// it appears that the first byte could be anywhere from 248-255 (0xf8 - 0xff)
// https://github.com/ethereum/wiki/wiki/RLP
// however, I've yet to see a block header where the first byte is not 249
const ethBlockHeaderFirstByte = 249

type Decompressor interface {
	Decompress(raw []byte) ([]byte, error)
}

type EthBlockDecompressor struct {
}

type decompressError struct {
	msg string
	err error
}

func (d decompressError) Error() string {
	if d.err != nil {
		return fmt.Sprintf("%s: %s", d.msg, d.err.Error())
	} else {
		return d.msg
	}
}

func NewDecompressError(msg string, err error) *decompressError {
	return &decompressError{
		msg: msg,
		err: err,
	}
}

func (d EthBlockDecompressor) Decompress(raw []byte) ([]byte, error) {
	if len(raw) == 0 {
		return nil, NewDecompressError("No input to decode", nil)
	}
	if raw[0] != ethBlockHeaderFirstByte {
		return nil, NewDecompressError(fmt.Sprintf("Received byte array not beginning with 249: %d", raw[0]), nil)
	}
	_, content, _, err := rlp.Split(raw)
	if err != nil {
		return nil, NewDecompressError("Error decoding block data:", err)
	}
	var result []byte
	keepIterating := true
	for keepIterating {
		if len(content) == 0 {
			keepIterating = false
			break
		}
		isCompressed, idx := leadsWithInvalidRLP(content[:2])
		if isCompressed {
			decompressed := getMatchingCommonRLP(idx)
			result = append(result, decompressed...)
			content = content[2:]
		} else {
			originalContent := content
			var contentToAppend []byte
			_, contentToAppend, content, err = rlp.Split(content)
			if err != nil {
				return nil, err
			}
			chunkHeader := getChunkHeaderData(originalContent, contentToAppend)
			result = append(append(result, chunkHeader...), contentToAppend...)
		}
	}
	rlpHeaderData := getRLPHeaderData(result)
	result = append(rlpHeaderData, result...)
	return result, nil
}

func getChunkHeaderData(original []byte, separator []byte) []byte {
	split := bytes.Split(original, separator)
	chunkHeader := split[0]
	if len(split[1]) > 0 {
		keepChecking := true
		remaining := split[1]
		for keepChecking {
			if remaining[0] == 0 {
				chunkHeader = append(chunkHeader, 0)
				if len(remaining) == 1 {
					keepChecking = false
				}
				remaining = remaining[1:]
			} else {
				keepChecking = false
			}
		}
	}
	return chunkHeader
}

func getRLPHeaderData(result []byte) []byte {
	rlpLength := len(result)
	var bs []byte
	bs = append(bs, ethBlockHeaderFirstByte)
	tmp := make([]byte, 2)
	binary.BigEndian.PutUint16(tmp, uint16(rlpLength))
	bs = append(bs, tmp...)
	return bs
}

func leadsWithInvalidRLP(data []byte) (bool, int) {
	for i := 0; i < len(InvalidRLPs); i++ {
		if bytes.Equal(data, InvalidRLPs[i]) {
			return true, i
		}
	}
	return false, 0
}

func getMatchingCommonRLP(idx int) []byte {
	if idx < len(CommonRLPs) {
		return CommonRLPs[idx]
	}
	return CommonRLPs[len(CommonRLPs)-1]
}
