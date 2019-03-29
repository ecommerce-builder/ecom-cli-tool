package configmgr

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"

	"bitbucket.org/andyfusniakteam/ecom-cli-tool/eclient"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

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
		return "", fmt.Errorf("failed user.Current(): %v", err)
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

func ensureConfigDirExists() (*string, error) {
	hd, err := homeDir()
	if err != nil {
		return nil, errors.Wrapf(err, "failed homeDir()")
	}
	configDir := filepath.Join(hd, configDir)

	exists, err := exists(configDir)
	if err != nil {
		return nil, errors.Wrapf(err, "failed exists(%s)", configDir)
	}
	if !exists {
		os.Mkdir(configDir, 0755)
		err = WriteCurrentProject("")
		if err != nil {
			return nil, errors.Wrapf(err, "failed write current project %q", "")
		}
	}
	return &configDir, nil
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

// ReadCurrentConfigName returns the contents of the CURRENT_PROJECT
// file. If the CURRENT_PROJECT file does not exists (for example, the
// first time the program is run), an empty file will be created.
func ReadCurrentConfigName() (string, error) {
	configDir, err := ensureConfigDirExists()
	if err != nil {
		return "", errors.Wrap(err, "ensure config dir exists failed")
	}
	cpf := filepath.Join(*configDir, "CURRENT_PROJECT")
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

// ReadConfig opens and read the .ecomrc file putting each section name in a map of EcomConfigEntrys.
func ReadConfig() (*EcomConfigurations, error) {
	err := ensureConfigFileExists()
	if err != nil {
		return nil, errors.Wrap(err, "failed to ensure config file exists")
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
		return nil, errors.Wrap(err, "failed to read in config file")
	}

	//configs := vp.GetStringMap("configurations")

	//fmt.Printf("%+v\n", configs)

	configurations := EcomConfigurations{}
	err = viper.Unmarshal(&configurations)
	if err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal")
	}

	//
	// for _, s := range cfg.Sections() {
	// 	if s.Name() == "DEFAULT" {
	// 		continue
	// 	}
	// 	ecomCfg[s.Name()] = EcomConfigEntry{
	// 		Name:           s.Key("name").String(),
	// 		FirebaseAPIKey: s.Key("firebase_api_key").String(),
	// 		DevKey:         s.Key("developer_key").String(),
	// 		Endpoint:       s.Key("endpoint").String(),
	// 	}
	// }
	return &configurations, nil
}

// WriteConfig writes the EcomConfigurations to the YAML config.
func WriteConfig(cfgs *EcomConfigurations) error {
	viper.Set("configurations", cfgs.Configurations)
	err := viper.WriteConfig()
	if err != nil {
		return errors.Wrap(err, "failed to write config file")
	}
	return nil
}

// WriteCurrentProject records the project API Key on the filesystem within the $HOME/.ecom directory in a file called CURRENT_API_KEY. The current API Key context is read between invocation of the command-line tool.
func WriteCurrentProject(name string) error {
	configDir, err := ensureConfigDirExists()
	if err != nil {
		return errors.Wrap(err, "couldn't ensure config dir exists")
	}

	cpf := filepath.Join(*configDir, "CURRENT_PROJECT")
	bs := []byte(name)
	err = ioutil.WriteFile(cpf, bs, 0644)
	if err != nil {
		return errors.Wrap(err, "failed to write CURRENT_PROJECT file")
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
		return false, fmt.Errorf("filename %q does not exist", filepath)
	}
	err = os.Remove(filepath)
	if err != nil {
		return false, errors.Wrapf(err, "failed to remove file %q", filepath)
	}
	return true, nil
}

// ReadTokenAndRefreshToken reads the token and refresh token from the filesystem
// or returns nil if the file has not yet been created.
func ReadTokenAndRefreshToken(filename string) (*eclient.TokenAndRefreshToken, error) {
	configDir, err := ensureConfigDirExists()
	if err != nil {
		return nil, errors.Wrap(err, "couldn't ensure config dir exists")
	}

	filepath := filepath.Join(*configDir, filename)
	exists, err := exists(filepath)
	if err != nil {
		return nil, errors.Wrapf(err, "failed exists(%s)", filepath)
	}

	if !exists {
		return nil, nil
	}

	f, err := os.Open(filepath)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to open file %q", filepath)
	}
	defer f.Close()

	var tar eclient.TokenAndRefreshToken
	json.NewDecoder(f).Decode(&tar)
	return &tar, nil
}

// WriteTokenAndRefreshToken writes a copy of the token and refresh token to the file system
func WriteTokenAndRefreshToken(filename string, tar *eclient.TokenAndRefreshToken) error {
	cfgDir, err := ensureConfigDirExists()
	if err != nil {
		return errors.Wrap(err, "couldn't ensure config dir exists")
	}

	filepath := filepath.Join(*cfgDir, filename)

	f, err := os.Create(filepath)
	if err != nil {
		return errors.Wrapf(err, "failed to open file %q", filepath)
	}
	defer f.Close()

	err = json.NewEncoder(f).Encode(tar)
	if err != nil {
		return errors.Wrapf(err, "failed to encode token")
	}

	return nil
}
