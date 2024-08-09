package main

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/bernardinorafael/kn-server/internal/config"
	"github.com/bernardinorafael/kn-server/internal/core/application/service"
	"github.com/bernardinorafael/kn-server/internal/infra/auth"
	"github.com/bernardinorafael/kn-server/internal/infra/database/gorm/gormrepo"
	database "github.com/bernardinorafael/kn-server/internal/infra/database/pg"
	"github.com/bernardinorafael/kn-server/internal/infra/http/middleware"
	"github.com/bernardinorafael/kn-server/internal/infra/http/route"
	"github.com/bernardinorafael/kn-server/internal/infra/s3client"
	"github.com/bernardinorafael/kn-server/internal/infra/twilioclient"
	"github.com/bernardinorafael/kn-server/pkg/logger"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

/*
 TODO: implement Graceful Shutdown
 TODO: implement Swagger
 TODO: resolve excess logic in the controller
*/

func main() {
	router := chi.NewRouter()
	router.Use(middleware.WithRecoverPanic)
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))
	l := slog.New(logger.NewLog(nil))

	env, err := config.NewConfig()
	if err != nil {
		l.Error("failed load env", "error", err)
		return
	}

	db, err := database.Connect(l, env.DSN)
	if err != nil {
		l.Error("error connecting db", "error", err)
		panic(err)
	}

	/*
	 TODO: change JWT implementation into a service
	*/
	jwtAuth, err := auth.NewJWTAuth(l, env.JWTSecret)
	if err != nil {
		l.Error("cannot initialize jwt auth")
		panic(err)
	}

	s3Client, err := s3client.New(env)
	if err != nil {
		l.Error("cannot initialize aws s3 client")
		panic(err)
	}

	twilioClient := twilioclient.New(env)

	/*
	 Repositories
	*/
	userRepo := gormrepo.NewUserRepo(db)
	productRepo := gormrepo.NewProductRepo(db)
	teamRepo := gormrepo.NewTeamRepo(db)
	/*
	 Services
	 TODO: make Option Pattern for services
	*/
	s3Service := service.NewS3Service(l, s3Client)
	twilioSMSService := service.NewTwilioSMSService(l, env.TwilioServiceID, twilioClient)
	twilioEmailService := service.NewTwilioEmailService(l, env.TwilioServiceID, twilioClient)
	authService := service.NewAuthService(l, twilioSMSService, twilioEmailService, userRepo)
	userService := service.NewUserService(l, userRepo)
	teamService := service.NewTeamService(l, teamRepo)
	productService := service.NewProductService(service.WithProductParams{
		Log:         l,
		Env:         env,
		ProductRepo: productRepo,
		FileService: s3Service,
	})
	/*
	 Controllers
	*/
	authHandler := route.NewAuthHandler(l, authService, jwtAuth, env)
	productHandler := route.NewProductHandler(l, productService, jwtAuth)
	userHandler := route.NewUserHandler(l, userService, jwtAuth)
	teamHandler := route.NewTeamHandler(l, teamService, jwtAuth)
	/*
	 Registering controllers
	*/
	authHandler.RegisterRoute(router)
	productHandler.RegisterRoute(router)
	userHandler.RegisterRoute(router)
	teamHandler.RegisterRoute(router)

	l.Info("server listening", "port", env.Port)

	err = http.ListenAndServe(":"+env.Port, router)
	if err != nil {
		l.Error("error connecting web server", "error", err)
		os.Exit(1)
	}
}
