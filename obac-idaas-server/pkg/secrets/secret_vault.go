package secrets

import (
	"github.com/magiconair/properties"
	"path/filepath"
)

var (
	defaultSecretsFile = "secrets.properties"
)

type SecretVault interface {
	GetSecret(secretKey string) (string, bool)
}

type EnvironmentSecretsVault struct {
	secrets *properties.Properties
}

func NewSecretVault(path string) (SecretVault, error) {
	vault := EnvironmentSecretsVault{}
	vault.secrets = properties.MustLoadFile(filepath.Join(path, defaultSecretsFile), properties.UTF8)
	return &vault, nil
}

func (e *EnvironmentSecretsVault) GetSecret(secretKey string) (string, bool) {
	return e.secrets.Get(secretKey)
}
