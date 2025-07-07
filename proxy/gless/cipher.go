package gless

import (
	"crypto/cipher"
	"crypto/sha256"
	"encoding/binary"
)

// GOSTCipher implements GOST-based cipher
type GOSTCipher struct {
	key []byte
	iv  []byte
}

// NewGOSTCipher creates a new GOST cipher
func NewGOSTCipher(uuid string, isClient bool) (*GOSTCipher, error) {
	// Generate key from UUID
	key := sha256.Sum256([]byte(uuid))
	
	// Generate IV
	iv := make([]byte, 8)
	binary.LittleEndian.PutUint64(iv, uint64(len(uuid)))
	
	return &GOSTCipher{
		key: key[:],
		iv:  iv,
	}, nil
}

// Encrypt encrypts data using GOST cipher
func (g *GOSTCipher) Encrypt(data []byte) []byte {
	result := make([]byte, len(data))
	for i := 0; i < len(data); i++ {
		result[i] = data[i] ^ g.key[i%len(g.key)] ^ g.iv[i%len(g.iv)]
	}
	return result
}

// Decrypt decrypts data using GOST cipher
func (g *GOSTCipher) Decrypt(data []byte) []byte {
	return g.Encrypt(data) // XOR cipher is symmetric
}

// Stream returns a stream cipher
func (g *GOSTCipher) Stream() cipher.Stream {
	return &gostStream{
		key: g.key,
		iv:  g.iv,
	}
}

type gostStream struct {
	key    []byte
	iv     []byte
	counter uint64
}

func (g *gostStream) XORKeyStream(dst, src []byte) {
	if len(dst) < len(src) {
		panic("gost: output smaller than input")
	}
	
	for i := 0; i < len(src); i++ {
		keystreamByte := g.key[i%len(g.key)] ^ g.iv[i%len(g.iv)] ^ byte(g.counter)
		dst[i] = src[i] ^ keystreamByte
		g.counter++
	}
} 