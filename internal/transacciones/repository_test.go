package transacciones

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

type DummyStore struct{}

func (d *DummyStore) Read(data interface{}) error {
	trans := data.(*[]Transaccion)
	*trans = []Transaccion{}
	return nil
}

func (s *DummyStore) Write(data interface{}) error {
	return nil
}

type ErrorReadStore struct {
	readWasCalled  bool
	writeWasCalled bool
}

func (d *ErrorReadStore) Read(data interface{}) error {
	d.readWasCalled = true
	return errors.New("error al leer la data dentro del store")
}

func (d *ErrorReadStore) Write(data interface{}) error {
	d.writeWasCalled = true
	return nil
}

type ErrorWriteStore struct {
	readWasCalled  bool
	writeWasCalled bool
}

func (d *ErrorWriteStore) Read(data interface{}) error {
	d.readWasCalled = true
	return nil
}

func (d *ErrorWriteStore) Write(data interface{}) error {
	d.writeWasCalled = true
	return errors.New("error al escribir la data dentro del store")
}

type StubStore struct{}

func (s *StubStore) Read(data interface{}) error {
	trans := data.(*[]Transaccion)
	*trans = []Transaccion{
		{
			Id:                1,
			CodigoTransaccion: "ctr1",
			Moneda:            "MXN",
			Monto:             4000,
			Emisor:            "Brandon",
			Receptor:          "Juan",
			FechaTransaccion:  "21/04/2022",
		},
		{
			Id:                2,
			CodigoTransaccion: "ctr2",
			Moneda:            "USD",
			Monto:             200,
			Emisor:            "Juan",
			Receptor:          "Brandon",
			FechaTransaccion:  "21/04/2022",
		},
	}
	return nil
}

func (s *StubStore) Write(data interface{}) error {
	return nil
}

type SpyStore struct {
	readWasCalled  bool
	writeWasCalled bool
}

func (s *SpyStore) Read(data interface{}) error {
	trans := data.(*[]Transaccion)
	*trans = []Transaccion{}
	s.readWasCalled = true
	return nil
}

func (s *SpyStore) Write(data interface{}) error {
	s.writeWasCalled = true
	return nil
}

type MockStore struct {
	readWasCalled  bool
	writeWasCalled bool
	Data           []Transaccion
}

func (s *MockStore) Read(data interface{}) error {
	trans := data.(*[]Transaccion)
	*trans = s.Data
	s.readWasCalled = true
	return nil
}

func (s *MockStore) Write(data interface{}) error {
	s.Data = data.([]Transaccion)
	s.writeWasCalled = true
	return nil
}

func TestRepositoryGetAll(t *testing.T) {
	// Arrange
	expected := []Transaccion{
		{
			Id:                1,
			CodigoTransaccion: "ctr1",
			Moneda:            "MXN",
			Monto:             4000,
			Emisor:            "Brandon",
			Receptor:          "Juan",
			FechaTransaccion:  "21/04/2022",
		},
		{
			Id:                2,
			CodigoTransaccion: "ctr2",
			Moneda:            "USD",
			Monto:             200,
			Emisor:            "Juan",
			Receptor:          "Brandon",
			FechaTransaccion:  "21/04/2022",
		},
	}
	stubStore := &StubStore{}
	repo := NewRepository(stubStore)

	// Act
	result, err := repo.GetAll()

	// Assert
	assert.Nil(t, err)
	assert.Equal(t, expected, result)
}

func TestRepositoryGetAllEmpty(t *testing.T) {
	// Arrange
	dummyStore := &DummyStore{}
	repo := NewRepository(dummyStore)

	// Act
	result, err := repo.GetAll()

	// Assert
	assert.NotNil(t, err)
	assert.Empty(t, result)
}

func TestRepositoryStore(t *testing.T) {
	// Arrange
	mockStore := &MockStore{
		Data: []Transaccion{{
			Id:                100,
			CodigoTransaccion: "Before Update",
			Moneda:            "MXN",
			Monto:             0,
			Emisor:            "Brandon",
			Receptor:          "Juan",
			FechaTransaccion:  "21/04/2022",
		}},
	}
	repo := NewRepository(mockStore)
	expected := Transaccion{
		Id:                101,
		CodigoTransaccion: "ctr",
		Moneda:            "MXN",
		Monto:             100,
		Emisor:            "Banamex",
		Receptor:          "Bancomer",
		FechaTransaccion:  "21/02/2022",
	}

	// Act
	result, err := repo.Store(expected.Id, expected.CodigoTransaccion, expected.Moneda,
		expected.Monto, expected.Emisor, expected.Receptor, expected.FechaTransaccion)

	// Assert
	assert.Nil(t, err)
	assert.True(t, mockStore.readWasCalled)
	assert.True(t, mockStore.writeWasCalled)
	assert.Equal(t, expected, result)
}

