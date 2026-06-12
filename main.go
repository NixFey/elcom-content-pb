// main.go
package main

import (
	"log"
	"os"
	"strings"

	"net/http"
	"net/url"

	"bytes"
	"encoding/json"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
)

func main() {
	app := pocketbase.New()

	callback_url := os.Getenv("CMS_CALLBACK_URL")
	if len(callback_url) > 0 {
		for _, v := range strings.Split(callback_url, ";") {
			_, err := url.ParseRequestURI(v)
			if err != nil {
				panic(err)
			}

			handleModelUpdate := func(e *core.ModelEvent) error {
				json, err := json.Marshal(e.Model)
				jsonBody := bytes.NewBuffer(json)
				_, err = http.Post(v, "application/json", jsonBody)
				if err != nil {
					log.Fatalln(err)
				}
				return nil
			}

			app.OnModelAfterCreate().Add(handleModelUpdate)
			app.OnModelAfterUpdate().Add(handleModelUpdate)
			app.OnModelAfterDelete().Add(handleModelUpdate)
		}
	}

	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}
