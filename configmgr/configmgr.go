package configmgr

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

type TokenAndRefreshToken struct {
	IDToken      string `json:"idToken"`
	RefreshToken string `json:"refreshToken"`
}

// EcomConfigEntry represents a single configuration set.
type EcomConfigEntry struct {
	FirebaseAPIKey string `mapstructure:"firebase-api-key" yaml:"firebase-api-key"`
	DevKey         string `mapstructure:"developer-key" yaml:"developer-key"`
	Endpoint       string `mapstructure:"endpoint" yaml:"endpoint"`
}

// EcomConfigurations contains the map of config entries.
type EcomConfigurations struct {
	Configurations map[string]EcomConfigEntry `mapstructure:"configurations" yaml:"configurations"`
}

const (
	configFile = ".ecomrc.yaml"
	configDir  = ".ecom"
)

func homeDir() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", fmt.Errorf("user.Current() failed: %v", err)
	}
	return usr.HomeDir, nil
}

func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}

func ensureConfigDirExists() error {
	hd, err := homeDir()
	if err != nil {
		return errors.Wrapf(err, "failed homeDir()")
	}
	cfgDir := filepath.Join(hd, configDir)

	exists, err := exists(cfgDir)
	if err != nil {
		return errors.Wrapf(err, "failed exists(%s)", configDir)
	}
	if !exists {
		os.Mkdir(cfgDir, 0755)
		err = WriteCurrentProject("")
		if err != nil {
			return errors.Wrapf(err, "failed write current project %q", "")
		}
	}
	return nil
}

func ensureConfigFileExists() error {
	hd, err := homeDir()
	if err != nil {
		return errors.Wrapf(err, "failed homeDir()")
	}

	cf := filepath.Join(hd, configFile)
	exists, err := exists(cf)
	if err != nil {
		return errors.Wrapf(err, "failed exists(%s)", configDir)
	}

	if !exists {
		f, err := os.Create(cf)
		if err != nil {
			return errors.Wrapf(err, "create(%q) failed", cf)
		}
		defer f.Close()
		_, err = f.WriteString("{}")
		if err != nil {
			return errors.Wrapf(err, "write string to file %q failed", cf)
		}
	}
	return nil
}

// URLToHostName converts a standard URL string to hostname replacing the dot character with underscores.
func URLToHostName(u string) (string, error) {
	url, err := url.Parse(u)
	if err != nil {
		return "", errors.Wrapf(err, "failed to parse url %q", u)
	}
	return strings.ReplaceAll(url.Hostname(), ".", "_"), nil
}

// TokenFilename returns the full filepath of the file corresponding
// to the given EcomConfigEntry.
func TokenFilename(e *EcomConfigEntry) (string, error) {
	hd, err := homeDir()
	if err != nil {
		return "", errors.Wrapf(err, "homeDir() failed")
	}

	hostname, err := URLToHostName(e.Endpoint)
	if err != nil {
		return "", errors.Wrapf(err, "url to hostname failed for %q", e.Endpoint)
	}
	file := fmt.Sprintf("%s-%s", e.FirebaseAPIKey, hostname)
	tokenFile := filepath.Join(hd, configDir, file)
	exists, err := exists(tokenFile)
	if err != nil {
		return "", errors.Wrapf(err, "exists(%s) failed", tokenFile)
	}

	if !exists {
		return "", fmt.Errorf("token file %q not found", tokenFile)
	}

	return tokenFile, nil
}

// ReadCurrentConfigName returns the contents of the CURRENT_PROJECT
// file. If the CURRENT_PROJECT file does not exists (for example, the
// first time the program is run), an empty file will be created.
func ReadCurrentConfigName() (string, error) {
	hd, err := homeDir()
	if err != nil {
		return "", errors.Wrapf(err, "homeDir() failed")
	}
	err = ensureConfigDirExists()
	if err != nil {
		return "", errors.Wrap(err, "ensure config dir exists failed")
	}
	cpf := filepath.Join(hd, configDir, "CURRENT_PROJECT")
	exists, err := exists(cpf)
	if err != nil {
		return "", errors.Wrapf(err, "exists(%s) failed", cpf)
	}
	if !exists {
		f, err := os.Create(cpf)
		if err != nil {
			return "", errors.Wrapf(err, "create file %q failed", cpf)
		}
		defer f.Close()
		return "", nil
	}
	bs, err := ioutil.ReadFile(cpf)
	if err != nil {
		return "", errors.Wrapf(err, "read file %q failed", cpf)
	}
	return string(bs), nil
}

