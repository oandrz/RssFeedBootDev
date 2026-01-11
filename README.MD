# Gator - RSS Feed Aggregator CLI

Gator is a command-line RSS feed aggregator written in Go. It allows you to follow RSS feeds, aggregate posts, and browse them from your terminal.

## Prerequisites

Before running Gator, you need to have the following installed:

### PostgreSQL

Gator uses PostgreSQL as its database. Install it using your package manager:

**macOS (Homebrew):**
```bash
brew install postgresql@15
brew services start postgresql@15
```

**Ubuntu/Debian:**
```bash
sudo apt update
sudo apt install postgresql postgresql-contrib
sudo systemctl start postgresql
```

**Windows:**
Download and install from [postgresql.org](https://www.postgresql.org/download/windows/)

After installing PostgreSQL, create a database for Gator:
```bash
createdb gator
```

### Go

Gator requires Go 1.24 or later. Install it from [go.dev/dl](https://go.dev/dl/) or using a package manager:

**macOS (Homebrew):**
```bash
brew install go
```

**Ubuntu/Debian:**
```bash
sudo apt install golang-go
```

## Installation

Install Gator directly using `go install`:

```bash
go install github.com/yourusername/gator@latest
```

Or clone the repository and build locally:

```bash
git clone https://github.com/yourusername/gator.git
cd gator
go build -o gator .
```

After building, move the binary to a directory in your PATH:
```bash
sudo mv gator /usr/local/bin/
```

## Configuration

Gator requires a configuration file at `~/.gatorconfig.json`. Create it with your database connection string:

```json
{
  "db_url": "postgres://username:password@localhost:5432/gator?sslmode=disable"
}
```

Replace `username` and `password` with your PostgreSQL credentials. If you're using local authentication without a password, you can use:

```json
{
  "db_url": "postgres://yourusername:@localhost:5432/gator?sslmode=disable"
}
```

## Database Migrations

Before using Gator, run the database migrations. If you have [goose](https://github.com/pressly/goose) installed:

```bash
cd sql/schema
goose postgres "postgres://username:@localhost:5432/gator" up
```

## Usage

### User Management

**Register a new user:**
```bash
gator register <username>
```

**Login as an existing user:**
```bash
gator login <username>
```

**List all users:**
```bash
gator users
```

### Feed Management

**Add a new RSS feed (must be logged in):**
```bash
gator addfeed "Feed Name" "https://example.com/feed.xml"
```

**List all available feeds:**
```bash
gator feeds
```

**Follow a feed (must be logged in):**
```bash
gator follow "https://example.com/feed.xml"
```

**View feeds you're following:**
```bash
gator following
```

**Unfollow a feed:**
```bash
gator unfollow "https://example.com/feed.xml"
```

### Aggregating and Browsing

**Start the feed aggregator:**
```bash
gator agg 1m
```
This will fetch new posts from feeds every 1 minute. You can adjust the interval (e.g., `30s`, `5m`, `1h`).

**Browse aggregated posts:**
```bash
gator browse        # Shows 2 most recent posts (default)
gator browse 10     # Shows 10 most recent posts
```

### Admin Commands

**Reset (delete all users):**
```bash
gator reset
```

## Development

For development, you can run the program directly with:
```bash
go run . <command> [args...]
```

Note: `go run .` is for development only. For production use, install the compiled binary with `go install` or `go build`.

## Example Workflow

```bash
# 1. Register and login
gator register alice
gator login alice

# 2. Add some RSS feeds
gator addfeed "Boot.dev Blog" "https://blog.boot.dev/index.xml"
gator addfeed "Lane's Blog" "https://www.wagslane.dev/index.xml"

# 3. Start aggregating (in a separate terminal)
gator agg 1m

# 4. Browse the latest posts
gator browse 5
```

## License

MIT
