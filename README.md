# sample-go-gqlgen

```shell
curl -X POST http://localhost:8080/graphql \
  --header 'Content-Type: application/json' \
  --header 'Authorization: Bearer test' \
  --data '{
    "query": "{__typename}",
    "extensions": {
      "persistedQuery": {
        "version": 1,
        "sha256Hash": "ecf4edb46db40b5132295c0291d62fb65d6759a9eedfa4d5d612dd5ec54a6b38"
      }
    }
  }'
```

```shell
curl -X POST http://localhost:8080/graphql \
  --header 'Content-Type: application/json' \
  --header 'Authorization: Bearer test' \
  --data '{
    "extensions": {
      "persistedQuery": {
        "version": 1,
        "sha256Hash": "ecf4edb46db40b5132295c0291d62fb65d6759a9eedfa4d5d612dd5ec54a6b38"
      }
    }
  }'
```
