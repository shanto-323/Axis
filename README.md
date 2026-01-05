# Chat-AI - Backend API

Backend service for authentication and multi-modal AI chat. Built with Go, PostgreSQL, and Docker.

## Quick Links

- **API Docs:** `http://localhost:8080/docs`
- **Base URL:** `/api/v1`
- **Auth:** JWT tokens (1-hour expiry)

## Tech Stack

Go • PostgreSQL • JWT • Docker

## Getting Started

### Prerequisites
- Go 1.20+
- PostgreSQL 12+
- Docker (optional)

### Installation

1. Clone the repository:
```bash
git clone https://github.com/yourusername/chat-ai.git
cd chat-ai
```

2. Install dependencies:
```bash
go mod download
```

3. Configure environment variables (see Configuration section below)

4. Run the application:
```bash
go run ./cmd/main
```

Visit `http://localhost:8080/docs` for interactive API documentation.

### Docker

```bash
docker-compose up --build
```

---

## Configuration

The application uses two configuration files:

### 1. `.env` - Application Configuration

Copy `.env.example` to `.env`:
```bash
cp .env.example .env
```

Edit `.env` with your preferences:

```dotenv
# Environment Type
ENV=local                          # Options: local, prod

# Database Type
DATABASE_TYPE=postgres             # Options: mock (testing), postgres (production)

# Logging Configuration
LEVEL=debug                        # Options: debug, info, warn, error
FORMAT=json                        # Options: json, text

# Server Configuration
SERVER_PORT=8080

# JWT Secret (change this in production!)
JWT_SECRET=your-super-secret-key-change-this

# AI Services
AI.LLM_PROVIDER=openrouter
AI.LLM_API_KEY=<your-openrouter-api-key>

AI.VLM_PROVIDER=mock
AI.VLM_API_KEY=mock-no-key-needed
```

#### Configuration Options

**ENV** (required)
- `local` - Development environment with verbose logging
- `prod` - Production environment with optimized settings

**DATABASE_TYPE** (required)
- `mock` - In-memory database for testing (data is not persisted)
- `postgres` - PostgreSQL database for production

**LEVEL** (required)
- `debug` - Verbose logging for development
- `info` - Standard information logging
- `warn` - Only warnings and errors
- `error` - Only error messages

**FORMAT** (required)
- `json` - Structured JSON logging for parsing
- `text` - Human-readable text logging

### 2. `.env.postgres` - Database Configuration

Copy `.env.postgres.example` to `.env.postgres`:
```bash
cp .env.postgres.example .env.postgres
```

Edit `.env.postgres` with your database details:

```dotenv
POSTGRES_USER=postgres
POSTGRES_PASSWORD=yourpassword
POSTGRES_DB=chat_ai
POSTGRES_HOST=localhost
POSTGRES_PORT=5432
```

This file is only used when `DATABASE_TYPE=postgres`.

---

## AI Services Setup

### LLM (Language Model) - OpenRouter

The application currently supports **OpenRouter** for language models.

**Setup Steps:**

1. Go to https://openrouter.ai/
2. Click "Sign Up" and create an account
3. Once logged in, visit https://openrouter.ai/settings/keys
4. Click "Create New Key" to generate an API key
5. Copy the API key
6. Paste it in `.env` file:
   ```dotenv
   AI.LLM_API_KEY=sk-or-xxxxxxxxxxxxx
   ```

### VLM (Vision Language Model)

Currently, only **mock** VLM provider is available. This is a placeholder that returns sample responses for image inputs.

**Configuration:**
```dotenv
AI.VLM_PROVIDER=mock
AI.VLM_API_KEY=mock-no-key-needed
```

Vision model support will be added in future updates.

---

## Authentication

All endpoints except `/auth/register` and `/auth/login` require a JWT token.

**Include token in requests:**
```
Authorization: Bearer <token>
```

---

## API Endpoints

