package ServerOpenApi

import (
	//	"fmt"
	"net/http"
	"log"
	"github.com/osbuild/osbuild-installer-backend/internal/api-openapitools"

)


type Handlers struct { }
func (s *Handlers) OsbuildInstallerViewsV1GetVersion() (interface{}, error) {
	return "1", nil
}

func Run(address string, routePrefix string) {
	var s Handlers
	DefaultApiController := openapi.NewDefaultApiController(&s)
	router := openapi.NewRouter(DefaultApiController)
	log.Fatal(http.ListenAndServe(address, router))
}
