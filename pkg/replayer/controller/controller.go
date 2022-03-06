package controller

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/crossedbot/common/golang/config"
	"github.com/crossedbot/common/golang/crypto"
	cdxjdb "github.com/crossedbot/hermes-archiver/pkg/database"
	"github.com/crossedbot/hermes-archiver/pkg/replayer"
	"github.com/crossedbot/hermes-archiver/pkg/replayer/models"
)

var (
	// Errors
	ErrorReplayNotFound = errors.New("replay not found")
)

// Controller represents an interface of an WARC-CDXJ record replayer's
// controller
type Controller interface {
	// Replay returns a WARC replay for the given record Id
	Replay(id string) (models.Replay, error)
}

// controller implements the Controller interface
type controller struct {
	ctx context.Context
	key []byte
	rpl replayer.Replayer
}

// Config represents the configuration of a Replayer's controller
type Config struct {
	DatabaseAddr        string `toml:"database_addr"`
	DropDatabaseOnStart bool   `toml:"drop_database_on_start"`

	// Encyption configuraiton
	EncryptionKey  string `toml:"encryption_key"`
	EncryptionSalt string `toml:"encryption_salt"`

	// IPFS configuration
	IpfsAddress string `toml:"ipfs_address"`
}

// control is a Contoller singleton
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
		db, err := cdxjdb.New(
			ctx,
			cfg.DatabaseAddr,
			cfg.DropDatabaseOnStart,
		)
		if err != nil {
			panic(fmt.Errorf(
				"Controller: failed to connect to database at "+
					"address ('%s') with error: %s",
				cfg.DatabaseAddr, err,
			))
		}
		rpl, err := replayer.New(ctx, cfg.IpfsAddress, db)
		if err != nil {
			panic(err)
		}
		control = New(ctx, extKey, rpl)
	})
	return control
}

// New returns a new Controller
func New(ctx context.Context, key []byte, rpl replayer.Replayer) Controller {
	return &controller{ctx: ctx, key: key, rpl: rpl}
}

// Replay returns a WARC replay for the given CDXJ record Id
func (c *controller) Replay(id string) (models.Replay, error) {
	return c.rpl.Replay(id, c.key)
}
