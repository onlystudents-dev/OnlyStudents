package helpers

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"math"
	"strconv"
	"strings"

	"golang.org/x/crypto/argon2"
)

func Argon2HashPassword(password string) (string, error) {
	salt := make([]byte, GetUintEnvFallback("ARGON2_SALTLEN", 16))
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}

	argon2_time := GetUint32EnvFallback("ARGON2_TIME", 1, 100)
	argon2_mem := GetUint32EnvFallback("ARGON2_MEMORY", 64*1024, 1<<24)
	argon2_threads := GetUint8EnvFallback("ARGON2_THREADS", 4, 32)
	argon2_len := GetUint32EnvFallback("ARGON2_KEYLEN", 32, 1024)

	hash := argon2.IDKey(
		[]byte(password),
		salt,
		argon2_time,
		argon2_mem,
		argon2_threads,
		argon2_len,
	)

	// this is good because if we change the parameters in the future it will continue working
	hash_string := fmt.Sprintf("$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s", argon2.Version, argon2_mem, argon2_time, argon2_threads, base64.RawStdEncoding.EncodeToString(salt), base64.RawStdEncoding.EncodeToString(hash))

	return hash_string, nil
}

func Argon2Verify(password string, hash_string string) bool {
	if !strings.HasPrefix(hash_string, "$argon2id$") {
		return false
	}

	parts := strings.Split(hash_string, "$")
	if len(parts) != 6 {
		return false
	}

	if parts[2] != fmt.Sprintf("v=%d", argon2.Version) {
		return false
	}

	var argon2_time, argon2_mem, argon2_threads uint64
	for param := range strings.SplitSeq(parts[3], ",") {
		pair := strings.SplitN(param, "=", 2)
		if len(pair) != 2 {
			return false
		}
		value, err := strconv.ParseUint(pair[1], 10, 32)
		if err != nil {
			return false
		}
		switch pair[0] {
		case "m":
			argon2_mem = value
		case "t":
			argon2_time = value
		case "p":
			argon2_threads = value
		default:
			return false
		}
	}

	if argon2_time == 0 || argon2_mem == 0 || argon2_threads == 0 {
		return false
	}

	salt, err := base64.RawStdEncoding.DecodeString(parts[4])
	if err != nil {
		return false
	}
	stored_hash, err := base64.RawStdEncoding.DecodeString(parts[5])
	if err != nil {
		return false
	}

	keyLen := len(stored_hash)
	if keyLen > math.MaxUint32 {
		keyLen = math.MaxUint32
	}

	hash := argon2.IDKey(
		[]byte(password),
		salt,
		uint32(argon2_time),   // #nosec G115 -- parsed with ParseUint bitSize 32, <= MaxUint32
		uint32(argon2_mem),    // #nosec G115 -- parsed with ParseUint bitSize 32, <= MaxUint32
		uint8(argon2_threads), // #nosec G115 -- parsed with ParseUint bitSize 32, <= MaxUint32
		uint32(keyLen),        // #nosec G115 -- keyLen bounded to MaxUint32 above
	)

	return subtle.ConstantTimeCompare(hash, stored_hash) == 1
}
