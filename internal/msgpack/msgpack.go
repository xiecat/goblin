package msgpack

import (
	"bytes"

	"github.com/ugorji/go/codec"
)

var msgpackHandler codec.MsgpackHandle

// Encode - Encodes object with msgpack.
func Encode(obj interface{}) ([]byte, error) {
	buff := new(bytes.Buffer)
	encoder := codec.NewEncoder(buff, &msgpackHandler)
	err := encoder.Encode(obj)

	return buff.Bytes(), err
}

// Decode - Decodes object with msgpack.
func Decode(b []byte, v interface{}) error {
	decoder := codec.NewDecoderBytes(b, &msgpackHandler)

	return decoder.Decode(v)
}
