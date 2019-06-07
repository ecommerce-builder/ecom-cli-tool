package products

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"

	"github.com/ecommerce-builder/ecom-cli-tool/configmgr"
	"github.com/ecommerce-builder/ecom-cli-tool/eclient"
	"github.com/spf13/cobra"
)

type productContainer struct {
	Product eclient.ProductApply `yaml:"product"`
}

// NewCmdProductsApply returns new initialized instance of the apply sub command
func NewCmdProductsApply() *cobra.Command {
	cfgs, curCfg, err := configmgr.GetCurrentConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%+v\n", err)
		os.Exit(1)
	}

	var cmd = &cobra.Command{
		Use:   "apply <product.yaml>|<dir>",
		Short: "Create or update an exising product",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			current := cfgs.Configurations[curCfg]
			client := eclient.New(current.Endpoint)
			if err := client.SetToken(&current); err != nil {
				log.Fatal(err)
			}

			isDir, err := isDirectory(args[0])
			if err != nil {
				fmt.Fprintf(os.Stderr, "%v\n", err)
				os.Exit(1)
			}
			if !isDir {
				if err := applyProduct(client, args[0]); err != nil {
					if err == errMissingEAN {
						fmt.Fprintf(os.Stderr, "Skipping %s as EAN is missing\n", args[0])
						os.Exit(1)
					}
					fmt.Fprintf(os.Stderr, "%v\n", err)
					os.Exit(1)
				}
				os.Exit(0)
			}

			matches, err := filepath.Glob(args[0] + "/*.yaml")
			if err != nil {
				fmt.Fprintf(os.Stderr, "%v\n", err)
				os.Exit(1)
			}

			for _, file := range matches {
				if err := applyProduct(client, file); err != nil {
					if err == errMissingEAN {
						fmt.Fprintf(os.Stderr, "Skipping %s as EAN is missing\n", file)
						continue
					}
					fmt.Fprintf(os.Stderr, "%v\n", err)
					os.Exit(1)
				}
			}
			os.Exit(0)
		},
	}
	return cmd
}

var errMissingEAN = errors.New("missing EAN")

func applyProduct(ec *eclient.EcomClient, filen string) error {
	file, err := os.Open(filen)
	if err != nil {
		return err
	}
	defer file.Close()
	p := productContainer{}
	dec := yaml.NewDecoder(file)
	if err := dec.Decode(&p); err != nil {
		return err
	}
	//if p.Product.EAN == "" {
	//	return errMissingEAN
	//}
	pu := eclient.ProductApply{
		EAN:     p.Product.EAN,
		Path:    p.Product.Path,
		Name:    p.Product.Name,
		Images:  p.Product.Images,
		Pricing: p.Product.Pricing,
		Content: p.Product.Content,
	}
	_, err = ec.ReplaceProduct(p.Product.SKU, &pu)
	if err != nil {
		return err
	}
	return nil
}

func isDirectory(path string) (bool, error) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false, err
	}
	return fileInfo.IsDir(), err
}
