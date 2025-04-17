# 🧩 1337b04rd — Team Project Division Plan

This document outlines the division of work between two teammates for the **1337b04rd** project — an anonymous imageboard implemented using Hexagonal Architecture in Go.

---

## 👤 mkaliyev — Domain & Infrastructure Lead

### 🔹 Responsibilities

#### 🧱 Domain Layer (Core Logic)
- Define core entities: `Post`, `Comment`, `UserSession`
- Define interfaces (ports):
  - `PostRepository`
  - `CommentRepository`
  - `SessionStore`
  - `ImageUploader`
  - `AvatarProvider`

#### 🗄️ PostgreSQL Integration
- Implement PostgreSQL-based adapters:
  - `pgPostRepository`
  - `pgCommentRepository`
  - `pgSessionStore`
- Create and manage SQL schema & migrations
- Ensure use of indexes, constraints, and foreign keys

#### ☁️ S3 Image Storage (MinIO or AWS)
- Configure at least 2 buckets (e.g., `posts-bucket`, `comments-bucket`)
- Validate image types and securely upload files
- Return public URLs for image retrieval

#### ⏱️ Time-Based Deletion Logic
- Implement timers for post expiration:
  - No comments: delete after **10 minutes**
  - With comments: delete after **15 minutes** since last comment
- Use `time.Ticker` or background goroutines for cleanup

#### ✅ Testing (Infrastructure Side)
- Unit test:
  - Repositories
  - ImageUploader adapter
  - Cleanup logic (deletion rules)

---

## 👤 dzhailan — API, Sessions & Integration Lead

### 🔹 Responsibilities

#### 🌐 HTTP Server & Handlers
- Set up HTTP server using `net/http`
- Serve frontend templates:
  - `catalog.html`, `archive.html`, `post.html`, `archive-post.html`, `create-post.html`, `error.html`
- Implement REST API endpoints:
  - Create post
  - Add comment
  - Load posts and comments

#### 🍪 Session & Cookie Management
- Implement secure, HTTP-only cookies
- Manage sessions with 1-week expiration
- Store session-based user identity and metadata

#### 🧑‍🚀 Rick & Morty API Integration
- Fetch unique avatar and name for each new session
- Handle name overrides (custom names replace API name)
- Avoid avatar duplication within a single post
- Reuse avatars if avatar pool is exhausted

#### 💬 Comment & Reply System
- Allow replies to posts and comments
- Show reply relationships via clickable IDs
- Assign user avatars and names to comments

#### 🛠️ Logging + CLI Interface
- Use `log/slog` for structured logging:
  - HTTP requests
  - Errors
  - DB interactions
- Implement CLI `--help` command with output:
    - ./1337b04rd --help


### Suggested Folder Structure
```
1337b04rd/
├── cmd/1337b04rd/         # Main entry point
├── internal/
│   ├── domain/            # Entities + interfaces (ports)
│   ├── usecase/           # Core services and business logic
│   ├── adapter/           # DB, S3, API adapters
│   ├── handler/           # HTTP handlers
│   ├── middleware/        # Session and auth middlewares
│   └── utils/             # Common utilities/helpers
├── templates/             # HTML templates (frontend)
├── s3/                    # S3/MinIO config or setup tools
├── migrations/            # SQL schema and init scripts
├── tests/                 # *_test.go files
└── go.mod
``` 
