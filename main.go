
package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/infisical/go-sdk"
	"github.com/spf13/cobra"
)

func main() {
	var env, path string

	rootCmd := &cobra.Command{
		Use:   "infisical-secrets-set [OPTIONS] SECRET_NAME",
		Short: "Writes STDIN into infisical with a secret name of SECRET_NAME.",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			secretName := args[0]

			client := infisical.NewInfisicalClient(context.Background(), infisical.Config{
				SiteUrl:          os.Getenv("INFISICAL_API_URL"),
				AutoTokenRefresh: true,
				SilentMode:       true,
			})

			_, err := client.Auth().UniversalAuthLogin("", "")
			if err != nil {
				fmt.Printf("Authentication failed: %v\n", err)
				os.Exit(1)
			}

			reader := bufio.NewReader(os.Stdin)
			secretValue, err := reader.ReadString('\n')
			if err != nil {
				fmt.Printf("Failed to read secret value: %v\n", err)
				os.Exit(1)
			}

			secretValue = strings.TrimSpace(secretValue)

			_, err = client.Secrets().Create(infisical.CreateSecretOptions{
				ProjectID:   os.Getenv("INFISICAL_PROJECT_ID"),
				Environment: env,
				SecretKey:   secretName,
				SecretValue: secretValue,
				SecretPath:  path,
			})

			if err != nil {
				fmt.Printf("Failed to set secret: %v\n", err)
				os.Exit(1)
			}

			fmt.Println("Secret set successfully.")
		},
	}

	rootCmd.Flags().StringVar(&env, "env", "prod", "Environment to use (dev|stage|prod). Default: prod")
	rootCmd.Flags().StringVar(&path, "path", "/", "Path to use. Default: /")

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
