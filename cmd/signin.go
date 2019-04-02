package cmd

import (
	"fmt"
	"os"

	"bitbucket.org/andyfusniakteam/ecom-cli-tool/configmgr"
	"bitbucket.org/andyfusniakteam/ecom-cli-tool/eclient"
	"github.com/spf13/cobra"
)

// signinCmd represents the signin command
var signinCmd = &cobra.Command{
	Use:   "signin",
	Short: "A brief description of your command",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		//current := rc.Projects[rc.CurrentProject.Name]
		// hostname, err := eclient.URLToHostName(rc.CurrentProject.Endpoint)
		// if err != nil {
		// 	fmt.Fprintf(os.Stderr, "%v\n", err)
		// 	os.Exit(1)
		// }
		// filename := fmt.Sprintf("%s-%s", rc.CurrentProject.FirebaseAPIKey, hostname)

		// token, err := configmgr.ReadTokenAndRefreshToken(filename)
		// if err != nil {
		// 	fmt.Fprintf(os.Stderr, "%v\n", err)
		// 	os.Exit(1)
		// }

		// var p jwt.Parser
		// t, _, err := p.ParseUnverified(token.IDToken, &jwt.StandardClaims{})

		// c := t.Claims.(*jwt.StandardClaims)
		// utcNow := time.Now().Unix()

		// // If the token has expired, use the refresh token to get another

		name, webKey, endpoint, devKey, err := promptAddConfiguration()
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
			os.Exit(1)
		}

		hostname, err := configmgr.URLToHostName(endpoint)
		filename := fmt.Sprintf("%s-%s", webKey, hostname)
		fmt.Println(filename)
		fmt.Println(name)
		// exists, err := configmgr.exists(filename)
		// if err != nil {
		// 	fmt.Fprintf(os.Stderr, "%v\n", err)
		// 	os.Exit(1)
		// }

		// if !exists {
		// 	fmt.Fprintf(os.Stderr, "A configuration already exists for these credentials. You should already be able to use the command line tool without needing to signin. If you experience problems, or wish to use a different developer key, use the signout commmand first before signin.")
		// 	os.Exit(1)
		// }

		ecomClient := eclient.NewEcomClient(webKey, endpoint, timeout)
		// if c.ExpiresAt-utcNow <= 0 {
		// 	tokens, err := ecomClient.ExchangeRefreshTokenForIDToken(token.RefreshToken)
		// 	if err != nil {
		// 		fmt.Fprintf(os.Stderr, "%+v\n", err)
		// 		os.Exit(1)
		// 	}
		// 	fmt.Printf("%+v\n", tokens)
		// }

		//jwt.Parse(token.IDToken)

		customToken, err := ecomClient.SignInWithDevKey(devKey)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
			os.Exit(1)
		}
		tar, err := ecomClient.ExchangeCustomTokenForIDAndRefreshToken(customToken)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
			os.Exit(1)
		}

		configmgr.WriteTokenAndRefreshToken(webKey, endpoint, tar)
	},
}

func init() {
	rootCmd.AddCommand(signinCmd)
}
