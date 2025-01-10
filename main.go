package main

import (
	"bufio"
	"context"
	"io"
	"log/slog"
	"os"
	"strings"

	"github.com/infisical/go-sdk"
	"github.com/spf13/cobra"
)

func main() {
	var env, path, logLevel string
	var overwrite bool

	var logLevelOption slog.Level
	switch logLevel {
	case "debug":
		logLevelOption = slog.LevelDebug
	case "info":
		logLevelOption = slog.LevelInfo
	case "warning":
		logLevelOption = slog.LevelWarn
	case "error":
		logLevelOption = slog.LevelError
	default:
		logLevelOption = slog.LevelInfo
	}

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: logLevelOption,
	}))

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
				logger.Error("Authentication failed", "error", err)
				os.Exit(1)
			}

			reader := bufio.NewReader(os.Stdin)
			secretBytes, err := io.ReadAll(reader)
			secretValue := string(secretBytes)
			if err != nil {
				logger.Error("Failed to read secret from STDIN", "error", err)
				os.Exit(1)
			}

			secretValue = strings.TrimSpace(secretValue)

			createSecret := func() error {
				_, err := client.Secrets().Create(infisical.CreateSecretOptions{
					ProjectID:   os.Getenv("INFISICAL_PROJECT_ID"),
					Environment: env,
					SecretKey:   secretName,
					SecretValue: secretValue,
					SecretPath:  path,
				})
				return err
			}

			err = createSecret()
			if err != nil {
				if strings.Contains(err.Error(), "Secret already exist") {
					if overwrite {
						logger.Debug("Secret already exists, overwriting due to --overwrite flag.")
						_, err = client.Secrets().Delete(infisical.DeleteSecretOptions{
							ProjectID:   os.Getenv("INFISICAL_PROJECT_ID"),
							Environment: env,
							SecretKey:   secretName,
							SecretPath:  path,
						})
						if err != nil {
							logger.Error("Failed to delete existing secret", "error", err)
							os.Exit(1)
						}
						err = createSecret()
					} else {
						logger.Warn("Secret exists. Use --overwrite to overwrite.")
						os.Exit(1)
					}
				}

				if err != nil {
					logger.Error("Failed to set secret", "error", err)
					os.Exit(1)
				}
			}
		},
	}

	rootCmd.Flags().StringVar(&logLevel, "log-level", "info", "Set log level to one of debug|info|warning|error (default \"info\")")
	rootCmd.Flags().BoolVarP(&overwrite, "overwrite", "o", false, "Overwrite secret if exists")
	rootCmd.Flags().StringVar(&env, "env", "prod", "Environment to use (dev|stage|prod). Default: prod")
	rootCmd.Flags().StringVar(&path, "path", "/", "Path to use. Default: /")

	if err := rootCmd.Execute(); err != nil {
		logger.Error("Command execution failed", "error", err)
		os.Exit(1)
	}
}
