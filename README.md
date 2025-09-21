# Minify

A simple, fast URL shortener written in Go with PostgreSQL and Prometheus monitoring.


## Features
- Minifies (shortens) URLs quickly with custom codes 
- User authorization (registration + login, JWT-based)  
- Tracks user metrics such as clicks, popular URLs, and usage stats with Prometheus  


## Local Setup

### Prerequisites

- Go 1.21+  
- PostgreSQL 12+  


### Quick Start

1. **Clone the repo**
```bash
git clone https://github.com/sameersaeed/minify.git
cd minify
```

2. **Install dependencies**
```bash
go mod tidy
```

3. **Setup PostgreSQL**
```bash
sudo -u postgres createdb minify
```

4. **Configure environment**
- For the backend, from the project root you can modify `.env`
- For the frontend, from it's project root (`frontend/`), you can modify `.env.local` to point to the backend's URL if you modified it

5. **Run backend**
```bash
make run
```
or
```bash
go run main.go
```

Or build the binary (from project root):
```bash
make
./minify
```
or
```bash
go build -o minify .
./minify
```

6. **Frontend**
```bashnpm install
npm run dev   # for development
npm run build # build for production
npm start     # start frontend
```


## API Endpoints
| Route                                          | What it does        |
|------------------------------------------------|---------------------|
| `POST /api/v1/users`                           | register            |
| `POST /api/v1/users/login`                     | login               |
| `POST /api/v1/minify`                          | create minified URL |
| `GET /api/v1/urls?user_id=X`                   | get user URLs       |
| `GET /{shortCode}`                             | redirect            |
| `GET /api/v1/analytics/overview`               | usage overview      |
| `GET /api/v1/analytics/popular`                | popular URLs        |
| `GET /api/v1/analytics/timeframe/{period}`     | timeframe stats     |
| `GET /metrics`                                 | Prometheus metrics  |
| `GET /health`                                  | health check        |


## Environment variables

| Variable       | Default                               | Description                      |
|----------------|---------------------------------------|----------------------------------|
| PORT           | 8080                                  | Backend port                     |
| DATABASE_URL   | postgres://...                        | PostgreSQL connection string     |
| BASE_URL       | http://localhost:8080                 | Base URL for short links         |
| JWT_SECRET     | your-secret-key                       | JWT signing secret               |


### Production
Update `.env` with real values based on your hosted server's endpoint:
```
env
BASE_URL=https://yourdomain.com
DATABASE_URL=postgres://user:pass@prod-db:5432/minify?sslmode=require
JWT_SECRET=<your-secret-key>
```

Build and run:
```bash
make
./minify
```
or
```bash
go build -o minify .
./minify
```


## Architecture
(Going to insert an image here)
- HTTP requests, clicks, and usage stats are handled by the Go backend  
- URL & user data are stored in PostgreSQL  
- Prometheus collects metrics from backend for monitoring
