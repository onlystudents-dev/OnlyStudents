package opaquepkg

import (
	"encoding/base64"
	"errors"

	"github.com/goccy/go-json"
)

type Message []byte

func (m Message) MarshalJSON() ([]byte, error) {
	return json.Marshal(base64.RawURLEncoding.EncodeToString(m))
}

func (m *Message) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	b, err := base64.RawURLEncoding.DecodeString(s)
	if err != nil {
		return errors.Join(errors.New("invalid opaque message encoding"), err)
	}

	*m = b

	return nil
}
