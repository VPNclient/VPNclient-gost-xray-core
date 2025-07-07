package encoding

import (
	"io"

	"github.com/xtls/xray-core/common/buf"
	"github.com/xtls/xray-core/common/errors"
	"github.com/xtls/xray-core/common/protocol"
	"github.com/xtls/xray-core/common/net"
	"github.com/xtls/xray-core/proxy/gless"
)

// DecodeRequestHeader decodes the request header from reader
func DecodeRequestHeader(reader io.Reader, validator gless.Validator) (*protocol.RequestHeader, error) {
	// Simplified implementation - in real implementation this would decode
	// the actual GLess protocol header
	buffer := buf.New()
	defer buffer.Release()

	// Read a small amount to simulate header reading
	if _, err := buffer.ReadFullFrom(reader, 1); err != nil {
		return nil, errors.New("failed to read header").Base(err)
	}

	// Create a mock request header
	request := &protocol.RequestHeader{
		Command: protocol.RequestCommandTCP,
		Address: net.IPAddress([]byte{127, 0, 0, 1}),
		Port:    80,
	}

	return request, nil
}

// EncodeRequestHeader encodes the request header to writer
func EncodeRequestHeader(writer io.Writer, request *protocol.RequestHeader, addons *Addons) error {
	// Simplified implementation - in real implementation this would encode
	// the actual GLess protocol header
	buffer := buf.New()
	defer buffer.Release()

	// Write a simple header
	buffer.WriteByte(0x01) // Version
	buffer.WriteByte(byte(request.Command))
	
	// Write address
	if request.Address != nil {
		buffer.Write(request.Address.IP())
	}
	
	// Write port
	buffer.WriteByte(byte(request.Port >> 8))
	buffer.WriteByte(byte(request.Port))

	_, err := writer.Write(buffer.Bytes())
	return err
}

// NewGOSTReader creates a new GOST reader
func NewGOSTReader(reader io.Reader, cipher *gless.GOSTCipher) io.Reader {
	return &gostReader{
		reader: reader,
		cipher: cipher,
	}
}

// NewGOSTWriter creates a new GOST writer
func NewGOSTWriter(writer io.Writer, cipher *gless.GOSTCipher) io.Writer {
	return &gostWriter{
		writer: writer,
		cipher: cipher,
	}
}

type gostReader struct {
	reader io.Reader
	cipher *gless.GOSTCipher
}

func (g *gostReader) Read(p []byte) (n int, err error) {
	n, err = g.reader.Read(p)
	if n > 0 {
		// Decrypt the data
		decrypted := g.cipher.Decrypt(p[:n])
		copy(p[:n], decrypted)
	}
	return n, err
}

type gostWriter struct {
	writer io.Writer
	cipher *gless.GOSTCipher
}

func (g *gostWriter) Write(p []byte) (n int, err error) {
	// Encrypt the data
	encrypted := g.cipher.Encrypt(p)
	return g.writer.Write(encrypted)
} 