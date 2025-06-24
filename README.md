# VietTick Backend

A comprehensive social media backend API built with Go, featuring user authentication, posts, comments, follows, and identity verification system for blue tick verification.

## Features

### 🔐 Authentication & Authorization
- JWT-based authentication with access and refresh tokens
- Email verification system
- Password change and reset functionality
- Rate limiting for security

### 👤 User Management
- User registration and login
- Profile management and updates
- User search functionality
- Username and email availability checking
- User statistics and recommendations

### 📱 Social Media Core Features
- Create, read, update, delete posts
- Image support for posts
- Like/unlike posts and comments
- Comment system with full CRUD operations
- User feed based on followed users
- Explore posts for discovery

### 👥 Follow System
- Follow/unfollow users
- Get followers and following lists
- Mutual follows detection
- Follow statistics and relationships
- Bulk follow/unfollow operations

### ✅ Verification System (Blue Tick)
- Identity verification through document upload
- Admin review system for verification requests
- Email notifications for verification status
- Verification requirements and guidelines

### 🛡️ Security & Quality
- Comprehensive middleware (CORS, rate limiting, security headers)
- Input validation and sanitization
- Error handling and logging
- Request ID tracking for debugging

## Tech Stack

- **Language**: Go 1.21+
- **Framework**: Gin (HTTP web framework)
- **Database**: MySQL 8.0+
- **Authentication**: JWT tokens
- **Email**: SMTP (configurable)
- **Password Hashing**: bcrypt

## Getting Started

### Prerequisites

- Go 1.21 or higher
- MySQL 8.0 or higher
- SMTP server access (for email features)

### Installation

1. **Clone the repository**
   ```bash
   git clone <repository-url>
   cd vietick-backend
   ```

2. **Install dependencies**
   ```bash
   go mod download
   ```

3. **Database Setup**
   ```bash
   # Create database
   mysql -u root -p -e "CREATE DATABASE vietick CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;"
   
   # Run migrations
   mysql -u root -p vietick < migrations/001_create_tables.sql
   ```

4. **Environment Configuration**
   ```bash
   cp .env.example .env
   # Edit .env with your configuration
   ```

5. **Run the server**
   ```bash
   go run cmd/main.go
   ```

The server will start on `http://localhost:8080`

## Configuration

### Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `SERVER_HOST` | Server host address | `0.0.0.0` |
| `SERVER_PORT` | Server port | `8080` |
| `DB_HOST` | MySQL host | `localhost` |
| `DB_PORT` | MySQL port | `3306` |
| `DB_USER` | MySQL username | `root` |
| `DB_PASSWORD` | MySQL password | - |
| `DB_NAME` | MySQL database name | `vietick` |
| `JWT_ACCESS_SECRET` | JWT access token secret | - |
| `JWT_REFRESH_SECRET` | JWT refresh token secret | - |
| `JWT_ACCESS_EXPIRY_HOUR` | Access token expiry in hours | `24` |
| `JWT_REFRESH_EXPIRY_DAY` | Refresh token expiry in days | `7` |
| `SMTP_HOST` | SMTP server host | `smtp.gmail.com` |
| `SMTP_PORT` | SMTP server port | `587` |
| `SMTP_USER` | SMTP username | - |
| `SMTP_PASSWORD` | SMTP password | - |
| `FROM_EMAIL` | From email address | `noreply@vietick.com` |
| `FROM_NAME` | From name | `VietTick` |

### SMTP Configuration

For Gmail, you need to:
1. Enable 2-factor authentication
2. Generate an app password
3. Use the app password as `SMTP_PASSWORD`

## API Documentation

### Base URL
```
http://localhost:8080/api/v1
```

### Authentication

Most endpoints require authentication. Include the JWT token in the Authorization header:

```
Authorization: Bearer <your-access-token>
```

### Rate Limits

- Global: 100 requests per minute
- Auth endpoints: 5 requests per minute
- API endpoints: 60 requests per minute

### Endpoints Overview

#### Authentication (`/auth`)
- `POST /auth/register` - Register new user
- `POST /auth/login` - Login user
- `POST /auth/refresh` - Refresh access token
- `POST /auth/logout` - Logout (invalidate refresh token)
- `POST /auth/logout-all` - Logout from all devices
- `POST /auth/verify-email` - Verify email address
- `POST /auth/resend-verification` - Resend verification email
- `POST /auth/change-password` - Change password
- `GET /auth/me` - Get current user info
- `GET /auth/check` - Check token validity