// ReadConfig opens and read the .ecomrc.yaml file putting each section
// name in a map of EcomConfigEntrys.
func ReadConfig() (*EcomConfigurations, error) {
	err := ensureConfigFileExists()
	if err != nil {
		return nil, errors.Wrap(err, "ensure config file exists failed")
	}
	viper.SetConfigName(".ecomrc")
	viper.SetConfigType("yaml")
	hd, err := homeDir()
	if err != nil {
		return nil, err
	}
	viper.AddConfigPath(hd)
	viper.ReadInConfig()
	if err != nil {
		return nil, errors.Wrap(err, "read in config file failed")
	}
	configurations := EcomConfigurations{}
	err = viper.Unmarshal(&configurations)
	if err != nil {
		return nil, errors.Wrap(err, "unmarshal configurations failed")
	}
	return &configurations, nil
}

// WriteConfig writes the EcomConfigurations to the YAML config.
func WriteConfig(cfgs *EcomConfigurations) error {
	viper.Set("configurations", cfgs.Configurations)
	err := viper.WriteConfig()
	if err != nil {
		return errors.Wrap(err, "write config file failed")
	}
	return nil
}

// WriteCurrentProject records the project API Key on the filesystem within the $HOME/.ecom directory in a file called CURRENT_API_KEY. The current API Key context is read between invocation of the command-line tool.
func WriteCurrentProject(name string) error {
	err := ensureConfigDirExists()
	if err != nil {
		return errors.Wrap(err, "couldn't ensure config dir exists")
	}
	hd, err := homeDir()
	if err != nil {
		return errors.Wrap(err, "failed to get home directory")
	}
	cpf := filepath.Join(hd, configDir, "CURRENT_PROJECT")
	bs := []byte(name)
	err = ioutil.WriteFile(cpf, bs, 0644)
	if err != nil {
		return errors.Wrap(err, "write CURRENT_PROJECT file failed")
	}
	return nil
}

// DeleteProject removes a project from the .ecom directory returning ok true if successful.
func DeleteProject(filename string) (bool, error) {
	hd, err := homeDir()
	if err != nil {
		return false, errors.Wrap(err, "failed to get home directory")
	}
	filepath := filepath.Join(hd, configDir, filename)

	exists, err := exists(filepath)
	if err != nil {
		return false, errors.Wrapf(err, "failed exists(%q)", filepath)
	}
	if !exists {
		return false, fmt.Errorf("filename %q not found", filepath)
	}
	err = os.Remove(filepath)
	if err != nil {
		return false, errors.Wrapf(err, "failed to remove file %q", filepath)
	}
	return true, nil
}

// ReadTokenAndRefreshToken reads the token and refresh token from the filesystem
// or returns nil if the file has not yet been created.
func ReadTokenAndRefreshToken(fp string) (*TokenAndRefreshToken, error) {
	exists, err := exists(fp)
	if err != nil {
		return nil, errors.Wrapf(err, "exists(%s) failed", fp)
	}

	if !exists {
		return nil, errors.Wrapf(err, "token file %q not found", fp)
	}

	f, err := os.Open(fp)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to open file %q", fp)
	}
	defer f.Close()

	var tar TokenAndRefreshToken
	json.NewDecoder(f).Decode(&tar)
	return &tar, nil
}

// WriteTokenAndRefreshToken writes a copy of the token and refresh token to the file system
func WriteTokenAndRefreshToken(webKey, endpoint string, tar *TokenAndRefreshToken) error {
	err := ensureConfigDirExists()
	if err != nil {
		return errors.Wrap(err, "couldn't ensure config dir exists")
	}

	hostname, err := URLToHostName(endpoint)
	filename := fmt.Sprintf("%s-%s", webKey, hostname)

	hd, err := homeDir()
	if err != nil {
		return errors.Wrap(err, "failed to get home directory")
	}
	filepath := filepath.Join(hd, configDir, filename)

	f, err := os.Create(filepath)
	if err != nil {
		return errors.Wrapf(err, "create file %q failed", filepath)
	}
	defer f.Close()

	err = json.NewEncoder(f).Encode(tar)
	if err != nil {
		return errors.Wrapf(err, "json encode token failed")
	}
	return nil
}
