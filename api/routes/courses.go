package routes

import (
	"os"

	"github.com/BetterGR/api-gateway/pkg/controllers"
	courseProtos "github.com/BetterGR/courses-microservice/protos"
	"github.com/gin-gonic/gin"
	"k8s.io/klog/v2"
)

// InitiateCoursesMicroservice initialize courses microservice
func InitiateCoursesMicroservice(router *gin.Engine) {
	_, err := RegisterCoursesRoutes(router)
	if err != nil {
		klog.Fatalf("Failed to register courses routes, %v", err)
	}
}

func RegisterCoursesRoutes(router *gin.Engine) (courseProtos.CoursesServiceClient, error) {
	// Initialize the gRPC client connection.
	coursesAddress := os.Getenv("COURSES_ADDRESS")
	klog.Infof("Courses address: %s", coursesAddress)
	grpcClient, err := controllers.InitCoursesGRPCClient(coursesAddress)
	if err != nil {
		klog.Fatalf("Failed to initialize gRPC client, %v", err)
	}

	// Rest endpoints.
	router.POST("/api/courses", func(c *gin.Context) {
		controllers.CreateCourseHandler(c, grpcClient)
	})

	router.GET("/api/courses/:courseID", func(c *gin.Context) {
		controllers.GetCourseHandler(c, grpcClient)
	})

	router.PUT("/api/courses/:courseID", func(c *gin.Context) {
		controllers.UpdateCourseHandler(c, grpcClient)
	})

	router.DELETE("/api/courses/:courseID", func(c *gin.Context) {
		controllers.DeleteCourseHandler(c, grpcClient)
	})

	router.POST("/api/courses/:courseID/students/:studentID", func(c *gin.Context) {
		controllers.AddStudentToCourseHandler(c, grpcClient)
	})

	router.DELETE("/api/courses/:courseID/students/:studentID", func(c *gin.Context) {
		controllers.RemoveStudentFromCourseHandler(c, grpcClient)
	})

	router.POST("/api/courses/:courseID/staff/:staffID", func(c *gin.Context) {
		controllers.AddStaffToCourseHandler(c, grpcClient)
	})

	router.DELETE("/api/courses/:courseID/staff/:staffID", func(c *gin.Context) {
		controllers.RemoveStaffFromCourseHandler(c, grpcClient)
	})

	router.GET("/api/courses/:courseID/students", func(c *gin.Context) {
		controllers.GetCourseStudentsHandler(c, grpcClient)
	})

	router.GET("/api/courses/:courseID/staff", func(c *gin.Context) {
		controllers.GetCourseStaffHandler(c, grpcClient)
	})

	router.GET("/api/students/:studentID", func(c *gin.Context) {
		controllers.GetStudentCoursesHandler(c, grpcClient)
	})

	router.GET("/api/staff/:staffID", func(c *gin.Context) {
		controllers.GetStaffCoursesHandler(c, grpcClient)
	})

	router.POST("/api/courses/:courseID/announcement", func(c *gin.Context) {
		controllers.AddAnnouncementToCourseHandler(c, grpcClient)
	})

	router.GET("/api/courses/:courseID/announcements", func(c *gin.Context) {
		controllers.GetCourseAnnouncementsHandler(c, grpcClient)
	})

	router.DELETE("/api/courses/:courseID/announcement/:announcementID", func(c *gin.Context) {
		controllers.DeleteAnnouncementFromCourseHandler(c, grpcClient)
	})

	return grpcClient, nil
}
