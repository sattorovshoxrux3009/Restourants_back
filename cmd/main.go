package main

import (
	"fmt"
	"log"

	"github.com/sattorovshoxrux3009/Restourants_back/config"
	_ "github.com/sattorovshoxrux3009/Restourants_back/docs" // Swagger docs
	"github.com/sattorovshoxrux3009/Restourants_back/seeds"
	"github.com/sattorovshoxrux3009/Restourants_back/server"
	"github.com/sattorovshoxrux3009/Restourants_back/storage"
	"github.com/sattorovshoxrux3009/Restourants_back/storage/repo"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

// @title           Restaurants API
// @version         1.0
// @description     This is a Restaurants server.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:3000
// @BasePath  /

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

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
	fmt.Println(mysqlUrl)
	// GORM bilan ulanish
	mysqlConn, err := gorm.Open(mysql.Open(mysqlUrl), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, // Jadvallarni ko‘plik shaklida yaratmasin
		},
		Logger: logger.Default.LogMode(logger.Silent),
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
		log.Fatal("Error running migrations: ", err)
	}
	log.Println("Migration muvaffaqiyatli yakunlandi!")

	// Storage yaratish
	strg := storage.NewStorage(mysqlConn)

	// Run database seeds
	seeder := seeds.NewSeeds(strg)
	seeder.RunAll()

	// Create and start server
	app := server.NewServer(&server.Options{
		Strg: strg,
	})

	log.Fatal(app.Listen(":3000"))
}
