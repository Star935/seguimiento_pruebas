package main

import (
	"context"
	. "inventory/handlers"
	"log"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

// Punto de entrada del sistema
func main() {
	// Instancia de Echo
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	client, err := mongo.Connect(options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal(err)
	}

	if err := client.Ping(context.Background(), nil); err != nil {
		log.Fatal(err)
	}

	// Define la base de datos y la coleccion
	Collection = client.Database("seguimiento2").Collection("inventories")

	// Rutas para la gestion de inventarios
	e.GET("/inventories", GetInventories)
	e.GET("/inventories/:id", GetInventoryById)
	e.POST("/inventories", CreateInventory)
	e.DELETE("/inventories/:id", DeleteInventory)

	e.Logger.Fatal(e.Start(":8080"))
}