// package tests

// import (
// 	"bytes"
// 	"context"
// 	"encoding/json"
// 	"inventory/handlers"
// 	"net/http"
// 	"os/exec"
// 	"testing"
// 	"time"

// 	"github.com/stretchr/testify/assert"
// 	"go.mongodb.org/mongo-driver/v2/bson"
// 	"go.mongodb.org/mongo-driver/v2/mongo"
// 	"go.mongodb.org/mongo-driver/v2/mongo/options"
// )

// const baseURL = "http://localhost:8080"

// var createdID string

// // Espera hasta que el servidor responda
// func waitServerReady() {
// 	for i := 0; i < 10; i++ {
// 		resp, err := http.Get(baseURL + "/inventories")
// 		if err == nil && resp.StatusCode < 500 {
// 			return
// 		}
// 		time.Sleep(500 * time.Millisecond)
// 	}
// 	panic("Servidor no arrancó a tiempo")
// }

// // Inicia el servidor en segundo plano
// func startServerForTest(t *testing.T) {
// 	cmd := exec.Command("go", "run", "main.go")
// 	if err := cmd.Start(); err != nil {
// 		t.Fatalf("no se pudo iniciar el servidor: %v", err)
// 	}
// 	time.Sleep(2 * time.Second)
// 	waitServerReady()
// }

// // Limpia la colección de prueba
// func clearMongoCollection(t *testing.T) {
// 	client, err := mongo.Connect(options.Client().ApplyURI("mongodb://localhost:27017"))
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	handlers.Collection = client.Database("test_db").Collection("inventories")
// 	_, _ = handlers.Collection.DeleteMany(context.Background(), bson.M{})
// }

// // Prueba 1: Crear inventario
// func TestE2E_CreateInventory(t *testing.T) {
// 	startServerForTest(t)
// 	clearMongoCollection(t)

// 	body := []byte(`{"name":"E2E Inventario","object":[{"name":"Proyector","description":"HD"}]}`)
// 	resp, err := http.Post(baseURL+"/inventories", "application/json", bytes.NewReader(body))

// 	assert.NoError(t, err)
// 	assert.Equal(t, http.StatusCreated, resp.StatusCode)

// 	var result map[string]interface{}
// 	_ = json.NewDecoder(resp.Body).Decode(&result)
// 	assert.NotNil(t, result["id"])

// 	createdID = result["id"].(string)
// }

// // Prueba 2: Obtener inventario por ID
// func TestE2E_GetInventoryById(t *testing.T) {
// 	waitServerReady()

// 	resp, err := http.Get(baseURL + "/inventories/" + createdID)
// 	assert.NoError(t, err)
// 	assert.Equal(t, http.StatusFound, resp.StatusCode)

// 	var result map[string]interface{}
// 	_ = json.NewDecoder(resp.Body).Decode(&result)
// 	assert.Equal(t, "E2E Inventario", result["name"])
// }

// // Prueba 3: Eliminar inventario
// func TestE2E_DeleteInventory(t *testing.T) {
// 	req, _ := http.NewRequest(http.MethodDelete, baseURL+"/inventories/"+createdID, nil)
// 	client := &http.Client{}
// 	resp, err := client.Do(req)

// 	assert.NoError(t, err)
// 	assert.Equal(t, http.StatusNoContent, resp.StatusCode)

// 	// Confirmar que ya no existe
// 	resp2, _ := http.Get(baseURL + "/inventories/" + createdID)
// 	assert.Equal(t, http.StatusNotFound, resp2.StatusCode)
// }
