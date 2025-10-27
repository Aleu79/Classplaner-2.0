package database

import (
	"classplanner/internal/infrastructure/database"
	"classplanner/internal/model"
	"classplanner/pkg/utils"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/sunary/sqlize"
)

type MigrationInstance struct {
	SQL *sqlize.Sqlize
}

var MgInstance *MigrationInstance

func migrateModels() []interface{} {
	return []interface{}{
		model.User{},
		model.Role{},
		model.Calendar{},
		model.Classes{},
		model.Comment{},
		model.Submission{},
		model.Tasks{},
		model.UserClass{},
		model.Address{},
	}
}

func PrepareSQL() *sqlize.Sqlize {
	utils.LoadEnv()

	migrationFolder := os.Getenv("MIGRATIONS_FOLDER")
	if migrationFolder == "" {
		log.Fatal("MIGRATIONS_FOLDER env var not set")
	}

	versionFile := os.Getenv("MIGRATIONS_VERSION")
	if versionFile == "" {
		log.Fatal("MIGRATIONS_VERSION env var not set")
	}

	if err := os.MkdirAll(migrationFolder, 0755); err != nil {
		log.Fatalf("error creating migrations folder: %v", err)
	}

	sqlz := sqlize.NewSqlize(
		sqlize.WithSqlTag("sql"),
		sqlize.WithMigrationFolder(migrationFolder),
		sqlize.WithPostgresql(),
	)

	if err := sqlz.FromObjects(migrateModels()...); err != nil {
		log.Fatalf("sqlize FromObjects: %v", err)
	}

	if err := sqlz.FromMigrationFolder(); err != nil {
		log.Fatalf("sqlize FromMigrationFolder: %v", err)
	}

	MgInstance = &MigrationInstance{SQL: sqlz}
	return sqlz
}

func getCurrentVersion(filePath string) int64 {
	data, err := os.ReadFile(filePath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return 0
		}
		log.Fatalf("error reading migration version file: %v", err)
	}

	if len(data) == 0 {
		return 0
	}

	version, err := strconv.ParseInt(string(data), 10, 64)
	if err != nil {
		log.Fatalf("error parsing migration version: %v", err)
	}

	return version
}

// writeMigrationVersion saves the current migration version hash to a file
func writeMigrationVersion(sqlz *sqlize.Sqlize, filePath string, version int64) {
	if err := os.MkdirAll(filepath.Dir(filePath), 0755); err != nil {
		log.Fatalf("error creating version folder: %v", err)
	}

	// Use WriteFile to create/overwrite the file with only the version number
	if err := os.WriteFile(filePath, []byte(fmt.Sprint(version)), 0644); err != nil {
		log.Fatalf("error writing migration version: %v", err)
	}
	log.Println("migration version saved:", version)
}

// write migrations up
func Migrate() {
	sqlz := PrepareSQL()
	db := database.Connect().DB

	// Get the UP SQL string for all pending changes
	migrationUpSQL := sqlz.StringUp()
	fmt.Println("Executing UP migration SQL:\n", migrationUpSQL)

	log.Println("executing migration up...")
	// Execute the entire migration script
	if _, err := db.Exec(migrationUpSQL); err != nil {
		log.Fatalf("pgsql DatabaseUp: %v", err)
	}

	// Save the new version *after* successful migration
	writeMigrationVersion(sqlz, os.Getenv("MIGRATIONS_VERSION"), sqlz.HashValue())
	log.Println("database successfully migrated")
}

// write migrations down
func Reset() {
	sqlz := PrepareSQL()
	db := database.Connect().DB

	// Get the DOWN SQL string to revert changes
	migrationDownSQL := splitInLines(sqlz.StringDown()) //split in lines
	fmt.Println("Executing DOWN migration SQL:\n", migrationDownSQL)

	log.Println("executing migration down...")
	for i := 0; i < len(migrationDownSQL); i++ {
		if _, err := db.Exec(migrationDownSQL[i]); err != nil {
			log.Fatalf("pgsql DatabaseDown: %v", err)
		}
	}
	// Reset the version record after successful revert
	writeMigrationVersion(sqlz, os.Getenv("MIGRATIONS_VERSION"), 0)
	log.Println("database successfully reset")
}

func splitInLines(sql string) []string {
	return strings.Split(sql, "\n")
}
