package main

import (
	"context"
	"fmt"
	"os"

	"github.com/weaviate/weaviate-go-client/v4/weaviate"
	"github.com/weaviate/weaviate-go-client/v4/weaviate/auth"
	"github.com/weaviate/weaviate-go-client/v4/weaviate/graphql"
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

	// Gibt eine Fehlermeldung aus, falls die Client-Initialisierung fehlschlägt.
	if err != nil {
		fmt.Println(err)
	}

	// Erstellt einen Kontext für die Anfrage.
	ctx := context.Background()

	// Führt eine GraphQL-Abfrage aus, um Daten aus der Klasse "Question" zu holen.
	response, err := client.GraphQL().Get().
		WithClassName("Question").
		WithFields(
			graphql.Field{Name: "question"},
			graphql.Field{Name: "answer"},
			graphql.Field{Name: "category"},
		).
		// Filtert die Ergebnisse basierend auf ihrer semantischen Nähe zum Konzept "biology".
		WithNearText(client.GraphQL().NearTextArgBuilder().
			WithConcepts([]string{"biology"})).
		// Es werden nur die relevantesten Ergebnisse zurückgegeben, begrenzt auf 1 Eintrag.
		WithLimit(1).
		// Führt die Abfrage aus und sendet sie an den Weaviate-Server
		Do(ctx)

	if err != nil {
		panic(err)
	}
	// Gibt die Antwort der GraphQL-Abfrage aus.
	fmt.Printf("%v", response)
}
