package helpers

import (
	// secure because over HTTPS and we are only sending the start of the hash
	"crypto/sha1" // nosec G401
	"encoding/hex"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/gofiber/fiber/v3"
)

func HaveIBeenPwnedCheck(pw string) (bool, error) {
	sum := sha1.Sum([]byte(pw))
	h := strings.ToUpper(hex.EncodeToString(sum[:]))

	resp, err := http.Get("https://api.pwnedpasswords.com/range/" + h[:5])
	if err != nil {
		return true, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return true, fmt.Errorf("pwnedpasswords: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return true, err
	}

	for line := range strings.FieldsSeq(string(body)) {
		if suffix, _, ok := strings.Cut(line, ":"); ok && suffix == h[5:] {
			return true, nil
		}
	}

	return false, nil
}

func PasswordChecks(c fiber.Ctx, password string, confirm_password string) (bool, error) {
	if password != confirm_password {
		return false, c.SendStatus(401)
	}

	if utf8.RuneCountInString(password) < GetIntEnvFallback("PASSWORD_MIN_LEN", 12, 128) {
		return false, c.Status(400).SendString(fmt.Sprintf("Your password is too weak: password must be atleast %d characters", GetIntEnvFallback("PASSWORD_MIN_LEN", 12, 128)))
	}

	if utf8.RuneCountInString(password) > GetIntEnvFallback("PASSWORD_MAX_LEN", 128, 256) {
		return false, c.Status(400).SendString(fmt.Sprintf("Your password is too weak: password must not exceed %d characters", GetIntEnvFallback("PASSWORD_MAX_LEN", 128, 256)))
	}

	var hasUpper, hasLower, hasDigit, hasSymbol bool
	for _, r := range password {
		switch {
		case unicode.IsUpper(r):
			hasUpper = true
		case unicode.IsLower(r):
			hasLower = true
		case unicode.IsDigit(r):
			hasDigit = true
		case unicode.IsPunct(r) || unicode.IsSymbol(r):
			hasSymbol = true
		}
	}

	if !hasUpper {
		return false, c.Status(400).SendString("Your password is too weak: password must contain an uppercase letter")
	}
	if !hasLower {
		return false, c.Status(400).SendString("Your password is too weak: password must contain a lowercase letter")
	}
	if !hasDigit {
		return false, c.Status(400).SendString("Your password is too weak: password must contain a digit")
	}
	if !hasSymbol {
		return false, c.Status(400).SendString("Your password is too weak: password must contain a symbol")
	}

	breached, err := HaveIBeenPwnedCheck(password)

	if err != nil {
		slog.Error("have i been pwned error", "err", err)
		return false, c.SendStatus(500)
	}

	if breached {
		return false, c.Status(400).SendString("Your password has previously been exposed in a data breach!")
	}

	return true, nil
}
