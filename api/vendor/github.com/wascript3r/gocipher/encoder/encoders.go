package encoder

import (
	"encoding/base64"
	"encoding/hex"
)

func HexEncode(data []byte) []byte {
	bs := make([]byte, hex.EncodedLen(len(data)))
	hex.Encode(bs, data)
	return bs
}

func HexDecode(data []byte) ([]byte, error) {
	bs := make([]byte, hex.DecodedLen(len(data)))

	_, err := hex.Decode(bs, data)
	if err != nil {
		return nil, err
	}

	return bs, nil
}

func Base64Encode(data []byte) []byte {
	bs := make([]byte, base64.RawStdEncoding.EncodedLen(len(data)))
	base64.RawStdEncoding.Encode(bs, data)
	return bs
}

func Base64Decode(data []byte) ([]byte, error) {
	bs := make([]byte, base64.RawStdEncoding.DecodedLen(len(data)))

	_, err := base64.RawStdEncoding.Decode(bs, data)
	if err != nil {
		return nil, err
	}

	return bs, nil
}
