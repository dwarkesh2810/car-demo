package test

import (
	"car_demo/conf"
	"testing"

	"github.com/beego/beego/v2/client/orm"
)

func TestInit(t *testing.T) {
	t.Run("Evironment", func(t *testing.T) {
		conf.LoadEnv(".")
		result := conf.EnvConfig.BaseStoragePath
		expected := "assets/img"

		if result != expected {
			t.Fatalf("Expected %s but got %s", expected, result)
		}
		t.Log("env load success")
	})

	t.Run("Driver registration", func(t *testing.T) {
		result := orm.RegisterDriver("postgres", orm.DRPostgres)

		if result != nil {
			t.Fatalf("Driver registration Failed err := %s", result)
		}
		t.Log("Driver registration success")
	})

	t.Run("Database registration", func(t *testing.T) {
		result := orm.RegisterDataBase("default1", "postgres", "user=postgres password=1234 dbname=postgres sslmode=disable")

		if result != nil {
			t.Fatalf("Database registration Failed err := %s", result)
		}
		t.Log("Database registration success")
	})

	t.Run("SyncDB", func(t *testing.T) {
		result := orm.RunSyncdb("default", false, true)

		if result != nil {
			t.Fatalf("SyncDB Failed err := %s", result)
		}
		t.Log("Database registration success")
	})

}
