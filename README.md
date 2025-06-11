# Server Template

A microservice project template based on the [Kratos](https://github.com/go-kratos/kratos) framework, integrating common middleware and tools, providing a complete development process and best practices.

## Features

- Based on Kratos v2 framework, providing a complete microservice architecture
- Using Protocol Buffers to define APIs, supporting gRPC and HTTP protocols
- Integrated MySQL master-slave read-write separation
- Integrated Redis cache
- Using Wire for dependency injection
- Using sqlc to generate type-safe database access code
- Using buf tool to manage Protocol Buffers
- Support for database transactions and error handling
- Integrated logging system
- Docker deployment support

## Tech Stack

- Go 1.24+
- Kratos v2
- Protocol Buffers
- gRPC/HTTP
- MySQL
- Redis
- Wire
- sqlc
- buf
- Docker

## Project Structure

```
.
├── api                 # API definitions
│   └── server          # Service API definitions
├── cmd                 # Application entries
│   └── server          # Service entry
├── internal            # Internal code
│   ├── biz             # Business logic layer
│   ├── conf            # Configuration definitions
│   ├── data            # Data access layer
│   │   ├── migration   # Database migrations
│   │   ├── queries     # Generated query code
│   │   └── query       # SQL query definitions
│   ├── server          # Server implementations
│   └── service         # Service interface implementations
├── pkg                 # Public packages
├── Dockerfile          # Docker build file
├── Makefile            # Build scripts
├── config.yaml         # Configuration file
└── sqlc.yaml           # sqlc configuration
```

## Development Environment Setup

### Install Dependencies

```bash
# Initialize development environment, install required tools
make init
```

This will install the following tools:
- kratos CLI
- protoc-gen-go-http
- wire
- sqlc
- buf
- protoc-gen-openapi
- protoc-gen-go
- protoc-gen-go-grpc
- migrate

## Development Guide

### Generate Code

```bash
# Generate API related code
make api

# Generate configuration related code
make config

# Generate all code
make all

# Generate dependency injection code
make wire

# Generate database access code
make sqlc
```

### Create Database Migration

```bash
# Create a new migration file
make new_migration name=migration_name
```

### Run Service

```bash
# Run the service
make server
```

### Build Project

```bash
# Build the project
make build
```

## Docker Deployment

```bash
# Build Docker image
docker build -t <your-docker-image-name> .

# Run Docker container
docker run --rm -p 8000:8000 -p 9000:9000 -v </path/to/your/configs>:/data/conf <your-docker-image-name>
```

## API Documentation

The project automatically generates an OpenAPI specification file `openapi.yaml`, which can be viewed using tools like Swagger UI.

## Database Design

The project includes the following database tables:

- `user`: Basic user information
- `user_detail`: Detailed user information

## Configuration

The configuration file is located at `config.yaml`, containing the following main configurations:

- Service configuration (HTTP/gRPC)
- Database configuration (master-slave)
- Redis configuration
- Logging configuration

## License

[MIT License](LICENSE)

