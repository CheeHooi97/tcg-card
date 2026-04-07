package config

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"log"
	"os"

	"github.com/joho/godotenv"
)

var (
	Env                      string
	DBHost                   string
	DBPort                   string
	DBUser                   string
	DBPassword               string
	DBName                   string
	SystemAesKey             string
	PokemonApiKey            string
	OSSEndpoint              string
	OSSAccessKeyID           string
	OSSAccessKeySecret       string
	OSSBucket                string
	AuthenticationPublicKey  *rsa.PublicKey
	AuthenticationPrivateKey *rsa.PrivateKey
)

// LoadConfig
func LoadConfig() {
	_ = godotenv.Load()

	Env = GetEnv("ENV")
	DBHost = GetEnv("MYSQL_HOST")
	DBPort = GetEnv("MYSQL_PORT")
	DBUser = GetEnv("MYSQL_USER")
	DBPassword = GetEnv("MYSQL_PASSWORD")
	DBName = GetEnv("MYSQL_DATABASE")
	SystemAesKey = GetEnv("SYSTEM_AES_KEY")
	PokemonApiKey = GetEnv("POKEMON_API_KEY") // Optional, so use os.Getenv directly
	OSSEndpoint = GetEnv("OSS_ENDPOINT")
	OSSAccessKeyID = GetEnv("OSS_ACCESS_KEY_ID")
	OSSAccessKeySecret = GetEnv("OSS_ACCESS_KEY_SECRET")
	OSSBucket = GetEnv("OSS_BUCKET")
	AuthPrivateKeyPath := GetEnv("AUTH_PRIVATE_KEY_PATH")
	AuthPublicKeyPath := GetEnv("AUTH_PUBLIC_KEY_PATH")

	err := loadRSAKeys(AuthPrivateKeyPath, AuthPublicKeyPath)
	if err != nil {
		log.Fatalf("Failed to load RSA keys: %v", err)
	}

}

func loadRSAKeys(privatePath, publicPath string) error {
	// Load private key
	privData, err := os.ReadFile(privatePath)
	if err != nil {
		return err
	}
	privBlock, _ := pem.Decode(privData)
	if privBlock == nil || privBlock.Type != "RSA PRIVATE KEY" {
		return errors.New("invalid private key PEM")
	}
	privateKey, err := x509.ParsePKCS1PrivateKey(privBlock.Bytes)
	if err != nil {
		return err
	}
	AuthenticationPrivateKey = privateKey

	// Load public key
	pubData, err := os.ReadFile(publicPath)
	if err != nil {
		return err
	}
	pubBlock, _ := pem.Decode(pubData)
	if pubBlock == nil || pubBlock.Type != "PUBLIC KEY" {
		return errors.New("invalid public key PEM")
	}
	pubInterface, err := x509.ParsePKIXPublicKey(pubBlock.Bytes)
	if err != nil {
		return err
	}
	publicKey, ok := pubInterface.(*rsa.PublicKey)
	if !ok {
		return errors.New("not an RSA public key")
	}
	AuthenticationPublicKey = publicKey

	return nil
}

func GetEnv(key string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		log.Fatalf("%s environment variable not set", key)
	}
	return value
}
