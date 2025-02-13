package storage

import (
	"context"
	"log"
	"url_shorter/internal/storage/memory"
	"url_shorter/internal/storage/postgres"
)

type Storage interface {
	Close()
	UrlInsert(origUrl string) (string, error)
	UrlRead(shortURL string) (string, error)
}

var store Storage

func InitStorage(storageType *string) {
	switch *storageType {
	case "memory":
		store = memory.Init()
	case "postgres":
		ctx := context.Background()
		store = postgres.Init(ctx)
	default:
		log.Fatal("Invalid storage type. Use 'memory' or 'postgres'")
	}
}
func GetStorage() *Storage {
	return &store
}
