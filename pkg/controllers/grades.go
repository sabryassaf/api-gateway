package controllers

import (
	"context"
	"net/http"
	"strings"
	"time"

	gradesProtos "github.com/BetterGR/grades-microservice/protos"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"k8s.io/klog/v2"
)

const (
	logLevelDebug = 5
)

// InitGradesGRPCClient initializes the grades-microservice gRPC client connection.
func InitGradesGRPCClient(address string) (gradesProtos.GradesServiceClient, error) {
	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	return gradesProtos.NewGradesServiceClient(conn), nil
}

// GetStudentCourseGradesHandler handles REST requests and calls the gRPC Grades Microservice.
func GetStudentCourseGradesHandler(c *gin.Context, grpcClient gradesProtos.GradesServiceClient) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "No authorization token provided"})
		return
	}

	token := strings.TrimPrefix(authHeader, "Bearer ")
	semester := c.Param("semester")
	studentId := c.Param("student_id")
	courseId := c.Param("courseId")

	// Log the parameters for debugging
	logger := klog.FromContext(c.Request.Context())
	logger.V(logLevelDebug).Info("Received request for student course grades",
		"course_id", courseId,
		"semester", semester,
		"student_id", studentId,
		"URI", c.Request.RequestURI,
		"params", c.Params)

	// Build gRPC request.
	request := &gradesProtos.GetStudentCourseGradesRequest{
		Token:     token,
		CourseID:  courseId,
		Semester:  semester,
		StudentID: studentId,
	}

	// Call the gRPC server.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	response, err := grpcClient.GetStudentCourseGrades(ctx, request)
	if err != nil {
		klog.Errorf("Error calling gRPC Grades Microservice: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch grades"})

		return
	}

	// Send response to the client.
	c.JSON(http.StatusOK, response.Grades)
}

// GetStudentSemesterGradesHandler handles REST requests and calls the gRPC Grades Microservice
// to return all the student grades for a specific semester.
func GetStudentSemesterGradesHandler(c *gin.Context, grpcClient gradesProtos.GradesServiceClient) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "No authorization token provided"})
		return
	}

	token := strings.TrimPrefix(authHeader, "Bearer ")
	semester := c.Param("semester")
	studentId := c.Param("student_id")

	// Log the parameters for debugging.
	logger := klog.FromContext(c.Request.Context())
	logger.V(logLevelDebug).Info("Received request for student grades",
		"semester", semester,
		"student_id", studentId,
		"URI", c.Request.RequestURI,
		"params", c.Params)

	// Build gRPC request.
	request := &gradesProtos.GetStudentSemesterGradesRequest{
		Token:     token,
		Semester:  semester,
		StudentID: studentId,
	}
	logger.V(logLevelDebug).Info("Request built with student ID", "studentID", request.StudentID)

	// Call the gRPC server.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	response, err := grpcClient.GetStudentSemesterGrades(ctx, request)
	if err != nil {
		klog.Errorf("Error calling gRPC Grades Microservice: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch grades"})
		return
	}

	// Return the entire response including the grades array
	c.JSON(http.StatusOK, response)
}

func AddSingleGradeHandler(c *gin.Context, grpcClient gradesProtos.GradesServiceClient) {
	c.JSON(http.StatusInternalServerError, gin.H{"error": "Not Implemented"})
}

func UpdateSingleGradeHandler(c *gin.Context, grpcClient gradesProtos.GradesServiceClient) {
	c.JSON(http.StatusInternalServerError, gin.H{"error": "Not Implemented"})
}

func DeleteSingleGradeHandler(c *gin.Context, grpcClient gradesProtos.GradesServiceClient) {
	c.JSON(http.StatusInternalServerError, gin.H{"error": "Not Implemented"})
}
