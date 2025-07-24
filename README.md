# leiia-api

A Go HTTP API that receives files via POST endpoint and extracts readable text content. Files and parsed content are stored in PostgreSQL database using GORM.

## Features

- HTTP POST endpoint for file uploads
- File text extraction using `ledongthuc/pdf` library
- PostgreSQL database integration for file and content storage
- GORM ORM for database operations
- Input validation (file type, size limits)
- JSON responses with extracted text content
- Retrieve previously uploaded files by ID
- Docker support for easy development setup

## Setup

### Prerequisites

1. Go 1.24+ installed
2. **Option A**: Docker and Docker Compose (recommended for development)
3. **Option B**: PostgreSQL database server running locally

### Option A: Docker Setup (Recommended)

This is the easiest way to get started. Docker will handle the PostgreSQL setup automatically.

1. **Install Docker and Docker Compose**:
   - [Docker Desktop](https://www.docker.com/products/docker-desktop/) (includes Docker Compose)
   - Or install Docker and Docker Compose separately

2. **Clone the repository**:
   ```bash
   git clone <repository-url>
   cd leiia-api
   ```

3. **Setup environment variables**:
   ```bash
   cp env.example .env
   # Edit .env file and set your database password
   nano .env  # or use your preferred editor
   ```

4. **Start PostgreSQL with Docker**:
   ```bash
   # Start only PostgreSQL
   docker-compose up -d postgres
   
   # Or start PostgreSQL with pgAdmin (optional database management UI)
   docker-compose --profile admin up -d
   ```

5. **Verify database is running**:
   ```bash
   docker-compose ps
   docker-compose logs postgres
   ```

6. **Install Go dependencies and run the API**:
   ```bash
   go mod download
   go run main.go
   ```

### Option B: Local PostgreSQL Setup

If you prefer to run PostgreSQL locally without Docker:

1. **Install PostgreSQL**:
   ```bash
   # macOS with Homebrew
   brew install postgresql
   brew services start postgresql
   
   # Ubuntu/Debian
   sudo apt-get install postgresql postgresql-contrib
   sudo systemctl start postgresql
   ```

2. **Create a database and user**:
   ```sql
   sudo -u postgres psql
   CREATE DATABASE leiia;
   CREATE USER leiia_user WITH PASSWORD 'your_password';
   GRANT ALL PRIVILEGES ON DATABASE leiia TO leiia_user;
   \q
   ```

3. **Setup environment variables**:
   ```bash
   cp env.example .env
   # Edit .env with your database configuration
   ```

4. **Run the application**:
   ```bash
   go mod download
   source .env  # Load environment variables
   go run main.go
   ```

## Docker Commands

### Basic Operations
```bash
# Start PostgreSQL
docker-compose up -d postgres

# Start PostgreSQL + pgAdmin
docker-compose --profile admin up -d

# View logs
docker-compose logs postgres
docker-compose logs pgadmin

# Stop services
docker-compose down

# Stop and remove volumes (⚠️ This will delete all data)
docker-compose down -v
```

### Database Management
```bash
# Access PostgreSQL CLI
docker-compose exec postgres psql -U postgres -d leiia

# Backup database
docker-compose exec postgres pg_dump -U postgres leiia > backup.sql

# Restore database
cat backup.sql | docker-compose exec -T postgres psql -U postgres -d leiia
```

### pgAdmin Access
If you started with the `admin` profile, you can access pgAdmin at:
- URL: http://localhost:5050
- Email: admin@leiia.com
- Password: admin123

To connect to the database in pgAdmin:
- Host: postgres (container name)
- Port: 5432
- Database: leiia (or your DB_NAME from .env)
- Username: postgres (or your DB_USER from .env)
- Password: (your DB_PASSWORD from .env)

## Environment Variables

Copy `env.example` to `.env` and update the values:

```bash
cp env.example .env
```

Required variables:
- `DB_HOST`: Database host (default: localhost)
- `DB_PORT`: Database port (default: 5432)
- `DB_USER`: Database username (default: postgres)
- `DB_PASSWORD`: Database password (required)
- `DB_NAME`: Database name (default: leiia)
- `DB_SSLMODE`: SSL mode (default: disable)

**Note**: When using Docker, set `DB_HOST=localhost` in your `.env` file (the Docker service will map to localhost).

## Application Startup

The application will:
1. Connect to PostgreSQL using the environment variables
2. Auto-migrate the database schema (creates tables if they don't exist)
3. Start the HTTP server on port 8080

Example startup output:
```
Database connection established and migrations completed
Server starting on port 8080...
File upload endpoint: POST http://localhost:8080/file/upload
Get file endpoint: GET http://localhost:8080/file/{id}
```

## Database Schema

The application creates a `files` table with the following structure:

```sql
CREATE TABLE files (
    id SERIAL PRIMARY KEY,
    filename VARCHAR NOT NULL,
    original_name VARCHAR NOT NULL,
    file_size BIGINT NOT NULL,
    content_type VARCHAR NOT NULL,
    file_data BYTEA NOT NULL,        -- Stores actual file content
    parsed_text TEXT,                -- Extracted text content (nullable)
    parse_error TEXT,                -- Error message if parsing failed (nullable)
    uploaded_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);
```

## API Endpoints

### POST /file/upload

Upload a file and extract its text content.

**Request:**
- Method: `POST`
- Content-Type: `multipart/form-data`
- Form field: `file` (file)

**Response:**
```json
{
  "id": 1,
  "message": "File uploaded and parsed successfully",
  "filename": "document.pdf",
  "size": 12345,
  "content": "Extracted text content here...",
  "uploaded_at": "2024-01-15T10:30:00Z"
}
```

### GET /file/{id}

Retrieve a previously uploaded file's information and content.

**Request:**
- Method: `GET`
- URL: `/file/{id}` where `{id}` is the file's database ID

**Response:**
```json
{
  "id": 1,
  "message": "File found and parsed successfully",
  "filename": "document.pdf", 
  "size": 12345,
  "content": "Extracted text content here...",
  "uploaded_at": "2024-01-15T10:30:00Z"
}
```

### GET /health

Health check endpoint that also verifies database connectivity.

**Response:**
```json
{
  "status": "healthy",
  "database": "connected"
}
```

## Testing with Postman

### Upload File
1. Create a new POST request to `http://localhost:8080/file/upload`
2. Set Body to `form-data`
3. Add a key named `file` with type `File`
4. Select your file
5. Send the request
6. Note the `id` in the response for retrieval

### Retrieve File
1. Create a new GET request to `http://localhost:8080/file/{id}`
2. Replace `{id}` with the ID from the upload response
3. Send the request

## Testing with cURL

### Upload
```bash
curl -X POST \
  -F "file=@/path/to/your/file.pdf" \
  http://localhost:8080/file/upload
```

### Retrieve
```bash
curl http://localhost:8080/file/1
```

## File Storage

- Files are stored as binary data (`BYTEA`) in PostgreSQL
- No temporary files are created on the filesystem
- Both original file content and parsed text are stored in the database
- Files can be retrieved by their database ID

## Error Handling

The API handles various error cases:
- Invalid file types (non-PDF)
- File size exceeding limits (10MB)
- Database connection failures
- File parsing errors (file is still saved but parsing failure is recorded)
- Record not found for retrieval requests

## Development

### Project Structure
```
leiia-api/
├── main.go                 # Main application file
├── docker-compose.yml      # Docker services configuration
├── env.example            # Environment variables template
├── .env                   # Your environment variables (not in git)
├── docker/
│   └── postgres/
│       └── init/
│           └── 01-init-db.sql  # Database initialization script
├── go.mod                 # Go module dependencies
├── go.sum                 # Go module checksums
└── README.md              # This file
```

### Troubleshooting

1. **Database connection issues**:
   ```bash
   # Check if PostgreSQL is running
   docker-compose ps
   
   # Check PostgreSQL logs
   docker-compose logs postgres
   
   # Test database connection
   docker-compose exec postgres psql -U postgres -d leiia -c "SELECT 1;"
   ```

2. **Port conflicts**:
   - If port 5432 is already in use, change `DB_PORT` in your `.env` file
   - Update the port mapping in `docker-compose.yml` if needed

3. **Permission issues**:
   ```bash
   # Reset Docker volumes if needed
   docker-compose down -v
   docker-compose up -d postgres
   ```

## Configuration

- Maximum file size: 10MB
- Supported format: PDF only
- Content-Type validation: `application/pdf`
- Database: PostgreSQL with GORM ORM
- Default ports: API (8080), PostgreSQL (5432), pgAdmin (5050)