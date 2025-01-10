# infisical-secrets-set

## Environment Variables
`infisical-secrets-set` consumes the following environment variables to initialize the infisical client:
```
INFISICAL_API_URL
INFISICAL_PROJECT_ID
INFISICAL_UNIVERSAL_AUTH_CLIENT_ID
INFISICAL_UNIVERSAL_AUTH_CLIENT_SECRET
```

It also disables telemetry.


## Usage
```
cat <<'EOF' > /tmp/secret.txt
line1
line2
EOF
infisical-secrets-set --env=dev --path=/foo my-secret < /tmp/secret.txt
```

Or without a file using Bash "Here Strings":
```
infisical-secrets-set --env=dev --path=/foo my-secret-2 <<< "content"
```


## Help
Run `infisical-secrets-set --help` to get the following Usage:
```
Reads STDIN and writes into infisical with a secret name of SECRET_NAME.

Usage:
  infisical-secrets-set [OPTIONS] SECRET_NAME [flags]

Flags:
      --env string         environment to use (dev|stage|prod) (default "prod")
  -h, --help               help for infisical-secrets-set
  -l, --log-level string   set log level to one of (debug|info|warning|error) (default "info")
  -o, --overwrite          overwrite secret if exists
      --path string        path to use (default "/")
```


## Libraries
This tool uses `spf13/cobra` for argument parsing.
