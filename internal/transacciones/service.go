package transacciones

import (
	"errors"
)

const (
	STRING_EMPTY = ""
	INT_ZERO     = 0
)

type Service interface {
	GetAll() ([]Transaccion, error)
	GetTransaccionFiltrada(id int, codigoTransaccion, moneda string, monto float64, emisor, receptor, fechaTransaccion string) ([]Transaccion, error)
	GetTransaccion(id int) (Transaccion, error)
	Store(codigoTransaccion, moneda string, monto float64, emisor, receptor, fechaTransaccion string) (Transaccion, error)
	Update(id int, codigoTransaccion, moneda string, monto float64, emisor, receptor, fechaTransaccion string) (Transaccion, error)
	Patch(id int, codigoTransaccion string, monto float64) (Transaccion, error)
	Delete(id int) error
}

type service struct {
	repository Repository
}

func NewService(r Repository) Service {
	return &service{repository: r}
}

func (s *service) GetAll() ([]Transaccion, error) {
	return s.repository.GetAll()
}

func (s *service) GetTransaccionFiltrada(id int, codigoTransaccion, moneda string, monto float64, emisor, receptor, fechaTransaccion string) ([]Transaccion, error) {
	transacciones, err := s.repository.GetAll()

	if err != nil {
		return []Transaccion{}, err
	}

	var transaccionesFiltradas []Transaccion

	for _, transaccion := range transacciones {
		if (id == INT_ZERO || transaccion.Id == id) &&
			(codigoTransaccion == STRING_EMPTY || transaccion.CodigoTransaccion == codigoTransaccion) &&
			(moneda == STRING_EMPTY || transaccion.Moneda == moneda) &&
			(monto == INT_ZERO || transaccion.Monto == monto) &&
			(emisor == STRING_EMPTY || transaccion.Emisor == emisor) &&
			(receptor == STRING_EMPTY || transaccion.Receptor == receptor) &&
			(fechaTransaccion == STRING_EMPTY || transaccion.FechaTransaccion == fechaTransaccion) {
			transaccionesFiltradas = append(transaccionesFiltradas, transaccion)
		}
	}

	if len(transaccionesFiltradas) == INT_ZERO {
		return []Transaccion{}, errors.New("ninguna transaccion fue encontrada")
	}

	return transaccionesFiltradas, nil
}

func (s *service) GetTransaccion(id int) (Transaccion, error) {
	transacciones, err := s.repository.GetAll()

	if err != nil {
		return Transaccion{}, err
	}

	for _, transaccion := range transacciones {
		if transaccion.Id == id {
			return transaccion, nil
		}
	}

	return Transaccion{}, errors.New("no se enconto la transaccion")
}

func (s *service) Store(codigoTransaccion, moneda string, monto float64, emisor, receptor, fechaTransaccion string) (Transaccion, error) {
	id, err := s.repository.LastID()
	if err != nil {
		return Transaccion{}, err
	}
	id++
	transaccion, err := s.repository.Store(id, codigoTransaccion, moneda, monto, emisor, receptor, fechaTransaccion)
	if err != nil {
		return Transaccion{}, err
	}
	return transaccion, nil
}

func (s *service) Update(id int, codigoTransaccion, moneda string, monto float64, emisor, receptor, fechaTransaccion string) (Transaccion, error) {
	return s.repository.Update(id, codigoTransaccion, moneda, monto, emisor, receptor, fechaTransaccion)
}

func (s *service) Patch(id int, codigoTransaccion string, monto float64) (Transaccion, error) {
	return s.repository.Patch(id, codigoTransaccion, monto)
}

func (s *service) Delete(id int) error {
	return s.repository.Delete(id)
}
