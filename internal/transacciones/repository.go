package transacciones

import (
	"errors"

	"github.com/BrandonICR/web_cl2_050422_8am/pkg/store"
)

type Transaccion struct {
	Id                int     `json:"id"`
	CodigoTransaccion string  `json:"codigo_transaccion"`
	Moneda            string  `json:"moneda"`
	Monto             float64 `json:"monto"`
	Emisor            string  `json:"emisor"`
	Receptor          string  `json:"receptor"`
	FechaTransaccion  string  `json:"fecha_transaccion"`
}

var transaccionesList []Transaccion

type Repository interface {
	GetAll() ([]Transaccion, error)
	Store(id int, codigoTransaccion, moneda string, monto float64, emisor, receptor, fechaTransaccion string) (Transaccion, error)
	Update(id int, codigoTransaccion, moneda string, monto float64, emisor, receptor, fechaTransaccion string) (Transaccion, error)
	Patch(id int, codigoTransaccion string, monto float64) (Transaccion, error)
	Delete(id int) error
	LastID() (int, error)
}

type repository struct {
	db store.Store
}

func NewRepository(db store.Store) Repository {
	return &repository{db: db}
}

func (r *repository) GetAll() ([]Transaccion, error) {
	if err := r.db.Read(&transaccionesList); err != nil {
		return []Transaccion{}, errors.New("error al leer del store")
	}

	if len(transaccionesList) == INT_ZERO {
		return []Transaccion{}, errors.New("ninguna transaccion fue encontrada")
	}

	return transaccionesList, nil
}

func (r *repository) Store(id int, codigoTransaccion, moneda string, monto float64, emisor, receptor, fechaTransaccion string) (Transaccion, error) {
	if err := r.db.Read(&transaccionesList); err != nil {
		return Transaccion{}, errors.New("error al leer del store")
	}

	transaccion := Transaccion{
		Id:                id,
		CodigoTransaccion: codigoTransaccion,
		Moneda:            moneda,
		Monto:             monto,
		Emisor:            emisor,
		Receptor:          receptor,
		FechaTransaccion:  fechaTransaccion,
	}

	transaccionesList = append(transaccionesList, transaccion)

	if err := r.db.Write(transaccionesList); err != nil {
		return Transaccion{}, err
	}

	return transaccion, nil
}

func (r *repository) Update(id int, codigoTransaccion, moneda string, monto float64, emisor, receptor, fechaTransaccion string) (Transaccion, error) {
	if err := r.db.Read(&transaccionesList); err != nil {
		return Transaccion{}, errors.New("error al leer del store")
	}
	transaccionUpdated := Transaccion{
		Id:                id,
		CodigoTransaccion: codigoTransaccion,
		Moneda:            moneda,
		Monto:             monto,
		Emisor:            emisor,
		Receptor:          receptor,
		FechaTransaccion:  fechaTransaccion,
	}

	var wasUpdated bool //Elegi con boolean en lugar de directo si no incrementaría la complejidad ciclomática por el writeRepository

	for index, transaccion := range transaccionesList {
		if transaccion.Id == transaccionUpdated.Id {
			transaccionesList[index] = transaccionUpdated
			wasUpdated = true
		}
	}

	if !wasUpdated {
		return Transaccion{}, errors.New("no se encontro la transaccion a actualizar")
	}

	if err := r.db.Write(transaccionesList); err != nil {
		return Transaccion{}, err
	}

	return transaccionUpdated, nil
}

func (r *repository) Patch(id int, codigoTransaccion string, monto float64) (Transaccion, error) {
	if err := r.db.Read(&transaccionesList); err != nil {
		return Transaccion{}, errors.New("error al leer del store")
	}
	var wasUpdated bool
	var transaccionUpdated Transaccion

	for index, transaccion := range transaccionesList {
		if transaccion.Id == id {
			transaccion.CodigoTransaccion = codigoTransaccion
			transaccion.Monto = monto
			transaccionUpdated = transaccion
			transaccionesList[index] = transaccion
			wasUpdated = true
		}
	}

	if !wasUpdated {
		return Transaccion{}, errors.New("no se encontro la transaccion a actualizar")
	}

	if err := r.db.Write(transaccionesList); err != nil {
		return Transaccion{}, err
	}

	return transaccionUpdated, nil
}

func (r *repository) LastID() (int, error) {
	if err := r.db.Read(&transaccionesList); err != nil {
		return 0, errors.New("error al leer del store")
	}
	var maxId int
	for _, transaccion := range transaccionesList {
		if maxId < transaccion.Id {
			maxId = transaccion.Id
		}
	}
	return maxId, nil
}

func (r *repository) Delete(id int) error {
	if err := r.db.Read(&transaccionesList); err != nil {
		return errors.New("error al leer del store")
	}
	var idIndex int

	for index, transaccion := range transaccionesList {
		if transaccion.Id == id {
			idIndex = index + 1
		}
	}

	if idIndex == 0 {
		return errors.New("la transaccion a eliminar no existe")
	}

	if idIndex > len(transaccionesList)-1 {
		transaccionesList = transaccionesList[:idIndex-1]
	} else {
		transaccionesList = append(transaccionesList[:idIndex-1], transaccionesList[idIndex:]...)
	}

	if err := r.db.Write(transaccionesList); err != nil {
		return err
	}

	return nil
}
