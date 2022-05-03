# SiriDB Enodo listener

**Before making a commit, make sure the linter succeeds:**

```
docker run --rm -v $(pwd):/app -w /app golangci/golangci-lint:v1.45.2 golangci-lint run -v
```

### Environment variable

Variable                      | Description
----------------------------- | -----------
ENODO_HUB_HOSTNAME            | Hostname/FQDN or IP address of the Enodo Hub, for example: enodohub
ENODO_HUB_PORT                | Connect to the Enodo Hub on this TCP port, for example: 9103.
ENODO_TCP_PORT                | Listen to this TCP port, for example: 9104.
ENODO_INTERNAL_SECURITY_TOKEN | (Optional) Security Token for connecting to the Hub.
