package DatabaseConnection

import (
	"github.com/alistairfink/Steak/Backend/Utilities"
	"github.com/go-pg/pg"
)

func Connect(config *Utilities.Config) *pg.DB {
	db := pg.Connect(&pg.Options{
		Addr:     config.DB.Address + ":" + config.DB.Port,
		User:     config.DB.Username,
		Password: config.DB.Password,
		Database: config.DB.DBName,
	})

	return db
}

func Close(db *pg.DB) {
	db.Close()
}
