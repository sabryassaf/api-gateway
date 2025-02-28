package controllers

import (
	"context"
	"net/http"
	"strings"
	"time"

	gradesProtos "github.com/BetterGR/grades-microservice/protos"
	"google.golang.org/grpc/credentials/insecure"
	"k8s.io/klog/v2"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
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
	studentId := c.Param("studentId")
	courseId := c.Param("courseId")

	// Build gRPC request.
	request := &gradesProtos.GetStudentCourseGradesRequest{
		Token:     token,
		CourseId:  courseId,
		Semester:  semester,
		StudentId: studentId,
	}

	logger := klog.FromContext(c.Request.Context())
	logger.V(logLevelDebug).Info("Received request for student course grades", "course_id", request.CourseId,
		"semester", request.Semester, "student_id", request.StudentId)

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

	logger := klog.FromContext(c.Request.Context())
	logger.V(logLevelDebug).Info("Received request for student grades", "semester", semester, "student_id", studentId)

	// Build gRPC request.
	request := &gradesProtos.GetStudentSemesterGradesRequest{
		Token:     token,
		Semester:  semester,
		StudentId: studentId,
	}
	logger.V(logLevelDebug).Info("Request built with student ID: '%s'", request.StudentId)

	// Call the gRPC server.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	response, err := grpcClient.GetStudentSemesterGrades(ctx, request)
	if err != nil {
		klog.Errorf("Error calling gRPC Grades Microservice: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch grades"})
		return
	}
	c.JSON(http.StatusOK, response)
}

func AddHomeworkGradeHandler(c *gin.Context, grpcClient gradesProtos.GradesServiceClient) {
	c.JSON(http.StatusInternalServerError, gin.H{"error": "Not Implemented"})
}

func AddExamGradeHandler(c *gin.Context, grpcClient gradesProtos.GradesServiceClient) {
	c.JSON(http.StatusInternalServerError, gin.H{"error": "Not Implemented"})
}

func UpdateHomeworkGradeHandler(c *gin.Context, grpcClient gradesProtos.GradesServiceClient) {
	c.JSON(http.StatusInternalServerError, gin.H{"error": "Not Implemented"})
}

func UpdateExamGradeHandler(c *gin.Context, grpcClient gradesProtos.GradesServiceClient) {
	c.JSON(http.StatusInternalServerError, gin.H{"error": "Not Implemented"})
}

func DeleteHomeworkGradeHandler(c *gin.Context, grpcClient gradesProtos.GradesServiceClient) {
	c.JSON(http.StatusInternalServerError, gin.H{"error": "Not Implemented"})
}

func DeleteExamGradeHandler(c *gin.Context, grpcClient gradesProtos.GradesServiceClient) {
	c.JSON(http.StatusInternalServerError, gin.H{"error": "Not Implemented"})
}
