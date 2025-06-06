package main

import (
	"context"
	"log"
	"net/http"
	"strings"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

// Estructura de objeto
type Object struct {
	ID	 		bson.ObjectID `json:"id,ompitempty" bson:"_id,ompitempty"`
	Name 		string 		  `json:"name" bson:"name"`
	Description string 		  `json:"description" bson:"description"`
}

// Estructura de inventario
type Inventory struct {
	ID	 		bson.ObjectID `json:"id,ompitempty" bson:"_id,ompitempty"`
	Name 		string 		  `json:"name" bson:"name"`
	Objects 	[]Object	  `json:"object" bson:"object"`
}

// Coleccion de Mongo
var collection *mongo.Collection

// Recupera todos los inventarios junto con sus objetos anidados
func getInventories(c echo.Context) error {
	// Recupera todos los inventarios
	cur, err := collection.Find(context.Background(), bson.M{})
	// Valuda si recupera los inventarios
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}
	
	// Lista de objetos
	var inventories []Inventory
	
	// Almacena en la lista de inventarios todos los inventarios recuperados y valida si la operacion es exitosa
	if err := cur.All(context.Background(), &inventories); err != nil{
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	// Retorna estado de respuesta ok y todos los inventarios recuperados
	return c.JSON(http.StatusOK, inventories)
}

// Recupera un inventario mediante su id
func getInventoryById(c echo.Context) error {
	// Recupera el parametro de consulta el id
	idParam := c.Param("id")

	// Convierte el id en ObjectID
	id, err := bson.ObjectIDFromHex(idParam)
	// Valida si la operacion fue exitosa
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Id invalido"})
	}

	// Instancia de Inventory
	var inventory Inventory

	// Recupera el inventario mediante su id y lo decodifica en el espacio de memoria de la instancia de inventario
	err = collection.FindOne(context.Background(), bson.M{"_id": id}).Decode(&inventory)
	// Valuda si no existe el documento
	if err == mongo.ErrNoDocuments {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "Inventario no encontrado"})
	} else if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	// Retorna estado de respuesta ok y todos los inventarios recuperados
	return c.JSON(http.StatusFound, inventory)
}

// Crea un nuevo inventario
func createInventory(c echo.Context) error {
	// Instancia de Inventory
	var inventory Inventory

	// Almacena el body recuperado de la peticion al inventory
	if err := c.Bind(&inventory); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Input invalido"})
	}

	// Valida la existencia del campo del nombre
	if strings.TrimSpace(inventory.Name) == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "El nombre es obligatorio"})
	}

	// Valida la existencia de objetos
	if inventory.Objects == nil || len(inventory.Objects) == 0 {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Los objetos son obligatorios"})
	}

	// Valida la lista de objetos
	for i, obj := range inventory.Objects {
		// Valida el nombre del objeto
		if strings.TrimSpace(obj.Name) == "" {
			return c.JSON(http.StatusBadRequest, echo.Map{"error": "El objeto en el indice " + strconv.Itoa(i) + " carece de nombre, es necesario asignarle nombre"})
		}

		// Valida la descripcion del objeto
		if strings.TrimSpace(obj.Description) == "" {
			return c.JSON(http.StatusBadRequest, echo.Map{"error": "El objeto en el indice " + strconv.Itoa(i) + " carece de descripcion, es necesario asignarle descripcion"})
		}
	}

	// Inserta el nuevo inventario a la base de datos
	res, err := collection.InsertOne(context.Background(), inventory)
	// Valida si la operacion fue exitosa
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	// Asigna el id generado del nuevo inventario insertado
	inventory.ID = res.InsertedID.(bson.ObjectID)

	return c.JSON(http.StatusCreated, inventory)
}

// Elimina un inventario
func deleteInventory(c echo.Context) error {
	// Recupera el parametro de consulta el id
	idParam := c.Param("id")

	// Convierte el id en ObjectID
	id, err := bson.ObjectIDFromHex(idParam)
	// Valida si la operacion fue exitosa
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Id invalido"})
	}
	
	// Realiza operacion de eliminado mediante el id recuperado del parametro de consulta
	res, err := collection.DeleteOne(context.Background(), bson.M{"_id": id})
	// Valida si la operacion fue exitosa
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}
	
	// Valida si se elimino algun documento
	if res.DeletedCount == 0 {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "Inventario no encontrado"})
	}

	return c.NoContent(http.StatusNoContent)
}

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
	collection = client.Database("seguimiento2").Collection("inventories")

	// Rutas para la gestion de inventarios
	e.GET("/inventories", getInventories)
	e.GET("/inventories/:id", getInventoryById)
	e.POST("/inventories", createInventory)
	e.DELETE("/inventories/:id", deleteInventory)

	e.Logger.Fatal(e.Start(":8080"))
}