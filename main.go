package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"strings"

	"github.com/infisical/go-sdk"
	"github.com/spf13/cobra"
)

func configureLogging(logLevel string) *slog.Logger {
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
	logger.Info("logging configured", "level", logLevel)
	return logger
}

func main() {
	var env, path, logLevel string
	var overwrite bool

	rootCmd := &cobra.Command{
		Use:   "infisical-secrets-set [OPTIONS] SECRET_NAME",
		Short: "Reads STDIN and writes into infisical with a secret name of SECRET_NAME.",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			var logger = configureLogging(logLevel)
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

	rootCmd.Flags().StringVarP(&logLevel, "log-level", "l", "info", "set log level to one of (debug|info|warning|error)")
	rootCmd.Flags().BoolVarP(&overwrite, "overwrite", "o", false, "overwrite secret if exists")
	rootCmd.Flags().StringVar(&env, "env", "prod", "environment to use (dev|stage|prod)")
	rootCmd.Flags().StringVar(&path, "path", "/", "path to use")

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "number of foo: %s", err)
		os.Exit(1)
	}
}
