package main

import (
	"femProject/internal/app"
	"femProject/internal/routes"
	"flag"
	"fmt"
	"net/http"
	"time"
)

func main() {

	var port int
	flag.IntVar(&port, "port", 8080, "Port to run the application on")
	flag.Parse()

	app, err := app.NewApplication()
	if err != nil {
		panic(err)
	}

	defer app.DB.Close()

	router := routes.SetupRoutes(app)

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", port),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
		Handler:      router,
	}

	app.Logger.Printf("Starting application on port %d\n", port)

	err = server.ListenAndServe()
	if err != nil {
		app.Logger.Fatal(err)
	}
}
