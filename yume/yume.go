package yume

import (
	"context"

	"encore.dev/beta/errs"
	"encore.dev/rlog"
	"encore.dev/storage/sqldb"
	"encore.dev/types/uuid"
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

// Notes/Docs /yume/documents
// for user-added documents
type Document struct {
	OWNER     string    // id of the user who added the document
	NAME      string    // name of the document
	embedding []float32 // embedding of the document, generated by the model
	// TODO: add permissions
}

// Namespace /yume/userns
// currently only in atari serverless pod [dimension: 768, metric: cosine, spec: aws us-east-1]
type UserNamespace struct {
	ID    uuid.UUID
	NAME  string
	OWNER uuid.UUID
}

type NewUserNamespaceParams struct {
	NAME  string // what
	OWNER uuid.UUID
}

// TODO: generate embedding and manage the pineconedb from typescript backend

//encore:api public path=/yume/userns/new
func NewUserNamespace(ctx context.Context, p *NewUserNamespaceParams) (*UserNamespace, error) {
	// add to db
	id, err := uuid.NewV4()
	if err != nil {
		return nil, err
	}
	db.Exec(ctx, `INSERT INTO namespace (id, name, owner) VALUES ($1, $2, $3)`, id, p.NAME, p.OWNER.String())
	return &UserNamespace{ID: id, NAME: p.NAME, OWNER: p.OWNER}, nil
}

type QueryUserNamespaceParams struct {
	ASDF string
}

func getUserNS(ctx context.Context, id string) error {
	// check if exist
	err := db.QueryRow(ctx, `SELECT (id) from namespace WHERE id = $1`, id).Scan(&id)
	if err != nil {
		if err.Error() != "not_found: sql: no rows in result set" { // TODO: change this to sqldb.Err..., currently not work somehow
			return &errs.Error{Code: errs.NotFound, Message: "namespace id not found"}
		}
		rlog.Error("sql fetch failed", "err", err)
		return &errs.Error{Code: errs.Internal, Message: "internal error"}
	}
	return nil
}

//encore:api public path=/yume/userns/insert/:id
func InsertToUserNamespace(ctx context.Context, id string, p *QueryUserNamespaceParams) error {
	err := getUserNS(ctx, id)
	return err
}

// Define a database named 'url', using the database
// migrations  in the "./migrations" folder.
// Encore provisions, migrates, and connects to the database.
// Learn more: https://encore.dev/docs/primitives/databases
var db = sqldb.NewDatabase("yume", sqldb.DatabaseConfig{
	Migrations: "./migrations",
})
