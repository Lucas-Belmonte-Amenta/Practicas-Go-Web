package storage

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

type storageJSON struct {
	fileName string
}

func NewStorageJSON(fileName string) (Storage, error) {

	file, err := os.Open(fileName)
	if err != nil {
		return nil, fmt.Errorf("Error al abrir el archivo Json: %v\n", err)
	}
	file.Close()

	storageJSON := &storageJSON{fileName: fileName}

	return storageJSON, nil
}

// función para leer un archivo JSON y deserializarlo en un slice de Product
func (sj *storageJSON) Read(emptyListEntity any) error {
	// Abrir el archivo
	file, err := os.Open(sj.fileName)
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

// función para escribir un slice de Product en un archivo JSON
func (sj *storageJSON) Write(emptyListEntity any) error {

	// Crear el archivo
	file, err := os.Create(sj.fileName)
	if err != nil {
		return fmt.Errorf("Error al crear el archivo Json: %s\n", err)
	}
	defer file.Close()

	// Serializar el slice de Product a JSON y escribirlo en el archivo
	json.NewEncoder(file).Encode(emptyListEntity)

	return nil

}
