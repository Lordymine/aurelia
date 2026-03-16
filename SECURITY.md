# SECURITY

## Reporting

Do not open a public issue for security-sensitive findings.

Report privately to the repository owner with:

- affected area
- impact
- reproduction steps
- whether a secret may have been exposed

## Sensitive Data Policy

The repository must not contain:

- API keys
- bot tokens
- local `.env` files
- local MCP config with real credentials
- local databases
- runtime memory artifacts
- debug output containing provider responses or headers

## If A Secret Is Exposed

Treat it as compromised.

Required response:

1. rotate or revoke the secret outside the repository
2. remove it from the working tree
3. ensure it is not present in the published repository
4. document any process lesson in [docs/LEARNINGS.md](docs/LEARNINGS.md) if it should prevent recurrence

## Security Review Areas

Review carefully when changes affect:

- configuration loading
- provider credentials
- Telegram auth boundaries
- command execution
- MCP configuration
- persistence and runtime state
- GitHub Actions and PR execution context
