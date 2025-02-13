package postgres

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pashagolub/pgxmock/v2"
	"log"
	"sync"
	"url_shorter/config"
	"url_shorter/inter/shorter"
)

type PgStorage struct {
	mu       sync.RWMutex
	pool     *pgxpool.Pool
	MockPool pgxmock.PgxPoolIface
}

func Init(ctx context.Context) *PgStorage {
	cfg := config.GetConfig()

	var err error
	pool, err := pgxpool.New(ctx, cfg.PostgresURL)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}

	log.Println("Postgres сonnect success")
	return &PgStorage{pool: pool}
}
func (p *PgStorage) Close() {
	if p.pool != nil {
		p.pool.Close()
		log.Println("Database connect closed")
	}
}

func (p *PgStorage) UrlInsert(origUrl string) (string, error) {

	p.mu.Lock()
	defer p.mu.Unlock()
	var exists bool
	var err error
	err = p.pool.QueryRow(context.Background(), "SELECT EXISTS (SELECT 1 FROM UrlShorter WHERE origUrl = $1)", origUrl).Scan(&exists)
	if err != nil {
		return "", err
	}
	var shortUrl string
	if exists == true {
		err = p.pool.QueryRow(context.Background(), "SELECT shortUrl FROM UrlShorter WHERE origUrl = $1 )", origUrl).Scan(&shortUrl)
		return shortUrl, err
	} else {
		shortUrl = shorter.UrlShorter(origUrl)
		var dbAnswer string
		err = p.pool.QueryRow(context.Background(), "INSERT INTO UrlShorter (origUrl, shortUrl) VALUES ($1, $2)", origUrl, shortUrl).Scan(&dbAnswer)
		return shortUrl, err
	}
}

func (p *PgStorage) UrlRead(shortUrl string) (string, error) {
	p.mu.Lock()
	defer p.mu.Unlock()
	var origUrl string
	err := p.pool.QueryRow(context.Background(), "SELECT origUrl FROM UrlShorter WHERE shortUrl = $1", shortUrl).Scan(&origUrl)
	if err != nil {
		return "", errors.New("short URL not found")
	}
	return origUrl, nil
}
