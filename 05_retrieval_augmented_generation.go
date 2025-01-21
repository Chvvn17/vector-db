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
	generatePrompt := "Write a tweet with emojis about these facts"

	// Führt eine GraphQL-Abfrage aus, um Daten aus der Klasse "Question" zu holen.
	gs := graphql.NewGenerativeSearch().GroupedResult(generatePrompt)

	response, err := client.GraphQL().Get().
		WithClassName("Question"). // Gibt an, dass die Abfrage auf der Klasse "Question" basiert.
		WithFields(
			graphql.Field{Name: "question"},
			graphql.Field{Name: "answer"},
			graphql.Field{Name: "category"},
		).
		// Filtert die Ergebnisse basierend auf ihrer semantischen Nähe zum Konzept "biology".
		WithGenerativeSearch(gs).
		WithNearText(client.GraphQL().NearTextArgBuilder().
			WithConcepts([]string{"nature"})).
		// Es werden nur die relevantesten Ergebnisse zurückgegeben, begrenzt auf 1 Eintrag.
		WithLimit(1).
		// Führt die Abfrage aus und sendet sie an den Weaviate-Server
		Do(ctx)

	if err != nil {
		// Beendet das Programm, falls die Abfrage fehlschlägt.
		panic(err)
	}
	// Gibt die Antwort der GraphQL-Abfrage aus.
	fmt.Printf("%v", response)
}
