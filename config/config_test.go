package config

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
	"time"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testConfig struct {
	Host    string        `mapstructure:"host"    env:"TEST_HOST"    validate:"required"`
	Port    int           `mapstructure:"port"    env:"TEST_PORT"    validate:"required,gte=1"`
	Timeout time.Duration `mapstructure:"timeout" env:"TEST_TIMEOUT"`
	Nested  nestedConfig  `mapstructure:"nested"`
	Ignored string        `mapstructure:"-"`
	Pointer *string       `mapstructure:"pointer" env:"TEST_POINTER"`
}

type nestedConfig struct {
	Name string `mapstructure:"name" env:"TEST_NESTED_NAME" validate:"required"`
}

func TestLoadFromFileAndEnv(t *testing.T) {
	t.Setenv("TEST_PORT", "9090")
	t.Setenv("TEST_NESTED_NAME", "env-name")

	path := filepath.Join(t.TempDir(), "config.yml")
	if err := os.WriteFile(path, []byte(`
host: localhost
port: 8080
timeout: 5s
nested:
  name: file-name
`), 0o600); err != nil {
		require.NoError(t, err)
	}

	cfg, err := Load[testConfig](path)
	require.NoError(t, err)
	assert.Equal(t, "localhost", cfg.Host)
	assert.Equal(t, 9090, cfg.Port)
	assert.Equal(t, 5*time.Second, cfg.Timeout)
	assert.Equal(t, "env-name", cfg.Nested.Name)
}

func TestLoadReadFileError(t *testing.T) {
	_, err := Load[testConfig](filepath.Join(t.TempDir(), "missing.yml"))
	require.Error(t, err)
	assert.ErrorContains(t, err, "config: read file")
}

func TestLoadValidationError(t *testing.T) {
	_, err := Load[testConfig]("")
	require.Error(t, err)
	assert.ErrorContains(t, err, "config:")
}

func TestLoadNonStruct(t *testing.T) {
	_, err := Load[int]("")
	require.Error(t, err)
	assert.ErrorContains(t, err, "config: unmarshal")
}

func TestLoadEnvCoercionFailure(t *testing.T) {
	t.Setenv("TEST_PORT", "not-a-number")
	t.Setenv("TEST_HOST", "localhost")
	t.Setenv("TEST_NESTED_NAME", "n")

	_, err := Load[testConfig]("")
	require.Error(t, err)
	assert.ErrorContains(t, err, "config: unmarshal")
}

func TestLoadMissingOptionalFieldsZeroValue(t *testing.T) {
	t.Setenv("TEST_HOST", "localhost")
	t.Setenv("TEST_PORT", "8080")
	t.Setenv("TEST_NESTED_NAME", "n")

	cfg, err := Load[testConfig]("")
	require.NoError(t, err)
	assert.Equal(t, time.Duration(0), cfg.Timeout)
	assert.Nil(t, cfg.Pointer)
	assert.Equal(t, "", cfg.Ignored)
}

func TestBindEnvsPointerType(t *testing.T) {
	t.Setenv("TEST_HOST", "localhost")
	t.Setenv("TEST_PORT", "8080")
	t.Setenv("TEST_NESTED_NAME", "pointer-name")

	v := viper.New()
	bindEnvs(v, reflect.TypeFor[*testConfig](), "")

	var cfg testConfig
	require.NoError(t, v.Unmarshal(&cfg))

	assert.Equal(t, "localhost", cfg.Host)
	assert.Equal(t, 8080, cfg.Port)
	assert.Equal(t, "pointer-name", cfg.Nested.Name)
}
