## Services

- **Order Service**: Manages order processing
- **Product Service**: Provides product and inventory management

## Technologies Used

- Go/Fiber
- Bun/Uptrace
- PostgreSQL
- RabbitMQ
- Docker / Docker Compose
- Vue.js

## Running Locally

#### Requirements

- Docker ve Docker Compose
- Go 1.21+

### Initialization Steps


1. Start PostgreSQL databases and RabbitMQ:

```bash
docker-compose up -d
```

2. Run migrations for each service:

```bash
cd backend/order-service 
make migrate-up

cd ../product-service
make migrate-up
```

3. Upload sample products

```bash
cd backend/product-service/
make seed
```

3. Run the services on separate terminals:

```bash
# Terminal 1
cd backend/order-service
go run server/cmd/main.go
```

```bash
# Terminal 2
cd backend/product-service
go run server/cmd/main.go
```

### Service APIs

#### Order Service

- **Endpoint**: http://localhost:3000
  - `POST /api/order` - Create new order
  - `GET /api/order/:id` - Get order details

#### Product Service

- **Endpoint**: http://localhost:3001
  - `POST /api/product` - Create new product
  - `GET /api/product` - Get all products
  - `GET /api/product/:id` - Get product details

### RabbitMQ Administrator Interface

- **URL**: http://localhost:15672
- **Kullanıcı adı**: guest
- **Şifre**: guest

### Frontend
```bash
cd frontend && yarn
yarn dev
```
Available from http://localhost:8080




## Sample Usage Flow

1. Create a new product with Product Service
2. Create a new order with Order Service
3. Order Service saves the order status as “pending”
4. Order Service sends a message to Product Service via RabbitMQ for stock reduction
5. Product Service checks stock status and drops stock if available
6. Product Service reports the transaction result to Order Service via RabbitMQ
7. Order Service updates the order status as “created” or “failed” according to the result

## Saga Pattern Usage

Saga Pattern's Choreography approach is used in the project. Communication between services is provided through RabbitMQ message queue.
