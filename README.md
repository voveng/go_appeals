# Go Appeals System

A Go-based appeals management system with SQLite storage and REST API.

## Features

- RESTful API for managing appeals
- SQLite database for data persistence
- Support for different appeal statuses (New, In Progress, Completed, Cancelled)
- Date-based filtering of appeals
- Automatic cancellation of in-progress appeals
- Comprehensive test coverage

## Prerequisites

- Go 1.19 or higher
- SQLite3

## Installation

1. Clone the repository:
```bash
git clone <repository-url>
cd go_appeals
```

2. Install dependencies:
```bash
go mod tidy
```

3. Build the application:
```bash
go build -o appeals cmd/server/main.go
```

## Usage

1. Run the server:
```bash
go run cmd/server/main.go
```

2. The server will start on `http://localhost:8080` by default.

## API Endpoints

- `POST /appeals` - Create a new appeal
- `GET /appeals` - Get all appeals
- `GET /appeals/:id` - Get an appeal by ID
- `PUT /appeals/:id` - Update an appeal
- `PUT /appeals/:id/solution` - Add a solution to an appeal
- `PUT /appeals/:id/cancel` - Cancel an appeal
- `GET /appeals/filter` - Filter appeals by date range
- `POST /appeals/cancel-all` - Cancel all in-progress appeals

## Database Schema

The application uses an SQLite database with a single `appeals` table containing:
- `id` - Unique identifier (UUID)
- `theme` - Appeal theme
- `message` - Appeal message
- `status` - Current status (New, In Progress, Completed, Cancelled)
- `solution` - Solution provided for the appeal
- `cansel_reason` - Reason for cancellation
- `created_at` - Creation timestamp
- `updated_at` - Last update timestamp

## Testing

Run all tests:
```bash
go test ./...
```

Run tests for a specific package:
```bash
go test ./repository/
go test ./models/
```

## Architecture

The application follows a layered architecture:
- `handlers` - HTTP request handlers
- `models` - Data structures and business logic
- `repository` - Database operations
- `services` - Business logic layer

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests for new functionality
5. Ensure all tests pass
6. Submit a pull request

## License

This project is licensed under the MIT License - see the LICENSE file for details.