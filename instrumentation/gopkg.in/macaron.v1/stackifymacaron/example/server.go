package main

import (
	"log"
	"net/http"

	"gopkg.in/macaron.v1"

	apm "github.com/stackify/stackify-go-apm"
	"github.com/stackify/stackify-go-apm/config"
	"github.com/stackify/stackify-go-apm/instrumentation/gopkg.in/macaron.v1/stackifymacaron"
)

func initStackifyTrace() (*apm.StackifyAPM, error) {
	return apm.NewStackifyAPM(
		config.WithApplicationName("Go Application"),
		config.WithEnvironmentName("Test"),
		config.WithDebug(true),
	)
}

func main() {
	stackifyAPM, err := initStackifyTrace()
	if err != nil {
		log.Fatalf("failed to initialize stackifyapm: %v", err)
	}
	defer stackifyAPM.Shutdown()

	m := macaron.Classic()
	m.Use(stackifymacaron.Middleware())

	m.Get("/index", func(ctx *macaron.Context) string {
		return "Hello World!"
	})

	http.ListenAndServe("0.0.0.0:8000", m)
}
