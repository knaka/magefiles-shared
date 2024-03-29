package psqldb

import (
	"bufio"
	"database/sql"
	"fmt"
	_ "github.com/jackc/pgx/v5/stdlib"
	. "github.com/knaka/go-utils"
	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
	"log"
	"net/url"
	"os"
	"strings"
)

//goland:noinspection GoUnusedExportedType, GoUnnecessarilyExportedIdentifiers
type Db mg.Namespace

func execMainDatabaseDdl(ddlWithDbNamePlaceholder string, confirms bool) error {
	adminDbUrl := os.Getenv("ADMIN_DB_URL")
	dbUrl := os.Getenv("DB_URL")
	urlDb, err := url.Parse(dbUrl)
	if err != nil {
		return err
	}
	dbName := strings.Replace(urlDb.Path, "/", "", 1)
	if confirms {
		reader := bufio.NewReader(os.Stdin)
		fmt.Printf("Are you sure you want to execute the query “%s” on the database “%s”? (yes/no): ", ddlWithDbNamePlaceholder, dbName)
		confirmation := V(reader.ReadString('\n'))
		confirmation = strings.TrimSpace(confirmation)
		if confirmation != "yes" {
			return nil
		}
		V0(fmt.Printf("Executing query “%s” on the database “%s”...", ddlWithDbNamePlaceholder, dbName))
	}
	if adminDbUrl == "" {
		urlDb.Path = "/template1"
		adminDbUrl = urlDb.String()
	}
	db, err := sql.Open("pgx", adminDbUrl)
	if err != nil {
		return err
	}
	defer (func() { _ = db.Close() })()
	_, err = db.Exec(fmt.Sprintf(ddlWithDbNamePlaceholder, dbName))
	if err != nil {
		return err
	}
	return nil
}

// Create creates a new PostgreSQL database.
//
//goland:noinspection GoUnusedExportedFunction, GoUnnecessarilyExportedIdentifiers
func (Db) Create() error {
	err := execMainDatabaseDdl("CREATE DATABASE %s", false)
	if err != nil {
		return err
	}
	return nil
}

// Drop drops a PostgreSQL database if exists.
//
//goland:noinspection GoUnusedExportedFunction, GoUnnecessarilyExportedIdentifiers
func (Db) Drop() error {
	return execMainDatabaseDdl("DROP DATABASE IF EXISTS %s", true)
}

func execDbQuery(query string) error {
	dbUrl := os.Getenv("DB_URL")
	if dbUrl == "" {
		return nil
	}
	db, err := sql.Open("pgx", dbUrl)
	if err != nil {
		return err
	}
	defer (func() { _ = db.Close() })()
	log.Println("Executing query:", query)
	result, err := db.Exec(query)
	if err != nil {
		return err
	}
	log.Println("Done.")
	n, err := result.RowsAffected()
	if err != nil {
		log.Panicf("panic 27d2d65 (%v)", err)
	}
	log.Println("Rows affected:", n)
	return nil
}

// Seed inserts seed data into the database.
//
//goland:noinspection GoUnusedExportedFunction, GoUnnecessarilyExportedIdentifiers
func (Db) Seed() error {
	return sh.RunWith(nil, mg.GoCmd(), "run", "./cmd/db-seed")
}