#### Users (`/users`)
- `GET /users/me` - Get current user profile
- `PUT /users/me` - Update profile
- `PUT /users/me/username` - Update username
- `PUT /users/me/email` - Update email
- `GET /users/{id}` - Get user profile by ID
- `GET /users/username/{username}` - Get user profile by username
- `GET /users/{id}/stats` - Get user statistics
- `GET /users/search` - Search users
- `GET /users/recommended` - Get recommended users
- `GET /users/check-username` - Check username availability
- `GET /users/check-email` - Check email availability

#### Posts (`/posts`)
- `POST /posts` - Create post
- `GET /posts/{id}` - Get post by ID
- `PUT /posts/{id}` - Update post
- `DELETE /posts/{id}` - Delete post
- `GET /posts/feed` - Get user feed
- `GET /posts/explore` - Get explore posts
- `GET /posts/search` - Search posts
- `GET /posts/user/{user_id}` - Get user posts
- `POST /posts/{id}/like` - Like post
- `POST /posts/{id}/unlike` - Unlike post
- `POST /posts/{id}/toggle-like` - Toggle like status
- `GET /posts/{id}/stats` - Get post statistics

#### Comments (`/comments` and `/posts/{id}/comments`)
- `POST /posts/{id}/comments` - Create comment
- `GET /posts/{id}/comments` - Get post comments
- `GET /comments/{id}` - Get comment by ID
- `PUT /comments/{id}` - Update comment
- `DELETE /comments/{id}` - Delete comment
- `POST /comments/{id}/like` - Like comment
- `POST /comments/{id}/unlike` - Unlike comment
- `POST /comments/{id}/toggle-like` - Toggle like status
- `GET /comments/{id}/stats` - Get comment statistics

#### Follow System (`/users/{id}/...` and `/follows`)
- `POST /users/{id}/follow` - Follow user
- `POST /users/{id}/unfollow` - Unfollow user
- `POST /users/{id}/toggle-follow` - Toggle follow status
- `GET /users/{id}/follow-status` - Get follow status
- `GET /users/{id}/followers` - Get user followers
- `GET /users/{id}/following` - Get users followed by user
- `GET /users/{id}/follow-counts` - Get follow counts
- `GET /users/{id}/mutual-follows` - Get mutual follows
- `GET /users/{id}/relationship` - Get follow relationship
- `GET /users/{id}/follow-stats` - Get follow statistics
- `POST /follows/bulk-follow` - Bulk follow users
- `POST /follows/bulk-unfollow` - Bulk unfollow users

#### Verification (`/verification`)
- `POST /verification/submit` - Submit identity verification
- `GET /verification/me` - Get user verification status
- `GET /verification/can-submit` - Check if can submit verification
- `GET /verification/requirements` - Get verification requirements
- `GET /verification/verified-users` - Get verified users list

##### Admin Only
- `GET /verification/pending` - Get pending verifications
- `GET /verification/all` - Get all verifications
- `GET /verification/{id}` - Get verification by ID
- `POST /verification/{id}/review` - Review verification
- `DELETE /verification/{id}` - Delete verification
- `GET /verification/stats` - Get verification statistics

### Response Format

#### Success Response
```json
{
  "data": {...},
  "message": "Success message"
}
```

#### Error Response
```json
{
  "error": "Error message",
  "message": "User-friendly message",
  "details": {...},
  "request_id": "unique-request-id"
}
```

### Example curl Requests

