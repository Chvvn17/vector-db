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
	cfg := weaviate.Config{
		Host:       os.Getenv("WCD_HOSTNAME"),
		Scheme:     "https",
		AuthConfig: auth.ApiKey{Value: os.Getenv("WCD_API_KEY")},
		Headers: map[string]string{
			"X-Cohere-Api-Key": os.Getenv("COHERE_APIKEY"),
		},
	}

	client, err := weaviate.NewClient(cfg)
	if err != nil {
		fmt.Println(err)
	}

	ctx := context.Background()

	// Definiert ein Prompt für die Generierung von Inhalten.
	generatePrompt := "How does a heart work."

	// Führt eine GraphQL-Abfrage aus, um Daten aus der Klasse "Question" zu holen.
	gs := graphql.NewGenerativeSearch().GroupedResult(generatePrompt)

	response, err := client.GraphQL().Get().
		WithClassName("Question"). // Gibt an, dass die Abfrage auf der Klasse "Question" basiert.
		WithFields(
			graphql.Field{Name: "question"}, // Holt das Feld "question".
			graphql.Field{Name: "answer"},   // Holt das Feld "answer".
			graphql.Field{Name: "category"}, // Holt das Feld "category".
		).
		WithGenerativeSearch(gs).                           // Fügt die generative Suche basierend auf dem Prompt hinzu.
		WithNearText(client.GraphQL().NearTextArgBuilder(). // Filtert die Abfrage mit "NearText".
									WithConcepts([]string{"biology"})). // Sucht nach Einträgen, die sich auf das Konzept "biology" beziehen.
		WithLimit(2).                                       // Beschränkt die Abfrage auf maximal 2 Ergebnisse.
		Do(ctx)                                             // Führt die Abfrage mit dem erstellten Kontext aus.

	if err != nil {
		// Beendet das Programm, falls die Abfrage fehlschlägt.
		panic(err)
	}
	// Gibt die Antwort der GraphQL-Abfrage aus.
	fmt.Printf("%v", response)
}
