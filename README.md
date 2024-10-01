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
- Idempotent Request Handling: Ensure repeated requests yield the same result.
- Docker Support: Pre-configured Dockerfile for easy deployment.
- CI/CD Integration with GitHub Actions: Streamline your deployment process.
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

4. run `cd docker && docker-compose up -d` to running monitoring tools.

5. Run the application: From your root app directory, run the command:

```bash
go run main.go
```

### API Documentation
API documentation is automatically generated and can be accessed at http://localhost:8081/swagger/doc.json.

if you want to login using seed data, you can try with this payload:
```json
{
    "email": "rijal.asep.nugroho@gmail.com",
    "password": "qwertyuiop!1Q"
}
```
  
## Folder Structure
The folder structure is organized to follow the principles of Clean Architecture, ensuring that the application remains maintainable and scalable:

```bash
go-rest-api-skeleton/
├── cmd/                # CLI command for migration and etc
├── docker/             # docker configuration
├── docs/               # Generated APi Doc for Swagger
├── internal/           # Application internals (domain logic, services)
│   ├── dto/            # Data transfer object to transform request into model and transform model into response
│   ├── handler/        # HTTP handlers (controllers)
│   ├── middleware/     # Middleware functions
│   ├── model/          # Data models and entities
│   ├── pkg/            # Utility functions
│   ├── repository/     # Database repository and interfaces
│   ├── route/          # API route definitions
│   └── usecase/        # Business logic and services
├── log/                # Directory for log application
├── migrations/         # Database migration scripts
├── tests/              # API tests
├── .env.example        # Example environment variables file
├── Dockerfile          # Docker configuration file
├── main.go             # Main Entrypoint for this app
└── README.md           # This README file
```

## Creating a New API
To create a new API endpoint, follow these steps:

1. Define the Model: Create a new model in the internal/model directory representing your data structure.

Example: `internal/model/product.go`

```go
package model

type Product struct {
    ID    int     `json:"id"`
    Name  string  `json:"name"`
    Price float64 `json:"price"`
}
```

2. Create the DTO: to transfer model with request and response payload

Example: `internal/dto/product_dto.go`

```go
package dto

import (
	"errors"
	"rest-skeleton/internal/model"
)

type ProductCreateRequest struct {
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

func (p *ProductCreateRequest) Validate() error {
	if len(p.Name) == 0 {
		return errors.New("name is required")
	}
	if p.Price <= 0 {
		return errors.New("price must be greater than zero")
	}
	return nil
}

func (p *ProductCreateRequest) ToEntity() model.Product {
	return model.Product{
		Name:  p.Name,
		Price: p.Price,
	}
}

type ProductResponse struct {
	ID    int64   `json:"id"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

func (p *ProductResponse) FromEntity(product model.Product) {
	p.ID = product.ID
	p.Name = product.Name
	p.Price = product.Price
}
```

3. Create the Handler: Implement the handler functions in the internal/handler directory.

Example: `internal/handler/product.go`

```go
package handler

import (
	"context"
	"database/sql"
	"net/http"
	"os"
	"rest-skeleton/internal/dto"
	"rest-skeleton/internal/model"
	"rest-skeleton/internal/pkg/httpresponse"
	"rest-skeleton/internal/pkg/logger"
	"rest-skeleton/internal/repository"

	"github.com/bytedance/sonic"
	"github.com/julienschmidt/httprouter"
	"go.opentelemetry.io/otel"
)

// Products handler
type Products struct {
	Log *logger.Logger
	DB  *sql.DB
}

