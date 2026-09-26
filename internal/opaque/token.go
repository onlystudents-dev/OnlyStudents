package opaquepkg

import (
	"encoding/hex"
	"onlystudents/internal/helpers"

	"github.com/bytemare/opaque"
)

const enrollTokenEnv = "OPAQUE_ENROLL_TOKEN_LEN"

func NewEnrollToken() string {
	return hex.EncodeToString(opaque.RandomBytes(helpers.GetIntEnvFallback(enrollTokenEnv, 32, 512)))
}

func ValidEnrollToken(t string) bool {
	return len(t) == helpers.GetIntEnvFallback(enrollTokenEnv, 32, 512)*2
}
