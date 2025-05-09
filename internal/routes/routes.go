package routes

import (
	"femProject/internal/app"

	"github.com/go-chi/chi/v5"
)

func SetupRoutes(app *app.Application) *chi.Mux {
	router := chi.NewRouter()

	router.Get("/health", app.HealthCheck)
	router.Get("/workouts/{id}", app.WorkoutHandler.HandleGetWorkoutByID)
	router.Post("/workouts", app.WorkoutHandler.HandleCreateWorkout)
	router.Put("/workouts/{id}", app.WorkoutHandler.HandleUpdateWorkoutByID)
	router.Delete("/workouts/{id}", app.WorkoutHandler.HandleDeleteWorkoutByID)

	router.Post("/users", app.UserHandler.HandleRegisterUser)

	return router
}
