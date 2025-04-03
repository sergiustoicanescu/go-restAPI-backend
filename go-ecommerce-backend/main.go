package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq"
	"github.com/sergiustoicanescu/go-restAPI-backend/go-ecommerce-backend/config"
	"github.com/sergiustoicanescu/go-restAPI-backend/go-ecommerce-backend/migrations"
	"github.com/sergiustoicanescu/go-restAPI-backend/go-ecommerce-backend/routes"
)

func main() {
	cfg, err := config.LoadConfig("../.env")
	if err != nil {
		log.Fatalf("failed to load configuration: %v", err)
	}

	db, err := sql.Open("postgres", cfg.DBProd)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	if err := migrations.RunMigrations(db, "ecommerce", "file://migrations"); err != nil {
		log.Fatalf("Migration failed: %v", err)
	}

	r := routes.SetupRoutes(db, cfg)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Server starting on port %s", port)
	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
