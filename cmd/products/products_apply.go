package products

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"

	"github.com/ecommerce-builder/ecom-cli-tool/configmgr"
	"github.com/ecommerce-builder/ecom-cli-tool/eclient"
	"github.com/spf13/cobra"

	"github.com/pkg/errors"
)

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

			// load all products
			pcontainer, err := client.GetProducts()
			if err != nil {
				fmt.Fprintf(os.Stderr, "%+v\n", err)
				os.Exit(1)
			}

			products := pcontainer.Data

			isDir, err := isDirectory(args[0])
			if err != nil {
				fmt.Fprintf(os.Stderr, "%+v\n", err)
				os.Exit(1)
			}
			if !isDir {
				if err := applyProduct(client, products, args[0]); err != nil {
					if err == errMissingEAN {
						fmt.Fprintf(os.Stderr, "Skipping %s as EAN is missing\n", args[0])
						os.Exit(1)
					}
					fmt.Fprintf(os.Stderr, "%+v\n", err)
					os.Exit(1)
				}
				os.Exit(0)
			}

			matches, err := filepath.Glob(args[0] + "/*.yaml")
			if err != nil {
				fmt.Fprintf(os.Stderr, "%+v\n", err)
				os.Exit(1)
			}

			for _, file := range matches {
				if err := applyProduct(client, products, file); err != nil {
					if err == errMissingEAN {
						fmt.Fprintf(os.Stderr, "Skipping %s as EAN is missing\n", file)
						continue
					}
					fmt.Fprintf(os.Stderr, "%+v\n", err)
					os.Exit(1)
				}
			}
			os.Exit(0)
		},
	}
	return cmd
}

var errMissingEAN = errors.New("missing EAN")

func findProductBySKU(products []*eclient.ProductResponse, sku string) *eclient.ProductResponse {
	for _, p := range products {
		if sku == p.SKU {
			return p
		}
	}
	return nil
}

func applyProduct(ec *eclient.EcomClient, products []*eclient.ProductResponse, filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return errors.Wrapf(err, "os.Open(%q) failed", filename)
	}
	defer file.Close()

	container := eclient.ProductContainerYAML{}
	dec := yaml.NewDecoder(file)
	if err := dec.Decode(&container); err != nil {
		return err
	}

	product := findProductBySKU(products, container.Product.SKU)
	request := eclient.ProductRequest{
		Path: container.Product.Path,
		SKU:  container.Product.SKU,
		Name: container.Product.Name,
	}
	if product != nil {
		_, err = ec.ReplaceProduct(product.ID, &request)
		if err != nil {
			return err
		}
		return nil
	}

	_, err = ec.CreateProduct(&request)
	if err != nil {
		return err
	}
	return nil
}

func isDirectory(path string) (bool, error) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false, errors.Wrapf(err, "os.Stat(%q) failed", path)
	}
	return fileInfo.IsDir(), err
}
