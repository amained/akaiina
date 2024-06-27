package yume

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"log"

	"encore.dev/storage/sqldb"
	"github.com/pinecone-io/go-pinecone/pinecone"
)

var secrets struct {
	pineconeKey string
}

//encore:service
type Service struct {
	pineconeClient  *pinecone.Client
	indexConnection *pinecone.IndexConnection
}

func initService() (*Service, error) {
	pc, err := pinecone.NewClient(pinecone.NewClientParams{ApiKey: secrets.pineconeKey})
	if err != nil {
		return nil, err
	}
	indx, err := pc.Index("https://multilingual-e5-large-qnilgk5.svc.aped-4627-b74a.pinecone.io")
	if err != nil {
		return nil, err
	}

	return &Service{pineconeClient: pc, indexConnection: indx}, nil
}

// Namespace /yume/userns
// currently only in atari serverless pod [dimension: 768, metric: cosine, spec: aws us-east-1]
type UserNamespace struct {
	ID    string
	NAME  string
	OWNER string
}

type NewUserNamespaceParams struct {
	NAME  string // what
	OWNER string // the fuck
}

//encore:api public path=/yume/userns/new
func NewUserNamespace(ctx context.Context, p *NewUserNamespaceParams) (*UserNamespace, error) {
	// add to db
	id, err := generateID()
	log.Println("we tf")
	if err != nil {
		return nil, err
	}
	db.Exec(ctx, `INSERT INTO namespace (id, name) VALUES ($1, $2)`, id, p.NAME)
	return &UserNamespace{ID: id, NAME: p.NAME}, nil
}

// generateID generates a random short ID.
func generateID() (string, error) {
	var data [6]byte // 6 bytes of entropy
	if _, err := rand.Read(data[:]); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(data[:]), nil
}

// Define a database named 'url', using the database
// migrations  in the "./migrations" folder.
// Encore provisions, migrates, and connects to the database.
// Learn more: https://encore.dev/docs/primitives/databases
var db = sqldb.NewDatabase("yume", sqldb.DatabaseConfig{
	Migrations: "./migrations",
})
