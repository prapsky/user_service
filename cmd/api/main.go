package main

import (
	"database/sql"
	"fmt"

	"github.com/indrasaputra/hashids"
	"github.com/labstack/echo/v4"
	"github.com/prapsky/user_service/common/config"
	"github.com/prapsky/user_service/common/logger/zerolog"
	handler_user "github.com/prapsky/user_service/internal/handler/user"
	repository "github.com/prapsky/user_service/internal/repository"
	service_auth "github.com/prapsky/user_service/service/auth"
	service_user_register "github.com/prapsky/user_service/service/user/register"
)

const (
	dbDriver = "postgres"
)

func main() {
	cfg := buildConfig(".env")
	db := buildDB(cfg)
	defer db.Close()

	hash, err := hashids.NewHashID(uint(cfg.HashIDMinLength), cfg.HashIDSalt)
	checkError(err)
	hashids.SetHasher(hash)

	logger := zerolog.NewZeroLog()

	userRepo := repository.NewUser(db, logger)
	authService := service_auth.NewJwtAuthService(service_auth.JwtAuthServiceOptions{
		PrivateKey: cfg.Auth.PrivateKey,
		PublicKey:  cfg.Auth.PublicKey,
	})
	userService := service_user_register.NewRegisterUserService(userRepo, authService, logger)
	userHandler := handler_user.NewRegisterUserHandler(&userService)

	e := echo.New()
	e.POST("/register", userHandler.Register)
	e.Logger.Fatal(e.Start(":8080"))
}
func buildConfig(env string) *config.Config {
	cfg, err := config.NewConfig(env)
	checkError(err)
	return cfg
}

func buildDB(cfg *config.Config) *sql.DB {
	sqlCfg := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.Username,
		cfg.Database.Password,
		cfg.Database.Name)

	db, err := sql.Open(dbDriver, sqlCfg)
	checkError(err)

	db.SetMaxOpenConns(cfg.Database.MaxOpenConns)
	db.SetMaxIdleConns(cfg.Database.MaxIdleConns)

	return db
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}
