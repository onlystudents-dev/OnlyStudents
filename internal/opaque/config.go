package opaquepkg

import (
	"encoding/hex"
	"fmt"

	"onlystudents/internal/helpers"

	"github.com/bytemare/ecc"
	"github.com/bytemare/opaque"
)

var Conf *opaque.Configuration

func CreateServerFromEnv() (*opaque.Server, error) {
	conf := opaque.DefaultConfiguration()
	Conf = conf
	group := ecc.Ristretto255Sha512

	seed, err := mustHexEnv("OPAQUE_OPRF_SEED_SECRET", 64)
	if err != nil {
		return nil, err
	}
	privBytes, err := mustHexEnv("OPAQUE_SERVER_PRIVATE_KEY", 32)
	if err != nil {
		return nil, err
	}

	priv := group.NewScalar()
	if err := priv.Decode(privBytes); err != nil {
		return nil, fmt.Errorf("OPAQUE_SERVER_PRIVATE_KEY: invalid ristretto255 scalar: %w", err)
	}

	pub := group.NewElement().Base().Multiply(priv)

	server, err := opaque.NewServer(conf)
	if err != nil {
		return nil, err
	}

	if err := server.SetKeyMaterial(&opaque.ServerKeyMaterial{
		PrivateKey:     priv,
		PublicKeyBytes: pub.Encode(),
		OPRFGlobalSeed: seed,
	}); err != nil {
		return nil, fmt.Errorf("set OPAQUE key material: %w", err)
	}

	return server, nil
}

func mustHexEnv(env string, want int) ([]byte, error) {
	raw := helpers.GetEnvFallback(env, "")
	if raw == "" {
		return nil, fmt.Errorf("%s is not set, generate it once and store it in sops/.env", env)
	}
	b, err := hex.DecodeString(raw)
	if err != nil {
		return nil, fmt.Errorf("%s: not valid hex: %w", env, err)
	}
	if len(b) != want {
		return nil, fmt.Errorf("%s: want %d bytes, got %d", env, want, len(b))
	}
	return b, nil
}