#### Auth
```bash
# Register
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "john_doe",
    "email": "john@example.com",
    "password": "securepassword",
    "full_name": "John Doe"
  }'

# Login
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john@example.com",
    "password": "securepassword"
  }'

# Refresh token
curl -X POST http://localhost:8080/api/v1/auth/refresh \
  -H "Content-Type: application/json" \
  -d '{
    "refresh_token": "<refresh_token>"
  }'

# Logout
curl -X POST http://localhost:8080/api/v1/auth/logout \
  -H "Authorization: Bearer <access_token>" \
  -H "Content-Type: application/json" \
  -d '{
    "refresh_token": "<refresh_token>"
  }'

# Logout all
curl -X POST http://localhost:8080/api/v1/auth/logout-all \
  -H "Authorization: Bearer <access_token>"

# Verify email
curl -X POST http://localhost:8080/api/v1/auth/verify-email \
  -H "Content-Type: application/json" \
  -d '{
    "token": "<verification_token>"
  }'

# Resend verification
curl -X POST http://localhost:8080/api/v1/auth/resend-verification \
  -H "Authorization: Bearer <access_token>"

# Change password
curl -X POST http://localhost:8080/api/v1/auth/change-password \
  -H "Authorization: Bearer <access_token>" \
  -H "Content-Type: application/json" \
  -d '{
    "current_password": "oldpass",
    "new_password": "newpass"
  }'

# Get current user info
curl -X GET http://localhost:8080/api/v1/auth/me \
  -H "Authorization: Bearer <access_token>"

# Check token
curl -X GET http://localhost:8080/api/v1/auth/check \
  -H "Authorization: Bearer <access_token>"
```

#### Users
```bash
# Get current profile
curl -X GET http://localhost:8080/api/v1/users/me \
  -H "Authorization: Bearer <access_token>"

# Update profile
curl -X PUT http://localhost:8080/api/v1/users/me \
  -H "Authorization: Bearer <access_token>" \
  -H "Content-Type: application/json" \
  -d '{
    "full_name": "New Name",
    "bio": "New bio"
  }'

# Update username
curl -X PUT http://localhost:8080/api/v1/users/me/username \
  -H "Authorization: Bearer <access_token>" \
  -H "Content-Type: application/json" \
  -d '{
    "username": "new_username"
  }'

# Update email
curl -X PUT http://localhost:8080/api/v1/users/me/email \
  -H "Authorization: Bearer <access_token>" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "new_email@example.com"
  }'

# Get user profile by ID
curl -X GET http://localhost:8080/api/v1/users/<user_id> \
  -H "Authorization: Bearer <access_token>"

# Get user profile by username
curl -X GET http://localhost:8080/api/v1/users/username/<username> \
  -H "Authorization: Bearer <access_token>"

# Get user stats
curl -X GET http://localhost:8080/api/v1/users/<user_id>/stats \
  -H "Authorization: Bearer <access_token>"

# Search users
curl -X GET "http://localhost:8080/api/v1/users/search?query=abc" \
  -H "Authorization: Bearer <access_token>"

# Get recommended users
curl -X GET http://localhost:8080/api/v1/users/recommended \
  -H "Authorization: Bearer <access_token>"

# Check username/email
curl -X GET "http://localhost:8080/api/v1/users/check-username?username=abc"
curl -X GET "http://localhost:8080/api/v1/users/check-email?email=abc@example.com"
```

#### Posts
```bash
# Create post
curl -X POST http://localhost:8080/api/v1/posts \
  -H "Authorization: Bearer <access_token>" \
  -H "Content-Type: application/json" \
  -d '{
    "content": "Hello, VietTick!",
    "image_urls": ["https://example.com/image.jpg"]
  }'

# Get post by ID
curl -X GET http://localhost:8080/api/v1/posts/<post_id> \
  -H "Authorization: Bearer <access_token>"

# Update post
curl -X PUT http://localhost:8080/api/v1/posts/<post_id> \
  -H "Authorization: Bearer <access_token>" \
  -H "Content-Type: application/json" \
  -d '{
    "content": "Nội dung mới"
  }'

# Delete post
curl -X DELETE http://localhost:8080/api/v1/posts/<post_id> \
  -H "Authorization: Bearer <access_token>"

# Feed/explore/search
curl -X GET http://localhost:8080/api/v1/posts/feed \
  -H "Authorization: Bearer <access_token>"
curl -X GET http://localhost:8080/api/v1/posts/explore \
  -H "Authorization: Bearer <access_token>"
curl -X GET "http://localhost:8080/api/v1/posts/search?query=abc" \
  -H "Authorization: Bearer <access_token>"

# Get user posts
curl -X GET http://localhost:8080/api/v1/posts/user/<user_id> \
  -H "Authorization: Bearer <access_token>"

# Like/Unlike/Toggle Like
curl -X POST http://localhost:8080/api/v1/posts/<post_id>/like \
  -H "Authorization: Bearer <access_token>"
curl -X POST http://localhost:8080/api/v1/posts/<post_id>/unlike \
  -H "Authorization: Bearer <access_token>"
curl -X POST http://localhost:8080/api/v1/posts/<post_id>/toggle-like \
  -H "Authorization: Bearer <access_token>"

# Get post stats
curl -X GET http://localhost:8080/api/v1/posts/<post_id>/stats \
  -H "Authorization: Bearer <access_token>"
```

