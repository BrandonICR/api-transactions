package handler

import (
	"fmt"
	"net/http"
	"os"
	"reflect"
	"strconv"

	"github.com/BrandonICR/web_cl2_050422_8am/internal/transacciones"
	"github.com/BrandonICR/web_cl2_050422_8am/pkg/web"
	"github.com/gin-gonic/gin"
)

type request struct {
	Id                int     `json:"id"`
	CodigoTransaccion string  `json:"codigo_transaccion" validation:"required"`
	Moneda            string  `json:"moneda" validation:"required"`
	Monto             float64 `json:"monto" validation:"required"`
	Emisor            string  `json:"emisor" validation:"required"`
	Receptor          string  `json:"receptor" validation:"required"`
	FechaTransaccion  string  `json:"fecha_transaccion" validation:"required"`
}

type patchRequest struct {
	CodigoTransaccion string  `json:"codigo_transaccion" validation:"required"`
	Monto             float64 `json:"monto" validation:"required"`
}

type Transaccion struct {
	service transacciones.Service
}

func NewTransaccion(s transacciones.Service) *Transaccion {
	return &Transaccion{service: s}
}

func ValidarTransaccion(request request) error {
	values := reflect.ValueOf(request)
	keys := reflect.TypeOf(request)
	var badParameters string
	for i := 0; i < values.NumField(); i++ {
		validation, errValidation := keys.Field(i).Tag.Lookup("validation")
		if !errValidation {
			continue
		}
		if values.Field(i).IsZero() && validation == "required" {
			tag, errTag := keys.Field(i).Tag.Lookup("json")
			if !errTag {
				badParameters += keys.Field(i).Name + ", "
				continue
			}
			badParameters += tag + ", "
		}
	}
	if badParameters == "" {
		return nil
	}
	return fmt.Errorf("el campo %s es requerido", badParameters[:len(badParameters)-2])
}

func ValidarToken() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if ctx.GetHeader("authorization") != os.Getenv("TOKEN") {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, web.NewResponse(http.StatusUnauthorized, "No tiene permisos", nil, "No tiene permisos para realizar la peticion solicitada"))
			return
		}
		ctx.Next()
	}
}

// Get all transactions
// @Summary Get all transactions
// @Tags Transaction
// @Description Get  alltransactions
// @Accept json
// @Produce json
// @Param authorization header string true "authorization"
// @Succes 200 {object} web.Response
// @Router /transacciones [GET]
func (t *Transaccion) GetAll() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		transacciones, err := t.service.GetAll()

		if err != nil {
			ctx.JSON(http.StatusInternalServerError, web.NewResponse(http.StatusInternalServerError, "Error al recuperar las transacciones", nil, err.Error()))
			return
		}

		ctx.JSON(http.StatusOK, web.NewResponse(http.StatusOK, "Transacciones recuperadas con exito", transacciones, ""))
	}
}

// Get transaction using a filter
// @Summary Get transaction using a filter
// @Tags Transaction
// @Description Get transaction using a filter
// @Accept json
// @Produce json
// @Param authorization header string true "authorization"
// @Param id query int false "id"
// @Param codigo_transaccion query string false "codigo_transaccion"
// @Param moneda query string false "moneda"
// @Param monto query float64 false "monto"
// @Param emisor query string false "emisor"
// @Param receptor query string false "receptor"
// @Param fecha_transaccion query string false "fecha_transaccion"
// @Succes 200 {object} web.Response
// @Router /transacciones/ [GET]
func (t *Transaccion) GetTransaccionFiltrada() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id, _ := strconv.Atoi(ctx.Query("id"))
		codigoTransaccion := ctx.Query("codigo_transaccion")
		moneda := ctx.Query("moneda")
		monto, _ := strconv.ParseFloat(ctx.Query("monto"), 64)
		emisor := ctx.Query("emisor")
		receptor := ctx.Query("receptor")
		fechaTransaccion := ctx.Query("fecha_transaccion")

		transacciones, err := t.service.GetTransaccionFiltrada(id, codigoTransaccion, moneda, monto, emisor,
			receptor, fechaTransaccion)

		if err != nil {
			ctx.JSON(http.StatusNotFound, web.NewResponse(http.StatusNotFound, "Error al tratar de recuperar las transacciones", nil, err.Error()))
			return
		}

		ctx.JSON(http.StatusOK, web.NewResponse(http.StatusOK, "Transacciones recuperadas con exito", transacciones, ""))
	}
}

// Get a specific transaction
// @Summary Get transaction
// @Tags Transaction
// @Description Get a specific transaction using the id
// @Accept json
// @Produce json
// @Param authorization header string true "authorization"
// @Param Id path int true "Id"
// @Succes 200 {object} web.Response
// @Router /transacciones/{Id} [GET]
func (t *Transaccion) GetTransaccion() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		idParam, err := strconv.Atoi(ctx.Param("Id"))
		if err != nil {
			ctx.JSON(http.StatusBadRequest, web.NewResponse(http.StatusOK, "No se selecciono la transaccion a recuperar", nil, err.Error()))
			return
		}

		transaccion, err := t.service.GetTransaccion(idParam)

		if err != nil {
			ctx.JSON(http.StatusNotFound, web.NewResponse(http.StatusNotFound, "Error al tratar de recuperar la transaccion", nil, err.Error()))
			return
		}

		ctx.JSON(http.StatusOK, web.NewResponse(http.StatusOK, "Transaccion recuperada con exito", transaccion, ""))
	}
}

