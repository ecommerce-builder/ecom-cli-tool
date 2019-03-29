package cmd

import (
	"github.com/spf13/cobra"
)

// sysinfoCmd represents the sysinfo command
var sysinfoCmd = &cobra.Command{
	Use:   "sysinfo",
	Short: "Prints system information from the running API service.",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		//current := rc.Configurations[rc.CurrentProject.Name]
		//ecomClient := eclient.NewEcomClient(current.FirebaseAPIKey, current.Endpoint, timeout)
		//ecomClient.SetJWT("eyJhbGciOiJSUzI1NiIsImtpZCI6IjBhZDdkNTY3ZWQ3M2M2NTEzZWQ0ZTE0ZTc4OGRjZWU4NjZlMzY3ZDMiLCJ0eXAiOiJKV1QifQ.eyJuYW1lIjoiQW5keSBhZG1pbmYiLCJjdXVpZCI6IjkzYzdlYWEyLTRmYzEtNGY3Yy1iNzJhLTY1MTc5Y2E4NTVmNSIsInJvbGUiOiJhZG1pbiIsImlzcyI6Imh0dHBzOi8vc2VjdXJldG9rZW4uZ29vZ2xlLmNvbS90ZXN0LWRhdGEtc3B5Y2FtZXJhY2N0diIsImF1ZCI6InRlc3QtZGF0YS1zcHljYW1lcmFjY3R2IiwiYXV0aF90aW1lIjoxNTUyOTk1OTY3LCJ1c2VyX2lkIjoiaHVXWGV6Q05kSk8yWHhGMlFybnppMTVZODFqMiIsInN1YiI6Imh1V1hlekNOZEpPMlh4RjJRcm56aTE1WTgxajIiLCJpYXQiOjE1NTI5OTU5NjcsImV4cCI6MTU1Mjk5OTU2NywiZW1haWwiOiJhbmR5K2FkbWluQGFuZHlmdXNuaWFrLmNvbSIsImVtYWlsX3ZlcmlmaWVkIjpmYWxzZSwiZmlyZWJhc2UiOnsiaWRlbnRpdGllcyI6eyJlbWFpbCI6WyJhbmR5K2FkbWluQGFuZHlmdXNuaWFrLmNvbSJdfSwic2lnbl9pbl9wcm92aWRlciI6InBhc3N3b3JkIn19.KcZ-7kdHhhAHsENxRMxIX3td196VjHvbUtD1UBqI26WJbeVOwt3SHL-dAzAVLb72LaP891fHvdJ-Hqe_w-77u1DWxPcMFxyPx9nc-i40z8nK3kv8ZrnmNfh1UBz1fcqHpenQFFMGUhXXVUN04xxlS7pUAK0-aUVay7tUDI1g1JLrhY0cBDvoGrJNM3x-Sc6N2vce9JZyBptBcN9QrIr2WD8CIWCHvyKvecpg1ywMbVgyMEtCZEb8tt2iTKY64zdsVKr7xYyjjczVLTaPG1-HRyzKSWtzGc5nGDoye8wtpQzL5dhe1bcfSn8gZByOn9d5tyHVsFRAJkKOQ-IEVGV8GQ")
		// sysInfo, err := ecomClient.SysInfo()
		// if err != nil {
		// 	fmt.Fprintf(os.Stderr, "%v\n", err)
		// 	os.Exit(1)
		// }

		// format := "%v\t%v\t\n"
		// tw := new(tabwriter.Writer).Init(os.Stdout, 0, 8, 2, ' ', 0)
		// fmt.Fprintf(tw, format, "API Service", "")
		// fmt.Fprintf(tw, format, "-----------", "")
		// fmt.Fprintf(tw, format, "Version", sysInfo.APIVersion)
		// fmt.Fprintf(tw, format, "", "")
		// fmt.Fprintf(tw, format, "Postgres", "")
		// fmt.Fprintf(tw, format, "--------", "")
		// fmt.Fprintf(tw, format, "Host", sysInfo.Env.Pg.Host)
		// fmt.Fprintf(tw, format, "Port", sysInfo.Env.Pg.Port)
		// fmt.Fprintf(tw, format, "Database", sysInfo.Env.Pg.Database)
		// fmt.Fprintf(tw, format, "User", sysInfo.Env.Pg.User)
		// fmt.Fprintf(tw, format, "SSLMode", sysInfo.Env.Pg.SslMode)
		// fmt.Fprintf(tw, format, "", "")
		// fmt.Fprintf(tw, format, "Google", "")
		// fmt.Fprintf(tw, format, "------", "")
		// fmt.Fprintf(tw, format, "Project ID", sysInfo.Env.Goog.ProjectID)
		// fmt.Fprintf(tw, format, "", "")
		// fmt.Fprintf(tw, format, "App", "")
		// fmt.Fprintf(tw, format, "---", "")
		// fmt.Fprintf(tw, format, "HTTP Port", sysInfo.Env.App.HTTPPort)
		// fmt.Fprintf(tw, format, "Root Email", sysInfo.Env.App.RootEmail)
		// tw.Flush()
	},
}

func init() {
	rootCmd.AddCommand(sysinfoCmd)
}
