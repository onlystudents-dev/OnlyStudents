package opaquepkg

import (
	"context"
	"errors"
	"log/slog"

	db_queries "onlystudents/internal/db/store"

	"github.com/bytemare/opaque"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

var DemoPassword = []byte("test")

func SeedDemo(ctx context.Context, pool *pgxpool.Pool, server *opaque.Server) error {
	queries := db_queries.New(pool)

	people := []struct {
		role string
		id   int32
	}{
		{role: "student", id: 1},
		{role: "teacher", id: 1},
		{role: "guardian", id: 1},
	}

	for _, person := range people {
		account, err := demoAccount(ctx, queries, person.role, person.id)
		if err != nil {
			slog.Warn("demo seed: account not found, skipping", "role", person.role, "id", person.id, "err", err)
			continue
		}

		if err := registerDemoAccount(ctx, queries, server, account); err != nil {
			return err
		}
	}

	return nil
}

func demoAccount(ctx context.Context, queries *db_queries.Queries, role string, id int32) (db_queries.Account, error) {
	v := pgtype.Int4{Int32: id, Valid: true}

	switch role {
	case "student":
		return queries.GetAccountByStudentID(ctx, v)
	case "teacher":
		return queries.GetAccountByTeacherID(ctx, v)
	case "guardian":
		return queries.GetAccountByGuardianID(ctx, v)
	default:
		return db_queries.Account{}, errors.New("unknown role")
	}
}

func registerDemoAccount(ctx context.Context, queries *db_queries.Queries, server *opaque.Server, account db_queries.Account) error {
	client, err := Conf.Client()
	if err != nil {
		return err
	}

	credID := []byte(account.ID.String())

	req, err := client.RegistrationInit(DemoPassword)
	if err != nil {
		return err
	}

	resp, err := server.RegistrationResponse(req, credID, nil)
	if err != nil {
		return err
	}

	// match JS client
	record, _, err := client.RegistrationFinalize(resp, nil, nil,
		&opaque.ClientOptions{KSFSalt: make([]byte, 16), KSFLength: 64})
	if err != nil {
		return err
	}

	if err := queries.CreateOpaqueRecord(ctx, db_queries.CreateOpaqueRecordParams{
		AccountUuid:          account.ID,
		RegistrationRecord:   record.Serialize(),
		CredentialIdentifier: credID,
	}); err != nil {
		return err
	}

	slog.Info("demo seed: enrolled account with OPAQUE", "role", account.Role, "account_uuid", account.ID.String())

	return nil
}
