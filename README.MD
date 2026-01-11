# Start Database
```bash
pg15-start
```

# Stop Database
```bash
pg15-stop
```

# SQL Query
```bash
sqlc generate
```

# Migration
do this in `/schema`
## Push Migration
```bash
goose postgres "postgres://andreas.oentoro:@localhost:5432/gator" up
```

## Drop Migration
```bash
goose postgres "postgres://andreas.oentoro:@localhost:5432/gator" down
```