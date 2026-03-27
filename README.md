# Visea Go Boilerplate

## Project Structure

```text
├── cmd/
│   └── server/          # Main application entrypoint (`main.go`)
├── internal/            # Private application code
│   ├── config/          # Environment configuration loading
│   ├── handler/         # HTTP request handlers (Controllers)
│   ├── middleware/      # Custom HTTP middlewares (e.g., auth, cors)
│   ├── model/           # Database models and data structures
│   ├── repository/      # Data access layer for DB operations
│   ├── request/         # Request DTOs and validation rules
│   ├── response/        # Response serialization structs
│   ├── route/           # Application HTTP routes registration
│   └── service/         # Core business logic layer
├── pkg/                 # Public/shared utility packages
│   ├── datatable/       # Data filtering and pagination helpers
│   ├── helpers/         # General helper functions
│   ├── logger/          # Custom logging implementation
│   ├── mail/            # Email sending notification service
│   ├── messages/        # Standardized API response messages
│   └── notifier/        # Slack/webhook notification dispatchers
├── .air.toml            # Air live-reload configuration
├── .env.example         # Example environment variables
├── go.mod               # Go module dependencies
└── README.md            # Project documentation
```

## Use Case: Product CRUD

This boilerplate includes a complete Product CRUD implementation to demonstrate the architectural pattern (Handler -> Service -> Repository -> Model).

### Architectural Layers
1. **Handler (`internal/handler/`)**: Handles HTTP routing (Gin) and request parsing.
2. **Service (`internal/service/`)**: Contains business logic and orchestrates data between handlers and repositories.
3. **Repository (`internal/repository/`)**: Manages database operations using GORM.
4. **Model (`internal/model/`)**: Defines the database schema and GORM tags.
5. **DTOs (`internal/request/` & `internal/response/`)**: Decouples the API contract from the database schema.

### Key Endpoints
- `POST /products`: Create a new product.
- `GET /products`: Retrieve all products.
- `GET /products/:id`: Get detailed information for a single product.
- `PUT /products/:id`: Update an existing product.
- `DELETE /products/:id`: Perform a soft delete of a product.
- `GET /products/datatable`: **Advanced Pagination** using the `datatable` package. Supports server-side searching, sorting, and filtering.

## System Features

### 1. Consistent Response Structure
All API responses use a standardized JSON envelope implemented in `pkg/helpers/response.go`:
```json
{
  "status": true,
  "message": "Product successfully retrieved",
  "data": { ... }
}
```

### 2. Multi-language (I18n) Translation
The system supports dual-language responses (Indonesian/English) via the `pkg/messages/` package.
- **Header Controlled**: Use `Accept-Language: id` or `Accept-Language: en` in your request headers.
- **Implementation**: Handlers use `middleware.GetLang(c)` and `messages.Translate()` to deliver localized feedback based on the user's preference.

### 3. Async Notifier System
A decoupled notification system located in `pkg/notifier/` for sending alerts to Slack, webhooks, or other providers.
- **Non-Blocking**: Uses the `Async` wrapper (goroutines) to ensure notifications never delay API responses.
- **Flexible Levels**: Supports `Info`, `Warning`, `Error`, and `Critical` levels with automatic emoji prefixing.
- **Providers**: Easily switch between `WebhookNotifier` (production) and `NoOpNotifier` (local development).

## Prerequisites

- **Go** (version 1.25.0 or later)
- **PostgreSQL** (Active database instance)
- **Redis** (Active cache instance)
- **[Air](https://github.com/cosmtrek/air)** (Optional) - for hot-reloading during development

## How to Start the Project

1. **Install Dependencies:**
   Ensure all Go module dependencies are downloaded and verified:
   ```bash
   go mod tidy
   ```

2. **Configure Environment Variables:**
   Copy the provided example environment template:
   ```bash
   cp .env.example .env
   ```
   *Note: Open `.env` and fill in correct values for `DB_HOST`, `DB_USER`, `DB_PASSWORD`, `REDIS_HOST`, etc. matching your local/remote instances.*

3. **Run the Application:**
   Start the application running via standard Go command:
   ```bash
   go run ./cmd/server/main.go
   ```

   **With Hot Reloading:**
   If you have `air` installed, run at the project root for automatic recompilation on file changes:
   ```bash
   air
   ```
   The application will start, usually accessible on `http://localhost:8080`, depending on the `APP_PORT` set in your `.env`.