func TestRepositoryUpdate(t *testing.T) {
	// Arrange
	mockStore := &MockStore{
		Data: []Transaccion{
			{
				Id:                1,
				CodigoTransaccion: "Before Update",
				Moneda:            "MXN",
				Monto:             100,
				Emisor:            "Banamex",
				Receptor:          "Bancomer",
				FechaTransaccion:  "22/04/2022",
			},
		},
	}
	repo := NewRepository(mockStore)
	expected := Transaccion{
		Id:                1,
		CodigoTransaccion: "After Update",
		Moneda:            "USD",
		Monto:             200,
		Emisor:            "Banregio",
		Receptor:          "Visa",
		FechaTransaccion:  "22/02/2022",
	}

	// Act
	result, err := repo.Update(expected.Id, expected.CodigoTransaccion, expected.Moneda,
		expected.Monto, expected.Emisor, expected.Receptor, expected.FechaTransaccion)

	// Assert
	assert.Nil(t, err)
	assert.True(t, mockStore.readWasCalled)
	assert.True(t, mockStore.writeWasCalled)
	assert.Equal(t, expected, result)
}

func TestRepositoryUpdateNotFound(t *testing.T) {
	// Arrange
	spyStore := &SpyStore{}
	repo := NewRepository(spyStore)
	data := Transaccion{
		Id:                1,
		CodigoTransaccion: "After Update",
		Moneda:            "USD",
		Monto:             200,
		Emisor:            "Banregio",
		Receptor:          "Visa",
		FechaTransaccion:  "22/02/2022",
	}

	// Act
	result, err := repo.Update(data.Id, data.CodigoTransaccion, data.Moneda,
		data.Monto, data.Emisor, data.Receptor, data.FechaTransaccion)

	// Assert
	assert.NotNil(t, err)
	assert.True(t, spyStore.readWasCalled)
	assert.False(t, spyStore.writeWasCalled)
	assert.Empty(t, result)
}

func TestRepositoryPatch(t *testing.T) {
	// Arrange
	mockStore := &MockStore{
		Data: []Transaccion{
			{
				Id:                1,
				CodigoTransaccion: "Before Update",
				Moneda:            "MXN",
				Monto:             0,
				Emisor:            "Brandon",
				Receptor:          "Juan",
				FechaTransaccion:  "21/04/2022",
			},
		}}
	repo := NewRepository(mockStore)
	id := 1
	newCodigoTransaction := "After Update"
	newMonto := 200.0

	// Act
	result, err := repo.Patch(id, newCodigoTransaction, newMonto)

	// Assert
	assert.Nil(t, err)
	assert.True(t, mockStore.readWasCalled)
	assert.Equal(t, newCodigoTransaction, result.CodigoTransaccion)
	assert.Equal(t, newMonto, result.Monto)
}

func TestRepositoryPatchNotFound(t *testing.T) {
	// Arrange
	spyStore := &SpyStore{}
	repo := NewRepository(spyStore)
	id := 1
	newCodigoTransaction := "After Update"
	newMonto := 200.0

	// Act
	result, err := repo.Patch(id, newCodigoTransaction, newMonto)

	// Assert
	assert.NotNil(t, err)
	assert.True(t, spyStore.readWasCalled)
	assert.False(t, spyStore.writeWasCalled)
	assert.Empty(t, result)
}

func TestRepositoryDelete(t *testing.T) {
	// Arrange
	mockStore := &MockStore{
		Data: []Transaccion{
			{
				Id:                1,
				CodigoTransaccion: "ctr",
				Moneda:            "MXN",
				Monto:             0,
				Emisor:            "Brandon",
				Receptor:          "Juan",
				FechaTransaccion:  "21/04/2022",
			},
		}}
	repo := NewRepository(mockStore)
	id := 1

	// Act
	err := repo.Delete(id)

	// Assert
	assert.Nil(t, err)
	assert.True(t, mockStore.readWasCalled)
	assert.True(t, mockStore.writeWasCalled)
}

func TestRepositoryLastID(t *testing.T) {
	// Arrange
	mockStore := &MockStore{
		Data: []Transaccion{
			{
				Id:                10,
				CodigoTransaccion: "ctr",
				Moneda:            "MXN",
				Monto:             0,
				Emisor:            "Brandon",
				Receptor:          "Juan",
				FechaTransaccion:  "21/04/2022",
			},
		}}
	repo := NewRepository(mockStore)
	expected := 10

	// Act
	result, err := repo.LastID()

	// Assert
	assert.Nil(t, err)
	assert.True(t, mockStore.readWasCalled)
	assert.False(t, mockStore.writeWasCalled)
	assert.Equal(t, expected, result)
}
