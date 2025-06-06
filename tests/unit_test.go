package tests

import (
	"inventory/handlers"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

// Utilidad para crear un contexto de prueba con cuerpo JSON simulado
func setupEchoTestContext(method, path, body string) (echo.Context, *httptest.ResponseRecorder) {
	e := echo.New()
	req := httptest.NewRequest(method, path, strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	return c, rec
}

// Prueba 1: Nombre de inventario vacío
func TestCreateInventory_MissingName(t *testing.T) {
	// Construye un json de un nuevo inventario pero sin nombre
	json := `{"name": "", "object": [{"name": "Objeto de un inventario sin nombre", "description": "No hay nombre para el inventario"}]}`
	c, rec := setupEchoTestContext(http.MethodPost, "/inventories", json)

	// Intenta guardar el inventario
	err := handlers.CreateInventory(c)

	// No debe haber error en la ejecucion del handler
	assert.NoError(t, err)
	// Compara que el codigo de respuesta http coincida con la esperada
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	// Compara que el mensaje de error coincida con el esperado
	assert.Contains(t, rec.Body.String(), "El nombre es obligatorio")
}

// Prueba 2: Lista de objetos vacía
func TestCreateInventory_MissingObjects(t *testing.T) {
	// Construye un json de un nuevo inventario pero sin con la lista vacia de object
	json := `{"name": "Inventario sin objetos", "object": []}`
	c, rec := setupEchoTestContext(http.MethodPost, "/inventories", json)

	// Intenta guardar el inventario
	err := handlers.CreateInventory(c)

	// No debe haber error en la ejecucion del handler
	assert.NoError(t, err)
	// Compara que el codigo de respuesta http coincida con la esperada
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	// Compara que el mensaje de error coincida con el esperado
	assert.Contains(t, rec.Body.String(), "Los objetos son obligatorios")
}

// Prueba 3: Objeto sin descripción
func TestCreateInventory_ObjectMissingDescription(t *testing.T) {
	// Construye un json de un nuevo inventario pero sin la descripcion del objeto
	json := `{"name": "Inventario con objeto sin descripcion", "object": [{"name": "Objeto sin descripcion", "description": "g"}]}`
	c, rec := setupEchoTestContext(http.MethodPost, "/inventories", json)
	
	// Intenta guardar el inventario
	err := handlers.CreateInventory(c)

	// No debe haber error en la ejecucion del handler
	assert.NoError(t, err)
	// Compara que el codigo de respuesta http coincida con la esperada
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	// Compara que el mensaje de error coincida con el esperado
	assert.Contains(t, rec.Body.String(), "carece de descripcion")
}