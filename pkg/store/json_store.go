package store

import (
	"encoding/json"
	"errors"
	"os"
)

type JsonFileStore struct {
	FileName string
}

func (s *JsonFileStore) Read(data interface{}) error {
	jsonData, err := os.ReadFile(s.FileName)
	if err != nil {
		return errors.New("archivo no encontrado")
	}
	serr := json.Unmarshal((jsonData), data)
	if serr != nil {
		return errors.New("archivo con formato no valido")
	}
	return nil
}

func (s *JsonFileStore) Write(data interface{}) error {
	content, err := json.Marshal(data)
	if err != nil {
		return errors.New("error al almacenar la nueva transaccion")
	}
	if err := os.WriteFile(s.FileName, content, 0644); err != nil {
		return errors.New("error al escribir en el archivo json")
	}
	return nil
}
