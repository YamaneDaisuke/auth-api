package db_test

import (
	"testing"

	"github.com/charakoba-com/auth-api/db"
)

func TestCreateUser(t *testing.T) {
	u := db.User{
		Name:     "Taro",
		Password: "none",
	}
	tx, err := db.BeginTx()
	if err != nil {
		t.Errorf("%s", err)
		return
	}
	err = u.Create(tx)
	if err != nil {
		t.Errorf("%s", err)
		return
	}
}

func TestLookupUser(t *testing.T) {
	u := db.User{}
	tx, err := db.BeginTx()
	if err != nil {
		t.Errorf("%s", err)
		return
	}
	err = u.Lookup(tx, "testuser")
	if err != nil {
		t.Errorf("%s", err)
		return
	}
	if u.Name != "testuser" {
		t.Errorf("%s != testuser", u.Name)
		return
	}
	if u.Password != "testpasswd" {
		t.Errorf("%s != testpasswd", u.Password)
		return
	}
}
