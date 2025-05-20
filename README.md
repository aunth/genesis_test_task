# Weather Update Service

A Go-based service that provides weather updates via email subscriptions. Users can subscribe to receive weather updates for their chosen cities at either hourly or daily intervals.

## Features

- Real-time weather data from OpenWeather API
- Email subscription system with confirmation flow
- Configurable update frequency (hourly/daily)
- Gmail integration for sending emails
- PostgreSQL database for storing subscriptions
- Docker support for easy deployment

## Prerequisites

- Go 1.21 or later
- Docker and Docker Compose
- PostgreSQL (if running locally)
- OpenWeather API key
- Gmail API credentials

## Environment Variables

Set the following environment variables:

```bash
# For macOS/Linux:
export OPENWEATHER_API_KEY=your_api_key
export GMAIL_CREDENTIALS=your_gmail_credentials
export GMAIL_FROM=your_email@gmail.com
export DB_HOST=localhost
export DB_PORT=5432
export DB_USER=postgres
export DB_PASSWORD=your_password
export DB_NAME=weather_service

# For Windows (Command Prompt):
set OPENWEATHER_API_KEY=your_api_key
set GMAIL_CREDENTIALS=your_gmail_credentials
set GMAIL_FROM=your_email@gmail.com
set DB_HOST=localhost
set DB_PORT=5432
set DB_USER=postgres
set DB_PASSWORD=your_password
set DB_NAME=weather_service

# For Windows (PowerShell):
$env:OPENWEATHER_API_KEY="your_api_key"
$env:GMAIL_CREDENTIALS="your_gmail_credentials"
$env:GMAIL_FROM="your_email@gmail.com"
$env:DB_HOST="localhost"
$env:DB_PORT="5432"
$env:DB_USER="postgres"
$env:DB_PASSWORD="your_password"
$env:DB_NAME="weather_service"
```

## Getting Started

1. Clone the repository:
```bash
git clone <repository-url>
cd weather-service
```

2. Set up Gmail API credentials:
   - Go to Google Cloud Console
   - Create a new project
   - Enable Gmail API
   - Create OAuth 2.0 credentials
   - Download credentials and set in GMAIL_CREDENTIALS
   - Use the auth tool to generate GMAIL_TOKEN:
   ```bash
   go run cmd/auth/main.go
   ```

3. Get OpenWeather API key:
   - Sign up at https://openweathermap.org/
   - Get your API key
   - Add it to OPENWEATHER_API_KEY in .env

4. Start the service:
```bash
# Run database migrations
./migrate

# Start the service
go run main.go
```

## API Endpoints

- `GET /weather?city={city}` - Get current weather for a city
- `POST /subscribe` - Subscribe to weather updates
- `GET /confirm/{token}` - Confirm subscription
- `POST /unsubscribe?token={token}` - Unsubscribe from updates

## Development

### Project Structure

```
.
├── cmd/
│   ├── auth/           # Gmail authentication tool
│   ├── migrate/        # Database migrations
│   └── test-email/     # Email testing tool
├── internal/
│   ├── database/       # Database connection and migrations
│   ├── handlers/       # HTTP request handlers
│   ├── models/         # Data models
│   ├── services/       # Business logic
│   ├── static/         # Static files (HTML, CSS)
│   └── server/         # HTTP server setup
├── docker-compose.yml
├── Dockerfile
└── main.go
```

### Database Schema

```sql
CREATE TABLE subscriptions (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) NOT NULL,
    city VARCHAR(255) NOT NULL,
    frequency VARCHAR(10) NOT NULL,
    confirmation_token UUID NOT NULL,
    confirmed BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

## Docker Support

Build and run with Docker:

```bash
# Build the image
docker build -t weather-service .

# Run the container
docker run -p 8080:8080 weather-service
```

Or use Docker Compose:

```bash
docker compose up -d
```


## Bonus task

My application is running on this ip address: http://146.190.33.131:8080




## Contributing

1. Fork the repository
2. Create your feature branch
3. Commit your changes
4. Push to the branch
5. Create a Pull Request

## License

This project is licensed under the MIT License.