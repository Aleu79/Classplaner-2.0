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

	"github.com/sunary/sqlize"
)

type MigrationInstance struct {
	SQL *sqlize.Sqlize
}

var MgInstance *MigrationInstance

// migrateModels defines models to be loaded
func migrateModels() []interface{} {
	return []interface{}{
		model.User{},
		model.Role{},
		model.Address{},
	}
}

// PrepareSQL sets up sqlize configuration and models
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

	migrationFile := os.Getenv("MIGRATIONS_FILE")
	if migrationFile == "" {
		log.Fatal("MIGRATIONS_FILE env var not set")
	}

	// create migrations folder if not exists
	if err := os.MkdirAll(migrationFolder, 0755); err != nil {
		log.Fatalf("error creating migrations folder: %v", err)
	}

	sqlz := sqlize.NewSqlize(
		sqlize.WithSqlTag("sql"),
		sqlize.WithMigrationFolder(migrationFolder),
		sqlize.WithPostgresql(),
	)

	// generate schemas
	if err := sqlz.FromObjects(migrateModels()...); err != nil {
		log.Fatalf("sqlize FromObjects: %v", err)
	}

	// load migration files
	if err := sqlz.FromMigrationFolder(); err != nil {
		log.Fatalf("sqlize FromMigrationFolder: %v", err)
	}

	// write migration files
	log.Println("writing migrations to", migrationFile)
	if err := sqlz.WriteFilesWithVersion(migrationFile, sqlz.HashValue(), false); err != nil {
		log.Fatalf("sqlize WriteFilesWithVersion: %v", err)
	}

	MgInstance = &MigrationInstance{SQL: sqlz}
	return sqlz
}

// create version saver for the migrations files
func writeMigrationVersion(sqlz *sqlize.Sqlize, filePath string, version int64) string {
	if err := os.MkdirAll(filepath.Dir(filePath), 0755); err != nil {
		log.Fatalf("error creating version folder: %v", err)
	}

	file, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil && !errors.Is(err, os.ErrExist) {
		log.Fatalf("error opening migration version file: %v", err)
	}
	defer file.Close()

	info, _ := file.Stat()
	if info.Size() == 0 {
		if _, err := file.WriteString(fmt.Sprint(version)); err != nil {
			log.Fatalf("error writing migration version: %v", err)
		}
		log.Println("initialized migration version:", version)
	}
	sql := sqlz.StringUp()
	fmt.Println(sql)
	return sql
}

// write migrations down
func Reset() {
	sqlz := PrepareSQL()
	db := database.Connect().DB

	log.Println("executing migration down...")
	if _, err := db.Exec(sqlz.StringDown()); err != nil {
		log.Fatalf("pgsql DatabaseDown: %v", err)
	}
	log.Println("database successfully reset")
}

// write migrations up
func Migrate() {
	sqlz := PrepareSQL()
	db := database.Connect().DB

	log.Println("executing migration up...")
	if _, err := db.Exec(writeMigrationVersion(sqlz, os.Getenv("MIGRATIONS_FILE"), sqlz.HashValue())); err != nil {
		log.Fatalf("pgsql DatabaseUp: %v", err)
	}
	log.Println("database successfully migrated")
}
