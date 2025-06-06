package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"inventory/handlers"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func setupEcho(t *testing.T) *echo.Echo {
	client, err := mongo.Connect(options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		t.Fatal(err)
	}
	handlers.Collection = client.Database("test_db").Collection("inventories")
	handlers.Collection.DeleteMany(context.Background(), bson.M{})
	e := echo.New()
	return e
}

// Prueba 1: Crear un inventario correctamente
func TestIntegration_CreateInventory_Success(t *testing.T) {
	e := setupEcho(t)

	payload := `{"name":"Inventario creado","object":[{"name":"Nombre del objeto","description":"Descripcion del objeto"}]}`
	req := httptest.NewRequest(http.MethodPost, "/inventories", bytes.NewBuffer([]byte(payload)))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := handlers.CreateInventory(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, rec.Code)

	var created handlers.Inventory
	_ = json.Unmarshal(rec.Body.Bytes(), &created)
	assert.Equal(t, "Inventario creado", created.Name)
	assert.Equal(t, "Nombre del objeto", created.Objects[0].Name)
	assert.Equal(t, "Descripcion del objeto", created.Objects[0].Description)
}

// Prueba 2: Obtener un inventario por su ID
func TestIntegration_GetInventoryById_Found(t *testing.T) {
	e := setupEcho(t)

	// Crea primero un inventario
	payload := `{"name":"Buscar este","object":[{"name":"Monitor","description":"24 pulgadas"}]}`
	req := httptest.NewRequest(http.MethodPost, "/inventories", bytes.NewBuffer([]byte(payload)))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	_ = handlers.CreateInventory(c)

	var inv handlers.Inventory
	_ = json.Unmarshal(rec.Body.Bytes(), &inv)
	id := inv.ID.Hex()

	// Luego buscarlo
	req2 := httptest.NewRequest(http.MethodGet, "/inventories/"+id, nil)
	rec2 := httptest.NewRecorder()
	c2 := e.NewContext(req2, rec2)
	c2.SetPath("/inventories/:id")
	c2.SetParamNames("id")
	c2.SetParamValues(id)

	err := handlers.GetInventoryById(c2)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusFound, rec2.Code)
	assert.Contains(t, rec2.Body.String(), "Buscar este")
}

// Prueba 3: Intentar eliminar un inventario inexistente
func TestIntegration_DeleteInventory_NotFound(t *testing.T) {
	e := setupEcho(t)

	// ID con formato válido pero no existente
	id := "60ddc9737c213e3d8c9e6a6f"
	req := httptest.NewRequest(http.MethodDelete, "/inventories/"+id, nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/inventories/:id")
	c.SetParamNames("id")
	c.SetParamValues(id)

	err := handlers.DeleteInventory(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, rec.Code)
	assert.Contains(t, rec.Body.String(), "Inventario no encontrado")
}