### Register
**POST** `/api/v1/auth/register`

```json
{
  "email": "user@example.com",
  "password": "password123"
}
```

Response:
```json
{
  "accessToken": "jwt_token_here"
}
```

### Login
**POST** `/api/v1/auth/login`

```json
{
  "email": "user@example.com",
  "password": "password123"
}
```

Response:
```json
{
  "accessToken": "jwt_token_here"
}
```

### Chat
**POST** `/api/v1/chat` (requires auth)

Send text and/or images. At least one must be provided.

```json
{
  "message_query": "What's in this image?",
  "image_s": [
    {
      "url": "https://example.com/image.jpg",
      "type": "image/jpeg"
    }
  ],
  "model_config": {
    "llm_model": "llama_70b",
    "vlm_model": "mock"
  }
}
```

**Image Format:**
```json
{
  "url": "string (optional - image URL)",
  "base64": "string (optional - base64 encoded image)",
  "type": "image/jpeg | image/png | image/webp"
}
```

Response:
```json
{
  "id": "uuid",
  "user_id": "uuid",
  "query": "What's in this image?",
  "response_text": "AI response here",
  "image_urls": [],
  "llm_model_name": "llama_70b",
  "vlm_model_name": "mock",
  "timestamp": "2025-12-14T10:30:00Z"
}
```

**Defaults:**
- `llm_model`: `llama_70b`
- `vlm_model`: `mock`

### Chat History
**GET** `/api/v1/chat/history` (requires auth)

Query Parameters:
- `page` (default: 1, min: 1)
- `limit` (default: 10, min: 1, max: 100)
- `order` (default: desc, options: asc, desc)

Example: `/api/v1/chat/history?page=1&limit=20&order=desc`

Response:
```json
{
  "data": [
    {
      "id": "uuid",
      "user_id": "uuid",
      "query": "What's in this image?",
      "response_text": "AI response here",
      "image_urls": [],
      "llm_model_name": "llama_70b",
      "vlm_model_name": "mock",
      "timestamp": "2025-12-14T10:30:00Z"
    }
  ],
  "page": 1,
  "limit": 20,
  "total": 45,
  "total_pages": 3
}
```

## Error Response

```json
{
  "error": "Error type",
  "message": "Detailed error message"
}
```

Common Status Codes:
- `200` - Success
- `201` - Created
- `400` - Bad request
- `401` - Unauthorized (invalid/expired token)
- `500` - Server error

## Project Structure

```
├── Makefile
├── README.md
├── .env.example
├── .env.postgres.example
├── cmd
│   └── main.go
├── config
│   ├── config.go
│   └── config_test.go
├── docker-compose.yaml
├── go.mod
├── go.sum
├── internal
│   ├── database
│   │   ├── database.go
│   │   ├── migrations
│   │   ├── postgres
│   │   └── mock
│   ├── server
│   │   ├── handler
│   │   ├── middleware
│   │   ├── router
│   │   └── manager
│   └── service
├── model
│   ├── dto
│   ├── entity
│   └── paginate.go
├── pkg
│   ├── hash.go
│   ├── jwt.go
│   └── logger
├── static
│   ├── openapi.html
│   └── openapi.json
└── sqlerr
```

## Environment Examples

```dotenv
ENV=local
DATABASE_TYPE=mock
LEVEL=debug
FORMAT=text
SERVER_PORT=8080
JWT_SECRET=dev-secret-key
AI.LLM_PROVIDER=openrouter
AI.LLM_API_KEY=sk-or-xxxxx
AI.VLM_PROVIDER=mock
AI.VLM_API_KEY=mock
```

## Notes

- Token expires after 1 hour - users must login again
- Images can be sent via URL or Base64
- Chat history is user-specific and paginated
- All timestamps are in ISO 8601 format
- Use `DATABASE_TYPE=mock` for testing without a database
- Change `JWT_SECRET` in production to a strong random string