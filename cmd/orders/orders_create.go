package orders

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/ecommerce-builder/ecom-cli-tool/cmdvalidate"
	"github.com/ecommerce-builder/ecom-cli-tool/configmgr"
	"github.com/ecommerce-builder/ecom-cli-tool/eclient"
	"github.com/spf13/cobra"
	"gopkg.in/AlecAivazis/survey.v1"
)

// NewCmdOrdersCreate returns new initialized instance of the create sub command
func NewCmdOrdersCreate() *cobra.Command {
	cfgs, curCfg, err := configmgr.GetCurrentConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%+v\n", err)
		os.Exit(1)
	}
	var cmd = &cobra.Command{
		Use:   "create <cart_id>",
		Short: "Place an order for a cart",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			current := cfgs.Configurations[curCfg]
			client := eclient.New(current.Endpoint)
			if err := client.SetToken(&current); err != nil {
				fmt.Fprintf(os.Stderr, "%+v\n", err)
				os.Exit(1)
			}

			cartID := args[0]
			if !cmdvalidate.IsValidUUID(cartID) {
				fmt.Fprintf(os.Stderr, "cart_id %q is not a valid v4 uuid\n",
					cartID)
				os.Exit(1)
			}

			req, err := promptCreateOrder()
			if err != nil {
				fmt.Fprintf(os.Stderr, "%+v\n", err)
				os.Exit(1)
			}
			req.CartID = &cartID
			fmt.Println(req)

			ctx := context.Background()
			order, err := client.PlaceOrder(ctx, req)
			if err == eclient.ErrCartNotFound {
				fmt.Fprintf(os.Stderr, "cart %q not found\n", cartID)
				os.Exit(1)
			}
			if err != nil {
				fmt.Fprintf(os.Stderr, "%+v\n", err)
				os.Exit(1)
			}
			showOrder(order)
		},
	}
	return cmd
}

func promptCreateOrder() (*eclient.OrderRequest, error) {
	var req eclient.OrderRequest

	// check if this order is a guest order or registered user order
	var guestOrder bool
	i := &survey.Confirm{
		Message: "Is this a guest order?",
		Default: false,
	}
	survey.AskOne(i, &guestOrder, survey.Required)

	if guestOrder {
		// contact_name
		var contactName string
		n := &survey.Input{
			Message: "Contact name:",
		}
		survey.AskOne(n, &contactName, survey.Required)
		req.ContactName = &contactName

		// email
		var email string
		e := &survey.Input{
			Message: "Email:",
		}
		survey.AskOne(e, &email, survey.Required)
		req.Email = &email

		billing := promptAddress("Billing")
		shipping := promptAddress("Shipping")
		req.Billing = billing
		req.Shipping = shipping
		return &req, nil
	}

	var usrID string
	uid := &survey.Input{
		Message: "User ID:",
	}
	survey.AskOne(uid, &usrID, survey.ComposeValidators(
		survey.Required,
		func(val interface{}) error {
			str, ok := val.(string)
			if !ok {
				return errors.New("invalid response")
			}

			if !cmdvalidate.IsValidUUID(str) {
				return errors.New("user id must be a valid v4 uuid")
			}
			return nil
		},
	))
	req.UserID = &usrID

	// billing_id
	var billingID string
	bid := &survey.Input{
		Message: "Billing ID:",
	}
	survey.AskOne(bid, &billingID, survey.ComposeValidators(
		survey.Required,
		func(val interface{}) error {
			str, ok := val.(string)
			if !ok {
				return errors.New("invalid response")
			}

			if !cmdvalidate.IsValidUUID(str) {
				return errors.New("billing id must be a valid v4 uuid")
			}
			return nil
		},
	))
	req.BillingID = &billingID

	// shipping_id
	var shippingID string
	sid := &survey.Input{
		Message: "Shipping ID:",
	}
	survey.AskOne(sid, &shippingID, survey.ComposeValidators(
		survey.Required,
		func(val interface{}) error {
			str, ok := val.(string)
			if !ok {
				return errors.New("invalid response")
			}

			if !cmdvalidate.IsValidUUID(str) {
				return errors.New("shipping id must be a valid v4 uuid")
			}
			return nil
		},
	))
	req.ShippingID = &shippingID

	return &req, nil
}

func promptAddress(title string) *eclient.OrderAddressRequest {
	var v eclient.OrderAddressRequest

	// contact_name
	var contactName string
	bcn := &survey.Input{
		Message: fmt.Sprintf("%s Contact name:", title),
	}
	survey.AskOne(bcn, &contactName, survey.Required)
	v.ContactName = &contactName

	// addr1
	var addr1 string
	ba1 := &survey.Input{
		Message: fmt.Sprintf("%s Address 1:", title),
	}
	survey.AskOne(ba1, &addr1, survey.Required)
	v.Addr1 = &addr1

	// addr2
	var addr2 string
	ba2 := &survey.Input{
		Message: fmt.Sprintf("%s Address 2:", title),
	}
	survey.AskOne(ba2, &addr2, nil)
	if addr2 != "" {
		v.Addr2 = &addr2
	}

	// city
	var city string
	ct := &survey.Input{
		Message: fmt.Sprintf("%s City:", title),
	}
	survey.AskOne(ct, &city, survey.Required)
	v.City = &city

	// county
	var county string
	co := &survey.Input{
		Message: fmt.Sprintf("%s County:", title),
	}
	survey.AskOne(co, &county, nil)
	if county != "" {
		v.County = &county
	}

	// postcode
	var postcode string
	pc := &survey.Input{
		Message: fmt.Sprintf("%s Postcode:", title),
	}
	survey.AskOne(pc, &postcode, survey.Required)
	v.Postcode = &postcode

	// country_code
	var countryCode string
	cc := &survey.Input{
		Message: fmt.Sprintf("%s Country code:", title),
	}
	survey.AskOne(cc, &countryCode, survey.Required)
	v.CountryCode = &countryCode

	return &v
}
