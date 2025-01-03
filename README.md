# userli-synapse-user-retention

[![Integration](https://github.com/systemli/userli-synapse-user-retention/actions/workflows/integration.yml/badge.svg)](https://github.com/systemli/userli-synapse-user-retention/actions/workflows/integration.yml) [![Maintainability Rating](https://sonarcloud.io/api/project_badges/measure?project=systemli_userli-synapse-user-retention&metric=sqale_rating)](https://sonarcloud.io/summary/new_code?id=systemli_userli-synapse-user-retention) [![Quality Gate Status](https://sonarcloud.io/api/project_badges/measure?project=systemli_userli-synapse-user-retention&metric=alert_status)](https://sonarcloud.io/summary/new_code?id=systemli_userli-synapse-user-retention) [![Coverage](https://sonarcloud.io/api/project_badges/measure?project=systemli_userli-synapse-user-retention&metric=coverage)](https://sonarcloud.io/summary/new_code?id=systemli_userli-synapse-user-retention)

This project is a simple extention for Userli to update users when they are active in Matrix.
It will also delete users that are not active anymore in Matrix.
You can run this as a cronjob or a systemd timer.

## Configuration

The following environment variables are required:

- `USERLI_URL`: The URL of the Userli instance.
- `USERLI_DOMAIN`: The domain you are using for Matrix.
- `USERLI_TOKEN`: The token to authenticate against the Userli instance.
- `SYNAPSE_URL`: The URL of the Synapse instance.
- `SYNAPSE_TOKEN`: The token to authenticate against the Synapse instance.

## Usage

We recommend to use the provided Docker image:

```bash
docker run -d --name userli-synapse-user-retention \
  -e USERLI_URL=https://userli \
  -e USERLI_DOMAIN=example.com \
  -e USERLI_TOKEN=secret \
  -e SYNAPSE_URL=https://synapse \
  -e SYNAPSE_TOKEN=secret \
  systemli/userli-synapse-user-retention:latest
```
