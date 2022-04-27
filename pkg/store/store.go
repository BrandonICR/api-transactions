package store

type Store interface {
	Read(data interface{}) error
	Write(data interface{}) error
}

type StoreType string

const (
	JsonFileType StoreType = "jsonFile"
)

func NewStore(storeType StoreType, filename string) Store {
	switch storeType {
	case JsonFileType:
		return &JsonFileStore{FileName: filename}
	}
	return nil
}
