package database

import (
	"github.com/jinzhu/gorm"

	"database/sql"
	"flag"
	"os"
	"strconv"
	"sync"

	seeds "ImaginatoGolangTestTask/seeder"
)

// IConnection ITransaction is
type IConnection interface {
	GetDB() *gorm.DB
}

// GormDB is
type connection struct {
	db *gorm.DB

	readonly      bool
	isTranscation bool
}

var db *gorm.DB

var connectionOnce sync.Once

func Init() {
	connectionOnce.Do(func() {
		var err error

		connectionString := os.Getenv("ConnectionString")
		dialect := os.Getenv("Dialect")
		logMode := os.Getenv("LogMode")
		db, err = gorm.Open(dialect, connectionString)
		if err != nil {
			panic(err)
		}
		boolLogMode, _ := strconv.ParseBool(logMode)
		db.LogMode(boolLogMode)

		//seeder code command execute
		flag.Parse()
		args := flag.Args()

		if len(args) >= 1 {
			switch args[0] {
			case "seed":
				conn, _ := sql.Open(dialect, connectionString)
				seeds.Execute(conn, args[1:]...)
				os.Exit(0)
			}
		}
	})
}

func NewConnection() IConnection {
	return &connection{
		db,
		false,
		true,
	}
}

func Close() error {
	if db != nil {
		return db.Close()
	}

	return nil
}

// GetDB is
func (selfConn *connection) GetDB() *gorm.DB {
	return selfConn.db
}
