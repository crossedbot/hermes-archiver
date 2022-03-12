package replayer

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/crossedbot/common/golang/crypto/aes"
	"github.com/crossedbot/simplewarc"
	"github.com/stretchr/testify/require"
)

func TestParseLocator(t *testing.T) {
	paths := []string{
		"QmdfTbBqBPQ7VNxZEYEj14VmRuZBkqFbiwReogJgS1zR1n",
		"QmcRD4wkPPi6dig81r5sLj9Zm1gDCL4zgpEj9CfuRrGbzF",
	}
	locator := strings.Join(append([]string{"ipfs"}, paths...), "/")
	parsed, err := parseLocator(locator)
	require.Nil(t, err)
	require.Len(t, parsed, len(paths))
	for i, path := range paths {
		expected := strings.Join([]string{"/ipfs", path}, "/")
		require.Equal(t, expected, parsed[i].String())
	}
}

func TestDecompress(t *testing.T) {
	msg := "hello world"
	compressed := bytes.NewBuffer([]byte{})
	gzw := gzip.NewWriter(compressed)
	_, err := gzw.Write([]byte(msg))
	require.Nil(t, err)
	gzw.Close()
	reader, err := decompress(compressed, simplewarc.GzipCompression)
	require.Nil(t, err)
	b, err := ioutil.ReadAll(reader)
	require.Nil(t, err)
	require.Equal(t, []byte(msg), b)
}

func TestDecode(t *testing.T) {
	key := []byte("supersecret")
	salt := []byte("somesalt")
	msg := "hello world"
	aead, nonce, err := aes.NewKey(key, salt)
	require.Nil(t, err)
	compressed := bytes.NewBuffer([]byte{})
	gzw := gzip.NewWriter(compressed)
	_, err = gzw.Write([]byte(msg))
	require.Nil(t, err)
	gzw.Close()
	// Test unencrypted data
	encoded := base64.URLEncoding.EncodeToString(compressed.Bytes())
	reader, err := decode(encoded, nil, nil)
	require.Nil(t, err)
	b, err := ioutil.ReadAll(reader)
	require.Nil(t, err)
	require.Equal(t, []byte(msg), b)
	// Test encrypted data
	encrypted := aead.Seal(nil, nonce, compressed.Bytes(), nil)
	encoded = base64.URLEncoding.EncodeToString(encrypted)
	reader, err = decode(encoded, aead, nonce)
	require.Nil(t, err)
	b, err = ioutil.ReadAll(reader)
	require.Nil(t, err)
	require.Equal(t, []byte(msg), b)
}
