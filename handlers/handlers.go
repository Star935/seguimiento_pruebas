package handlers

import (
	"context"
	"net/http"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

// Estructura de objeto
type Object struct {
	ID	 		primitive.ObjectID `json:"id,ompitempty" bson:"_id,ompitempty"`
	Name 		string 		  `json:"name" bson:"name"`
	Description string 		  `json:"description" bson:"description"`
}

// Estructura de inventario
type Inventory struct {
	ID	 		primitive.ObjectID `json:"id,ompitempty" bson:"_id,ompitempty"`
	Name 		string 		  `json:"name" bson:"name"`
	Objects 	[]Object	  `json:"object" bson:"object"`
}

// Coleccion de Mongo
var Collection *mongo.Collection

// Recupera todos los inventarios junto con sus objetos anidados
func GetInventories(c echo.Context) error {
	// Valida la conexion a la coleccion
	if Collection == nil {
		return c.JSON(http.StatusNotFound, echo.Map{"data": nil, "status": http.StatusNotFound, "message": "Sin conexion a la colección"})
	}

	// Recupera todos los inventarios
	cur, err := Collection.Find(context.Background(), bson.M{})
	// Valuda si recupera los inventarios
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"data": nil, "status": http.StatusInternalServerError, "message": err.Error()})
	}
	
	// Lista de objetos
	var inventories []Inventory
	
	// Almacena en la lista de inventarios todos los inventarios recuperados y valida si la operacion es exitosa
	if err := cur.All(context.Background(), &inventories); err != nil{
		return c.JSON(http.StatusInternalServerError, echo.Map{"data": nil, "status": http.StatusInternalServerError, "message": err.Error()})
	}

	if len(inventories) == 0 {
		return c.JSON(http.StatusNotFound, echo.Map{"data": nil, "status": http.StatusNotFound, "message": "No se encontraron inventarios"})
	}

	// Retorna estado de respuesta ok y todos los inventarios recuperados
	return c.JSON(http.StatusOK, inventories)
}

// Recupera un inventario mediante su id
func GetInventoryById(c echo.Context) error {
	// Recupera el parametro de consulta el id
	idParam := c.Param("id")

	// Convierte el id en ObjectID
	id, err := primitive.ObjectIDFromHex(idParam)
	// Valida si la operacion fue exitosa
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"data": nil, "status": http.StatusBadRequest, "message": "Id invalido"})
	}

	// Instancia de Inventory
	var inventory Inventory

	// Valida la conexion a la coleccion
	if Collection == nil {
		return c.JSON(http.StatusNotFound, echo.Map{"data": nil, "status": http.StatusNotFound, "message": "Sin conexion a la colección"})
	}

	// Recupera el inventario mediante su id y lo decodifica en el espacio de memoria de la instancia de inventario
	err = Collection.FindOne(context.Background(), bson.M{"_id": id}).Decode(&inventory)
	// Valuda si no existe el documento
	if err == mongo.ErrNoDocuments {
		return c.JSON(http.StatusNotFound, echo.Map{"data": nil, "status": http.StatusNotFound, "message": "Inventario no encontrado"})
	} else if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"data": nil, "status": http.StatusInternalServerError, "message": err.Error()})
	}

	// Retorna estado de respuesta ok y todos los inventarios recuperados
	return c.JSON(http.StatusFound, inventory)
}

// Crea un nuevo inventario
func CreateInventory(c echo.Context) error {
	var inventory Inventory

	if err := c.Bind(&inventory); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"data": nil, "status": http.StatusBadRequest, "message": "Input invalido"})
	}

	if strings.TrimSpace(inventory.Name) == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"data": nil, "status": http.StatusBadRequest, "message": "El nombre es obligatorio"})
	}

	if inventory.Objects == nil || len(inventory.Objects) == 0 {
		return c.JSON(http.StatusBadRequest, echo.Map{"data": nil, "status": http.StatusBadRequest, "message": "Los objetos son obligatorios"})
	}

	// Generar nuevo ID para el inventario
	inventory.ID = primitive.NewObjectID()

	for i := range inventory.Objects {
		if strings.TrimSpace(inventory.Objects[i].Name) == "" {
			return c.JSON(http.StatusBadRequest, echo.Map{"data": nil, "status": http.StatusBadRequest, "message": "El objeto en el indice " + strconv.Itoa(i) + " carece de nombre"})
		}
		if strings.TrimSpace(inventory.Objects[i].Description) == "" {
			return c.JSON(http.StatusBadRequest, echo.Map{"data": nil, "status": http.StatusBadRequest, "message": "El objeto en el indice " + strconv.Itoa(i) + " carece de descripcion"})
		}
		// Generar nuevo ID para cada objeto
		inventory.Objects[i].ID = primitive.NewObjectID()
	}

	if Collection == nil {
		return c.JSON(http.StatusNotFound, echo.Map{"data": nil, "status": http.StatusNotFound, "message": "Sin conexion a la colección"})
	}

	_, err := Collection.InsertOne(context.Background(), inventory)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"data": nil, "status": http.StatusInternalServerError, "message": err.Error()})
	}

	return c.JSON(http.StatusCreated, inventory)
}
// Elimina un inventario
func DeleteInventory(c echo.Context) error {
	// Recupera el parametro de consulta el id
	idParam := c.Param("id")

	// Convierte el id en ObjectID
	id, err := primitive.ObjectIDFromHex(idParam)
	// Valida si la operacion fue exitosa
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"data": nil, "status": http.StatusBadRequest, "message": "Id invalido"})
	}

	// Valida la conexion a la coleccion
	if Collection == nil {
		return c.JSON(http.StatusNotFound, echo.Map{"data": nil, "status": http.StatusNotFound, "message": "Sin conexion a la colección"})
	}

	// Realiza operacion de eliminado mediante el id recuperado del parametro de consulta
	res, err := Collection.DeleteOne(context.Background(), bson.M{"_id": id})
	// Valida si la operacion fue exitosa
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"data": nil, "status": http.StatusInternalServerError, "message": err.Error()})
	}
	
	// Valida si se elimino algun documento
	if res.DeletedCount == 0 {
		return c.JSON(http.StatusNotFound, echo.Map{"data": nil, "status": http.StatusNotFound, "message": "Inventario no encontrado"})
	}

	return c.NoContent(http.StatusNoContent)
}