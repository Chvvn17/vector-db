package main

import (
	"context"
	"fmt"
	"os"

	"github.com/weaviate/weaviate-go-client/v4/weaviate"
	"github.com/weaviate/weaviate-go-client/v4/weaviate/auth"
)

func main() {
	// Erstellt eine Konfigurationsstruktur für die Verbindung mit Weaviate.
	// Die Host-Adresse und der API-Schlüssel werden aus Umgebungsvariablen gelesen.
	cfg := weaviate.Config{
		Host:       os.Getenv("WCD_HOSTNAME"),
		Scheme:     "https",
		AuthConfig: auth.ApiKey{Value: os.Getenv("WCD_API_KEY")},
	}

	// Erstellt einen neuen Weaviate-Client basierend auf der Konfiguration.
	client, err := weaviate.NewClient(cfg)
	if err != nil {
		// Gibt eine Fehlermeldung aus, falls die Client-Initialisierung fehlschlägt.
		fmt.Println(err)
	}

	// Überprüft die Verbindung mit dem Weaviate-Server.
	ready, err := client.Misc().ReadyChecker().Do(context.Background())
	if err != nil {
		// Beendet das Programm, wenn die Überprüfung fehlschlägt.
		panic(err)
	}
	// Gibt den Verbindungsstatus aus (true/false).
	fmt.Printf("%v", ready)
}
