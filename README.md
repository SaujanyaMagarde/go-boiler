# Go RSS Feed Aggregator Backend

This project is a powerful RSS feed aggregator backend built in Go. It allows users to register, add RSS feeds, follow feeds, and aggregates posts automatically using background goroutines.

This `README.md` is designed not just as project documentation, but as a **comprehensive revision guide** for building a production-ready API in Go. It covers everything from setting up the environment to understanding the core packages and architecture.

---

## 🚀 Features

- **User Management**: Create users and authenticate using API Keys.
- **Feed Management**: Add new RSS feeds to the database.
- **Feed Following**: Users can follow and unfollow feeds.
- **Background Scraper**: A concurrent goroutine constantly fetches new posts from the followed RSS feeds and stores them.
- **RESTful API**: Standardized JSON responses and error handling using the `chi` router.
- **Database Integration**: Postgres database driven by type-safe queries generated via `sqlc`.

---

## 🛠️ Tech Stack & Packages Used

This project relies on the standard library and a few highly-regarded open-source packages:

### 1. `github.com/go-chi/chi`
- **What it is**: A lightweight, idiomatic, and composable router for building Go HTTP services.
- **Why we use it**: It provides powerful routing (like `/v1/...`), URL parameters extraction (`chi.URLParam`), and excellent middleware support (like CORS and Auth) without reinventing the standard `net/http` wheel.

### 2. `github.com/go-chi/cors`
- **What it is**: Cross-Origin Resource Sharing (CORS) middleware for `chi`.
- **Why we use it**: To allow frontend applications (running on different ports/domains) to securely make requests to our API.

### 3. `github.com/lib/pq`
- **What it is**: A pure Go Postgres driver for the standard `database/sql` package.
- **Why we use it**: It registers the "postgres" driver under the hood, allowing `sql.Open("postgres", dbURL)` to establish our database connection.

### 4. `github.com/joho/godotenv`
- **What it is**: A port of the Ruby dotenv project.
- **Why we use it**: It reads the `.env` file in the root directory and loads those variables into the OS environment so we can access them securely via `os.Getenv()`.

### 5. `github.com/google/uuid`
- **What it is**: A library to generate and inspect UUIDs based on RFC 4122.
- **Why we use it**: We use UUIDs as primary keys for our database entities (Users, Feeds, Posts) for better security and distributed scalability.

### 6. `sqlc` & `goose` (CLI Tools)
- **`sqlc`**: Generates fully type-safe Go code from pure SQL queries (found in `sql/queries`). It outputs to `internal/database`, keeping our data access layer clean and bug-free.
- **`goose`**: A database migration tool. We use it to manage our Postgres schema (found in `sql/schema`) systematically.

---

## 📚 Project Architecture & Revision Guide

When reviewing this project, keep the following Go patterns in mind:

### 1. Directory Structure
- `main.go`: The entry point. Initializes the DB, router, background worker, and starts the HTTP server.
- `internal/database/`: Auto-generated code by `sqlc`. **Do not edit manually.**
- `sql/schema/`: Contains `.sql` migration files to create/alter tables (e.g., users, feeds, posts).
- `sql/queries/`: Contains raw SQL queries. `sqlc` reads these and generates Go functions.
- `handler_*.go`: Contains the HTTP handler functions for our endpoints.
- `models.go`: Contains functions to convert raw database structs into JSON-friendly structs.

### 2. Dependency Injection via Struct Methods
Instead of using global variables for the database connection, we use a struct to hold our state:
```go
type apiConfig struct {
	DB *database.Queries
}
```
All HTTP handlers are methods on `apiConfig`. This allows every handler (e.g., `apiCfg.handlerCreateUser`) to access the `DB` connection safely.

### 3. Middleware and Authentication
We created a custom middleware `middelwareAuth` in `middelware_auth.go`. 
It intercepts the request, extracts the `Authorization: ApiKey ...` header, verifies the user from the database, and injects the `database.User` object into the handler. This keeps our handlers clean from repetitive auth logic!

### 4. Concurrency (Goroutines & Channels)
In `main.go`, we spin up a background worker:
```go
go startScrapping(apiCfg.DB, 10, time.Minute)
```
The `go` keyword spins off `startScrapping` into its own lightweight thread (goroutine). It periodically checks the database for outdated feeds, fetches their XML via HTTP, parses the RSS, and saves new posts back to the database concurrently without blocking our main API server!

### 5. HTTP Handlers (`w` and `r`)
In Go, an HTTP handler function always takes two parameters: a place to put the incoming request (`r *http.Request`) and a tool to write the response back (`w http.ResponseWriter`).
- **`r *http.Request`**: This holds everything the client sent you.
  - `r.Method`: Is it a `POST`? A `GET`?
  - `r.Body`: The actual JSON data they sent (e.g., `{"name": "saujanya"}`).
  - `r.Header`: Contains security tokens, API keys, or browser information.
