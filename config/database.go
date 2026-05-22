package config

import (
	"fmt"
	"log"
	"os"

	"donation-api/models"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Bypass file .env, menggunakan variabel sistem (Railway).")
	}

	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Jakarta",
			os.Getenv("DB_HOST"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"),
			os.Getenv("DB_NAME"), os.Getenv("DB_PORT"),
		)
	}

	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Gagal terhubung ke database!\n", err)
	}

	// === AUTO MIGRATE ===
	err = database.AutoMigrate(
		&models.User{},
		&models.AppSetting{},
		&models.Category{},
		&models.Campaign{},
		&models.Transaction{},
		&models.Article{},
	)
	
	if err != nil {
		log.Fatal("Gagal migrasi database!\n", err)
	}

	DB = database
	fmt.Println("Koneksi Database dan Migrasi Tabel Berhasil!")
}