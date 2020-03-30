package main

import (
	"bytes"
	"text/template"
	"fmt"
	"os"

	"github.com/osbuild/osbuild-installer-backend/internal/server-openapitools"

	"github.com/getkin/kin-openapi/openapi3"
)

type Prefix struct {
	PathPrefix  string
	AppName     string
}

func main() {
	// PATH_PREFIX
	// APP_NAME

	pathPrefix, ok := os.LookupEnv("PATH_PREFIX")
	if !ok {
		pathPrefix = "api"
	}
	appName, ok := os.LookupEnv("APP_NAME")
	if !ok {
		appName = "osbuild_installer"
	}
	prefix := Prefix{pathPrefix, appName}

	// e.group pathPrefix, appName

	tmpl, err := template.ParseFiles("osbuild_installer/openapi/api.spec.yaml")
	if err != nil {
		panic(err)
	}

	var buffer bytes.Buffer
	err = tmpl.Execute(&buffer, prefix)
	if err != nil {
		panic(err)
	}

	swagger, err := openapi3.NewSwaggerLoader().LoadSwaggerFromData(buffer.Bytes())
	if err != nil {
		panic(err)
	}
	fmt.Println("URL: localhost:8086/version")
	ServerOpenApi.Run("localhost:8086", swagger.Servers[0].URL)
}
