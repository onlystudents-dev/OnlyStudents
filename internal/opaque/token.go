package opaquepkg

import (
	"encoding/hex"
	"onlystudents/internal/helpers"

	"github.com/bytemare/opaque"
)

func NewEnrollToken() string {
	return hex.EncodeToString(opaque.RandomBytes(helpers.GetIntEnvFallback("OPAQUE_ENROLL_TOKEN_LEN", 32, 512)))
}

func ValidEnrollToken(t string) bool {
	return len(t) == helpers.GetIntEnvFallback("OPAQUE_ENROLL_TOKEN_LEN", 32, 512)*2
}
