package main

import (
	"log"

	"gorm.io/gen"

	"gosync/internal/config"
	"gosync/internal/database"
	"gosync/internal/models"
)

func main() {
	err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	db, err := database.Connect()
	if err != nil {
		log.Fatal(err)
	}

	g := gen.NewGenerator(gen.Config{
		OutPath: "internal/api/repository",
		Mode:    gen.WithDefaultQuery | gen.WithQueryInterface,
	})
	g.UseDB(db)

	g.ApplyBasic(
		models.User{},
		models.Device{},
		models.Directory{},
		models.File{},
		models.FileVersion{},
		models.SyncEvent{},
		models.RefreshToken{},
	)

	g.Execute()
}
