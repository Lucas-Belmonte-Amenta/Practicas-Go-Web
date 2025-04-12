package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

// función para leer un archivo JSON y deserializarlo en un slice de Product
func ReadJSONFile[T any](fileName string, emptyListEntity *[]T) error {
	// Abrir el archivo
	file, err := os.Open(fileName)
	if err != nil {
		return fmt.Errorf("Error al abrir el archivo Json: %v\n", err)
	}
	defer file.Close()

	// Leer el contenido del archivo
	byteValue, err := io.ReadAll(file)
	if err != nil {
		return fmt.Errorf("Error al leer el archivo Json: %v\n", err)
	}

	// Deserializar el JSON en el slice de Product
	err = json.Unmarshal(byteValue, emptyListEntity)
	if err != nil {
		return fmt.Errorf("Error al deserializar el JSON: %v\n", err)
	}

	return nil

}
