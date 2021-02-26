package controller

import (
	"context"
	"errors"
	"sync"

	"github.com/crossedbot/common/golang/config"
	"github.com/crossedbot/common/golang/crypto"
	cdxjdb "github.com/crossedbot/hermes-archiver/pkg/database"
	"github.com/crossedbot/hermes-archiver/pkg/replayer"
	"github.com/crossedbot/hermes-archiver/pkg/replayer/models"
)

var (
	ErrorReplayNotFound = errors.New("replay not found")
)

type Controller interface {
	Replay(id string) (models.Replay, error)
}

type controller struct {
	ctx context.Context
	key []byte
	rpl replayer.Replayer
}

type Config struct {
	DatabaseAddr string `toml:"database_addr"`

	// Encyption configuraiton
	EncryptionKey  string `toml:"encryption_key"`
	EncryptionSalt string `toml:"encryption_salt"`

	// IPFS configuration
	IpfsAddress string `toml:"ipfs_address"`
}

var control Controller
var controllerOnce sync.Once
var V1 = func() Controller {
	controllerOnce.Do(func() {
		var cfg Config
		if err := config.Load(&cfg); err != nil {
			panic(err)
		}
		ctx := context.Background()
		extKey := []byte{}
		if cfg.EncryptionKey != "" {
			extKey = crypto.ExtendKey(
				[]byte(cfg.EncryptionKey),
				[]byte(cfg.EncryptionSalt),
			)
		}
		db := cdxjdb.New(ctx, cfg.DatabaseAddr)
		rpl, err := replayer.New(ctx, cfg.IpfsAddress, db)
		if err != nil {
			panic(err)
		}
		control = New(ctx, extKey, rpl)
	})
	return control
}

func New(ctx context.Context, key []byte, rpl replayer.Replayer) Controller {
	return &controller{ctx: ctx, key: key, rpl: rpl}
}

func (c *controller) Replay(id string) (models.Replay, error) {
	return c.rpl.Replay(id, c.key)
}