#### Comments
```bash
# Create comment
curl -X POST http://localhost:8080/api/v1/posts/<post_id>/comments \
  -H "Authorization: Bearer <access_token>" \
  -H "Content-Type: application/json" \
  -d '{
    "content": "Bình luận của bạn"
  }'

# Get post comments
curl -X GET http://localhost:8080/api/v1/posts/<post_id>/comments \
  -H "Authorization: Bearer <access_token>"

# Get/Update/Delete comment
curl -X GET http://localhost:8080/api/v1/comments/<comment_id> \
  -H "Authorization: Bearer <access_token>"
curl -X PUT http://localhost:8080/api/v1/comments/<comment_id> \
  -H "Authorization: Bearer <access_token>" \
  -H "Content-Type: application/json" \
  -d '{
    "content": "Nội dung mới"
  }'
curl -X DELETE http://localhost:8080/api/v1/comments/<comment_id> \
  -H "Authorization: Bearer <access_token>"

# Like/Unlike/Toggle Like comment
curl -X POST http://localhost:8080/api/v1/comments/<comment_id>/like \
  -H "Authorization: Bearer <access_token>"
curl -X POST http://localhost:8080/api/v1/comments/<comment_id>/unlike \
  -H "Authorization: Bearer <access_token>"
curl -X POST http://localhost:8080/api/v1/comments/<comment_id>/toggle-like \
  -H "Authorization: Bearer <access_token>"

# Get comment stats
curl -X GET http://localhost:8080/api/v1/comments/<comment_id>/stats \
  -H "Authorization: Bearer <access_token>"
```

#### Follow System
```bash
# Follow/Unfollow/Toggle
curl -X POST http://localhost:8080/api/v1/users/<user_id>/follow \
  -H "Authorization: Bearer <access_token>"
curl -X POST http://localhost:8080/api/v1/users/<user_id>/unfollow \
  -H "Authorization: Bearer <access_token>"
curl -X POST http://localhost:8080/api/v1/users/<user_id>/toggle-follow \
  -H "Authorization: Bearer <access_token>"

# Get follow status, followers, following, counts, mutual, relationship, stats
curl -X GET http://localhost:8080/api/v1/users/<user_id>/follow-status \
  -H "Authorization: Bearer <access_token>"
curl -X GET http://localhost:8080/api/v1/users/<user_id>/followers \
  -H "Authorization: Bearer <access_token>"
curl -X GET http://localhost:8080/api/v1/users/<user_id>/following \
  -H "Authorization: Bearer <access_token>"
curl -X GET http://localhost:8080/api/v1/users/<user_id>/follow-counts \
  -H "Authorization: Bearer <access_token>"
curl -X GET http://localhost:8080/api/v1/users/<user_id>/mutual-follows \
  -H "Authorization: Bearer <access_token>"
curl -X GET http://localhost:8080/api/v1/users/<user_id>/relationship \
  -H "Authorization: Bearer <access_token>"
curl -X GET http://localhost:8080/api/v1/users/<user_id>/follow-stats \
  -H "Authorization: Bearer <access_token>"

# Bulk follow/unfollow
curl -X POST http://localhost:8080/api/v1/follows/bulk-follow \
  -H "Authorization: Bearer <access_token>" \
  -H "Content-Type: application/json" \
  -d '{
    "user_ids": ["id1", "id2"]
  }'
curl -X POST http://localhost:8080/api/v1/follows/bulk-unfollow \
  -H "Authorization: Bearer <access_token>" \
  -H "Content-Type: application/json" \
  -d '{
    "user_ids": ["id1", "id2"]
  }'
```

