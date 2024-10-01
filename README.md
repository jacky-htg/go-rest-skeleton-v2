# Go REST API Skeleton
A robust and scalable RESTful API skeleton built with Go, featuring essential tools and practices for modern web applications.

## Features
- Concurrency Limit: Control the maximum number of concurrent requests.
- Rate Limiter: Protect your API from abuse by limiting request rates.
- JWT Authentication: Secure your API with JSON Web Tokens.
- RBAC Authorization: Implement role-based access control for fine-grained permissions.
- Dependency Injection Pattern: Promote modular and testable code.
- Structured Logging: Enhanced logging for errors and information.
- Environment Configuration: Option to use OS environment variables or a .env file for configuration.
- Redis Caching: Improve performance with caching.
- Graceful Shutdown: Ensure all requests complete before shutting down the server.
- CORS Handling: Manage Cross-Origin Resource Sharing.
- Clean Architecture: Maintainable and organized code structure.
- Panic Recovery Handling: Safeguard against server crashes.
- Context Error Handling: Manage request timeouts and cancellations.
- Database Migrations: Version control your database schema.
- API Testing: Ensure your API functions as expected.
- Swagger Documentation: Auto-generate API documentation for easy reference.
- Log Monitoring with Loki: Monitor logs efficiently.
- Tracing with OpenTelemetry: Track and analyze performance with Jaeger and otel-collector.
- Business Metrics with OpenTelemetry: Collect metrics relevant to business logic.
- Common Golang Metrics with Prometheus: Utilize Prometheus for golang server metrics.
- Idempotent Request Handling: (To be implemented) Ensure repeated requests yield the same result.
- Docker Support: Pre-configured Dockerfile for easy deployment.
- CI/CD Integration with GitHub Actions: (To be implemented) Streamline your deployment process.
- Example CRUD Operations: Included examples for user management and authentication.

## Getting Started
### Prerequisites
- Go version 1.23.1 or later
- Docker (optional, for containerization)
- PostgreSQL for the database

### Installation
1. Clone the repository:

```bash
git clone https://github.com/yourusername/go-rest-api-skeleton.git
cd go-rest-api-skeleton
```

2. Install dependencies:

```bash
go mod tidy
```

3. Create a .env file or set environment variables based on the provided configuration template.

4. Run the application:

```bash
go run main.go
```

### API Documentation
API documentation is automatically generated and can be accessed at http://localhost:8081/swagger/doc.json.

## Running Tests
To run the API tests, use the following command:

```bash
go test ./...
```

## Contributing
Contributions are welcome! Please fork the repository and submit a pull request for any enhancements or fixes.

## License
This project is licensed under the MIT License. See the [LICENSE](./LICENSE) file for details.