package gradesUnitTests

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/BetterGR/api-gateway/pkg/controllers"
	gradesProtos "github.com/BetterGR/grades-microservice/protos"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// MockGradesServiceClient is a mock of the GradesServiceClient
type MockGradesServiceClient struct {
	mock.Mock
}

// Implement all the methods from the GradesServiceClient interface
func (m *MockGradesServiceClient) GetCourseGrades(ctx context.Context, in *gradesProtos.GetCourseGradesRequest, opts ...grpc.CallOption) (*gradesProtos.GetCourseGradesResponse, error) {
	args := m.Called(ctx, in, opts)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*gradesProtos.GetCourseGradesResponse), args.Error(1)
}

func (m *MockGradesServiceClient) GetStudentCourseGrades(ctx context.Context, in *gradesProtos.GetStudentCourseGradesRequest, opts ...grpc.CallOption) (*gradesProtos.GetStudentCourseGradesResponse, error) {
	args := m.Called(ctx, in, opts)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*gradesProtos.GetStudentCourseGradesResponse), args.Error(1)
}

func (m *MockGradesServiceClient) AddSingleGrade(ctx context.Context, in *gradesProtos.AddSingleGradeRequest, opts ...grpc.CallOption) (*gradesProtos.AddSingleGradeResponse, error) {
	args := m.Called(ctx, in, opts)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*gradesProtos.AddSingleGradeResponse), args.Error(1)
}

func (m *MockGradesServiceClient) UpdateSingleGrade(ctx context.Context, in *gradesProtos.UpdateSingleGradeRequest, opts ...grpc.CallOption) (*gradesProtos.UpdateSingleGradeResponse, error) {
	args := m.Called(ctx, in, opts)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*gradesProtos.UpdateSingleGradeResponse), args.Error(1)
}

func (m *MockGradesServiceClient) RemoveSingleGrade(ctx context.Context, in *gradesProtos.RemoveSingleGradeRequest, opts ...grpc.CallOption) (*gradesProtos.RemoveSingleGradeResponse, error) {
	args := m.Called(ctx, in, opts)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*gradesProtos.RemoveSingleGradeResponse), args.Error(1)
}

func (m *MockGradesServiceClient) GetStudentSemesterGrades(ctx context.Context, in *gradesProtos.GetStudentSemesterGradesRequest, opts ...grpc.CallOption) (*gradesProtos.GetStudentSemesterGradesResponse, error) {
	args := m.Called(ctx, in, opts)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*gradesProtos.GetStudentSemesterGradesResponse), args.Error(1)
}

// Helper function to create a test Gin context with the given path and method
func setupTestContext(method, path string, params map[string]string, authToken string) (*gin.Context, *httptest.ResponseRecorder) {
	gin.SetMode(gin.TestMode)

	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)

	req, _ := http.NewRequest(method, path, nil)

	// Add authorization header if provided
	if authToken != "" {
		req.Header.Set("Authorization", "Bearer "+authToken)
	}

	ctx.Request = req

	// Set URL parameters
	for key, value := range params {
		ctx.Params = append(ctx.Params, gin.Param{Key: key, Value: value})
	}

	return ctx, recorder
}

// Test GetStudentCourseGradesHandler with valid input
func TestGetStudentCourseGradesHandler_Success(t *testing.T) {
	// Create a mock client
	mockClient := new(MockGradesServiceClient)

	// Define test data
	semesterStr := "SEMESTER_TERM_WINTER_2024"
	studentID := "student123"
	courseID := "course456"
	token := "valid-token"

	// Setup the expected response
	mockGrades := []*gradesProtos.SingleGrade{
		{
			GradeID:    "grade1",
			StudentID:  studentID,
			CourseID:   courseID,
			Semester:   semesterStr,
			GradeType:  "Exam",
			ItemID:     "midterm",
			GradeValue: "A",
			GradedBy:   "Professor X",
			Comments:   "Excellent work",
		},
	}

	// Configure mock to return these grades when called with appropriate parameters
	mockClient.On("GetStudentCourseGrades", mock.Anything,
		&gradesProtos.GetStudentCourseGradesRequest{
			Token:     token,
			CourseID:  courseID,
			Semester:  semesterStr,
			StudentID: studentID,
		},
		mock.Anything).Return(&gradesProtos.GetStudentCourseGradesResponse{
		Grades: mockGrades,
	}, nil)

	// Setup test context
	params := map[string]string{
		"semester":  semesterStr,
		"studentId": studentID,
		"courseId":  courseID,
	}
	ctx, recorder := setupTestContext("GET", "/api/grades/"+semesterStr+"/"+studentID+"/"+courseID, params, token)

	// Call the handler
	controllers.GetStudentCourseGradesHandler(ctx, mockClient)

	// Assert response
	require.Equal(t, http.StatusOK, recorder.Code)

	// Parse the response body
	var responseGrades []*gradesProtos.SingleGrade
	err := json.Unmarshal(recorder.Body.Bytes(), &responseGrades)
	require.NoError(t, err)

	// Verify the response
	assert.Len(t, responseGrades, 1)
	assert.Equal(t, mockGrades[0].GradeID, responseGrades[0].GradeID)
	assert.Equal(t, mockGrades[0].GradeValue, responseGrades[0].GradeValue)

	// Verify that the mock was called as expected
	mockClient.AssertExpectations(t)
}

