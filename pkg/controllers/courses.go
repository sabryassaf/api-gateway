package controllers

import (
	"context"
	"net/http"
	"strings"
	"time"

	coursesProtos "github.com/BetterGR/courses-microservice/protos"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"k8s.io/klog/v2"
)

// InitCoursesGRPCClient initializes the course-microservice gRPC client connection.
func InitCoursesGRPCClient(address string) (coursesProtos.CoursesServiceClient, error) {
	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	return coursesProtos.NewCoursesServiceClient(conn), nil
}

func CreateCourseHandler(c *gin.Context, grpcClient coursesProtos.CoursesServiceClient) {
	var course coursesProtos.Course
	if err := c.ShouldBindJSON(&course); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if _, err := grpcClient.CreateCourse(ctx, &coursesProtos.CreateCourseRequest{Course: &course}); err != nil {
		klog.Errorf("Failed to create course: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create course"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Course created successfully"})
}

func GetCourseHandler(c *gin.Context, grpcClient coursesProtos.CoursesServiceClient) {
	courseID := c.Param("courseID")
	if courseID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "courseID is required"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := grpcClient.GetCourse(ctx, &coursesProtos.GetCourseRequest{CourseID: courseID})
	if err != nil {
		klog.Errorf("Failed to get course: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get course"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"courseID":    resp.GetCourse().GetCourseID(),
		"courseName":  resp.GetCourse().GetCourseName(),
		"semester":    resp.GetCourse().GetSemester(),
		"description": resp.GetCourse().GetDescription(),
	})
}

func UpdateCourseHandler(c *gin.Context, grpcClient coursesProtos.CoursesServiceClient) {
	var course coursesProtos.Course
	if err := c.ShouldBindJSON(&course); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if _, err := grpcClient.UpdateCourse(ctx, &coursesProtos.UpdateCourseRequest{Course: &course}); err != nil {
		klog.Errorf("Failed to update course: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update course"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Course updated successfully"})
}

func DeleteCourseHandler(c *gin.Context, grpcClient coursesProtos.CoursesServiceClient) {
	courseID := c.Param("courseID")
	if courseID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "courseID is required"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if _, err := grpcClient.DeleteCourse(ctx, &coursesProtos.DeleteCourseRequest{CourseID: courseID}); err != nil {
		klog.Errorf("Failed to delete course: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete course"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Course deleted successfully"})
}

func AddStudentToCourseHandler(c *gin.Context, grpcClient coursesProtos.CoursesServiceClient) {
	courseID := c.Param("courseID")
	studentID := c.Param("studentID")
	if courseID == "" || studentID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "courseID and studentID are required"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if _, err := grpcClient.AddStudentToCourse(ctx, &coursesProtos.AddStudentRequest{CourseID: courseID, StudentID: studentID}); err != nil {
		klog.Errorf("Failed to add student to course: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add student to course"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Student added to course successfully"})
}

func RemoveStudentFromCourseHandler(c *gin.Context, grpcClient coursesProtos.CoursesServiceClient) {
	courseID := c.Param("courseID")
	studentID := c.Param("studentID")
	if courseID == "" || studentID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "courseID and studentID are required"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if _, err := grpcClient.RemoveStudentFromCourse(ctx, &coursesProtos.RemoveStudentRequest{CourseID: courseID, StudentID: studentID}); err != nil {
		klog.Errorf("Failed to remove student from course: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove student from course"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Student removed from course successfully"})
}

func AddStaffToCourseHandler(c *gin.Context, grpcClient coursesProtos.CoursesServiceClient) {
	courseID := c.Param("courseID")
	staffID := c.Param("staffID")
	if courseID == "" || staffID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "courseID and staffID are required"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if _, err := grpcClient.AddStaffToCourse(ctx, &coursesProtos.AddStaffRequest{CourseID: courseID, StaffID: staffID}); err != nil {
		klog.Errorf("Failed to add staff to course: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add staff to course"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Staff added to course successfully"})
}

func RemoveStaffFromCourseHandler(c *gin.Context, grpcClient coursesProtos.CoursesServiceClient) {
	courseID := c.Param("courseID")
	staffID := c.Param("staffID")
	if courseID == "" || staffID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "courseID and staffID are required"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if _, err := grpcClient.RemoveStaffFromCourse(ctx, &coursesProtos.RemoveStaffRequest{CourseID: courseID, StaffID: staffID}); err != nil {
		klog.Errorf("Failed to remove staff from course: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove staff from course"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Staff removed from course successfully"})
}

func GetCourseStudentsHandler(c *gin.Context, grpcClient coursesProtos.CoursesServiceClient) {
	courseID := c.Param("courseID")
	if courseID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "courseID is required"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := grpcClient.GetCourseStudents(ctx, &coursesProtos.GetCourseStudentsRequest{CourseID: courseID})
	if err != nil {
		klog.Errorf("Failed to get course students: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get course students"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"students": resp.GetStudentsIDs()})
}

func GetCourseStaffHandler(c *gin.Context, grpcClient coursesProtos.CoursesServiceClient) {
	courseID := c.Param("courseID")
	if courseID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "courseID is required"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := grpcClient.GetCourseStaff(ctx, &coursesProtos.GetCourseStaffRequest{CourseID: courseID})
	if err != nil {
		klog.Errorf("Failed to get course staff: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get course staff"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"staff": resp.GetStaffIDs()})
}

func GetStudentCoursesHandler(c *gin.Context, grpcClient coursesProtos.CoursesServiceClient) {
	studentID := c.Param("studentID")

	klog.Infof("Received studentID: %s", studentID)

	if studentID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "studentID is required"})
		return
	}

	// Extract token from the Authorization header
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "No authorization token provided"})
		return
	}

	token := strings.TrimPrefix(authHeader, "Bearer ")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := grpcClient.GetStudentCourses(ctx, &coursesProtos.GetStudentCoursesRequest{
		StudentID: studentID,
		Token:     token,
	})
	if err != nil {
		klog.Errorf("Failed to get student courses: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get student courses"})
		return
	}

	// Get course details for each course ID
	var coursesList []gin.H
	courseIDs := resp.GetCoursesIDs()

	if len(courseIDs) > 0 {
		for _, courseID := range courseIDs {
			// Create a new context for each course request
			courseCtx, courseCancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer courseCancel()
			courseResp, _ := grpcClient.GetCourse(courseCtx, &coursesProtos.GetCourseRequest{CourseID: courseID, Token: token})
			// Add the course with its details to the list
			coursesList = append(coursesList, gin.H{
				"id":          courseID,
				"name":        courseResp.GetCourse().GetCourseName(),
				"description": courseResp.GetCourse().GetDescription(),
				"semester":    courseResp.GetCourse().GetSemester(),
			})
		}
	}

	klog.Infof("Courses list: %v", coursesList)
	klog.Infof("Course IDs: %v", courseIDs)

	c.JSON(http.StatusOK, gin.H{
		"courses":   coursesList,
		"courseIds": courseIDs, // Keep the original IDs array for backward compatibility
	})
}

func GetStaffCoursesHandler(c *gin.Context, grpcClient coursesProtos.CoursesServiceClient) {
	staffID := c.Param("staffID")
	if staffID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "staffID is required"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := grpcClient.GetStaffCourses(ctx, &coursesProtos.GetStaffCoursesRequest{StaffID: staffID})
	if err != nil {
		klog.Errorf("Failed to get staff courses: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get staff courses"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"courses": resp.GetCoursesIDs()})
}

// func AddAnnouncementToCourseHandler(c *gin.Context, grpcClient coursesProtos.CoursesServiceClient) {
// 	courseID := c.Param("courseID")
// 	announcement := c.Param("announcement")
// 	if courseID == "" || announcement == "" {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "courseID and announcement are required"})
// 		return
// 	}

// 	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
// 	defer cancel()

// 	if _, err := grpcClient.AddAnnouncementToCourse(ctx, &coursesProtos.AddAnnouncementRequest{CourseID: courseID, Announcement: announcement}); err != nil {
// 		klog.Errorf("Failed to add announcement to course: %v", err)
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add announcement to course"})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{"message": "Announcement added to course successfully"})
// }

func GetCourseAnnouncementsHandler(c *gin.Context, grpcClient coursesProtos.CoursesServiceClient) {
	courseID := c.Param("courseID")
	if courseID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "courseID is required"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "No authorization token provided"})
		return
	}

	token := strings.TrimPrefix(authHeader, "Bearer ")

	resp, err := grpcClient.GetCourseAnnouncements(ctx, &coursesProtos.GetCourseAnnouncementsRequest{CourseID: courseID, Token: token})
	if err != nil {
		klog.Errorf("Failed to get course announcements: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get course announcements"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"announcements": resp.GetAnnouncements()})
}

func DeleteAnnouncementFromCourseHandler(c *gin.Context, grpcClient coursesProtos.CoursesServiceClient) {
	courseID := c.Param("courseID")
	announcementID := c.Param("announcementID")
	if courseID == "" || announcementID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "courseID and announcementID are required"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if _, err := grpcClient.RemoveAnnouncementFromCourse(ctx, &coursesProtos.RemoveAnnouncementRequest{CourseID: courseID, AnnouncementID: announcementID}); err != nil {
		klog.Errorf("Failed to delete announcement from course: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete announcement from course"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Announcement deleted from course successfully"})
}
