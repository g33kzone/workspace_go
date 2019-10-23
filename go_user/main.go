package main

import (
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

func main() {
	fmt.Println("Hello World !!!")
	/*
		r := gin.Default()
		r.GET("/", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"message": "Hello World!",
			})
		})

		r.Run(":8080")
	*/
	const (
		host     = "odj-devops.clzebwtlh4xa.us-east-1.rds.amazonaws.com"
		port     = 5432
		user     = "odj-user"
		password = "godevops2019"
		dbname   = "odjdb"
	)

	type Employee struct {
		gorm.Model
		Birthday time.Time
		Age      int
		Name     string `gorm:"size:255"` // Default size for string is 255, reset it with this tag
		Num      int    `gorm:"AUTO_INCREMENT"`
	}

	type Person struct {
		ID        uint   `json:"id"`
		FirstName string `json:"firstname"`
		LastName  string `json:"lastname"`
	}

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+"password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	db, err := gorm.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	//db.CreateTable(&Employee{})

	db.AutoMigrate(&Person{})

	p1 := Person{FirstName: "John", LastName: "Doe"}
	p2 := Person{FirstName: "Jane", LastName: "Smith"}

	db.Create(&p1)
	db.Create(&p2)

	var p3 Person
	db.First(&p3)

	emp := Employee{Name: "Manish", Age: 18, Birthday: time.Now()}
	db.Create(&emp)

}
