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
	rec, err := conf.GetFakeRecord([]byte("fake-credential-id-pad-to-64-bytes-00000000000000000000000"))
	if err != nil {
		return err
	}
	FakeRecord = rec
	return nil
}

type EnrollState struct {
	AccountUUID string `json:"account_uuid"`
}

type loginState struct {
	AccountUUID string `json:"account_uuid"`
	ClientMAC   []byte `json:"client_mac"`
}

func EnrollInit(ctx context.Context, s *opaque.Server, rdb *redis.Client, enroll_token string, reqBytes []byte) (respBytes []byte, err error) {
	if enroll_token == "" || !ValidEnrollToken(enroll_token) {
		return []byte{}, errors.New("Wrong args")
	}

	registration_req, err := s.Deserialize.RegistrationRequest(reqBytes)

	if err != nil {
		return []byte{}, err
	}

	raw, err := rdb.Get(ctx, "enroll:"+enroll_token).Bytes()

	var enroll_state EnrollState

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

	var enroll_state EnrollState

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

func ChangePasswordInit(ctx context.Context, s *opaque.Server, rdb *redis.Client, record *opaque.ClientRecord, ke1Bytes []byte) (ke2Bytes []byte, loginhandle string, err error) {
	ke1, err := s.Deserialize.KE1(ke1Bytes)

	if err != nil {
		return []byte{}, "", err
	}

	ke2, serverOutput, err := s.GenerateKE2(ke1, record)

	if err != nil {
		return []byte{}, "", err
	}

	change_password_handle := rand.Text()

	change_password_state_json, err := json.Marshal(loginState{
		AccountUUID: string(record.CredentialIdentifier),
		ClientMAC:   serverOutput.ClientMAC,
	})

	if err != nil {
		return []byte{}, "", err
	}

	ttl := time.Duration(helpers.GetInt64EnvFallback("OPAQUE_CHANGE_PASSWORD_TTL", 120, 600)) * time.Second

	rdb_err := rdb.Set(ctx, fmt.Sprintf("opaque_change_password:%s", change_password_handle), change_password_state_json, ttl).Err()

	if rdb_err != nil {
		return []byte{}, "", rdb_err
	}

	return ke2.Serialize(), change_password_handle, nil
}

func ChangePasswordFinish(ctx context.Context, s *opaque.Server, rdb *redis.Client, pool *pgxpool.Pool, change_password_handle string, ke3Bytes []byte) (string, string, error) {
	if change_password_handle == "" {
		return "", "", errors.New("Missing args")
	}

	ke3, err := s.Deserialize.KE3(ke3Bytes)

	if err != nil {
		return "", "", err
	}

	raw, err := rdb.GetDel(ctx, fmt.Sprintf("opaque_change_password:%s", change_password_handle)).Bytes()

	var change_password_state loginState

	if err != nil {
		return "", "", err
	}

	json_err := json.Unmarshal(raw, &change_password_state)

	if json_err != nil {
		return "", "", json_err
	}

	if err := s.LoginFinish(ke3, change_password_state.ClientMAC); err != nil {
		return "", "", err
	}

	parsed_account_uuid, err := uuid.Parse(change_password_state.AccountUUID)

	if err != nil {
		return "", "", err
	}

	queries := db_queries.New(pool)

	account, err := queries.GetAccountByUUID(ctx, pgtype.UUID{Bytes: parsed_account_uuid, Valid: true})

	if err != nil {
		return "", "", err
	}

	db_revoke_err := queries.RevokeAllSessions(ctx, pgtype.UUID{Bytes: parsed_account_uuid, Valid: true})

	if db_revoke_err != nil {
		return "", "", db_revoke_err
	}

	serr := helpers.SessionRevokeAll(ctx, rdb, parsed_account_uuid.String())

	if serr != nil {
		return "", "", serr
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
		return "", "", errors.New("invalid role")
	}

	session_uuid, err := uuid.NewRandom()

	if err != nil {
		return "", "", err
	}

	create_session_err := queries.CreateSession(ctx, db_queries.CreateSessionParams{
		ID:          pgtype.UUID{Bytes: session_uuid, Valid: true},
		AccountUuid: pgtype.UUID{Bytes: parsed_account_uuid, Valid: true},
	})

	if create_session_err != nil {
		return "", "", create_session_err
	}

	enroll_state_bytes, err := json.Marshal(EnrollState{
		AccountUUID: parsed_account_uuid.String(),
	})

	if err != nil {
		return "", "", err
	}

	crypto_grant_handle := rand.Text()
	ttl := time.Duration(helpers.GetInt64EnvFallback("OPAQUE_CRYPTO_GRANT_TTL", 300, 900)) * time.Second
	rdb_err := rdb.Set(ctx, fmt.Sprintf("crypto_grant:%s", crypto_grant_handle), enroll_state_bytes, ttl).Err()

	if rdb_err != nil {
		return "", "", rdb_err
	}

	session_token, serr := helpers.SessionCreate(ctx, rdb, account_id, change_password_state.AccountUUID, session_uuid.String(), account.Role)

	if serr != nil {
		// we dont know if a session was actually created, so we cannot delete it. Should be random enough that it doesnt cause issues.
		return "", "", serr
	}

	return session_token, crypto_grant_handle, nil
}

func RekeyInit(ctx context.Context, s *opaque.Server, account_uuid string, reqBytes []byte) (respBytes []byte, err error) {
	registration_req, err := s.Deserialize.RegistrationRequest(reqBytes)

	if err != nil {
		return []byte{}, err
	}

	resp, err := s.RegistrationResponse(registration_req, []byte(account_uuid), nil)

	if err != nil {
		return []byte{}, err
	}

	return resp.Serialize(), nil
}

func RekeyFinish(ctx context.Context, s *opaque.Server, rdb *redis.Client, pool *pgxpool.Pool, account_uuid string, crypto_grant_handle string, wrapped_mk_password []byte, recBytes []byte) error {
	if crypto_grant_handle == "" {
		return errors.New("Missing args")
	}

	rec, err := s.Deserialize.RegistrationRecord(recBytes)

	if err != nil {
		return err
	}

	raw, err := rdb.Get(ctx, "crypto_grant:"+crypto_grant_handle).Bytes()

	var crypto_grant EnrollState

	if err != nil {
		return err
	}

	unmarshal_err := json.Unmarshal(raw, &crypto_grant)

	if unmarshal_err != nil {
		return unmarshal_err
	}

	if crypto_grant.AccountUUID != account_uuid {
		return errors.New("session and crypto grant mismatch")
	}

	tx, err := pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	queries := db_queries.New(pool).WithTx(tx)

	parsed_account_uuid, err := uuid.Parse(crypto_grant.AccountUUID)

	if err != nil {
		return err
	}

	rows, opaque_record_update_err := queries.UpdateOpaqueRecord(ctx, db_queries.UpdateOpaqueRecordParams{
		AccountUuid:          pgtype.UUID{Bytes: parsed_account_uuid, Valid: true},
		RegistrationRecord:   rec.Serialize(),
		CredentialIdentifier: []byte(crypto_grant.AccountUUID),
	})

	if rows != 1 {
		return errors.New("Opaque record update did not match any rows in rekey finish")
	}

	if opaque_record_update_err != nil {
		return opaque_record_update_err
	}

	rows, account_crypto_update_err := queries.UpdateWrappedMKPassword(ctx, db_queries.UpdateWrappedMKPasswordParams{
		WrappedMkPassword: wrapped_mk_password,
		AccountUuid:       pgtype.UUID{Bytes: parsed_account_uuid, Valid: true},
	})

	if rows != 1 {
		return errors.New("Account Crypto update did not match any rows in rekey finish")
	}

	if account_crypto_update_err != nil {
		return account_crypto_update_err
	}

	commit_err := tx.Commit(ctx)

	if commit_err != nil {
		return commit_err
	}

	if err := rdb.Del(ctx, "crypto_grant:"+crypto_grant_handle).Err(); err != nil {
		return err
	}

	return nil
}
