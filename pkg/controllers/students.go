package controllers

import (
	"context"
	"net/http"
	"time"

	studentsProtos "github.com/BetterGR/students-microservice/protos"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"k8s.io/klog/v2"
)

// InitStudentsGRPCClient initializes the students-microservice gRPC client connection.
func InitStudentsGRPCClient(address string) (studentsProtos.StudentsServiceClient, error) {
	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	return studentsProtos.NewStudentsServiceClient(conn), nil
}

func CreateStudentHandler(c *gin.Context, grpcClient studentsProtos.StudentsServiceClient) {
	var student studentsProtos.Student
	if err := c.ShouldBindJSON(&student); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if _, err := grpcClient.CreateStudent(ctx, &studentsProtos.CreateStudentRequest{Student: &student}); err != nil {
		klog.Errorf("Failed to create student: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create student"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Student created successfully"})
}

func GetStudentHandler(c *gin.Context, grpcClient studentsProtos.StudentsServiceClient) {
	studentID := c.Param("studentID")
	if studentID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "studentID is required"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := grpcClient.GetStudent(ctx, &studentsProtos.GetStudentRequest{StudentID: studentID})
	if err != nil {
		klog.Errorf("Failed to get student: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch student"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"studentID":   resp.GetStudent().GetStudentID(),
		"firstName":   resp.GetStudent().GetFirstName(),
		"lastName":    resp.GetStudent().GetLastName(),
		"email":       resp.GetStudent().GetEmail(),
		"phoneNumber": resp.GetStudent().GetPhoneNumber(),
	})
}

func UpdateStudentHandler(c *gin.Context, grpcClient studentsProtos.StudentsServiceClient) {
	var student studentsProtos.Student
	if err := c.ShouldBindJSON(&student); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if _, err := grpcClient.UpdateStudent(ctx, &studentsProtos.UpdateStudentRequest{Student: &student}); err != nil {
		klog.Errorf("Failed to update student: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update student"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Student updated successfully"})
}

func DeleteStudentHandler(c *gin.Context, grpcClient studentsProtos.StudentsServiceClient) {
	studentID := c.Param("studentID")
	if studentID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "studentID is required"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if _, err := grpcClient.DeleteStudent(ctx, &studentsProtos.DeleteStudentRequest{StudentID: studentID}); err != nil {
		klog.Errorf("Failed to delete student: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete student"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Student deleted successfully"})
}
