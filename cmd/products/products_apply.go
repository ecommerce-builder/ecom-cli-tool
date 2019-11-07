package products

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"

	"github.com/ecommerce-builder/ecom-cli-tool/configmgr"
	"github.com/ecommerce-builder/ecom-cli-tool/eclient"
	"github.com/spf13/cobra"
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
				fmt.Fprintf(os.Stderr, "%+v\n", err)
				os.Exit(1)
			}

			// load all price lists
			ctx := context.Background()
			priceLists, err := client.GetPriceLists(ctx)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%+v\n", err)
				os.Exit(1)
			}

			// load all products
			products, err := client.GetProducts(context.TODO())
			if err != nil {
				fmt.Fprintf(os.Stderr, "%+v\n", err)
				os.Exit(1)
			}

			isDir, err := isDirectory(args[0])
			if err != nil {
				fmt.Fprintf(os.Stderr, "%+v\n", err)
				os.Exit(1)
			}
			if !isDir {
				if err := applyProduct(client, products, priceLists, args[0]); err != nil {
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
				if err := applyProduct(client, products, priceLists, file); err != nil {
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

func applyProduct(ec *eclient.EcomClient, products []*eclient.ProductResponse, priceLists []*eclient.PriceList, filename string) error {
	// create a map of priceListCode -> priceListID
	priceListCodeToID := make(map[string]string)
	for _, p := range priceLists {
		priceListCodeToID[p.PriceListCode] = p.ID
	}

	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("os.Open(%q) failed: %w", filename, err)
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

		if err := ec.DeleteProductImages(product.ID); err != nil {
			fmt.Fprintf(os.Stderr, "failed to delete images for product sku=%s\n", product.SKU)
			os.Exit(1)
		}
	} else {
		product, err = ec.CreateProduct(&request)
		if err != nil {
			return err
		}
	}

	for _, i := range container.Product.Images {
		ir := eclient.ImageRequest{
			ProductID: product.ID,
			Path:      i.Path,
		}
		_, err := ec.CreateImage(ir)
		if err != nil {
			fmt.Fprintf(os.Stderr, "create image failed: %v", err)
			os.Exit(1)
		}
	}

	// set prices
	newPrices := make([]*eclient.PriceRequest, 0)
	for priceListCode, prices := range container.Product.Prices {
		for _, price := range prices {
			pr := eclient.PriceRequest{
				Break:     price.Break,
				UnitPrice: price.UnitPrice,
			}
			newPrices = append(newPrices, &pr)
		}
		ec.SetPrices(product.ID, priceListCodeToID[priceListCode], newPrices)
	}

	return nil
}

func isDirectory(path string) (bool, error) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false, fmt.Errorf("os.Stat(%q) failed: %w", path, err)
	}
	return fileInfo.IsDir(), err
}
