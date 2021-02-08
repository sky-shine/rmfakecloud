package config

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"os"
	"path"
	"path/filepath"

	log "github.com/sirupsen/logrus"
)

const (
	DefaultPort     = "3000"
	DefaultDataDir  = "data"
	defaultTrashDir = "trash"

	// DefaultHost fake url
	DefaultHost = "local.appspot.com"

	EnvLogLevel   = "LOGLEVEL"
	EnvDataDir    = "DATADIR"
	EnvPort       = "PORT"
	EnvStorageURL = "STORAGE_URL"

	// smtp
	envJWTSecretKey     = "JWT_SECRET_KEY"
	envRegistrationOpen = "OPEN_REGISTRATION"
	EnvSmtpServer       = "RM_SMTP_SERVER"
	EnvSmtpUsername     = "RM_SMTP_USERNAME"
	EnvSmtpPassword     = "RM_SMTP_PASSWORD"
	EnvSmtpHelo         = "RM_SMTP_HELO"
	EnvSmtpInsecureTLS  = "RM_SMTP_INSECURE_TLS"
	EnvSmtpFrom         = "RM_SMTP_FROM"

	// myScript hwr api keys
	EnvHwrApplicationKey = "RMAPI_HWR_APPLICATIONKEY"
	EnvHwrHmac           = "RMAPI_HWR_HMAC"
)

// Config config
type Config struct {
	Port             string
	StorageURL       string
	DataDir          string
	TrashDir         string
	JWTSecretKey     []byte
	RegistrationOpen bool
}

// FromEnv config from environment values
func FromEnv() *Config {
	var err error
	var dataDir string
	data := os.Getenv(EnvDataDir)
	if data != "" {
		dataDir = data
	} else {
		dataDir, err = filepath.Abs(DefaultDataDir)
		if err != nil {
			panic(err)
		}
	}
	trashDir := path.Join(dataDir, defaultTrashDir)
	err = os.MkdirAll(trashDir, 0700)
	if err != nil {
		panic(err)
	}

	port := os.Getenv(EnvPort)
	if port == "" {
		port = DefaultPort
	}

	uploadURL := os.Getenv(EnvStorageURL)
	if uploadURL == "" {
		host, err := os.Hostname()
		if err != nil {
			log.Warn("cannot get hostname")
			host = DefaultHost
		}
		uploadURL = fmt.Sprintf("http://%s:%s", host, port)
	}

	jwtSecretKey, err := hex.DecodeString(os.Getenv(envJWTSecretKey))
	if err != nil || len(jwtSecretKey) == 0 {
		jwtSecretKey = make([]byte, 32)
		_, err := rand.Read(jwtSecretKey)
		if err != nil {
			panic(err)
		}
		log.Warnf("You have to set %s with some content. Eg: %s='%X'", envJWTSecretKey, envJWTSecretKey, jwtSecretKey)
		log.Warn("  without this variable set, you'll be disconnected after this program restart")
	}

	//openRegistration := os.Getenv(envRegistrationOpen)

	cfg := Config{
		Port:             port,
		StorageURL:       uploadURL,
		DataDir:          dataDir,
		TrashDir:         trashDir,
		JWTSecretKey:     jwtSecretKey,
		RegistrationOpen: true,
	}
	return &cfg
}
