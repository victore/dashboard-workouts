package app

import (
	"database/sql"
	"femProject/internal/api"
	"femProject/internal/store"
	"femProject/migrations"
	"fmt"
	"log"
	"net/http"
	"os"
)

type Application struct {
	Logger         *log.Logger
	WorkoutHandler *api.WorkoutHandler
	UserHandler    *api.UserHandler
	DB             *sql.DB
	TokenHandler   *api.TokenHandler
}

func NewApplication() (*Application, error) {
	pgDB, err := store.Open()
	if err != nil {
		return nil, err
	}

	err = store.MigrateFS(pgDB, migrations.FS, ".")
	if err != nil {
		panic(err)
	}

	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	// stores go here
	workoutStore := store.NewPostgresWorkoutStore(pgDB)
	userStore := store.NewPostgresUserStore(pgDB)
	tokenStore := store.NewPostgresTokenStore(pgDB)

	// handlers go here
	workoutHandler := api.NewWorkoutHandler(workoutStore, logger)
	userHandler := api.NewUserHandler(userStore, logger)
	tokenHandler := api.NewTokenHandler(tokenStore, userStore, logger)

	app := &Application{
		Logger:         logger,
		WorkoutHandler: workoutHandler,
		UserHandler:    userHandler,
		DB:             pgDB,
		TokenHandler:   tokenHandler,
	}

	return app, nil
}

func (app *Application) HealthCheck(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Status is available")
}
