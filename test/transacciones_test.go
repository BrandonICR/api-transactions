package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/BrandonICR/web_cl2_050422_8am/cmd/server/engine"
	"github.com/stretchr/testify/assert"
)

const FILE_STORE = "transacciones.json"

type transaccion struct {
	Id                int     `json:"id"`
	CodigoTransaccion string  `json:"codigo_transaccion"`
	Moneda            string  `json:"moneda"`
	Monto             float64 `json:"monto"`
	Emisor            string  `json:"emisor"`
	Receptor          string  `json:"receptor"`
	FechaTransaccion  string  `json:"fecha_transaccion"`
}

func TestUpdate(t *testing.T) {
	tempFileName := "transacciones_update_temp.json"
	router := engine.GetEngine(FILE_STORE, tempFileName, "./../.env")

	type response struct {
		Code    string      `json:"code"`
		Message string      `json:"message"`
		Data    transaccion `json:"data,omitempty"`
		Error   string      `json:"error,omitempty"`
	}
	var resBody response

	reqBody := transaccion{
		Id:                2,
		CodigoTransaccion: "ctr new",
		Moneda:            "USD",
		Monto:             900,
		Emisor:            "Banamex",
		Receptor:          "Banxico",
		FechaTransaccion:  "23/04/2022",
	}

	reqBytesBody, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPut,
		fmt.Sprintf("/api/v1/transacciones/%d", reqBody.Id),
		bytes.NewBuffer(reqBytesBody))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("authorization", "12345")

	res := httptest.NewRecorder()

	router.ServeHTTP(res, req)

	assert.Equal(t, http.StatusOK, res.Code)
	err := json.Unmarshal(res.Body.Bytes(), &resBody)
	assert.Nil(t, err)
	assert.Equal(t, reqBody, resBody.Data)

	os.Remove(tempFileName)
}

func TestDelete(t *testing.T) {
	tempFileName := "transacciones_delete_temp.json"
	router := engine.GetEngine(FILE_STORE, tempFileName, "./../.env")

	type response struct {
		Code    string      `json:"code"`
		Message string      `json:"message"`
		Data    transaccion `json:"data,omitempty"`
		Error   string      `json:"error,omitempty"`
	}
	var resBody response

	id := 2

	req := httptest.NewRequest(http.MethodDelete,
		fmt.Sprintf("/api/v1/transacciones/%d", id),
		nil)
	req.Header.Add("authorization", "12345")

	res := httptest.NewRecorder()

	router.ServeHTTP(res, req)

	assert.Equal(t, http.StatusOK, res.Code)
	err := json.Unmarshal(res.Body.Bytes(), &resBody)
	assert.Nil(t, err)

	os.Remove(tempFileName)
}
