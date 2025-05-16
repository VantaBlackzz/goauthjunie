# JWT Authentication Server

A Go-based JWT authentication server built with the Gin framework. This server provides endpoints for user registration, login, token refresh, and logout.

## Project Structure

```
├── .env                    # Environment variables
├── go.mod                  # Go module file
├── go.sum                  # Go module checksum file
├── main.go                 # Entry point
└── internal                # Internal packages
    ├── config              # Configuration
    ├── models              # Data models
    ├── repository          # Data access layer
    ├── service             # Business logic
    ├── handlers            # HTTP handlers
    ├── middleware          # Middleware
    └── utils               # Utility functions
```

## Features

- User registration and login
- JWT-based authentication
- Access and refresh tokens
- Token refresh
- Logout functionality
- Protected routes
- API documentation with Swagger

## API Endpoints

### Authentication

- `POST /auth/register` - Register a new user
- `POST /auth/login` - Login and get tokens
- `POST /auth/refresh` - Refresh access token
- `POST /auth/logout` - Logout (invalidate refresh token)

### User

- `GET /user/profile` - Get user profile (protected route)

## Getting Started

1. Clone the repository
2. Configure the `.env` file
3. Run the server:

```bash
go run main.go
```

## API Usage Examples

### Register a new user

```bash
curl -X POST http://localhost:8080/auth/register \
  -H "Content-Type: application/json" \
  -d '{"username": "testuser", "email": "test@example.com", "password": "password123"}'
```

### Login

```bash
curl -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username": "testuser", "password": "password123"}'
```

### Get user profile (protected)

```bash
curl -X GET http://localhost:8080/user/profile \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

### Refresh token

```bash
curl -X POST http://localhost:8080/auth/refresh \
  -H "Content-Type: application/json" \
  -d '{"refresh_token": "YOUR_REFRESH_TOKEN"}'
```

### Logout

```bash
curl -X POST http://localhost:8080/auth/logout \
  -H "Content-Type: application/json" \
  -d '{"refresh_token": "YOUR_REFRESH_TOKEN"}'
```

## API Documentation

The API is documented using Swagger. You can access the Swagger UI at:

```
http://localhost:8080/swagger/index.html
```

This provides an interactive documentation where you can:
- View all available endpoints
- See request and response schemas
- Test the API directly from the browser

## Profiling with pprof

The application includes Go's built-in profiling tool, pprof, which helps analyze performance and identify bottlenecks.

### Accessing pprof

The pprof server runs on port 6060. You can access the various profiling endpoints:

```
http://localhost:6060/debug/pprof/
```

### Common pprof Commands

You can use the Go pprof tool to analyze the profiles:

#### CPU Profiling

```bash
# Collect a 30-second CPU profile
go tool pprof http://localhost:6060/debug/pprof/profile?seconds=30

# Once in the pprof interactive console, you can use commands like:
(pprof) top10      # Show top 10 functions by CPU usage
(pprof) web        # Generate a graph visualization (requires graphviz)
(pprof) list <func> # Show source code for a function with CPU usage
```

#### Memory Profiling

```bash
# Analyze heap allocations
go tool pprof http://localhost:6060/debug/pprof/heap

# Analyze memory allocations
go tool pprof http://localhost:6060/debug/pprof/allocs
```

#### Goroutine Profiling

```bash
# Analyze goroutines
go tool pprof http://localhost:6060/debug/pprof/goroutine
```

#### Block Profiling

```bash
# Analyze blocking operations
go tool pprof http://localhost:6060/debug/pprof/block
```

### Visualizing Profiles in Browser

You can also use the interactive web UI:

```bash
go tool pprof -http=:8081 http://localhost:6060/debug/pprof/profile
```

This will open a web browser with an interactive visualization of the profile.

## Security Considerations

- The JWT secret key in the `.env` file should be changed in production
- In a production environment, use a proper database instead of the in-memory store
- Consider adding rate limiting to prevent brute force attacks
- Use HTTPS in production
