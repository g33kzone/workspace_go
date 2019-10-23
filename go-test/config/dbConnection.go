package config

import (
	"fmt"

	"github.com/jinzhu/gorm"
)

// DBGormConn -
func DBGormConn(dbConn DBConfig) *gorm.DB {

	fmt.Println(fmt.Sprintf("Attempting to connect %s Database on %s:%d", dbConn.database, dbConn.server, dbConn.port))

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s sslmode=disable",
		dbConn.server, dbConn.port, dbConn.username, dbConn.database,
		dbConn.password)

	dbConnection, err := gorm.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer dbConnection.Close()

	fmt.Println("Database connection successful...")

	return dbConnection
}
