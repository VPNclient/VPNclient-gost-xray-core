package encoding

import (
	"io"

	"github.com/xtls/xray-core/common/buf"
	"github.com/xtls/xray-core/common/errors"
	"github.com/xtls/xray-core/common/protocol"
	"github.com/xtls/xray-core/proxy/gless"
	"google.golang.org/protobuf/proto"
)

func EncodeHeaderAddons(buffer *buf.Buffer, addons *Addons) error {
	switch addons.Flow {
	case gless.GRV:
		bytes, err := proto.Marshal(addons)
		if err != nil {
			return errors.New("failed to marshal addons protobuf value").Base(err)
		}
		if err := buffer.WriteByte(byte(len(bytes))); err != nil {
			return errors.New("failed to write addons protobuf length").Base(err)
		}
		if _, err := buffer.Write(bytes); err != nil {
			return errors.New("failed to write addons protobuf value").Base(err)
		}
	default:
		if err := buffer.WriteByte(0); err != nil {
			return errors.New("failed to write addons protobuf length").Base(err)
		}
	}

	return nil
}

func DecodeHeaderAddons(buffer *buf.Buffer, reader io.Reader) (*Addons, error) {
	addons := new(Addons)
	buffer.Clear()
	if _, err := buffer.ReadFullFrom(reader, 1); err != nil {
		return nil, errors.New("failed to read addons protobuf length").Base(err)
	}

	if length := int32(buffer.Byte(0)); length != 0 {
		buffer.Clear()
		if _, err := buffer.ReadFullFrom(reader, length); err != nil {
			return nil, errors.New("failed to read addons protobuf value").Base(err)
		}

		if err := proto.Unmarshal(buffer.Bytes(), addons); err != nil {
			return nil, errors.New("failed to unmarshal addons protobuf value").Base(err)
		}

		// Verification.
		switch addons.Flow {
		default:
		}
	}

	return addons, nil
}

// EncodeBodyAddons returns a Writer that auto-encrypt content written by caller.
func EncodeBodyAddons(writer io.Writer, request *protocol.RequestHeader, addons *Addons) buf.Writer {
	switch addons.Flow {
	case gless.GRV:
		return buf.NewWriter(writer)
	default:
		return buf.NewWriter(writer)
	}
}

// DecodeBodyAddons returns a Reader from which caller read and decrypts content.
func DecodeBodyAddons(reader io.Reader, request *protocol.RequestHeader, addons *Addons) buf.Reader {
	switch addons.Flow {
	case gless.GRV:
		return buf.NewReader(reader)
	default:
		return buf.NewReader(reader)
	}
} 