// Test GetStudentCourseGradesHandler with missing authentication
func TestGetStudentCourseGradesHandler_NoAuth(t *testing.T) {
	// Create a mock client
	mockClient := new(MockGradesServiceClient)

	// Define test data
	semesterStr := "SEMESTER_TERM_WINTER_2024"
	studentID := "student123"
	courseID := "course456"

	// Setup test context without auth token
	params := map[string]string{
		"semester":  semesterStr,
		"studentId": studentID,
		"courseId":  courseID,
	}
	ctx, recorder := setupTestContext("GET", "/api/grades/"+semesterStr+"/"+studentID+"/"+courseID, params, "")

	// Call the handler
	controllers.GetStudentCourseGradesHandler(ctx, mockClient)

	// Assert response - should be unauthorized
	require.Equal(t, http.StatusUnauthorized, recorder.Code)

	// Verify the mock was not called
	mockClient.AssertNotCalled(t, "GetStudentCourseGrades")
}

// Test GetStudentCourseGradesHandler with gRPC error
func TestGetStudentCourseGradesHandler_GRPCError(t *testing.T) {
	// Create a mock client
	mockClient := new(MockGradesServiceClient)

	// Define test data
	semesterStr := "SEMESTER_TERM_WINTER_2024"
	studentID := "student123"
	courseID := "course456"
	token := "valid-token"

	// Configure mock to return an error
	mockClient.On("GetStudentCourseGrades", mock.Anything, mock.Anything, mock.Anything).
		Return(nil, status.Error(codes.NotFound, "grades not found"))

	// Setup test context
	params := map[string]string{
		"semester":  semesterStr,
		"studentId": studentID,
		"courseId":  courseID,
	}
	ctx, recorder := setupTestContext("GET", "/api/grades/"+semesterStr+"/"+studentID+"/"+courseID, params, token)

	// Call the handler
	controllers.GetStudentCourseGradesHandler(ctx, mockClient)

	// Assert response - should be internal server error
	require.Equal(t, http.StatusInternalServerError, recorder.Code)

	// Verify that the mock was called
	mockClient.AssertExpectations(t)
}

// Test GetStudentSemesterGradesHandler with valid input
func TestGetStudentSemesterGradesHandler_Success(t *testing.T) {
	// Create a mock client
	mockClient := new(MockGradesServiceClient)

	// Define test data
	semesterStr := "SEMESTER_TERM_WINTER_2024"
	studentID := "student123"
	token := "valid-token"

	// Setup the expected response
	mockGrades := []*gradesProtos.SingleGrade{
		{
			GradeID:    "grade1",
			StudentID:  studentID,
			CourseID:   "course1",
			Semester:   semesterStr,
			GradeType:  "Exam",
			ItemID:     "midterm",
			GradeValue: "A",
		},
		{
			GradeID:    "grade2",
			StudentID:  studentID,
			CourseID:   "course2",
			Semester:   semesterStr,
			GradeType:  "Homework",
			ItemID:     "hw1",
			GradeValue: "B+",
		},
	}

	// Configure mock to return these grades
	mockClient.On("GetStudentSemesterGrades", mock.Anything,
		&gradesProtos.GetStudentSemesterGradesRequest{
			Token:     token,
			Semester:  semesterStr,
			StudentID: studentID,
		},
		mock.Anything).Return(&gradesProtos.GetStudentSemesterGradesResponse{
		Grades: mockGrades,
	}, nil)

	// Setup test context
	params := map[string]string{
		"semester":  semesterStr,
		"studentId": studentID,
	}
	ctx, recorder := setupTestContext("GET", "/api/grades/"+semesterStr+"/"+studentID, params, token)

	// Call the handler
	controllers.GetStudentSemesterGradesHandler(ctx, mockClient)

	// Assert response
	require.Equal(t, http.StatusOK, recorder.Code)

	// Verify that the mock was called as expected
	mockClient.AssertExpectations(t)
}

// Test AddHomeworkGradeHandler
func TestAddHomeworkGradeHandler(t *testing.T) {
	// Create a mock client
	mockClient := new(MockGradesServiceClient)

	// Setup test context
	ctx, recorder := setupTestContext("POST", "/api/grades/homework", nil, "valid-token")

	// Call the handler
	controllers.AddSingleGradeHandler(ctx, mockClient)

	// Since this is not implemented yet, it should return Internal Server Error
	require.Equal(t, http.StatusInternalServerError, recorder.Code)

	// Parse response to check error message
	var response map[string]string
	err := json.Unmarshal(recorder.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, "Not Implemented", response["error"])
}

// Additional tests could be added for other handlers once they are implemented

// This test would verify the timeout handling
func TestHandlerTimeoutHandling(t *testing.T) {
	// Create a mock client that simulates a timeout
	mockClient := new(MockGradesServiceClient)

	// Define test data
	semesterStr := "SEMESTER_TERM_WINTER_2024"
	studentID := "student123"
	courseID := "course456"
	token := "valid-token"

	// Configure mock to delay longer than the timeout
	mockClient.On("GetStudentCourseGrades", mock.Anything, mock.Anything, mock.Anything).
		Run(func(args mock.Arguments) {
			// Sleep longer than the context timeout in the handler
			time.Sleep(2 * time.Second)
		}).
		Return(&gradesProtos.GetStudentCourseGradesResponse{}, nil)

	// Setup test context
	params := map[string]string{
		"semester":  semesterStr,
		"studentId": studentID,
		"courseId":  courseID,
	}
	ctx, recorder := setupTestContext("GET", "/api/grades/"+semesterStr+"/"+studentID+"/"+courseID, params, token)

	// Call the handler
	controllers.GetStudentCourseGradesHandler(ctx, mockClient)

	// Assert response - should be internal server error due to timeout
	require.Equal(t, http.StatusInternalServerError, recorder.Code)
}
