package main

import (
	"fmt"
	"log"

	"github.com/sattorovshoxrux3009/Restourants_back/config"
	"github.com/sattorovshoxrux3009/Restourants_back/server"
	"github.com/sattorovshoxrux3009/Restourants_back/storage"
	"github.com/sattorovshoxrux3009/Restourants_back/storage/repo" // Modellarni chaqiramiz
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	// "gorm.io/gorm/logger"
)

func main() {
	cfg := config.Load(".")
	// fmt.Println(cfg)

	mysqlUrl := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.Mysql.User,     // Foydalanuvchi nomi
		cfg.Mysql.Password, // Parol
		cfg.Mysql.Host,     // Host (masalan, "localhost")
		cfg.Mysql.Port,     // Port (masalan, "3306")
		cfg.Mysql.Database, // Ma'lumotlar bazasi nomi
	)

	// GORM bilan ulanish
	mysqlConn, err := gorm.Open(mysql.Open(mysqlUrl), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, // Jadvallarni ko‘plik shaklida yaratmasin
		},
		// Logger: logger.Default.LogMode(logger.Silent),
	})

	if err != nil {
		log.Fatal("Error opening database connection: ", err)
	}

	log.Println("Connection success!")

	// AutoMigrate qilish
	err = mysqlConn.AutoMigrate(
		&repo.Restaurant{},
		&repo.Admin{},
		&repo.EventPrice{},
		&repo.Menu{},
		&repo.Token{},
		&repo.SuperAdmin{},
		&repo.AdminRestaurantLimit{},
	)
	if err != nil {
		log.Fatal("Migrationda xatolik:", err)
	}

	fmt.Println("Migration muvaffaqiyatli yakunlandi!")

	strg := storage.NewStorage(mysqlConn)

	router := server.NewServer(&server.Options{
		Strg: strg,
	})

	if err := router.Listen(cfg.Port); err != nil {
		log.Fatal("Error starting server: ", err)
	}
}