- **`w http.ResponseWriter`**: This is how you reply to the client.
  - `w.WriteHeader(201)`: Tells the client the status code (e.g., "201 Created").
  - `w.Write([]byte("User created!"))`: Sends the actual text or JSON payload back to the screen.

### 6. Goroutines vs. OS Threads
When we initialize our database connection (`sql.Open`), Go automatically creates a connection pool managed by goroutines.
- **Why Goroutines?** Traditional OS threads consume roughly **1MB of memory per thread**. Goroutines, on the other hand, consume only **~1KB to 2KB per routine**. This means a Go server can handle thousands (or even millions) of concurrent tasks (like database queries or HTTP requests) simultaneously without crashing or running out of memory!

### 7. The Custom HTTP Server
In `main.go`, we explicitly define a custom server object:
```go
srv := &http.Server{
    Handler: router, // Gives requests to the chi router to figure out what to do.
    Addr:    ":" + portString, // Tells the server which port to listen on.
}
err = srv.ListenAndServe() // This blocks and constantly waits for network requests.
```

### 8. How Middleware *Actually* Works
Sometimes, before executing an HTTP handler, we want to extract user info or check authentication. Instead of writing that logic in *every single handler*, we use **Middleware**.
- **Input**: A middleware function takes a standard HTTP handler function that we *want* to run.
- **Output**: It returns a *brand new* standard HTTP function to the `chi` router.
- **Under the hood**: The `chi` router calls this new function. The middleware runs its logic (e.g., checking the `Authorization` header), and if everything is okay, it executes the original input function. If not, it writes an error response using `w` and stops right there.

---

## 💻 Step-by-Step Setup Guide

Follow these steps to recreate the environment from scratch.

### Step 1: Prerequisites
- Install **Go** (v1.26 or latest)
- Install **PostgreSQL** and ensure it's running.
- Install `sqlc`: `go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest`
- Install `goose`: `go install github.com/pressly/goose/v3/cmd/goose@latest`

### Step 2: Clone & Install Dependencies
```bash
git clone <your-repo-url>
cd go-server
go mod download
```

### Step 3: Setup Environment Variables
Create a `.env` file in the root directory:
```env
PORT=8080
DB_URL="postgres://<username>:<password>@localhost:5432/rssagg?sslmode=disable"
```
*(Replace `<username>` and `<password>` with your Postgres credentials, and make sure to create an empty database named `rssagg` in Postgres first).*

### Step 4: Run Database Migrations
We need to set up our tables (Users, Feeds, Posts, etc.). Use `goose` to run the migrations located in the `sql/schema` directory.
```bash
cd sql/schema
goose postgres "postgres://<username>:<password>@localhost:5432/rssagg?sslmode=disable" up
cd ../..
```

### Step 5: Generate Database Code (Optional/If you modify SQL)
If you ever change the queries in `sql/queries` or the schema, run:
```bash
sqlc generate
```
This updates the `internal/database` folder.

### Step 6: Start the Server
```bash
go run .
# or
go build -o go-server
./go-server
```
You should see: `server is running on port: 8080`.

---

## 🧪 API Endpoints & Testing

You can use tools like Postman, Insomnia, or Thunder Client to test these endpoints. 

*Note: For authenticated routes, you must provide the header: `Authorization: ApiKey <your_api_key>`*

### 1. General
- **GET** `/v1/healthz`: Check if the server is alive.
- **GET** `/v1/error`: Test the error response format.

### 2. Users
- **POST** `/v1/users`: Create a new user.
  - Body: `{ "name": "Your Name" }`
  - *Returns your API Key. Keep it safe!*
- **GET** `/v1/users`: Get your user details. *(Requires Auth)*

### 3. Feeds
- **POST** `/v1/createfeed`: Add a new RSS feed to the system. *(Requires Auth)*
  - Body: `{ "name": "Go Blog", "url": "https://go.dev/blog/feed.atom" }`
- **GET** `/v1/getfeeds`: List all feeds in the system. *(Requires Auth)*
- **GET** `/v1/getfeed`: Get a specific feed. *(Requires Auth)*

### 4. Feed Follows
- **POST** `/v1/feed_follows`: Follow a specific feed. *(Requires Auth)*
  - Body: `{ "feed_id": "<uuid-of-feed>" }`
- **GET** `/v1/feed_follows`: Get a list of feeds you are currently following. *(Requires Auth)*
- **DELETE** `/v1/feed_follows/{feedFollowID}`: Unfollow a feed. *(Requires Auth)*

### 5. Posts
- **GET** `/v1/posts`: Retrieve a list of the latest aggregated posts from the feeds you follow. *(Requires Auth)*
  - Query Param (Optional): `?limit=10`

---
*Happy Coding! Use this project as a sandbox to experiment with Go concurrency, middleware, and database interactions.*
