package main

import (
	"ImaginatoGolangTestTask/model"
	"ImaginatoGolangTestTask/routers"
	"ImaginatoGolangTestTask/shared/database"
	"ImaginatoGolangTestTask/shared/log"
	"ImaginatoGolangTestTask/shared/utils"
	"context"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"os"
	"strconv"
	"time"
)

//Execution starts from main function
func main() {
	godotenv.Load()

	level, _ := strconv.ParseUint(os.Getenv("Level"), 10, 32)
	maxAge, _ := strconv.ParseUint(os.Getenv("MaxAge"), 10, 32)

	log.Init("admin", os.Getenv("path"), logrus.Level(level), time.Duration(maxAge))

	database.Init()
	log.GetLog().Info("", "DB connected")

	model.AutoMigrate()

	rt := routers.NewRouter()
	rt.Setup()

	go rt.Run()

	utils.GracefulStop(log.GetLog(), func(ctx context.Context) error {
		var err error
		if err = rt.Close(ctx); err != nil {
			return err
		}
		if err = database.Close(); err != nil {
			return err
		}
		return nil
	})
}
