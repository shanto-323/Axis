# Axis - A Backend-Focused Web-Interface for interacting with multiple llm agent

## Quick Links

- **API Docs:** `http://localhost:8080/docs`
- **Base URL:** `http://localhost:8080/api/v1` <br><br>
- **WEB:**      `http://localhost:8000` 

## Tech Stack

- **Frontend:**  Go • PostgreSQL • Docker • Nginx
- **Backend:**  React • Type-Script

## Configuration

- **Make sure to create a .env and .env.postgres**

### 2. `.env` - Server Configuration

```dotenv
# Server Configuration
SERVER.PORT=8080
SERVER.CORS_ALLOWED_ORIGINS=http://localhost:8000

# JWT 
SERVER.JWT_KEY=secret_key   

# AI Services
AI_MANAGER.PROVIDER=https://openrouter.ai/api/v1
AI_MANAGER.API_KEY=api_key
```


### 2. `.env.postgres` - Database Configuration

```dotenv
POSTGRES_USER=postgres
POSTGRES_PASSWORD=yourpassword
POSTGRES_DB=chat_ai
```
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
   AI_MANAGER.API_KEY=sk-or-xxxxxxxxxxxxx
   ```

## Authentication

All endpoints except `/auth/register` and `/auth/login` require a JWT token as cookie

```
Key:<token>
```

## API Endpoints

### Register
**POST** `/api/v1/auth/register`

```json
{
  "email": "user@example.com",
  "password": "password123"
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

### Chat
**POST** `/api/v1/chat` (requires auth)

```json
{
  "message": "hello",
  "model": "llama-70b"
}
```

Response:
```json
{
  "id": "95539e01-21fc-44ca-9540-00d314ae0b12",
  "llm_model_name": "llama-70b",
  "timestamp": "2025-12-28T18:51:53.391628Z",
  "user_id": "2be4cf6b-4b5b-43fa-9bed-ad51911cefcf",
  "query": "hello",
  "response_text": "Hello! How can I assist you today?"
}
```

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
      "id": "string",
      "llm_model_name": "string",
      "timestamp": "2026-01-05T15:29:14.023Z",
      "user_id": "string",
      "query": "string",
      "response_text": "string"
    }
  ],
  "page": 1,
  "limit": 10,
  "total": 1,
  "total_pages": 1
}
```

## Error Response

```json
{
  "error": "Error type",
  "message": "Detailed error message"
}
```