// Store a specific transaction
// @Summary Store transaction
// @Tags Transaction
// @Description Store a specific transaction using the body
// @Accept json
// @Produce json
// @Param authorization header string true "authorization"
// @Param transaction body request true "transaction"
// @Succes 200 {object} web.Response
// @Router /transacciones [POST]
func (t *Transaccion) Store() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var request request

		if err := ctx.ShouldBindJSON(&request); err != nil {
			ctx.JSON(http.StatusBadRequest, web.NewResponse(http.StatusBadRequest, "Peticion no valida", nil, err.Error()))
			return
		}

		if err := ValidarTransaccion(request); err != nil {
			ctx.JSON(http.StatusBadRequest, web.NewResponse(http.StatusBadRequest, "Peticio no valida", nil, err.Error()))
			return
		}

		transaccion, err := t.service.Store(request.CodigoTransaccion, request.Moneda,
			request.Monto, request.Emisor, request.Receptor, request.FechaTransaccion)

		if err != nil {
			ctx.JSON(http.StatusInternalServerError, web.NewResponse(http.StatusInternalServerError, "Error al tratar de almacenar la transaccion", nil, err.Error()))
			return
		}

		ctx.JSON(http.StatusOK, web.NewResponse(http.StatusOK, "Transaccion almacenada con exito", transaccion, ""))
	}
}

// Update a specific transaction
// @Summary Update transaction
// @Tags Transaction
// @Description Update a specific transaction using the id and body
// @Accept json
// @Produce json
// @Param authorization header string true "authorization"
// @Param Id path int true "Id"
// @Param transaction body request true "transaction"
// @Succes 200 {object} web.Response
// @Router /transacciones/{Id} [PUT]
func (t *Transaccion) Update() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var request request

		id, err := strconv.Atoi(ctx.Param("Id"))

		if err != nil {
			ctx.JSON(http.StatusBadRequest, web.NewResponse(http.StatusBadRequest, "No se selecciono la transaccion a actualizar", nil, err.Error()))
			return
		}

		if err := ctx.ShouldBindJSON(&request); err != nil {
			ctx.JSON(http.StatusBadRequest, web.NewResponse(http.StatusBadRequest, "La peticion no es valida", nil, err.Error()))
			return
		}

		if err := ValidarTransaccion(request); err != nil {
			ctx.JSON(http.StatusBadRequest, web.NewResponse(http.StatusBadRequest, "La peticion no es valida", nil, err.Error()))
			return
		}

		transaccion, err := t.service.Update(id, request.CodigoTransaccion, request.Moneda,
			request.Monto, request.Emisor, request.Receptor, request.FechaTransaccion)

		if err != nil {
			ctx.JSON(http.StatusNotFound, web.NewResponse(http.StatusNotFound, "Error al tratar de eliminar la transaccion", nil, err.Error()))
			return
		}

		ctx.JSON(http.StatusOK, web.NewResponse(http.StatusOK, "Transaccion actualizada con exito", transaccion, ""))
	}
}

// Update partiality a specific transaction
// @Summary Patch transaction
// @Tags Transaction
// @Description Patch an specific transaction using the id and body
// @Accept json
// @Produce json
// @Param authorization header string true "authorization"
// @Param Id path int true "Id"
// @Param transaction body patchRequest true "transaction"
// @Succes 200 {object} web.Response
// @Router /transacciones/{Id} [PATCH]
func (t *Transaccion) Patch() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var request patchRequest

		id, err := strconv.Atoi(ctx.Param("Id"))

		if err != nil {
			ctx.JSON(http.StatusBadRequest, web.NewResponse(http.StatusBadRequest, "No se selecciono la transaccion a actualizar", nil, err.Error()))
			return
		}

		if err := ctx.ShouldBindJSON(&request); err != nil {
			ctx.JSON(http.StatusBadRequest, web.NewResponse(http.StatusBadRequest, "El request no es valido", nil, err.Error()))
			return
		}

		if request.CodigoTransaccion == "" || request.Monto <= 0 {
			ctx.JSON(http.StatusBadRequest, web.NewResponse(http.StatusBadRequest, "El request no es valido", nil, ""))
			return
		}

		transaccion, err := t.service.Patch(id, request.CodigoTransaccion, request.Monto)

		if err != nil {
			ctx.JSON(http.StatusNotFound, web.NewResponse(http.StatusNotFound, "Error al tratar de actualizar la transaccion", nil, err.Error()))
			return
		}

		ctx.JSON(http.StatusOK, web.NewResponse(http.StatusOK, "Transaccion actualizada con exito", transaccion, ""))
	}
}

// Delete a specific transaction
// @Summary Delete transaction
// @Tags Transaction
// @Description Delete an specific transaction using the id
// @Accept json
// @Produce json
// @Param authorization header string true "authorization"
// @Param Id path int true "Id"
// @Succes 200 {object} web.Response
// @Router /transacciones/{Id} [DELETE]
func (t *Transaccion) Delete() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id, err := strconv.Atoi(ctx.Param("Id"))

		if err != nil {
			ctx.JSON(http.StatusBadRequest, web.NewResponse(http.StatusBadRequest, "No se selecciono la transaccion a eliminar", nil, err.Error()))
			return
		}

		if err := t.service.Delete(id); err != nil {
			ctx.JSON(http.StatusNotFound, web.NewResponse(http.StatusNotFound, "Ocurrio un error al eliminar la transaccion", nil, err.Error()))
			return
		}
		ctx.JSON(http.StatusOK, web.NewResponse(http.StatusOK, "Transaccion eliminada con exito", nil, ""))
	}
}
