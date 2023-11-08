package testing_helpers

import (
	"database/sql"
	"testing"

	"desprit/hammertime/src/config"
	"desprit/hammertime/src/db"
)

func PrepareDB(t *testing.T) (*sql.DB, error) {
	c := config.GetConfig()
	d, err := db.Open(c)
	if err != nil {
		return nil, err
	}
	if err = db.InitDB(d, c); err != nil {
		return nil, err
	}
	return d, err
}
