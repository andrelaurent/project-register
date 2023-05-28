package database

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/andrelaurent/project-register/config"
	"github.com/andrelaurent/project-register/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Dbinstance struct {
	Db *gorm.DB
}

var DB Dbinstance

func Connect() {
	p := config.Config("DB_PORT")
	port, err := strconv.ParseUint(p, 10, 32)

	if err != nil {
		fmt.Println("error parsing")
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=Asia/Shanghai", config.Config("DB_HOST"), config.Config("DB_USER"), config.Config("DB_PASSWORD"), config.Config("DB_NAME"), port)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		log.Fatal("failed to connect to database.\n", err)
		os.Exit(2)
	}

	log.Println("connected")
	db.Logger = logger.Default.LogMode(logger.Info)
	log.Println("running migations")
	db.AutoMigrate(
		&models.Company{},
		&models.ProjectType{},
		&models.User{},
		&models.City{},
		&models.Province{},
		&models.Prospect{},
		&models.Project{},
		&models.Contact{},
		&models.Client{},
		&models.Locations{},
		&models.Employment{},
		&models.ClientContact{},
	)

	DB = Dbinstance{
		Db: db,
	}
}
