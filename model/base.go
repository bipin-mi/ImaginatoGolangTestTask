package model

import (
	"database/sql"
	"fmt"
	"github.com/benthor/clustersql"
	"github.com/go-sql-driver/mysql"
	"time"

	"ImaginatoGolangTestTask/shared/database"
)

type Model struct {
	Id        int64      `gorm:"primary_key" json:"id"`
	CreatedAt time.Time  `json:"created_at" sql:"DEFAULT:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time  `json:"updated_at" sql:"DEFAULT:CURRENT_TIMESTAMP"`
	DeletedAt *time.Time `json:"deleted_at,omitempty" sql:"index"`
}

//AutoMigrate for auto migrate the model tables in DB
func AutoMigrate() {
	conn := database.NewConnection()
	conn.GetDB().AutoMigrate(
		&Admin{},
	)
}

func ClusterDB() *sql.DB {
	mysqlDriver := mysql.MySQLDriver{}

	clusterDriver := clustersql.NewDriver(mysqlDriver)

	clusterDriver.AddNode("galera2", "root:some_characters@tcp(192.168.1.23:3306)/goapi?parseTime=true")
	clusterDriver.AddNode("galera1", "root:root@tcp(127.0.0.1:3306)/goapi?parseTime=true")

	sql.Register("myCluster", clusterDriver)
	db, err := sql.Open("myCluster", "galera")
	if err != nil {
		fmt.Println("Error(open) : ", err.Error())
	}

	return db
}
