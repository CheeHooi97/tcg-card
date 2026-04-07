package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	r "math/rand"
	config "pkm/config"
	"pkm/errcode"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

const (
	minSeed = int64(100000000)
	maxSeed = int64(999999999)
)

var numberRegex = regexp.MustCompile(`\b\d+(?:/\d+)?\b`)

func ValidateRequest(c echo.Context, i any) (string, error) {
	if err := c.Bind(&i); err != nil {
		return errcode.InvalidRequest.Message, err
	}

	if err := c.Validate(i); err != nil {
		return errcode.ValidationError.Message, err
	}

	return "", nil
}

func EncryptAES(plaintext string) (string, error) {
	keyBytes, err := base64.StdEncoding.DecodeString(config.SystemAesKey)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(keyBytes)
	if err != nil {
		return "", err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, aesGCM.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	ciphertext := aesGCM.Seal(nil, nonce, []byte(plaintext), nil)

	encrypted := append(nonce, ciphertext...)

	return base64.StdEncoding.EncodeToString(encrypted), nil
}

func DecryptAES(encryptedText string) (string, error) {
	keyBytes, err := base64.StdEncoding.DecodeString(config.SystemAesKey)
	if err != nil {
		return "", err
	}

	data, err := base64.StdEncoding.DecodeString(encryptedText)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(keyBytes)
	if err != nil {
		return "", err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonceSize := aesGCM.NonceSize()
	if len(data) < nonceSize {
		return "", errors.New(errcode.InvalidEncryptedText.Message)
	}

	nonce, ciphertext := data[:nonceSize], data[nonceSize:]

	plaintext, err := aesGCM.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}

// random
var randomSeed = r.New(r.NewSource(time.Now().UnixNano()))

func Alphanumeric(n int) string {
	var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890-")
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[randomSeed.Intn(len(letterRunes))]
	}
	return string(b)
}

func Numeric(n int) string {
	var letterRunes = []rune("1234567890")
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[randomSeed.Intn(len(letterRunes))]
	}

	s := string(b)
	return s
}

// GenerateUNIQ :
func UniqueID() string {
	strID := strconv.FormatInt(time.Now().Unix(), 10) + strconv.FormatInt(r.Int63n(maxSeed-minSeed)+minSeed, 10)
	return strID
}

func CapitalizeFirst(s string) string {
	if s == "" {
		return s
	}
	c := cases.Title(language.Und)
	return c.String(s[:1]) + s[1:]
}

func AddHashToNumbers(input string) string {
	return numberRegex.ReplaceAllStringFunc(input, func(match string) string {
		if strings.HasPrefix(match, "#") {
			return match // already has #
		}
		return "#" + match
	})
}

func SplitNameAndRarity(input string, rarities []string) (name string, rarity string) {
	for _, r := range rarities {
		if strings.HasSuffix(input, r) {
			name = strings.TrimSpace(strings.TrimSuffix(input, r))
			rarity = r
			return
		}
	}
	return input, "" // no rarity found
}

func SplitCardName(s string) (part1, part2, part3 string) {
	s = strings.TrimSpace(s)

	// Case 1: has [] and #
	if strings.Contains(s, "[") && strings.Contains(s, "]") && strings.Contains(s, "#") {
		// split name + rest
		beforeBracket, afterBracket, _ := strings.Cut(s, "[")

		part1 = strings.TrimSpace(beforeBracket)

		// get content inside []
		bracketContent, after, _ := strings.Cut(afterBracket, "]")
		part2 = strings.TrimSpace(bracketContent)

		// get number after #
		if _, num, ok := strings.Cut(after, "#"); ok {
			part3 = strings.TrimSpace(num)
		}

		return
	}

	// Case 2: only has #
	if strings.Contains(s, "#") {
		name, num, _ := strings.Cut(s, "#")
		part1 = strings.TrimSpace(name)
		part2 = strings.TrimSpace(num)
		return
	}

	// Fallback: only name
	part1 = s
	return
}

func UniqueUint(input []uint) []uint {
	seen := make(map[uint]struct{})
	result := make([]uint, 0, len(input))

	for _, v := range input {
		if _, ok := seen[v]; ok {
			continue
		}
		seen[v] = struct{}{}
		result = append(result, v)
	}
	return result
}

func OrderByField(field string, ids []string) string {
	values := make([]string, len(ids))
	for i, id := range ids {
		values[i] = id
	}
	return fmt.Sprintf("FIELD(%s, %s)", field, strings.Join(values, ","))
}

func ParsePrice(price string) float64 {
	price = strings.TrimSpace(price)

	// Remove currency symbols and commas
	price = strings.ReplaceAll(price, "$", "")
	price = strings.ReplaceAll(price, ",", "")

	val, err := strconv.ParseFloat(price, 64)
	if err != nil {
		return 0
	}

	return val
}
