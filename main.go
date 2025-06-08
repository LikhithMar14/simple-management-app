package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/LikhithMar14/management/database"
	carHandler "github.com/LikhithMar14/management/handler/car"
	engineHandler "github.com/LikhithMar14/management/handler/engine"
	"github.com/LikhithMar14/management/handler/login"
	"github.com/LikhithMar14/management/middleware"
	"github.com/LikhithMar14/management/migrations"
	carService "github.com/LikhithMar14/management/service/car"
	engineService "github.com/LikhithMar14/management/service/engine"
	carStore "github.com/LikhithMar14/management/store/car"
	engineStore "github.com/LikhithMar14/management/store/engine"
	"github.com/go-chi/chi/v5"
	"github.com/pressly/goose/v3"


	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	database.InitDB()
	defer database.CloseDB()

	db := database.GetDB()
	goose.SetBaseFS(migrations.FS)
	if err := goose.Up(db, "."); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}
	log.Println("Migrations applied successfully.")

	carStore := carStore.NewCarStore(db)
	carService := carService.NewCarService(carStore)
	carHandler := carHandler.NewCarHandler(carService)

	engineStore := engineStore.NewEngineStore(db)
	engineService := engineService.NewEngineService(engineStore)
	engineHandler := engineHandler.NewEngineHandler(engineService)

	router := chi.NewRouter()
	login.InitGoogleOauthConfig()
	login.InitGitHubOauthConfig()
	router.Get("/auth/google", login.GoogleLoginHandler)
	router.Get("/auth/google/callback", login.GoogleCallbackHandler)
	router.Get("/auth/github", login.GitHubLoginHandler)
	router.Get("/auth/github/callback", login.GitHubCallbackHandler)

	router.Post("/login", login.LoginHandler)
	router.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		response := map[string]string{"status": "ok"}
		json.NewEncoder(w).Encode(response)
	})


	router.Route("/", func(r chi.Router) {
		r.Use(middleware.AuthMiddleware)
		
		
		r.Get("/cars/{id}", carHandler.GetCarByID)
		r.Get("/cars", carHandler.GetCarsByBrand)
		r.Post("/cars", carHandler.CreateCar)
		r.Put("/cars/{id}", carHandler.UpdateCar)
		r.Delete("/cars/{id}", carHandler.DeleteCar)

		
		r.Get("/engine/{id}", engineHandler.GetEngineByID)
		r.Post("/engine", engineHandler.CreateEngine)
		r.Put("/engine/{id}", engineHandler.UpdateEngine)
		r.Delete("/engine/{id}", engineHandler.DeleteEngine)
	})

	log.Println("Server starting on :8080")
	err = http.ListenAndServe(":8080", router)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}



