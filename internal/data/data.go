package data

import (
	"os"
	"teslatrack/internal/conf"
	"teslatrack/internal/data/ent"

	"github.com/go-kratos/kratos/v2/log"
	_ "github.com/go-sql-driver/mysql"
	"github.com/google/wire"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(NewData, NewGreeterRepo)

// Data .
type Data struct {
	// db *ent.Edge
	db *ent.Client
}

// NewData .
func NewData(c *conf.Data, logger log.Logger) (*Data, func(), error) {
	db := mustNewMysqlSqlClient(c, logger)

	cleanup := func() {
		log.NewHelper(logger).Info("closing the data resources")
		_ = db.Close()
	}

	return &Data{db: db}, cleanup, nil
}

func mustNewMysqlSqlClient(c *conf.Data, logger log.Logger) *ent.Client {
	helper := log.NewHelper(logger)

	var databaseSource = os.Getenv("DATABASE_SOURCE")
	if databaseSource == "" {
		databaseSource = c.Database.Source
	}

	var databaseDriver = os.Getenv("DATABASE_DRIVER")
	if databaseDriver == "" {
		databaseDriver = c.Database.Driver
	}

	helper.Debugw("msg", "mysql connecting", "databaseDriver", databaseDriver, "databaseSource", databaseSource)
	clt, err := ent.Open(databaseDriver, databaseSource)
	if err != nil {
		panic("must be new mysql client")
	}

	helper.Info("mysql connection success")
	return clt
}
