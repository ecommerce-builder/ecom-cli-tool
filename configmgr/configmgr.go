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

	"github.com/spf13/viper"
)

// TokenAndRefreshToken contains a pair of JTW and refresh token for Firebase.
type TokenAndRefreshToken struct {
	IDToken      string `json:"idToken"`
	RefreshToken string `json:"refreshToken"`
}

// Customer details
type Customer struct {
	ID        string `mapstructure:"id" yaml:"id"`
	UID       string `mapstructure:"uid" yaml:"uid"`
	Role      string `mapstructure:"role" yaml:"role"`
	Email     string `mapstructure:"email" yaml:"email"`
	Firstname string `mapstructure:"firstname" yaml:"firstname"`
	Lastname  string `mapstructure:"lastname" yaml:"lastname"`
}

// EcomConfigEntry represents a single configuration set.
type EcomConfigEntry struct {
	Endpoint string   `mapstructure:"endpoint" yaml:"endpoint"`
	DevKey   string   `mapstructure:"developer-key" yaml:"developer-key"`
	Customer Customer `mapstructure:"user" yaml:"user"`
}

// EcomConfigurations contains the map of config entries.
type EcomConfigurations struct {
	Configurations map[string]EcomConfigEntry `mapstructure:"configurations" yaml:"configurations"`
}

// func (e EcomConfigurations) CurrentConfigEntry() EcomConfigEntry {
// 	for _, c := range e.Configurations {
// 		c.Endpoint
// 	}
// }

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
		return fmt.Errorf("failed homeDir(): %w", err)
	}
	cfgDir := filepath.Join(hd, configDir)

	exists, err := exists(cfgDir)
	if err != nil {
		return fmt.Errorf("failed exists(%s): %w", configDir, err)
	}
	if !exists {
		os.Mkdir(cfgDir, 0755)
		err = WriteCurrentProject("")
		if err != nil {
			return fmt.Errorf("failed write current project %q: %w", "", err)
		}
	}
	return nil
}

func ensureConfigFileExists() error {
	hd, err := homeDir()
	if err != nil {
		return fmt.Errorf("failed homeDir(): %w", err)
	}

	cf := filepath.Join(hd, configFile)
	exists, err := exists(cf)
	if err != nil {
		return fmt.Errorf("failed exists(%s): %w", configDir, err)
	}

	if !exists {
		f, err := os.Create(cf)
		if err != nil {
			return fmt.Errorf("create(%q) failed: %w", cf, err)
		}
		defer f.Close()
		_, err = f.WriteString("{}")
		if err != nil {
			return fmt.Errorf("write string to file %q failed: %w", cf, err)
		}
	}
	return nil
}

// URLToHostName converts a standard URL string to hostname replacing the dot character with underscores.
func URLToHostName(u string) (string, error) {
	url, err := url.Parse(u)
	if err != nil {
		return "", fmt.Errorf("failed to parse url %q: %w", u, err)
	}
	return strings.ReplaceAll(url.Hostname(), ".", "_"), nil
}

// TokenFilename returns the full filepath of the file corresponding
// to the given EcomConfigEntry.
func TokenFilename(e *EcomConfigEntry) (string, error) {
	hd, err := homeDir()
	if err != nil {
		return "", fmt.Errorf("homeDir() failed: %w", err)
	}

	hostname, err := URLToHostName(e.Endpoint)
	if err != nil {
		return "", fmt.Errorf("url to hostname failed for %q: %w", e.Endpoint, err)
	}
	filename := fmt.Sprintf("%s-%s", hostname, e.DevKey[:6])
	tokenFile := filepath.Join(hd, configDir, filename)
	exists, err := exists(tokenFile)
	if err != nil {
		return "", fmt.Errorf("exists(%s) failed: %w", tokenFile, err)
	}
	if !exists {
		return "", fmt.Errorf("token file %q not found: %w", tokenFile, err)
	}
	return tokenFile, nil
}

// ReadCurrentConfigName returns the contents of the CURRENT_PROJECT
// file. If the CURRENT_PROJECT file does not exists (for example, the
// first time the program is run), an empty file will be created.
func ReadCurrentConfigName() (string, error) {
	hd, err := homeDir()
	if err != nil {
		return "", fmt.Errorf("homeDir() failed: %w", err)
	}
	err = ensureConfigDirExists()
	if err != nil {
		return "", fmt.Errorf("ensure config dir exists failed: %w", err)
	}
	cpf := filepath.Join(hd, configDir, "CURRENT_PROJECT")
	exists, err := exists(cpf)
	if err != nil {
		return "", fmt.Errorf("exists(%s) failed: %w", cpf, err)
	}
	if !exists {
		f, err := os.Create(cpf)
		if err != nil {
			return "", fmt.Errorf("create file %q failed: %w", cpf, err)
		}
		defer f.Close()
		return "", nil
	}
	bs, err := ioutil.ReadFile(cpf)
	if err != nil {
		return "", fmt.Errorf("read file %q failed: %w", cpf, err)
	}
	return string(bs), nil
}

