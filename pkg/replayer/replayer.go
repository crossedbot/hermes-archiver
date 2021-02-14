package replayer

import (
	"bytes"
	"compress/gzip"
	"context"
	"crypto/cipher"
	"encoding/base64"
	"io"
	"io/ioutil"
	"strings"
	"time"

	ipfsfiles "github.com/ipfs/go-ipfs-files"
	ipfshttpapi "github.com/ipfs/go-ipfs-http-client"
	ipfspath "github.com/ipfs/go-path"
	ipfscorepath "github.com/ipfs/interface-go-ipfs-core/path"
	ma "github.com/multiformats/go-multiaddr"

	"github.com/crossedbot/common/golang/crypto/aes"
	cdxjdb "github.com/crossedbot/hermes-archiver/pkg/indexer/database"
	"github.com/crossedbot/hermes-archiver/pkg/replayer/models"
	"github.com/crossedbot/simplewarc"
)

type Replayer interface {
	Replay(id string, key []byte) (models.Replay, error)
}

type replayer struct {
	ctx    context.Context
	db     cdxjdb.CdxjRecords
	client *ipfshttpapi.HttpApi
}

func New(
	ctx context.Context,
	ipfsAddr string,
	db cdxjdb.CdxjRecords,
) (Replayer, error) {
	multiaddr, err := ma.NewMultiaddr(ipfsAddr)
	if err != nil {
		return nil, err
	}
	client, err := ipfshttpapi.NewApi(multiaddr)
	return &replayer{
		ctx:    ctx,
		db:     db,
		client: client,
	}, nil
}

func (rp *replayer) Replay(id string, key []byte) (models.Replay, error) {
	rec, err := rp.db.Get(id)
	if err != nil {
		return models.Replay{}, err
	}
	loc, err := parseLocator(rec.Content.Locator)
	if err != nil {
		return models.Replay{}, err
	}
	nonce, err := base64.URLEncoding.DecodeString(rec.Content.EncryptionNonce)
	if err != nil {
		return models.Replay{}, err
	}
	aead, err := aes.AesGcmKey(key)
	if err != nil {
		return models.Replay{}, err
	}
	payload := []byte{}
	if len(loc) > 1 {
		msg, err := pull(rp.ctx, rp.client, loc[1].String())
		if err != nil {
			return models.Replay{}, err
		}
		decoded, err := decode(msg, aead, nonce)
		if err != nil {
			return models.Replay{}, err
		}
		payload, err = ioutil.ReadAll(decoded)
		if err != nil {
			return models.Replay{}, err
		}
	}
	return models.Replay{
		Uri:       rec.Content.Uri,
		Title:     rec.Content.Title,
		Sha:       rec.Content.Sha,
		Timestamp: rec.Timestamp.Format(time.RFC3339),
		Content:   string(payload),
	}, nil
}

func parseLocator(locator string) ([]ipfspath.Path, error) {
	paths := []ipfspath.Path{}
	parts := strings.Split(locator, "/")
	for _, part := range parts[1:] {
		path, err := ipfspath.ParseCidToPath(part)
		if err != nil {
			return nil, err
		}
		paths = append(paths, path)
	}
	return paths, nil
}

func pull(
	ctx context.Context,
	client *ipfshttpapi.HttpApi,
	cid string,
) (string, error) {
	fn, err := client.Unixfs().Get(ctx, ipfscorepath.New(cid))
	if err != nil {
		return "", err
	}
	sz, err := fn.Size()
	if err != nil {
		return "", err
	}
	b := make([]byte, sz)
	if _, err := ipfsfiles.ToFile(fn).Read(b); err != nil {
		return "", err
	}
	return string(b), nil
}

func decompress(src io.Reader, c simplewarc.CompressionType) (io.Reader, error) {
	r, err := gzip.NewReader(src)
	if err != nil {
		return nil, err
	}
	defer r.Close()
	dst := bytes.NewBuffer([]byte{})
	_, err = io.Copy(dst, r)
	return dst, err
}

func decode(msg string, key cipher.AEAD, nonce []byte) (io.Reader, error) {
	decoded, err := base64.URLEncoding.DecodeString(msg)
	if err != nil {
		return nil, err
	}
	decrypted, err := key.Open(nil, nonce, decoded, nil)
	if err != nil {
		return nil, err
	}
	return decompress(bytes.NewReader(decrypted), simplewarc.GzipCompression)
}