// @Security Bearer
// @Summary Create Product
// @Description Create Product
// @Tags Products
// @Accept  json
// @Produce  json
// @Param product body dto.ProductCreateRequest true "Product to add"
// @Param Authorization header string true "Bearer token"
// @Success 201 {object} dto.ProductResponse
// @Router /products [post]
func (h *Products) Create(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var ctx = r.Context()
	ctx, span := otel.Tracer(os.Getenv("APP_NAME")).Start(ctx, "CreateProductHandler")
	defer span.End()

	defer r.Body.Close()
	var httpres = httpresponse.Response{}
	var productRequest dto.ProductCreateRequest
	err := sonic.ConfigDefault.NewDecoder(r.Body).Decode(&productRequest)
	if err != nil {
		h.Log.Error(ctx, err)
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	if err := productRequest.Validate(); err != nil {
		h.Log.Error(ctx, err)
		http.Error(w, "Invalid input: "+err.Error(), http.StatusBadRequest)
		return
	}

	var productRepo = repository.ProductRepository{Log: h.Log, Db: h.DB}
	productRepo.ProductEntity = productRequest.ToEntity()
	if err := productRepo.Save(ctx); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	var response dto.ProductResponse
	response.FromEntity(productRepo.ProductEntity)
	httpres.SetMarshal(ctx, w, http.StatusCreated, response, "")
}

```

4. Create Respository

Example: `internal/repository/product_repository.go`

```go
package repository

import (
	"context"
	"database/sql"
	"rest-skeleton/internal/model"
)

type ProductRepository struct {
	Log           *logger.Logger
	Db            *sql.DB
	ProductEntity model.Product
}

func (r *ProductRepository) Save(ctx context.Context) error {
	// Implementasi simpan produk ke database
	return nil
}
```

5. Define Routes: Add new routes in the internal/route directory to map HTTP requests to your handler functions.

Example: `internal/route/route.go`

```go
func ApiRoute(log *logger.Logger, db *database.Database, cache *redis.Cache, latencyMetric metric.Int64Histogram) *httprouter.Router {
    // .... existing code
    productHandler := handler.Products{Log: log, DB: db.Conn}
    router.POST("/products", mid.WrapMiddleware(privateMiddlewares, productHandler.Create))
    // .... existing code
}
```

6. Create swagger documentation with command `swag init`. Attention to install `go install github.com/swaggo/swag/cmd/swag@latest` before you run `swag init`. 

7. Testing: Write tests for your new API endpoint in the tests directory to ensure it behaves as expected.

## Creating Migration Scripts
When adding a new database migration for your product model, follow these steps and naming conventions:

1. Naming Convention: Use the following prefixes for your migration scripts:

- For functions: `1.001_fn_random_bigint.sql`
- For tables: `2.001_t_access.sql`
- For seeding: `3.001_seed.sql`

After the prefix, include a three-digit serial number followed by a descriptive name for the migration.

2. Create a Migration for Products: For the products table, create a migration file in the migrations directory named `2.006_t_products.sql`.

Example: `migrations/2.006_t_products.sql`

```sql
CREATE TABLE products (
    id int8 DEFAULT int64_id('products'::text, 'id'::text) NOT NULL,
    name VARCHAR(255) NOT NULL,
    price NUMERIC(10, 2) NOT NULL,
    created_at TIMESTAMPTZ DEFAULT timezone('utc'::text, now()) NULL,
);
```

3. Functions or Seeds: If you need to create a function or a seed, follow the same naming conventions for your migration files.

Example function file: `1.001_fn_add_tax.sql`

```sql
CREATE OR REPLACE FUNCTION add_tax(price NUMERIC)
RETURNS NUMERIC AS $$
BEGIN
    RETURN price * 1.1;  -- Adds a 10% tax
END;
$$ LANGUAGE plpgsql;
```

Example seed file: `3.001_seed.sql`

```sql
INSERT INTO products (name, price) VALUES
('Sample Product', 19.99),
('Another Product', 29.99);
```

4. Running migration command using `go run cmd/main.go migrate`

## Running Tests
To run the API tests, use the following command:

```bash
go test ./...
```

## Contributing
Contributions are welcome! Please fork the repository and submit a pull request for any enhancements or fixes.

## License
This project is licensed under the MIT License. See the [LICENSE](./LICENSE) file for details.