// ReadConfig opens and read the .ecomrc.yaml file putting each section
// name in a map of EcomConfigEntrys.
func ReadConfig() (*EcomConfigurations, error) {
	err := ensureConfigFileExists()
	if err != nil {
		return nil, fmt.Errorf("ensure config file exists failed: %w", err)
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
		return nil, fmt.Errorf("read in config file failed: %w", err)
	}
	configurations := EcomConfigurations{}
	err = viper.Unmarshal(&configurations)
	if err != nil {
		return nil, fmt.Errorf("unmarshal configurations failed: %w", err)
	}
	return &configurations, nil
}

// WriteConfig writes the EcomConfigurations to the YAML config.
func WriteConfig(cfgs *EcomConfigurations) error {
	viper.Set("configurations", cfgs.Configurations)
	err := viper.WriteConfig()
	if err != nil {
		return fmt.Errorf("write config file failed: %w", err)
	}
	return nil
}

// WriteCurrentProject records the project API Key on the filesystem within the $HOME/.ecom
// directory in a file called CURRENT_API_KEY. The current API Key context is read
// between invocation of the command-line tool.
func WriteCurrentProject(name string) error {
	err := ensureConfigDirExists()
	if err != nil {
		return fmt.Errorf("couldn't ensure config dir exists: %w", err)
	}
	hd, err := homeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}
	cpf := filepath.Join(hd, configDir, "CURRENT_PROJECT")
	bs := []byte(name)
	err = ioutil.WriteFile(cpf, bs, 0644)
	if err != nil {
		return fmt.Errorf("write CURRENT_PROJECT file failed: %w", err)
	}
	return nil
}

// DeleteProject removes a project from the .ecom directory returning ok true if successful.
func DeleteProject(filename string) (bool, error) {
	hd, err := homeDir()
	if err != nil {
		return false, fmt.Errorf("failed to get home directory: %w", err)
	}
	filepath := filepath.Join(hd, configDir, filename)

	exists, err := exists(filepath)
	if err != nil {
		return false, fmt.Errorf("exists(path=%q) failed: %w", filepath, err)
	}
	if !exists {
		return false, fmt.Errorf("filename %q not found: %w", filepath, err)
	}
	err = os.Remove(filepath)
	if err != nil {
		return false, fmt.Errorf("failed to remove file %q: %w", filepath, err)
	}
	return true, nil
}

// ReadTokenAndRefreshToken reads the token and refresh token from the filesystem
// or returns nil if the file has not yet been created.
func ReadTokenAndRefreshToken(fp string) (*TokenAndRefreshToken, error) {
	exists, err := exists(fp)
	if err != nil {
		return nil, fmt.Errorf("exists(path=%q) failed: %w", fp, err)
	}

	if !exists {
		return nil, fmt.Errorf("token file %q not found: %w", fp, err)
	}

	f, err := os.Open(fp)
	if err != nil {
		return nil, fmt.Errorf("failed to open file %q: %w", fp, err)
	}
	defer f.Close()

	var tar TokenAndRefreshToken
	json.NewDecoder(f).Decode(&tar)
	return &tar, nil
}

// WriteTokenAndRefreshToken writes a copy of the token and refresh token to file.
func WriteTokenAndRefreshToken(filename string, tar *TokenAndRefreshToken) error {
	err := ensureConfigDirExists()
	if err != nil {
		return fmt.Errorf("couldn't ensure config dir exists: %w", err)
	}
	hd, err := homeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}
	filepath := filepath.Join(hd, configDir, filename)

	f, err := os.Create(filepath)
	if err != nil {
		return fmt.Errorf("create file %q failed: %w", filepath, err)
	}
	defer f.Close()

	err = json.NewEncoder(f).Encode(tar)
	if err != nil {
		return fmt.Errorf("json encode token failed: %w", err)
	}
	return nil
}

// GetCurrentConfig returns a EcomConfigurations struct containing a map
// of all configurations (known as profiles to the user) along with a string
// key mapping to the current EcomConfigEntry.
func GetCurrentConfig() (cfgs *EcomConfigurations, curCfg string, err error) {
	cfgs, err = ReadConfig()
	if err != nil {
		return nil, "", fmt.Errorf("ReadConfig failed: %w", err)
	}
	curCfg, err = ReadCurrentConfigName()
	if err != nil {
		return nil, "", fmt.Errorf("ReadCurrentConfigName failed: %w", err)
	}
	return cfgs, curCfg, nil
}
