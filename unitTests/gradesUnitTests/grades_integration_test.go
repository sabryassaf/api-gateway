package gradesUnitTests

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"testing"
	"time"

	"github.com/BetterGR/api-gateway/api/routes"
	"github.com/BetterGR/api-gateway/pkg/controllers"
	gradesProtos "github.com/BetterGR/grades-microservice/protos"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Integration test that tests the actual interaction between the API gateway and the grades microservice
// Note: This test requires both the API Gateway and Grades Microservice to be running
// You can run the microservices using docker-compose or manually start them before running the test
func TestGradesMicroserviceIntegration(t *testing.T) {
	// Skip this test if we're not in integration test mode
	if os.Getenv("RUN_INTEGRATION_TESTS") != "true" {
		t.Skip("Skipping integration test. Set RUN_INTEGRATION_TESTS=true to run")
	}

	// Initialize gin router
	gin.SetMode(gin.TestMode)
	router := gin.New()

	// Get the grades microservice address from environment or use default
	gradesAddress := os.Getenv("GRADES_ADDRESS")
	if gradesAddress == "" {
		gradesAddress = "localhost:50052"
		os.Setenv("GRADES_ADDRESS", gradesAddress)
	}

	// Register the grades routes
	grpcClient, err := routes.RegisterGradesRoutes(router)
	require.NoError(t, err)
	require.NotNil(t, grpcClient)

	// Test connectivity to the grades microservice
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Try to make a simple gRPC call to verify connectivity
	// We'll use a dummy token since we're just testing connectivity
	dummyToken := "test-token"
	dummyRequest := &gradesProtos.GetCourseGradesRequest{
		Token:    dummyToken,
		CourseID: "test-course",
		Semester: "SEMESTER_TERM_WINTER_2024",
	}

	// This call will likely fail with an authentication error, but that's OK
	// We just want to make sure we can connect to the service
	_, err = grpcClient.GetCourseGrades(ctx, dummyRequest)
	// The error should be an authentication error, not a connection error
	if err != nil {
		t.Logf("Expected authentication error: %v", err)
		assert.Contains(t, err.Error(), "authentication",
			"Error should be about authentication, not connectivity")
	}

	t.Log("Successfully connected to the grades microservice")
}

// Helper function to set up the test environment
// This function starts the grades microservice in a separate process
func startGradesMicroservice(t *testing.T) (*exec.Cmd, error) {
	// Set necessary environment variables for the grades microservice
	env := []string{
		"GRPC_PORT=50052",
		"DB_HOST=localhost",
		"DB_PORT=5432",
		"DB_USER=postgres",
		"DB_PASSWORD=postgres",
		"DB_NAME=grades",
	}

	// Start the grades microservice
	cmd := exec.Command("go", "run", "../../../grades-microservice/server/server.go")
	cmd.Env = append(os.Environ(), env...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Start()
	if err != nil {
		return nil, fmt.Errorf("failed to start grades microservice: %w", err)
	}

	// Wait for the service to start
	time.Sleep(2 * time.Second)

	return cmd, nil
}

// This test requires docker and docker-compose to be installed
// It will start the grades microservice using docker-compose
func TestGradesMicroserviceIntegrationWithDocker(t *testing.T) {
	// Skip this test if we're not in integration test mode with docker
	if os.Getenv("RUN_DOCKER_INTEGRATION_TESTS") != "true" {
		t.Skip("Skipping docker integration test. Set RUN_DOCKER_INTEGRATION_TESTS=true to run")
	}

	// Start the grades microservice using docker-compose
	cmd := exec.Command("docker-compose", "-f", "../../../docker-compose.yml", "up", "-d", "grades-microservice")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	require.NoError(t, err, "Failed to start grades microservice with docker-compose")

	// Ensure we clean up after ourselves
	defer func() {
		cmd := exec.Command("docker-compose", "-f", "../../../docker-compose.yml", "down")
		cmd.Run()
	}()

	// Wait for the service to start
	time.Sleep(5 * time.Second)

	// Initialize the gRPC client
	gradesAddress := "localhost:50052"
	grpcClient, err := controllers.InitGradesGRPCClient(gradesAddress)
	require.NoError(t, err)
	require.NotNil(t, grpcClient)

	// Test connectivity to the grades microservice
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Try to make a simple gRPC call to verify connectivity
	dummyToken := "test-token"
	dummyRequest := &gradesProtos.GetCourseGradesRequest{
		Token:    dummyToken,
		CourseID: "test-course",
		Semester: "SEMESTER_TERM_WINTER_2024",
	}

	// This call will likely fail with an authentication error, but that's OK
	_, err = grpcClient.GetCourseGrades(ctx, dummyRequest)
	// The error should be an authentication error, not a connection error
	if err != nil {
		t.Logf("Expected authentication error: %v", err)
		assert.Contains(t, err.Error(), "authentication",
			"Error should be about authentication, not connectivity")
	}

	t.Log("Successfully connected to the grades microservice in Docker")
}

// Test with a mock HTTP server
func TestGradesMicroserviceHTTPEndpoints(t *testing.T) {
	// Skip if not in integration test mode
	if os.Getenv("RUN_INTEGRATION_TESTS") != "true" {
		t.Skip("Skipping HTTP endpoints test. Set RUN_INTEGRATION_TESTS=true to run")
	}

	// Initialize gin router
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	// Register the grades routes
	os.Setenv("GRADES_ADDRESS", "localhost:50052")
	_, err := routes.RegisterGradesRoutes(router)
	require.NoError(t, err)

	// Test the /api/grades/:semester/:student_id/:courseId endpoint
	// Create an HTTP test server
	server := http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	// Start the server in a goroutine
	go func() {
		server.ListenAndServe()
	}()
	defer server.Shutdown(context.Background())

	// Wait for the server to start
	time.Sleep(1 * time.Second)

	// Make an HTTP request to the server
	client := &http.Client{}
	req, err := http.NewRequest("GET", "http://localhost:8080/api/grades/SEMESTER_TERM_WINTER_2024/student123/course456", nil)
	require.NoError(t, err)

	// Add authorization header
	req.Header.Add("Authorization", "Bearer test-token")

	resp, err := client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	// We expect a 500 error since the grades microservice will return an auth error
	// The important thing is that the API Gateway properly forwarded the request
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)

	// Decode the response
	var response map[string]string
	err = json.NewDecoder(resp.Body).Decode(&response)
	require.NoError(t, err)

	// We expect an error message about failing to fetch grades
	assert.Contains(t, response["error"], "Failed to fetch grades")
}
