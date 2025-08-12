package data

import (
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
	clt, err := ent.Open(c.Database.Driver, c.Database.Source)
	if err != nil {
		panic("must be new mysql client")
	}
	log.NewHelper(logger).Info("mysql connection success")
	return clt
}
