package databasemanager

import (
	"github.com/luiz-pv9/dixte-analytics/dixteconfig"
	"io/ioutil"
	"log"
	"path/filepath"
	"testing"
)

func TestConnection(t *testing.T) {
	dc, err := dixteconfig.LoadFromFile(filepath.Join("..", "config.json"))
	if err != nil {
		t.Error(err)
	}
	dc.AssignDefaults()
	db, err := Connect(dc)
	if err != nil {
		t.Error(err)
	}
	err = db.Conn.Ping()
	if err != nil {
		t.Error(err)
	}
}

func TestTablesNames(t *testing.T) {
	dc, _ := dixteconfig.LoadFromFile(filepath.Join("..", "config.json"))
	db, _ := Connect(dc)
	_, err := db.TablesNames()
	if err != nil {
		t.Error(err)
	}
}

func TestResetDatabase(t *testing.T) {
	dc, _ := dixteconfig.LoadFromFile(filepath.Join("..", "config.json"))
	db, _ := Connect(dc)
	err := db.Reset()
	if err != nil {
		t.Error(err)
	}

	tablesNames, err := db.TablesNames()
	if err != nil {
		t.Error(nil)
	}

	if len(tablesNames) > 0 {
		t.Error("Didn't delete all tables")
	}
}

func TestCreateMigrationsTable(t *testing.T) {
	dc, _ := dixteconfig.LoadFromFile(filepath.Join("..", "config.json"))
	db, _ := Connect(dc)
	err := db.Reset()
	if err != nil {
		t.Error(err)
	}

	err = db.CreateMigrationsTable()
	if err != nil {
		t.Error(err)
	}

	tablesNames, _ := db.TablesNames()
	if len(tablesNames) != 1 || tablesNames[0] != "migrations" {
		t.Error("Didn't create migrations table")
	}
}

func TestMigrationTableExists(t *testing.T) {
	dc, _ := dixteconfig.LoadFromFile(filepath.Join("..", "config.json"))
	db, _ := Connect(dc)
	err := db.Reset()
	if err != nil {
		t.Error(err)
	}

	if db.HasMigrationsTable() != false {
		t.Error("Shouldn't have a migrations table")
	}

	err = db.CreateMigrationsTable()
	if err != nil {
		t.Error(err)
	}

	if db.HasMigrationsTable() != true {
		t.Error("Didn't detect migrations table")
	}
}

func TestMigrate(t *testing.T) {
	dc, _ := dixteconfig.LoadFromFile(filepath.Join("..", "config.json"))
	db, _ := Connect(dc)
	db.Reset()
	migratedFiles, err := db.Migrate(filepath.Join("..", "migrations"))
	if err != nil {
		t.Error(err)
	}
	if len(migratedFiles) < 1 {
		t.Error("Didn't run any migration")
	}
	tables, _ := db.TablesNames()
	if len(tables) < 1 {
		t.Error("Didn't create any tables")
	}
	files, _ := ioutil.ReadDir(filepath.Join("..", "migrations"))
	latestMigratedFiles, _ := db.MigratedFiles()
	if len(files) != len(latestMigratedFiles) {
		t.Error("Didn't store the migrations in the database")
	}

	// Running migrations again shouldn't change anything
	migratedFiles, err = db.Migrate(filepath.Join("..", "migrations"))
	if err != nil {
		t.Error(err)
	}
	if len(migratedFiles) != 0 {
		t.Error("Migrations should not be runned")
	}
}

func TestMigrateNewFiles(t *testing.T) {
	dc, _ := dixteconfig.LoadFromFile(filepath.Join("..", "config.json"))
	db, _ := Connect(dc)
	db.Reset()
	migratedFiles, err := db.Migrate(filepath.Join("..", "migrations"))
	if err != nil {
		t.Error(err)
	}
	log.Print("Migrated files")
	log.Print(migratedFiles)
	if len(migratedFiles) < 1 {
		t.Error("Didn't migrate any files")
	}

	migration := "SELECT schema_name FROM information_schema.schemata"
	err = ioutil.WriteFile(filepath.Join("..", "migrations", "999-test.sql"),
		[]byte(migration), 0644)

	if err != nil {
		t.Error(err)
	}

	// Run migrations and create another file
	// Run migrations again and see changes
}
