# 🌐 API Conventions

## 1. REST & Content Types
- **General:** The API follows REST principles.
- **JSON Default:** Requests and Responses default to `application/json`.
- **File Uploads:**
  - **Request:** Use `multipart/form-data`.
  - **Response:** Return a JSON object containing the `file_path` and metadata. Do NOT return the binary content directly.

## 2. Response Envelopes
All endpoints must return a consistent structure:

```json
// Success
{
  "data": { ... }
}

// Error
{
  "error": {
    "code": "RESOURCE_NOT_FOUND",
    "msg": "User not found"
  }
}
```

## 3. Middleware Stack
Requests flow through this pipeline:
1.  **Recovery:** Catch panics -> 500 JSON.
2.  **RequestID:** Propagate `X-Request-ID`.
3.  **Logger:** Structured logging (slog).
4.  **Auth:** Validate JWT -> Inject `UserClaims` into Context.
5.  **CORS:** Strict policy for authorized origins.

## 4. Context usage
- **Data Passing:** Use `context.Context` to pass Request-scoped data (UserID, TenantID) from Middleware to Handlers.
- **Keys:** Use custom types for Context keys to avoid collisions.
