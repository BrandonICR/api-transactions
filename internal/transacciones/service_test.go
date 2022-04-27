package transacciones

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestServiceGetAll(t *testing.T) {
	// Arrange
	expected := []Transaccion{{
		Id:                1,
		CodigoTransaccion: "ctr1",
		Moneda:            "MXN",
		Monto:             0,
		Emisor:            "Bancomer",
		Receptor:          "Banamex",
		FechaTransaccion:  "21/04/2022",
	}, {
		Id:                2,
		CodigoTransaccion: "ctr2",
		Moneda:            "USD",
		Monto:             100,
		Emisor:            "Banamex",
		Receptor:          "Bancomer",
		FechaTransaccion:  "22/04/2022",
	}}
	mock := MockStore{
		Data: expected,
	}
	repo := NewRepository(&mock)
	service := NewService(repo)

	// Act
	result, err := service.GetAll()

	// Assert
	assert.Nil(t, err)
	assert.True(t, mock.readWasCalled)
	assert.False(t, mock.writeWasCalled)
	assert.Equal(t, expected, result)
}

func TestServiceGetTransaccionFiltrada(t *testing.T) {
	// Arrange
	lenExpected := 1
	resultExpected := []Transaccion{{
		Id:                2,
		CodigoTransaccion: "ctr2",
		Moneda:            "USD",
		Monto:             100,
		Emisor:            "Banamex",
		Receptor:          "Bancomer",
		FechaTransaccion:  "22/04/2022",
	}}
	filter := Transaccion{
		CodigoTransaccion: "ctr2",
	}

	mock := MockStore{
		Data: []Transaccion{{
			Id:                1,
			CodigoTransaccion: "ctr1",
			Moneda:            "MXN",
			Monto:             0,
			Emisor:            "Bancomer",
			Receptor:          "Banamex",
			FechaTransaccion:  "21/04/2022",
		}, resultExpected[0]},
	}
	repo := NewRepository(&mock)
	service := NewService(repo)

	// Act
	result, err := service.GetTransaccionFiltrada(filter.Id, filter.CodigoTransaccion, filter.Moneda,
		filter.Monto, filter.Emisor, filter.Receptor, filter.FechaTransaccion)

	// Assert
	assert.Nil(t, err)
	assert.True(t, mock.readWasCalled)
	assert.False(t, mock.writeWasCalled)
	assert.Len(t, result, lenExpected)
	assert.Equal(t, resultExpected, result)
}

func TestServiceGetTransaccion(t *testing.T) {
	// Arrange
	mock := MockStore{
		Data: []Transaccion{{
			Id:                1,
			CodigoTransaccion: "ctr1",
			Moneda:            "MXN",
			Monto:             0,
			Emisor:            "Brandon",
			Receptor:          "Juan",
			FechaTransaccion:  "21/04/2022",
		}},
	}
	repo := NewRepository(&mock)
	service := NewService(repo)
	id := 1

	// Act
	result, err := service.GetTransaccion(id)

	// Assert
	assert.Nil(t, err)
	assert.True(t, mock.readWasCalled)
	assert.False(t, mock.writeWasCalled)
	assert.Equal(t, id, result.Id)
}

func TestServiceStore(t *testing.T) {
	// Arrange
	mock := MockStore{
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
	repo := NewRepository(&mock)
	service := NewService(repo)
	expected := Transaccion{
		Id:                101,
		CodigoTransaccion: "After Update",
		Moneda:            "USD",
		Monto:             100,
		Emisor:            "Juan",
		Receptor:          "Pedro",
		FechaTransaccion:  "22/04/2022",
	}

	// Act
	result, err := service.Store(expected.CodigoTransaccion, expected.Moneda,
		expected.Monto, expected.Emisor, expected.Receptor, expected.FechaTransaccion)

	// Assert
	assert.Nil(t, err)
	assert.True(t, mock.readWasCalled)
	assert.True(t, mock.writeWasCalled)
	assert.Equal(t, expected, result)
}

func TestServiceUpdate(t *testing.T) {
	// Arrange
	mock := MockStore{
		Data: []Transaccion{{
			Id:                1,
			CodigoTransaccion: "Before Update",
			Moneda:            "MXN",
			Monto:             0,
			Emisor:            "Brandon",
			Receptor:          "Juan",
			FechaTransaccion:  "21/04/2022",
		}},
	}
	repo := NewRepository(&mock)
	service := NewService(repo)
	expected := Transaccion{
		Id:                1,
		CodigoTransaccion: "After Update",
		Moneda:            "USD",
		Monto:             100,
		Emisor:            "Juan",
		Receptor:          "Pedro",
		FechaTransaccion:  "22/04/2022",
	}

	// Act
	result, err := service.Update(expected.Id, expected.CodigoTransaccion, expected.Moneda,
		expected.Monto, expected.Emisor, expected.Receptor, expected.FechaTransaccion)

	// Assert
	assert.Nil(t, err)
	assert.True(t, mock.readWasCalled)
	assert.True(t, mock.writeWasCalled)
	assert.Equal(t, expected, result)
}

func TestServicePatch(t *testing.T) {
	// Arrange
	mock := MockStore{
		Data: []Transaccion{{
			Id:                1,
			CodigoTransaccion: "Before Update",
			Moneda:            "MXN",
			Monto:             0,
			Emisor:            "Brandon",
			Receptor:          "Juan",
			FechaTransaccion:  "21/04/2022",
		}},
	}
	repo := NewRepository(&mock)
	service := NewService(repo)

	id := 1
	codigoTransaccion := "After Update"
	monto := 200.0

	// Act
	result, err := service.Patch(id, codigoTransaccion, monto)

	// Assert
	assert.Nil(t, err)
	assert.True(t, mock.readWasCalled)
	assert.True(t, mock.writeWasCalled)
	assert.Equal(t, id, result.Id)
	assert.Equal(t, codigoTransaccion, result.CodigoTransaccion)
	assert.Equal(t, monto, result.Monto)
}

func TestServiceDelete(t *testing.T) {
	// Arrange
	mock := MockStore{
		Data: []Transaccion{{
			Id:                1,
			CodigoTransaccion: "ctr1",
			Moneda:            "MXN",
			Monto:             100,
			Emisor:            "Banxico",
			Receptor:          "Banamex",
			FechaTransaccion:  "21/04/2022",
		}, {
			Id:                2,
			CodigoTransaccion: "ctr2",
			Moneda:            "MXN",
			Monto:             200,
			Emisor:            "Bancomer",
			Receptor:          "Banxico",
			FechaTransaccion:  "22/04/2022",
		}, {
			Id:                3,
			CodigoTransaccion: "ctr3",
			Moneda:            "MXN",
			Monto:             300,
			Emisor:            "Banamex",
			Receptor:          "Bancomer",
			FechaTransaccion:  "23/04/2022",
		}},
	}
	repo := NewRepository(&mock)
	service := NewService(repo)
	id := 1
	lenExpected := len(mock.Data) - 1

	// Act
	err := service.Delete(id)
	err2 := service.Delete(id)

	// Assert
	assert.Nil(t, err)
	assert.True(t, mock.readWasCalled)
	assert.True(t, mock.writeWasCalled)
	assert.NotNil(t, err2)
	assert.Len(t, mock.Data, lenExpected)
}
