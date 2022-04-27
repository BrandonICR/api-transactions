package engine

import (
	"os"

	"github.com/BrandonICR/web_cl2_050422_8am/cmd/server/handler"
	"github.com/BrandonICR/web_cl2_050422_8am/cmd/server/route"
	"github.com/BrandonICR/web_cl2_050422_8am/docs"
	"github.com/BrandonICR/web_cl2_050422_8am/pkg/store"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
)

func copyFileStore(fileStore string, tempFileStore string) error {
	data, err := os.ReadFile(fileStore)
	if err != nil {
		return err
	}

	if err := os.WriteFile(tempFileStore, data, 0666); err != nil {
		return err
	}

	return nil
}

func GetEngine(fileStore string, tempFileStore string, fileEnv string) *gin.Engine {
	if fileEnv != "" {
		if err := godotenv.Load(fileEnv); err != nil {
			panic("error: no se lograron cargar las variables de entorno")
		}
	} else {
		if err := godotenv.Load(); err != nil {
			panic("error: no se lograron cargar las variables de entorno")
		}
	}

	if tempFileStore != "" {
		if err := copyFileStore(fileStore, tempFileStore); err != nil {
			panic("error al crear file store temporal")
		}
		fileStore = tempFileStore
	}

	store := store.NewStore(store.JsonFileType, fileStore)

	router := gin.Default()

	docs.SwaggerInfo.Host = os.Getenv("HOST")
	router.GET("docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	router.Use(handler.ValidarToken())
	routes := route.NewRouter(router, &store)
	routes.MapRoutes()

	return router
}
