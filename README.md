# URL Shortener Service

A high-performance URL shortener service built with Go, Fiber, and Redis. This service allows you to shorten long URLs and provides rate limiting, custom short URLs, and expiry management.

## Features

- 🔗 **URL Shortening**: Convert long URLs into short, manageable links
- ⚡ **High Performance**: Built with Go and Fiber for fast response times
- 🚦 **Rate Limiting**: IP-based rate limiting to prevent abuse
- 🎯 **Custom Short URLs**: Option to create custom short URL identifiers
- ⏰ **Expiry Management**: Set expiration times for shortened URLs
- 🐳 **Docker Support**: Easy deployment with Docker Compose
- 📊 **Request Logging**: Built-in request logging middleware

## Tech Stack

- **Go 1.23+**: Programming language
- **Fiber v2**: Web framework
- **Redis**: In-memory data store for URL mapping and rate limiting
- **Docker**: Containerization
- **Docker Compose**: Multi-container orchestration

## Prerequisites

- Docker and Docker Compose installed
- Go 1.23+ (for local development)

## Quick Start

### Using Docker Compose (Recommended)

1. Clone the repository:
```bash
git clone https://github.com/harsh2971/shorten-url-go-redis.git
cd shorten-url-go-redis
```

2. Start the services:
```bash
docker compose up -d
```

3. The API will be available at `http://localhost:3300`

### Local Development

1. Navigate to the API directory:
```bash
cd api
```

2. Install dependencies:
```bash
go mod download
```

3. Set up environment variables (create a `.env` file):
```env
APP_PORT=:3300
DB_ADDR=localhost:6379
DB_PASS=
DOMAIN=http://localhost:3300
API_QUOTA=10
```

4. Make sure Redis is running locally:
```bash
docker run -d -p 6379:6379 redis:alpine
```

5. Run the application:
```bash
go run main.go
```

## API Endpoints

### Shorten URL

**POST** `/api/v1/shorten`

Shortens a long URL into a short link.

**Request Body:**
```json
{
  "url": "https://www.example.com/very/long/url",
  "custom_short": "my-link",  // Optional
  "expiry": 24                 // Optional (in hours, default: 24)
}
```

**Response:**
```json
{
  "url": "https://www.example.com/very/long/url",
  "custom_short": "http://localhost:3300/my-link",
  "expiry": 24,
  "rate_limit": 9,
  "rate_limit_reset": 29
}
```

### Resolve URL

**GET** `/:url`

Redirects to the original URL from a short URL.

**Example:**
```
GET http://localhost:3300/my-link
→ Redirects to the original URL
```

## Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `APP_PORT` | Port for the API server | `:3300` |
| `DB_ADDR` | Redis server address | `db:6379` (Docker) or `localhost:6379` (local) |
| `DB_PASS` | Redis password | `` (empty) |
| `DOMAIN` | Domain for short URLs | `http://localhost:3300` |
| `API_QUOTA` | Rate limit quota per IP | `10` |

## Docker Configuration

The project uses Docker Compose to orchestrate two services:

- **api**: Go API server (port 3300)
- **db**: Redis database (port 6379)

### Updating Environment Variables

Edit `docker-compose.yml` to change environment variables:

```yaml
environment:
  - APP_PORT=:3300
  - DB_ADDR=db:6379
  - DB_PASS=
  - DOMAIN=http://localhost:3300
  - API_QUOTA=10
```

After updating, restart the services:
```bash
docker compose up -d --force-recreate api
```

## Rate Limiting

The API implements IP-based rate limiting:
- Each IP address gets a quota (default: 10 requests)
- Quota resets after 30 minutes
- Rate limit information is included in the response

## Custom Short URLs

You can create custom short URLs by providing a `custom_short` field:
- Minimum length: 6 characters
- Maximum length: 15 characters
- Must be unique

If not provided, a random 6-character ID will be generated.

## Project Structure

```
shorten-url-go-redis/
├── api/
│   ├── database/       # Redis connection
│   ├── helpers/        # Helper functions
│   ├── routes/         # API route handlers
│   ├── main.go         # Application entry point
│   ├── Dockerfile      # API container definition
│   └── go.mod          # Go dependencies
├── db/
│   ├── Dockerfile      # Redis container definition
│   └── docker-compose.yml
├── docker-compose.yml  # Service orchestration
└── README.md
```

## Development

### Building the API

```bash
cd api
go build -o main .
```

### Running Tests

```bash
cd api
go test ./...
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is open source and available under the MIT License.

## Author

Built by [harsh2971](https://github.com/harsh2971)

