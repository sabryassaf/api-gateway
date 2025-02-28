# Grades Microservice Unit Tests

This directory contains unit and integration tests for the Grades microservice API integration in the API Gateway.

## Test Types

1. **Unit Tests** - Test the API Gateway's controllers without making actual gRPC calls
2. **Integration Tests** - Test the actual interaction between the API Gateway and the Grades microservice

## Running Unit Tests

To run the unit tests:

```bash
cd api-gateway
go test ./unitTests/gradesUnitTests -run "^Test[^Integration]" -v
```

These tests use mocks to simulate the Grades microservice, so you don't need the actual microservice running.

## Running Integration Tests

### Prerequisites

Before running integration tests, make sure:

1. The Grades microservice is running
2. You have the correct environment variables set (especially `GRADES_ADDRESS`)

### Option 1: Start services manually

1. Start the Grades microservice:
   ```bash
   cd grades-microservice
   # Set required environment variables
   export GRPC_PORT=50052
   export DB_HOST=localhost
   export DB_PORT=5432
   export DB_USER=postgres
   export DB_PASSWORD=postgres
   export DB_NAME=grades
   
   go run server/server.go
   ```

2. In another terminal, run the integration tests:
   ```bash
   cd api-gateway
   export RUN_INTEGRATION_TESTS=true
   export GRADES_ADDRESS=localhost:50052
   go test ./unitTests/gradesUnitTests -run "Integration" -v
   ```

### Option 2: Use Docker Compose

1. Run the integration tests with Docker Compose:
   ```bash
   cd api-gateway
   export RUN_DOCKER_INTEGRATION_TESTS=true
   go test ./unitTests/gradesUnitTests -run "IntegrationWithDocker" -v
   ```

   This test will automatically start the Grades microservice using Docker Compose.

## Test Coverage

To see test coverage:

```bash
cd api-gateway
go test ./unitTests/gradesUnitTests -cover
```

For a detailed coverage report:

```bash
cd api-gateway
go test ./unitTests/gradesUnitTests -coverprofile=coverage.out
go tool cover -html=coverage.out
```

## Troubleshooting

### Connection Issues

If you encounter connectivity issues during integration tests:

1. Verify the Grades microservice is running: `ps aux | grep grades-microservice`
2. Check the correct port is being used (default is 50052)
3. Ensure no firewall is blocking the connection

### Authentication Errors

Authentication errors in integration tests are expected since we're using dummy tokens.
The important thing is that the API Gateway can connect to the Grades microservice. 