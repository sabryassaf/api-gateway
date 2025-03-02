package controllers

import (
	"context"
	"net/http"
	"time"

	staffProtos "github.com/BetterGR/staff-microservice/protos"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"k8s.io/klog/v2"
)

// InitStaffGRPCClient initializes the staff-microservice gRPC client connection.
func InitStaffGRPCClient(address string) (staffProtos.StaffServiceClient, error) {
	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	return staffProtos.NewStaffServiceClient(conn), nil
}

func CreateStaffMemberHandler(c *gin.Context, grpcClient staffProtos.StaffServiceClient) {
	var staffMember staffProtos.StaffMember
	if err := c.ShouldBindJSON(&staffMember); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if _, err := grpcClient.CreateStaffMember(ctx, &staffProtos.CreateStaffMemberRequest{StaffMember: &staffMember}); err != nil {
		klog.Errorf("Failed to create staff member: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create staff member"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Staff member created successfully"})
}

func GetStaffMemberHandler(c *gin.Context, grpcClient staffProtos.StaffServiceClient) {
	staffID := c.Param("staffID")
	if staffID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "staffID is required"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := grpcClient.GetStaffMember(ctx, &staffProtos.GetStaffMemberRequest{StaffID: staffID})
	if err != nil {
		klog.Errorf("Failed to get staff member: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get staff member"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"staffID":     resp.GetStaffMember().GetStaffID(),
		"firstName":   resp.GetStaffMember().GetFirstName(),
		"lastName":    resp.GetStaffMember().GetLastName(),
		"email":       resp.GetStaffMember().GetEmail(),
		"phoneNumber": resp.GetStaffMember().GetPhoneNumber(),
		"title":       resp.GetStaffMember().GetTitle(),
		"office":      resp.GetStaffMember().GetOffice(),
	})
}

func UpdateStaffMemberHandler(c *gin.Context, grpcClient staffProtos.StaffServiceClient) {
	var staffMember staffProtos.StaffMember
	if err := c.ShouldBindJSON(&staffMember); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if _, err := grpcClient.UpdateStaffMember(ctx, &staffProtos.UpdateStaffMemberRequest{StaffMember: &staffMember}); err != nil {
		klog.Errorf("Failed to update staff member: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update staff member"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Staff member updated successfully"})
}

func DeleteStaffMemberHandler(c *gin.Context, grpcClient staffProtos.StaffServiceClient) {
	staffID := c.Param("staffID")
	if staffID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "staffID is required"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if _, err := grpcClient.DeleteStaffMember(ctx, &staffProtos.DeleteStaffMemberRequest{StaffID: staffID}); err != nil {
		klog.Errorf("Failed to delete staff member: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete staff member"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Staff member deleted successfully"})
}
