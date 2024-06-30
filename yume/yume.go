package yume

import (
	"context"
	"encoding/base64"
	"errors"

	"encore.dev/beta/errs"
	"encore.dev/config"
	"encore.dev/rlog"
	"encore.dev/storage/sqldb"
	"encore.dev/types/uuid"
	"github.com/pinecone-io/go-pinecone/pinecone"
)

var secrets struct {
	pineconeKey string
}

type YumeConfig struct {
	Prod            bool
	LogPotentialBug bool
}

var yumeconf = config.Load[*YumeConfig]()

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

// User /yume/user
// user management, TODO: connect to auth0 or something
type User struct {
	ID       uuid.UUID `encore:"sensitive"` // will return 400 if not valid uuid
	UID      string    `encore:"sensitive"` // user id from auth0, currently not used
	USERNAME string    `encore:"sensitive"` // username
}

type NewUserParams struct {
	USERNAME string `encore:"sensitive"`
}

//encore:api private path=/yume/user/new
func NewUser(ctx context.Context, p *NewUserParams) (*User, error) {
	id, err := uuid.NewV4()
	if err != nil {
		rlog.Error("uuid generation failed", "err", err)
		return nil, err
	}
	db.Exec(ctx, `INSERT INTO yuser (id, username) VALUES ($1, $2)`, id, p.USERNAME)
	return &User{ID: id, USERNAME: p.USERNAME}, nil
}

func checkUser(ctx context.Context, id string) error {
	// check if exist
	err := db.QueryRow(ctx, `SELECT (id) FROM yuser WHERE id = $1`, id).Scan(&id)
	if err != nil {
		if errors.Is(err, sqldb.ErrNoRows) {
			return errors.New("user not found")
		}
		rlog.Error("sql bug", "err", err)
		return err
	}
	return nil
}

// Notes/Docs /yume/documents
// for user-added documents
type Document struct {
	ID        uuid.UUID `encore:"sensitive"`
	OWNER     string    `encore:"sensitive"` // id of the user who added the document
	NAME      string    `encore:"sensitive"` // name of the document
	embedding []float32 `encore:"sensitive"` // embedding of the document, generated by the model
	// TODO: add permissions
}

// TODO: client-side encryption for content
type NewDocumentParams struct {
	OWNER          string `encore:"sensitive"`
	NAME           string `encore:"sensitive"`
	CONTENT_BASE64 string `encore:"sensitive"` // please base64 encode the content
}

func IsBase64(s string) bool {
	_, err := base64.StdEncoding.DecodeString(s)
	return err == nil
}

//encore:api public path=/yume/documents/new
func NewDocument(ctx context.Context, p *NewDocumentParams) (*Document, error) {
	if !IsBase64(p.CONTENT_BASE64) {
		if yumeconf.LogPotentialBug {
			rlog.Error("content in NewDocument is not valid base64, frontend bug?")
		}
		return nil, &errs.Error{Code: errs.InvalidArgument, Message: "content is not base64 encoded"}
	}
	id, err := uuid.NewV4()
	if err != nil {
		rlog.Error("uuid generation failed", "err", err)
		return nil, err
	}
	if checkUser(ctx, p.OWNER) != nil {
		if yumeconf.LogPotentialBug {
			rlog.Error("user not found", "id", p.OWNER)
		}
		return nil, &errs.Error{Code: errs.NotFound, Message: "failed to find user with id"}
	}
	db.Exec(ctx, `INSERT INTO document (id, namespace_id, name, b64_content) VALUES ($1, $2, $3, $4)`, id, p.OWNER, p.NAME, p.CONTENT_BASE64)
	return &Document{ID: id, OWNER: p.OWNER, NAME: p.NAME}, nil
}

func getDocument(ctx context.Context, id string) error {
	// check if exist
	err := db.QueryRow(ctx, `SELECT (id) from document WHERE id = $1`, id).Scan(&id)
	if err != nil {
		return err
	}
	return nil
}

// Namespace /yume/userns
// currently only in atari serverless pod [dimension: 768, metric: cosine, spec: aws us-east-1]
type UserNamespace struct {
	ID    uuid.UUID
	NAME  string    `encore:"sensitive"`
	OWNER uuid.UUID `encore:"sensitive"`
}

type NewUserNamespaceParams struct {
	NAME  string    `encore:"sensitive"`
	OWNER uuid.UUID `encore:"sensitive"`
}

// TODO: generate embedding and manage the pineconedb from typescript backend

//encore:api public path=/yume/userns/new
func NewUserNamespace(ctx context.Context, p *NewUserNamespaceParams) (*UserNamespace, error) {
	id, err := uuid.NewV4()
	if err != nil {
		return nil, err
	}
	if checkUser(ctx, p.OWNER.String()) != nil {
		return nil, &errs.Error{Code: errs.NotFound, Message: "failed to find user with id"}
	}
	db.Exec(ctx, `INSERT INTO namespace (id, name, owner) VALUES ($1, $2, $3)`, id, p.NAME, p.OWNER.String())
	return &UserNamespace{ID: id, NAME: p.NAME, OWNER: p.OWNER}, nil
}

type QueryUserNamespaceParams struct {
	DOCUMENTID uuid.UUID `encore:"sensitive"` // id of the document
}

func getUserNS(ctx context.Context, id string) error {
	// check if exist
	err := db.QueryRow(ctx, `SELECT (id) from namespace WHERE id = $1`, id).Scan(&id)
	if err != nil {
		if errors.Is(err, sqldb.ErrNoRows) {
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
	if err != nil {
		if err.(*errs.Error).Code == errs.NotFound {
			rlog.Error("invalid namespace id, frontend bug?", "id", id)
		}
		return err
	}
	// check if document exist
	err = getDocument(ctx, p.DOCUMENTID.String())
	if err != nil {
		rlog.Error("document not found", "err", err)
		return &errs.Error{Code: errs.NotFound, Message: "document not found"}
	}
	db.Exec(ctx, `INSERT INTO namespace_document (namespace_id, document_id) VALUES ($1, $2)`, id, p.DOCUMENTID.String())
	return err
}

// Define a database named 'url', using the database
// migrations  in the "./migrations" folder.
// Encore provisions, migrates, and connects to the database.
// Learn more: https://encore.dev/docs/primitives/databases
var db = sqldb.NewDatabase("yume", sqldb.DatabaseConfig{
	Migrations: "./migrations",
})
