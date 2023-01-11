# trello-app

## Compile docs

```bash
swagger-cli bundle -o docs/swagger.json docs/swagger.yml
```

## How run swagger docs

```bash
swagger serve docs/swagger.json -p 8080
```