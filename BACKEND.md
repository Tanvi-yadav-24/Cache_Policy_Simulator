# Cache Simulator Backend

Complete Go backend with REST API, database persistence via Supabase, and fallback to local TypeScript implementations.

## Architecture

```
Frontend (React + TypeScript)
    |
    |-- API calls to Go backend (if available)
    |-- Falls back to local TS implementations (if backend unavailable)
    |
Go Backend (Gin + Supabase)
    |
    |-- REST API: /api/simulate, /api/compare, /api/generate, /api/benchmark
    |-- History: /api/history/simulations, /api/history/comparisons, /api/history/benchmarks
    |
Supabase Database
    |-- simulations, comparisons, benchmarks, traces tables
```

## Database Schema

Tables created in Supabase:

- `simulations` - Individual policy simulation runs
- `comparisons` - Multi-policy comparison runs
- `benchmarks` - Benchmark runs across cache sizes
- `traces` - Generated request traces

All tables have RLS enabled with `anon, authenticated` access for this single-tenant demo app.

## API Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/health` | Health check |
| GET | `/api/policies` | List cache policies |
| POST | `/api/simulate` | Run single simulation |
| POST | `/api/generate` | Generate request trace |
| POST | `/api/compare` | Compare all policies |
| POST | `/api/benchmark` | Run benchmark |
| GET | `/api/history/simulations` | List recent simulations |
| GET | `/api/history/comparisons` | List recent comparisons |
| GET | `/api/history/benchmarks` | List recent benchmarks |

## Running Locally

### 1. Prerequisites

- Go 1.21+
- Supabase project (already provisioned)

### 2. Set Environment Variables

```bash
export SUPABASE_URL="https://your-project.supabase.co"
export SUPABASE_SERVICE_ROLE_KEY="your-service-role-key"
# or
export VITE_SUPABASE_URL="https://your-project.supabase.co"
export VITE_SUPABASE_ANON_KEY="your-anon-key"
```

### 3. Install Go Dependencies

```bash
go mod tidy
```

### 4. Run the Backend

```bash
go run cmd/server/main.go
```

Server starts on port 8080 (or PORT env var).

### 5. Run the Frontend (in another terminal)

```bash
npm install
npm run dev
```

The frontend will proxy `/api` requests to `http://localhost:8080` automatically via Vite config.

## Frontend-Backend Connection

The frontend automatically detects if the backend is running:

- **Green "Backend Connected"**: All operations go to Go backend + database
- **Amber "Local Mode"**: Frontend uses local TypeScript implementations (no persistence)

The frontend tries the backend first, and if unavailable, falls back to local implementations seamlessly.
