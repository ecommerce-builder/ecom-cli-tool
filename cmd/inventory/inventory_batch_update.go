package inventory

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"text/tabwriter"

	"github.com/ecommerce-builder/ecom-cli-tool/configmgr"
	"github.com/ecommerce-builder/ecom-cli-tool/eclient"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

// NewCmdInventoryBatchUpdate returns new initialized instance of batch-update sub command
func NewCmdInventoryBatchUpdate() *cobra.Command {
	cfgs, curCfg, err := configmgr.GetCurrentConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%+v\n", err)
		os.Exit(1)
	}

	var cmd = &cobra.Command{
		Use:   "batch-update <inventory.yaml>",
		Short: "Batch update inventory",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			current := cfgs.Configurations[curCfg]
			client := eclient.New(current.Endpoint)
			if err := client.SetToken(&current); err != nil {
				fmt.Fprintf(os.Stderr, "%+v\n", err)
				os.Exit(1)
			}

			data, err := ioutil.ReadFile(args[0])
			if err != nil {
				fmt.Fprintf(os.Stderr, "%v\n", err)
				os.Exit(1)
			}

			ctx := context.Background()
			products, err := client.GetProducts(ctx)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%+v\n", err)
				os.Exit(1)
			}

			productMap := make(map[string]string, len(products))
			for _, p := range products {
				productMap[p.SKU] = p.ID
			}

			var invYAML eclient.InventoryBatchContainerYAML
			if err = yaml.Unmarshal([]byte(data), &invYAML); err != nil {
				fmt.Fprintf(os.Stderr, "%+v\n", err)
				os.Exit(1)
			}

			req := buildRequest(productMap, &invYAML)
			inv, err := client.UpdateInventoryBatch(ctx, req)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%+v", err)
				os.Exit(1)
			}

			format := "%s\t%s\t%v\t%v\t%v\t%v\t\n"
			tw := new(tabwriter.Writer).Init(os.Stdout, 0, 8, 2, ' ', 0)
			fmt.Fprintf(tw, "%s\t%s\t%s\t%s\t%s\t%s\t\n",
				"Inventory ID",
				"Product SKU",
				"Onhand",
				"Overselling",
				"Created",
				"Modified")
			fmt.Fprintf(tw, "%s\t%s\t%s\t%s\t%s\t%s\t\n",
				"------------",
				"-----------",
				"------",
				"-----------",
				"-------",
				"--------")
			for _, v := range inv {
				fmt.Fprintf(tw, format,
					v.ID,
					v.ProductSKU,
					v.Onhand,
					v.Overselling,
					v.Created.In(location).Format(timeDisplayFormat),
					v.Modified.In(location).Format(timeDisplayFormat))
			}
			tw.Flush()
		},
	}
	return cmd
}

func buildRequest(productMap map[string]string, invYAML *eclient.InventoryBatchContainerYAML) []*eclient.InventoryBatchUpdateRequest {
	list := make([]*eclient.InventoryBatchUpdateRequest, 0, 16)
	for _, v := range invYAML.Inventory {
		productID := productMap[v.SKU]
		onhand := v.Onhand
		overselling := v.Overselling
		inv := eclient.InventoryBatchUpdateRequest{
			ProductID:   &productID,
			Onhand:      &onhand,
			Overselling: &overselling,
		}
		list = append(list, &inv)
	}
	return list
}
