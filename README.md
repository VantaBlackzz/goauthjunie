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

## Security Considerations

- The JWT secret key in the `.env` file should be changed in production
- In a production environment, use a proper database instead of the in-memory store
- Consider adding rate limiting to prevent brute force attacks
- Use HTTPS in production
