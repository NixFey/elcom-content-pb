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
		for v := range strings.SplitSeq(callback_url, ";") {
			_, err := url.ParseRequestURI(v)
			if err != nil {
				panic(err)
			}

			handleModelUpdate := func(e *core.RecordEvent) error {
				json, err := json.Marshal(e.Record)
				jsonBody := bytes.NewBuffer(json)
				_, err = http.Post(v, "application/json", jsonBody)
				if err != nil {
					log.Fatalln(err)
				}
				return e.Next()
			}

			app.OnRecordAfterCreateSuccess().BindFunc(handleModelUpdate)
			app.OnRecordAfterUpdateSuccess().BindFunc(handleModelUpdate)
			app.OnRecordAfterDeleteSuccess().BindFunc(handleModelUpdate)
		}
	}

	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}
