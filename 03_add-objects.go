package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/weaviate/weaviate-go-client/v4/weaviate"
	"github.com/weaviate/weaviate-go-client/v4/weaviate/auth"
	"github.com/weaviate/weaviate/entities/models"
)

func main() {
	// Erstellt die Konfiguration für die Verbindung mit Weaviate.
	cfg := weaviate.Config{
		Host:       os.Getenv("WCD_HOSTNAME"),
		Scheme:     "https",
		AuthConfig: auth.ApiKey{Value: os.Getenv("WCD_API_KEY")},
		Headers: map[string]string{
			"X-Cohere-Api-Key": os.Getenv("COHERE_APIKEY"),
		},
	}
	// Initialisiert den Weaviate-Client mit der oben definierten Konfiguration.
	client, err := weaviate.NewClient(cfg)
	if err != nil {
		// Gibt eine Fehlermeldung aus, falls die Client-Initialisierung fehlschlägt.
		fmt.Println(err)
	}

	// Ruft Daten von einer angegebenen URL ab.
	data, err := http.DefaultClient.Get("https://raw.githubusercontent.com/weaviate-tutorials/quickstart/main/data/jeopardy_tiny.json")
	if err != nil {
		// Beendet das Programm, wenn die Daten nicht abgerufen werden können.
		panic(err)
	}
	// Stellt sicher, dass der Body geschlossen wird, nachdem er verarbeitet wurde.
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(data.Body)

	// Dekodiert die JSON-Daten in eine Variable `items`.
	var items []map[string]string
	if err := json.NewDecoder(data.Body).Decode(&items); err != nil {
		// Beendet das Programm, falls die Dekodierung fehlschlägt.
		panic(err)
	}

	// Konvertiert die Daten in eine Slice von `models.Object`, die von Weaviate verarbeitet werden kann.
	objects := make([]*models.Object, len(items))
	for i := range items {
		objects[i] = &models.Object{
			Class: "Tiere",
			Properties: map[string]any{
				"category": items[i]["Category"],
				"question": items[i]["Question"],
				"answer":   items[i]["Answer"],
			},
		}
	}

	// Führt ein Batch-Schreiben der Datenobjekte zu Weaviate aus.
	batchRes, err := client.Batch().ObjectsBatcher().WithObjects(objects...).Do(context.Background())
	if err != nil {
		// Beendet das Programm, falls das Batch-Schreiben fehlschlägt.
		panic(err)
	}

	// Überprüft die Ergebnisse des Batch-Schreibens auf Fehler.
	for _, res := range batchRes {
		if res.Result.Errors != nil {
			// Beendet das Programm, falls Fehler im Batch-Ergebnis auftreten
			panic(res.Result.Errors.Error)
		}
	}
}