#### Verification
```bash
# Submit identity verification
curl -X POST http://localhost:8080/api/v1/verification/submit \
  -H "Authorization: Bearer <access_token>" \
  -H "Content-Type: application/json" \
  -d '{
    "full_name": "Tên",
    "id_number": "123456789",
    "id_type": "CCCD",
    "front_image_url": "...",
    "back_image_url": "...",
    "selfie_image_url": "..."
  }'

# Get verification status, requirements, verified users
curl -X GET http://localhost:8080/api/v1/verification/me \
  -H "Authorization: Bearer <access_token>"
curl -X GET http://localhost:8080/api/v1/verification/can-submit \
  -H "Authorization: Bearer <access_token>"
curl -X GET http://localhost:8080/api/v1/verification/requirements
curl -X GET http://localhost:8080/api/v1/verification/verified-users

# Admin: pending, all, get, review, delete, stats
curl -X GET http://localhost:8080/api/v1/verification/pending \
  -H "Authorization: Bearer <admin_token>"
curl -X GET http://localhost:8080/api/v1/verification/all \
  -H "Authorization: Bearer <admin_token>"
curl -X GET http://localhost:8080/api/v1/verification/<id> \
  -H "Authorization: Bearer <admin_token>"
curl -X POST http://localhost:8080/api/v1/verification/<id>/review \
  -H "Authorization: Bearer <admin_token>" \
  -H "Content-Type: application/json" \
  -d '{
    "status": "approved",
    "admin_notes": "OK"
  }'
curl -X DELETE http://localhost:8080/api/v1/verification/<id> \
  -H "Authorization: Bearer <admin_token>"
curl -X GET http://localhost:8080/api/v1/verification/stats \
  -H "Authorization: Bearer <admin_token>"
```

## Database Schema

### Key Tables

- **users** - User accounts and profiles
- **posts** - User posts/status updates
- **comments** - Comments on posts
- **post_likes** - Post likes
- **comment_likes** - Comment likes
- **follows** - Follow relationships
- **refresh_tokens** - JWT refresh tokens
- **identity_verifications** - Identity verification requests

## Development

### Project Structure

```
vietick-backend/
├── cmd/
│   └── main.go                 # Application entry point
├── internal/
│   ├── config/                 # Configuration management
│   ├── handler/               # HTTP handlers/controllers
│   ├── middleware/            # HTTP middleware
│   ├── model/                 # Data models
│   ├── repository/            # Data access layer
│   ├── service/               # Business logic layer
│   └── utils/                 # Utility functions
├── pkg/
│   ├── database/              # Database connection
│   ├── email/                 # Email service
│   └── jwt/                   # JWT management
├── migrations/                # Database migrations
├── docs/                      # Documentation
├── go.mod                     # Go module file
├── go.sum                     # Go dependencies
├── .env.example              # Environment variables example
└── README.md                 # This file
```

### Adding New Features

1. **Model** - Define data structures in `internal/model/`
2. **Repository** - Add data access methods in `internal/repository/`
3. **Service** - Implement business logic in `internal/service/`
4. **Handler** - Create HTTP handlers in `internal/handler/`
5. **Routes** - Add routes in `cmd/main.go`

### Code Guidelines

- Follow Go conventions and best practices
- Use meaningful variable and function names
- Add comments for exported functions
- Handle errors appropriately
- Write tests for critical functionality
- Use dependency injection for services

## Security Considerations

- Change default JWT secrets in production
- Use HTTPS in production
- Implement proper input validation
- Use prepared statements for database queries
- Enable CORS only for trusted domains
- Monitor rate limits and adjust as needed
- Regularly update dependencies
- Use strong passwords for database and SMTP

## Deployment

### Production Checklist

- [ ] Set strong JWT secrets
- [ ] Configure production database
- [ ] Set up HTTPS/TLS
- [ ] Configure production SMTP
- [ ] Set appropriate CORS origins
- [ ] Enable logging and monitoring
- [ ] Set up database backups
- [ ] Configure rate limiting
- [ ] Review security headers
- [ ] Test all endpoints

### Docker Deployment

```dockerfile
# Dockerfile example
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY go.* ./
RUN go mod download
COPY . .
RUN go build -o main cmd/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/main .
EXPOSE 8080
CMD ["./main"]
```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Commit your changes
4. Push to the branch
5. Create a Pull Request

## License

This project is licensed under the MIT License.

## Support

For questions or issues, please open an issue on the repository.
