package route

import (
	"github.com/BrandonICR/web_cl2_050422_8am/cmd/server/handler"
	"github.com/BrandonICR/web_cl2_050422_8am/internal/transacciones"
	"github.com/BrandonICR/web_cl2_050422_8am/pkg/store"
	"github.com/gin-gonic/gin"
)

type Router interface {
	MapRoutes()
}

type router struct {
	r  *gin.Engine
	rg *gin.RouterGroup
	db *store.Store
}

func NewRouter(r *gin.Engine, db *store.Store) Router {
	return &router{r: r, db: db}
}

func (r *router) MapRoutes() {
	r.setGroup()
	r.buildTransactionRoutes()
}

func (r *router) setGroup() {
	r.rg = r.r.Group("/api/v1/transacciones")
}

func (r *router) buildTransactionRoutes() {
	repository := transacciones.NewRepository(*r.db)
	service := transacciones.NewService(repository)
	transacciones := handler.NewTransaccion(service)

	r.rg.GET("", transacciones.GetAll())
	r.rg.GET("/", transacciones.GetTransaccionFiltrada())
	r.rg.POST("/:Id", transacciones.Store())
	r.rg.GET("/:Id", transacciones.GetTransaccion())
	r.rg.PUT("/:Id", transacciones.Update())
	r.rg.PATCH("/:Id", transacciones.Patch())
	r.rg.DELETE("/:Id", transacciones.Delete())
}
