package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"syscall"

	"github.com/crossedbot/common/golang/config"
	"github.com/crossedbot/common/golang/logger"
	"github.com/crossedbot/common/golang/server"
	"github.com/crossedbot/common/golang/service"

	"github.com/crossedbot/hermes-archiver/cmd"
	cdxjdb "github.com/crossedbot/hermes-archiver/pkg/database"
	"github.com/crossedbot/hermes-archiver/pkg/indexer"
	"github.com/crossedbot/hermes-archiver/pkg/indexer/controller"
)

const (
	// Exit codes
	FatalExitCode = iota + 1
)

var (
	// Build variables
	Version = "-"
	Build   = "-"
)

// Config represents an indexer's configuration
type Config struct {
	Host                string `toml:"host"`
	Port                int    `toml:"port"`
	ReadTimeout         int    `toml:"read_timeout"`  // in seconds
	WriteTimeout        int    `toml:"write_timeout"` // in seconds
	WarcDirectory       string `toml:"warc_directory"`
	DatabaseAddr        string `toml:"database_addr"`
	DropDatabaseOnStart bool   `toml:"drop_database_on_start"`

	// Encyption configuraiton
	EncryptionKey  string `toml:"encryption_key"`
	EncryptionSalt string `toml:"encryption_salt"`

	// IPFS configuration
	IpfsAddress string `toml:"ipfs_address"`
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	svc := service.New(ctx)
	err := svc.Run(run, syscall.SIGINT, syscall.SIGTERM)
	// cancel irregardless of error state
	cancel()
	if err != nil {
		fatal("Error: %s", err)
	}
}

// fatal exits with fatal after printing the given formatted message
func fatal(format string, a ...interface{}) {
	logger.Error(fmt.Errorf(format, a...))
	os.Exit(FatalExitCode)
}

// newServer returns a new indexer server using the given configuration
func newServer(c Config) server.Server {
	hostport := net.JoinHostPort(c.Host, strconv.Itoa(c.Port))
	srv := server.New(
		hostport,
		c.ReadTimeout,
		c.WriteTimeout,
	)
	for _, route := range controller.Routes {
		srv.Add(
			route.Handler,
			route.Method,
			route.Path,
			route.ResponseSettings...,
		)
	}
	controller.V1()
	return srv
}

// run serves the main entry point into the program
func run(ctx context.Context) error {
	f := cmd.ParseFlags()
	if f.Version {
		fmt.Printf(
			"%s version %s, build %s\n",
			filepath.Base(os.Args[0]),
			Version, Build,
		)
		return nil
	}
	config.Path(f.ConfigFile)
	var c Config
	if err := config.Load(&c); err != nil {
		return err
	}
	srv := newServer(c)
	db, err := cdxjdb.New(ctx, c.DatabaseAddr, c.DropDatabaseOnStart)
	if err != nil {
		return err
	}
	in, err := indexer.New(ctx, c.IpfsAddress, c.WarcDirectory, db)
	if err != nil {
		return err
	}
	in.SetEncryptionKey([]byte(c.EncryptionKey), []byte(c.EncryptionSalt))
	if err := in.Start(); err != nil {
		return err
	}
	logger.Info(fmt.Sprintf("Started indexer on %s", c.WarcDirectory))
	if err := srv.Start(); err != nil {
		return err
	}
	logger.Info(fmt.Sprintf("Listening on %s:%d", c.Host, c.Port))
	<-ctx.Done()
	logger.Info("Received signal, shutting down...")
	return nil
}
