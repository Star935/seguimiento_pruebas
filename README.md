# API de Inventarios - Seguimiento Corte 2

Esta API permite la gestión de inventarios y objetos anidados usando Go, Echo y MongoDB.

## Ejecución del Proyecto

### Requisitos

- Go 1.23+
- MongoDB corriendo en `mongodb://localhost:27017`

### Ejecutar el servidor

```bash
go run main.go

## Endpoints

### Obtener todos los inventarios
**GET** `/inventories`  
- **Descripción:** Lista todos los inventarios registrados.
- **Respuesta:** `200 OK` con arreglo de inventarios.

---

### Obtener inventario por ID
**GET** `/inventories/:id`  
- **Descripción:** Devuelve un inventario por su ID.
- **Parámetros de ruta:**  
  - `id` (string): ID del inventario (formato ObjectID).
- **Respuestas:**  
  - `200 OK` con el inventario.  
  - `400 Bad Request` si el ID es inválido.  
  - `404 Not Found` si no se encuentra.

---

### Crear inventario
**POST** `/inventories`  
- **Descripción:** Crea un nuevo inventario con al menos un objeto.
- **Body JSON:**
```json
{
  "name": "Inventario de laboratorio de ingenieria",
  "objects": [
    {
      "name": "Proyector",
      "description": "Proyector Epson HD"
    },
    {
      "name": "Computadora",
      "description": "PC de escritorio con 16GB RAM"
    }
  ]
}
