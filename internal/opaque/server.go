package opaquepkg

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	db_queries "onlystudents/internal/db/store"
	"onlystudents/internal/helpers"
	"time"

	"github.com/bytemare/opaque"
	"github.com/goccy/go-json"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

var FakeRecord *opaque.ClientRecord

func InitFakeRecord(conf *opaque.Configuration) error {
	rec, err := conf.GetFakeRecord([]byte("onlystudents-fake-credential-id-pad-to-64-bytes-0000000000"))
	if err != nil {
		return err
	}
	FakeRecord = rec
	return nil
}

type enrollState struct {
	AccountUUID string `json:"account_uuid"`
}

type loginState struct {
	AccountUUID string `json:"account_uuid"`
	ClientMAC   []byte `json:"client_mac"`
}

func EnrollInit(ctx context.Context, s *opaque.Server, rdb *redis.Client, enroll_token string, reqBytes []byte) (respBytes []byte, err error) {
	if enroll_token == "" {
		return []byte{}, errors.New("Missing args")
	}

	registration_req, err := s.Deserialize.RegistrationRequest(reqBytes)

	if err != nil {
		return []byte{}, err
	}

	raw, err := rdb.Get(ctx, "enroll:"+enroll_token).Bytes()

	var enroll_state enrollState

	if err != nil {
		return []byte{}, err
	}

	unmarshal_err := json.Unmarshal(raw, &enroll_state)

	if unmarshal_err != nil {
		return []byte{}, unmarshal_err
	}

	resp, err := s.RegistrationResponse(registration_req, []byte(enroll_state.AccountUUID), nil)

	if err != nil {
		return []byte{}, err
	}

	return resp.Serialize(), nil
}

func EnrollFinish(ctx context.Context, s *opaque.Server, rdb *redis.Client, pool *pgxpool.Pool, enroll_token string, recBytes []byte) (string, error) {
	if enroll_token == "" {
		return "", errors.New("Missing args")
	}

	rec, err := s.Deserialize.RegistrationRecord(recBytes)

	if err != nil {
		return "", err
	}

	raw, err := rdb.GetDel(ctx, "enroll:"+enroll_token).Bytes()

	var enroll_state enrollState

	if err != nil {
		return "", err
	}

	unmarshal_err := json.Unmarshal(raw, &enroll_state)

	if unmarshal_err != nil {
		return "", unmarshal_err
	}

	queries := db_queries.New(pool)

	parsed_account_uuid, err := uuid.Parse(enroll_state.AccountUUID)

	if err != nil {
		return "", err
	}

	opaque_err := queries.CreateOpaqueRecord(ctx, db_queries.CreateOpaqueRecordParams{
		AccountUuid:          pgtype.UUID{Bytes: parsed_account_uuid, Valid: true},
		RegistrationRecord:   rec.Serialize(),
		CredentialIdentifier: []byte(enroll_state.AccountUUID),
	})

	if opaque_err != nil {
		return "", opaque_err
	}

	account, err := queries.GetAccountByUUID(ctx, pgtype.UUID{Bytes: parsed_account_uuid, Valid: true})

	if err != nil {
		return "", err
	}

	var account_id int32

	switch account.Role {
	case "student":
		account_id = account.StudentID.Int32
	case "teacher":
		account_id = account.TeacherID.Int32
	case "guardian":
		account_id = account.GuardianID.Int32
	default:
		return "", errors.New("invalid role")
	}

	session_uuid, err := uuid.NewRandom()

	if err != nil {
		return "", err
	}

	create_session_err := queries.CreateSession(ctx, db_queries.CreateSessionParams{
		ID:          pgtype.UUID{Bytes: session_uuid, Valid: true},
		AccountUuid: pgtype.UUID{Bytes: parsed_account_uuid, Valid: true},
	})

	if create_session_err != nil {
		return "", create_session_err
	}

	return helpers.SessionCreate(ctx, rdb, account_id, enroll_state.AccountUUID, session_uuid.String(), account.Role)
}

func LoginInit(ctx context.Context, s *opaque.Server, rdb *redis.Client, record *opaque.ClientRecord, ke1Bytes []byte) (ke2Bytes []byte, loginhandle string, err error) {
	ke1, err := s.Deserialize.KE1(ke1Bytes)

	if err != nil {
		return []byte{}, "", err
	}

	ke2, serverOutput, err := s.GenerateKE2(ke1, record)

	if err != nil {
		return []byte{}, "", err
	}

	login_handle := rand.Text()

	login_state_json, err := json.Marshal(loginState{
		AccountUUID: string(record.CredentialIdentifier),
		ClientMAC:   serverOutput.ClientMAC,
	})

	if err != nil {
		return []byte{}, "", err
	}

	ttl := time.Duration(helpers.GetInt64EnvFallback("OPAQUE_LOGIN_TTL", 120, 600)) * time.Second

	rdb_err := rdb.Set(ctx, fmt.Sprintf("opaque_login:%s", login_handle), login_state_json, ttl).Err()

	if rdb_err != nil {
		return []byte{}, "", rdb_err
	}

	return ke2.Serialize(), login_handle, nil
}

func LoginFinish(ctx context.Context, s *opaque.Server, rdb *redis.Client, pool *pgxpool.Pool, login_handle string, ke3Bytes []byte) (string, error) {
	if login_handle == "" {
		return "", errors.New("Missing args")
	}

	ke3, err := s.Deserialize.KE3(ke3Bytes)

	if err != nil {
		return "", err
	}

	raw, err := rdb.GetDel(ctx, fmt.Sprintf("opaque_login:%s", login_handle)).Bytes()

	var login_state loginState

	if err != nil {
		return "", err
	}

	json_err := json.Unmarshal(raw, &login_state)

	if json_err != nil {
		return "", json_err
	}

	if err := s.LoginFinish(ke3, login_state.ClientMAC); err != nil {
		return "", err
	}

	queries := db_queries.New(pool)

	parsed_account_uuid, err := uuid.Parse(login_state.AccountUUID)

	if err != nil {
		return "", err
	}

	account, err := queries.GetAccountByUUID(ctx, pgtype.UUID{Bytes: parsed_account_uuid, Valid: true})

	if err != nil {
		return "", err
	}

	var account_id int32

	switch account.Role {
	case "student":
		account_id = account.StudentID.Int32
	case "teacher":
		account_id = account.TeacherID.Int32
	case "guardian":
		account_id = account.GuardianID.Int32
	default:
		return "", errors.New("invalid role")
	}

	session_uuid, err := uuid.NewRandom()

	if err != nil {
		return "", err
	}

	create_session_err := queries.CreateSession(ctx, db_queries.CreateSessionParams{
		ID:          pgtype.UUID{Bytes: session_uuid, Valid: true},
		AccountUuid: pgtype.UUID{Bytes: parsed_account_uuid, Valid: true},
	})

	if create_session_err != nil {
		return "", create_session_err
	}

	return helpers.SessionCreate(ctx, rdb, account_id, login_state.AccountUUID, session_uuid.String(), account.Role)
}
