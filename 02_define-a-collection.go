package main

import (
	"fmt"
	"github.com/weaviate/weaviate-go-client/v4/weaviate"
	"github.com/weaviate/weaviate-go-client/v4/weaviate/auth"
	"github.com/weaviate/weaviate/entities/models"
	"golang.org/x/net/context"
	"os"
)

func main() {

	// Erstellt die Konfiguration für die Verbindung mit Weaviate.
	cfg := weaviate.Config{
		Host:       os.Getenv("WCD_HOSTNAME"),                    // Hostname der Weaviate-Instanz
		Scheme:     "https",                                      // Kommunikationsschema (HTTPS)
		AuthConfig: auth.ApiKey{Value: os.Getenv("WCD_API_KEY")}, // Authentifizierung mit API-Schlüssel
	}
	// Initialisiert den Weaviate-Client mit der oben definierten Konfiguration.
	client, err := weaviate.NewClient(cfg)
	// Gibt eine Fehlermeldung aus, falls die Client-Initialisierung fehlschlägt.
	if err != nil {
		fmt.Println(err)
	}

	// Definiert eine neue Klasse (Collection) für Weaviate.
	// Define the collection
	classObj := &models.Class{
		Class:      "Question",
		Vectorizer: "text2vec-cohere",
		ModuleConfig: map[string]interface{}{
			"text2vec-cohere":   map[string]interface{}{},
			"generative-cohere": map[string]interface{}{},
		},
	}

	// Fügt die definierte Klasse (Collection) zur Weaviate-Schema hinzu.
	err = client.Schema().ClassCreator().WithClass(classObj).Do(context.Background())
	// Beendet das Programm, falls das Hinzufügen der Klasse fehlschlägt.
	if err != nil {
		panic(err)
	}
}
