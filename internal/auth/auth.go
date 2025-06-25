package auth

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/user"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

type TokenData struct {
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
}

type TokenManager struct {
	tokenDir string
	filename string
	key      []byte
}

func NewTokenManager() *TokenManager {
	randomBytes := make([]byte, 8)
	rand.Read(randomBytes)

	currentUser, _ := user.Current()
	userInfo := currentUser.Username
	if userInfo == "" {
		userInfo = "passenger"
	}
	key := deriveKey(userInfo)

	hash := sha256.Sum256([]byte("passenger-go-cli-" + userInfo))
	filename := fmt.Sprintf("pass_%x_%x.tmp",
		hash[:4],
		randomBytes[:4])

	return &TokenManager{
		tokenDir: getSecureTokenDir(),
		filename: filename,
		key:      key,
	}
}

func getSecureTokenDir() string {
	if runtime.GOOS == "darwin" || runtime.GOOS == "linux" {
		if currentUser, err := user.Current(); err == nil {
			userTempDir := filepath.Join("/tmp", "passenger-"+currentUser.Username)
			if err := os.MkdirAll(userTempDir, 0700); err == nil {
				return userTempDir
			}
		}
	}
	return os.TempDir()
}

func deriveKey(userInfo string) []byte {
	hash := sha256.Sum256([]byte("passenger-secret-" + userInfo))
	return hash[:]
}

func (tokenManager *TokenManager) encrypt(data []byte) ([]byte, error) {
	block, err := aes.NewCipher(tokenManager.key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	ciphertext := gcm.Seal(nonce, nonce, data, nil)
	return ciphertext, nil
}

func (tokenManager *TokenManager) decrypt(data []byte) ([]byte, error) {
	block, err := aes.NewCipher(tokenManager.key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return nil, fmt.Errorf("ciphertext too short")
	}

	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}

func (tokenManager *TokenManager) getTokenPath() string {
	return filepath.Join(tokenManager.tokenDir, tokenManager.filename)
}

func secureWipe(data []byte) {
	for i := range data {
		data[i] = 0
	}
}

func (tokenManager *TokenManager) StoreToken(token string) error {
	now := time.Now()
	tokenData := TokenData{
		Token:     token,
		ExpiresAt: now.Add(5 * time.Minute),
		CreatedAt: now,
	}

	jsonData, err := json.Marshal(tokenData)
	if err != nil {
		return fmt.Errorf("failed to marshal token data: %w", err)
	}
	defer secureWipe(jsonData)

	encryptedData, err := tokenManager.encrypt(jsonData)
	if err != nil {
		return fmt.Errorf("failed to encrypt token data: %w", err)
	}

	encodedData := base64.StdEncoding.EncodeToString(encryptedData)

	tokenPath := tokenManager.getTokenPath()
	err = os.WriteFile(tokenPath, []byte(encodedData), 0600)
	if err != nil {
		return fmt.Errorf("failed to write token file: %w", err)
	}

	return nil
}

func (tokenManager *TokenManager) GetValidToken() (string, error) {
	tokenPath := tokenManager.getTokenPath()

	if _, err := os.Stat(tokenPath); os.IsNotExist(err) {
		return "", fmt.Errorf("no token found")
	}

	encodedData, err := os.ReadFile(tokenPath)
	if err != nil {
		return "", fmt.Errorf("failed to read token file: %w", err)
	}

	encryptedData, err := base64.StdEncoding.DecodeString(string(encodedData))
	if err != nil {
		return "", fmt.Errorf("failed to decode token data: %w", err)
	}

	jsonData, err := tokenManager.decrypt(encryptedData)
	if err != nil {
		return "", fmt.Errorf("failed to decrypt token data: %w", err)
	}
	defer secureWipe(jsonData)

	var tokenData TokenData
	err = json.Unmarshal(jsonData, &tokenData)
	if err != nil {
		return "", fmt.Errorf("failed to parse token data: %w", err)
	}

	if time.Now().After(tokenData.ExpiresAt) {
		tokenManager.ClearToken()
		return "", fmt.Errorf("token expired")
	}

	return tokenData.Token, nil
}

func (tokenManager *TokenManager) ClearToken() error {
	tokenPath := tokenManager.getTokenPath()
	if _, err := os.Stat(tokenPath); os.IsNotExist(err) {
		return nil
	}

	// Simple file shredding
	if file, err := os.OpenFile(tokenPath, os.O_WRONLY, 0); err == nil {
		stat, _ := file.Stat()
		randomBytes := make([]byte, stat.Size())
		rand.Read(randomBytes)
		file.WriteAt(randomBytes, 0)
		file.Sync()
		file.Close()
	}

	err := os.Remove(tokenPath)
	if err != nil {
		return fmt.Errorf("failed to remove token file: %w", err)
	}

	return nil
}

func (tokenManager *TokenManager) cleanupExpiredTokens() error {
	tokenDir := tokenManager.tokenDir

	files, err := os.ReadDir(tokenDir)
	if err != nil {
		return nil // If directory doesn't exist or can't be read, just continue
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		fileName := file.Name()
		// Only check files that match our token naming pattern
		if !strings.HasPrefix(fileName, "pass_") || !strings.HasSuffix(fileName, ".tmp") {
			continue
		}

		filePath := filepath.Join(tokenDir, fileName)

		// Try to read and decrypt the file to check if it's expired
		encodedData, err := os.ReadFile(filePath)
		if err != nil {
			continue // Skip files we can't read
		}

		encryptedData, err := base64.StdEncoding.DecodeString(string(encodedData))
		if err != nil {
			continue // Skip files we can't decode
		}

		jsonData, err := tokenManager.decrypt(encryptedData)
		if err != nil {
			continue // Skip files we can't decrypt
		}

		var tokenData TokenData
		err = json.Unmarshal(jsonData, &tokenData)
		if err != nil {
			continue // Skip files we can't parse
		}

		// Secure wipe of sensitive data
		secureWipe(jsonData)

		// If token is expired, remove the file
		if time.Now().After(tokenData.ExpiresAt) {
			// Securely overwrite the file before removing
			if f, err := os.OpenFile(filePath, os.O_WRONLY, 0); err == nil {
				stat, _ := f.Stat()
				zeros := make([]byte, stat.Size())
				f.WriteAt(zeros, 0)
				f.Sync()
				f.Close()
			}
			os.Remove(filePath)
		}
	}

	return nil
}

var defaultTokenManager = NewTokenManager()

func StoreToken(token string) error {
	// Do not check errors here, it's not critical
	defaultTokenManager.cleanupExpiredTokens()

	return defaultTokenManager.StoreToken(token)
}

func GetValidToken() (string, error) {
	return defaultTokenManager.GetValidToken()
}

func ClearToken() error {
	return defaultTokenManager.ClearToken()
}
