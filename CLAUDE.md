# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is an HTMX-based todo list application built with Go, deployed to AWS using Elastic Beanstalk and managed via AWS CDK. The application uses AWS Cognito for authentication and RDS MySQL for data persistence. CloudFront provides CDN capabilities fronting the Elastic Beanstalk environment.

## Architecture

The application consists of three main components:

1. **GoApp/**: Go backend serving HTMX-enhanced HTML templates
   - Entry point: `cmd/server/main.go` (initializes db, auth, router in sequence)
   - Router: Gorilla mux handles routes for login, home, and task CRUD operations (port 5000)
   - Database: MySQL via sqlc for type-safe queries
   - Views: templ templates for server-side HTML generation
   - Auth: AWS Cognito integration for user authentication

2. **cdk/**: AWS CDK TypeScript stack for infrastructure
   - CloudFront distribution with custom domain (gohtmxtodo.grantstarkman.com)
   - Route53 DNS configuration
   - ACM certificate management
   - Cognito User Pool provisioning
   - Origin: Elastic Beanstalk environment (HTTP only on port 80)

3. **Elastic Beanstalk**: Go application hosting (manually configured)
   - Deployment managed via EB CLI
   - Environment must be restarted after deployment

## Development Commands

### Go Application (from GoApp/ directory)

**Local development with hot reload:**
```bash
air
# Watches .go, .templ, .sql files and rebuilds on changes
# Runs sqlc generate, templ generate, then builds to tmp/main
# Access via http://localhost:8090 (proxy) or http://localhost:8080 (direct)
```

**Build for production:**
```bash
TEMPL_EXPERIMENT=rawgo templ generate && go build -o bin/application cmd/server/main.go
```

**Generate database code (after modifying SQL files):**
```bash
sqlc generate
# Reads sqlc.yaml and generates Go code from:
# - internal/db/store/schema/schema.sql (schema)
# - internal/db/store/queries/query.sql (queries)
# Output: internal/db/store/sqlc/
```

**Generate templ templates:**
```bash
templ generate
# Converts .templ files to .templ.go files
# Required before building or running
```

**Install dependencies:**
```bash
go mod tidy
```

### CDK Infrastructure (from cdk/ directory)

**Install dependencies:**
```bash
npm install
```

**Synthesize CloudFormation:**
```bash
cdk synth
```

**Deploy infrastructure:**
```bash
cdk deploy --require-approval never
```

**Compare deployed vs. current state:**
```bash
cdk diff
```

### Full Deployment

**Complete build and deploy (from root):**
```bash
./recursiveBuild.sh
# Executes full pipeline:
# 1. Builds Go app (templ + go build)
# 2. Deploys to Elastic Beanstalk via EB CLI
# 3. Restarts EB environment
# 4. Deploys CDK stack
# 5. Invalidates CloudFront cache
```

## Key Technical Details

### Database Layer

- **sqlc** generates type-safe Go code from SQL
- Schema: `internal/db/store/schema/schema.sql` (tasks table)
- Queries: `internal/db/store/queries/query.sql` (CRUD operations)
- Generated code: `internal/db/store/sqlc/`
- Connection: RDS MySQL via environment variables (RDS_HOSTNAME, RDS_DB_NAME, RDS_USERNAME, RDS_PASSWORD)
- Singleton pattern: `db.Init()` creates global `dbInstance` accessible via `db.GetStore()`

### Authentication Flow

- **AWS Cognito** for user authentication
- Environment variables required: COGNITO_USER_POOL_ID, COGNITO_APP_CLIENT_ID, COGNITO_APP_CLIENT_SECRET
- Region: us-east-1 (hardcoded in auth/auth.go:14)
- App struct in `internal/app/app.go` holds Cognito client and configuration
- Initialized in `auth.Init()` and passed to login handlers

### View Layer

- **templ** for type-safe Go templates (https://templ.guide/)
- Templates located in `internal/views/` with .templ extension
- Must run `templ generate` to create _templ.go files before building
- HTMX attributes in templates enable dynamic updates without full page reloads

### Routing Structure

Routes defined in `internal/router/router.go`:
- `/login` - GET/POST for authentication
- `/home` - GET for main todo page
- `/tasks` - GET (fetch tasks), POST (add task)
- `/completed`, `/incomplete` - GET for filtered views
- `/checked` - POST to toggle task completion
- `/delete/{id}` - POST to delete task
- `/` - Static file server for assets in ./static/

### Infrastructure Details

- CloudFront cache policy: UseOriginCacheControlHeaders-QueryStrings (ID: 4cc15a8a-d715-48a4-82b8-cc0b614638fe)
- Cache invalidation required after deployments: distribution ID EBED81QRRI7NP
- Elastic Beanstalk environment ID: e-r36cz5jntw (region: us-east-2)
- Domain managed in Route53: grantstarkman.com hosted zone

### Environment Variables Required

**Runtime (Elastic Beanstalk):**
- RDS_HOSTNAME, RDS_DB_NAME, RDS_USERNAME, RDS_PASSWORD
- COGNITO_USER_POOL_ID, COGNITO_APP_CLIENT_ID, COGNITO_APP_CLIENT_SECRET

**Build/Deploy:**
- ELASTIC_BEANSTALK_ENVIRONMENT_DOMAIN (for CDK)

## CI/CD

GitHub Actions workflow (.github/workflows/main.yml) runs on push to main:
1. Sets up Go 1.22, Node.js 14, Python 3.10
2. Installs templ, AWS CLI, CDK, EB CLI
3. Generates templ files and installs Go dependencies
4. Executes recursiveBuild.sh for full deployment

## Important Notes

- Always regenerate templ files after modifying .templ templates
- Always regenerate sqlc code after modifying schema or queries
- CloudFront cache must be invalidated after code deployments to see changes
- EB environment restart is necessary for application updates to take effect
- air.toml configures hot reload for local development (excludes _test.go files)
