package main

import (
	"log"
	"net/http"
	"project/car-zone/db"
	"project/car-zone/store/car"
	"project/car-zone/store/engine"

	carService "project/car-zone/service/car"
	engineService "project/car-zone/service/engine"

	carHandler "project/car-zone/handler/car"
	engineHandler "project/car-zone/handler/engine"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}

	db, err := db.InitDB()

	if err != nil {
		log.Fatalf("failed initialise database %v", err)
	}

	carStore := car.New(db)
	engineStore := engine.New(db)

	carService := carService.New(carStore)
	engineService := engineService.New(engineStore)

	carHandler := carHandler.New(carService)
	engineHandler := engineHandler.New(engineService)

	router := mux.NewRouter()
	router.HandleFunc("/cars/{id}", carHandler.GetCarById).Methods("GET")
	router.HandleFunc("/cars", carHandler.GetCarByBrand).Methods("GET")
	router.HandleFunc("/cars", carHandler.CreateCar).Methods("POST")

	router.HandleFunc("/engine/{id}", engineHandler.GetEngineById).Methods("GET")
	router.HandleFunc("/engine", engineHandler.CreateEngine).Methods("POST")

	log.Println("Server Running on port 8080")
	http.ListenAndServe(":8080", router)
}
