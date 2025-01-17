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

	cfg := weaviate.Config{
		Host:       os.Getenv("WCD_HOSTNAME"),
		Scheme:     "https",
		AuthConfig: auth.ApiKey{Value: os.Getenv("WCD_API_KEY")},
	}

	client, err := weaviate.NewClient(cfg)
	if err != nil {
		fmt.Println(err)
	}

	// Define the collection
	classObj := &models.Class{
		Class:      "Question",
		Vectorizer: "text2vec-cohere",
		ModuleConfig: map[string]interface{}{
			"text2vec-cohere":   map[string]interface{}{},
			"generative-cohere": map[string]interface{}{},
		},
	}

	// add the collection
	err = client.Schema().ClassCreator().WithClass(classObj).Do(context.Background())
	if err != nil {
		panic(err)
	}
}
