package routes

import (
	"os"

	"github.com/BetterGR/api-gateway/pkg/controllers"
	gradesProtos "github.com/BetterGR/grades-microservice/protos"
	"github.com/gin-gonic/gin"
	"k8s.io/klog/v2"
)

// InitiateGradesMicroservice initialize grades microservice
func InitiateGradesMicroservice(router *gin.Engine) {
	_, err := RegisterGradesRoutes(router)
	if err != nil {
		klog.Fatalf("Failed to register grades routes, %v", err)
	}
}

func RegisterGradesRoutes(router *gin.Engine) (gradesProtos.GradesServiceClient, error) {
	// Initialize the gRPC client connection.
	gradesAddress := os.Getenv("GRADES_ADDRESS")
	klog.Infof("Grades address: %s", gradesAddress)
	grpcClient, err := controllers.InitGradesGRPCClient(gradesAddress)
	if err != nil {
		klog.Fatalf("Failed to initialize gRPC client, %v", err)
	}

	// Rest endpoints.

	// Get all grades for a student in a specific course for a specific semester.
	router.GET("/api/grades/:semester/:student_id/:courseId", func(c *gin.Context) {
		controllers.GetStudentCourseGradesHandler(c, grpcClient)
	})

	// Get all grades for a student in a specific semester.
	router.GET("/api/grades/:semester/:student_id", func(c *gin.Context) {
		controllers.GetStudentSemesterGradesHandler(c, grpcClient)
	})

	router.POST("/api/grades/homework", func(c *gin.Context) {
		controllers.AddHomeworkGradeHandler(c, grpcClient)
	})

	router.POST("/api/grades/exam", func(c *gin.Context) {
		controllers.AddExamGradeHandler(c, grpcClient)
	})

	router.PUT("/api/grades/homework", func(c *gin.Context) {
		controllers.UpdateHomeworkGradeHandler(c, grpcClient)
	})

	router.PUT("/api/grades/exam", func(c *gin.Context) {
		controllers.UpdateExamGradeHandler(c, grpcClient)
	})

	router.DELETE("/api/grades/homework", func(c *gin.Context) {
		controllers.DeleteHomeworkGradeHandler(c, grpcClient)
	})

	router.DELETE("/api/grades/exam", func(c *gin.Context) {
		controllers.DeleteExamGradeHandler(c, grpcClient)
	})

	return grpcClient, nil
}
