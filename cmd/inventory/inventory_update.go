package inventory

import (
	"context"
	"fmt"
	"os"
	"strconv"

	"github.com/ecommerce-builder/ecom-cli-tool/cmdvalidate"
	"github.com/ecommerce-builder/ecom-cli-tool/configmgr"
	"github.com/ecommerce-builder/ecom-cli-tool/eclient"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"gopkg.in/AlecAivazis/survey.v1"
)

// NewCmdInventoryUpdate returns new initialized instance of the update sub command
func NewCmdInventoryUpdate() *cobra.Command {
	cfgs, curCfg, err := configmgr.GetCurrentConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%+v\n", err)
		os.Exit(1)
	}

	var cmd = &cobra.Command{
		Use:   "update <inventory_id>",
		Short: "Update an individual product inventory",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			current := cfgs.Configurations[curCfg]
			client := eclient.New(current.Endpoint)
			if err := client.SetToken(&current); err != nil {
				fmt.Fprintf(os.Stderr, "%+v\n", err)
				os.Exit(1)
			}

			// inventory command parameter
			invID := args[0]
			if !cmdvalidate.IsValidUUID(invID) {
				fmt.Fprintf(os.Stderr,
					"inventory_id %q is not a valid v4 uuid\n",
					invID)
				os.Exit(1)
			}

			ctx := context.Background()
			existing, err := client.GetInventory(ctx, invID)
			if err == eclient.ErrInventoryNotFound {
				fmt.Fprintf(os.Stderr, "inventory %q not found\n", invID)
				os.Exit(1)
			}
			if err != nil {
				fmt.Fprintf(os.Stderr, "%v\n", err)
				os.Exit(1)
			}

			req, err := promptUpdateInventory(existing)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%v\n", err)
				os.Exit(1)
			}

			// if nothing has changed, then no need to update.
			if req.Onhand == nil && req.Overselling == nil {
				os.Exit(0)
			}

			inventory, err := client.UpdateInventory(ctx, invID, req)
			if err == eclient.ErrInventoryNotFound {
				fmt.Fprintf(os.Stderr, "inventory %q not found\n", invID)
				os.Exit(1)
			}
			if err != nil {
				fmt.Fprintf(os.Stderr, "%v\n", err)
				os.Exit(1)
			}

			showInventory(inventory)
		},
	}
	return cmd
}

func promptUpdateInventory(inv *eclient.Inventory) (*eclient.UpdateInventoryRequest, error) {
	var req eclient.UpdateInventoryRequest

	// onhand
	var onhand string
	o := &survey.Input{
		Message: "Onhand:",
		Default: strconv.Itoa(inv.Onhand),
	}
	survey.AskOne(o, &onhand, survey.ComposeValidators(
		survey.Required,
		func(val interface{}) error {
			str, ok := val.(string)
			if !ok {
				return errors.New("invalid response")
			}

			v, err := strconv.Atoi(str)
			if err != nil {
				return err
			}

			if v < 0 {
				return errors.New("amount must be a positive number")
			}
			return nil
		},
	))

	onh, err := strconv.Atoi(onhand)
	if err != nil {
		return nil, errors.Wrapf(err, "atoi: onhand=%d", onhand)
	}
	// only set { "onhand": v } in the request
	// if the value has changed.
	if onh != inv.Onhand {
		req.Onhand = &onh
	}

	// use date a date range?
	var overselling bool
	d := &survey.Confirm{
		Message: "Enable overselling?",
		Default: inv.Overselling,
	}
	survey.AskOne(d, &overselling, survey.Required)

	// only set { "overselling": v "} in the request
	// if the value has changed.
	if overselling != inv.Overselling {
		req.Overselling = &overselling
	}

	return &req, nil
}
