/*
 * Osbuild-installer backend service
 *
 * Service that relays image build requests
 *
 * API version: 1.0
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package openapi

import (
	_ "encoding/json"
	"net/http"
	"strings"

	_ "github.com/gorilla/mux"
)

// A DefaultApiController binds http requests to an api service and writes the service results to the http response
type DefaultApiController struct {
	service DefaultApiServicer
}

// NewDefaultApiController creates a default api controller
func NewDefaultApiController(s DefaultApiServicer) Router {
	return &DefaultApiController{ service: s }
}

// Routes returns all of the api route for the DefaultApiController
func (c *DefaultApiController) Routes() Routes {
	return Routes{ 
		{
			"OsbuildInstallerViewsV1GetVersion",
			strings.ToUpper("Get"),
			"/version",
			c.OsbuildInstallerViewsV1GetVersion,
		},
	}
}

// OsbuildInstallerViewsV1GetVersion - get the service version
func (c *DefaultApiController) OsbuildInstallerViewsV1GetVersion(w http.ResponseWriter, r *http.Request) {
	result, err := c.service.OsbuildInstallerViewsV1GetVersion()
	if err != nil {
		w.WriteHeader(500)
		return
	}
	
	EncodeJSONResponse(result, nil, w)
}
