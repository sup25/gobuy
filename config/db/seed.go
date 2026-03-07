package db

import (
	"context"
	"fmt"
	"time"

	"github.com/sup25/gobuy/internal/permissions"
	"github.com/sup25/gobuy/internal/user/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

// SeedAdminAndManager inserts default system admin and a sample merchant admin
func SeedAdminAndManager(client *mongo.Client) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	userCollection := client.Database(dbName).Collection("users")

	// Create a demo merchant ObjectID (you can generate a new one or use a fixed one)
	demoMerchantID := primitive.NewObjectID()

	// System Admin permissions - ONLY platform management
	systemAdminPermissions := []string{
		// Platform permissions only
		permissions.PlatformManageMerchants.String(),
		permissions.PlatformViewAnalytics.String(),
		permissions.PlatformManageUsers.String(),
		permissions.PlatformManageRoles.String(),
		// Permission management
		permissions.PermissionManage.String(),
	}

	// Merchant Admin permissions - FULL control over their shop
	merchantAdminPermissions := []string{
		// Products - full access
		permissions.ProductCreate.String(),
		permissions.ProductRead.String(),
		permissions.ProductUpdate.String(),
		permissions.ProductDelete.String(),
		permissions.ProductList.String(),
		// Orders - full access
		permissions.OrderCreate.String(),
		permissions.OrderRead.String(),
		permissions.OrderUpdate.String(),
		permissions.OrderDelete.String(),
		permissions.OrderList.String(),
		// Inventory - full access
		permissions.InventoryCreate.String(),
		permissions.InventoryRead.String(),
		permissions.InventoryUpdate.String(),
		permissions.InventoryDelete.String(),
		// Users - full access (their staff)
		permissions.UserCreate.String(),
		permissions.UserRead.String(),
		permissions.UserUpdate.String(),
		permissions.UserDelete.String(),
		permissions.UserList.String(),
		// Payment - full access
		permissions.PaymentProcess.String(),
		permissions.PaymentRefund.String(),
		permissions.PaymentView.String(),
	}

	// Manager permissions - limited operational access
	managerPermissions := []string{
		// Products - full access
		permissions.ProductCreate.String(),
		permissions.ProductRead.String(),
		permissions.ProductUpdate.String(),
		permissions.ProductDelete.String(),
		permissions.ProductList.String(),
		// Orders - full access
		permissions.OrderCreate.String(),
		permissions.OrderRead.String(),
		permissions.OrderUpdate.String(),
		permissions.OrderDelete.String(),
		permissions.OrderList.String(),
		// Inventory - full access
		permissions.InventoryCreate.String(),
		permissions.InventoryRead.String(),
		permissions.InventoryUpdate.String(),
		permissions.InventoryDelete.String(),
		// Users - read only
		permissions.UserRead.String(),
		permissions.UserList.String(),
		// Payment - view and process only (no refunds)
		permissions.PaymentProcess.String(),
		permissions.PaymentView.String(),
	}

	usersToSeed := []struct {
		User        models.User
		Password    string
		Permissions []string
	}{
		{
			User: models.User{
				Name:       "System Admin",
				Email:      "admin@gobuy.com",
				Role:       models.RoleSystemAdmin,
				MerchantID: nil, // System admin doesn't belong to any merchant
			},
			Password:    "@Hala_madrid123",
			Permissions: systemAdminPermissions,
		},
		{
			User: models.User{
				Name:       "Demo Merchant",
				Email:      "merchant@example.com",
				Role:       models.RoleMerchantAdmin,
				MerchantID: &demoMerchantID, // Pointer to ObjectID
			},
			Password:    "@Hala_madrid123",
			Permissions: merchantAdminPermissions,
		},
		{
			User: models.User{
				Name:       "Demo Manager",
				Email:      "manager@example.com",
				Role:       models.RoleManager,
				MerchantID: &demoMerchantID, // Same merchant as the merchant admin
			},
			Password:    "@Hala_madrid123",
			Permissions: managerPermissions,
		},
	}

	for _, seedData := range usersToSeed {
		// Check if user already exists
		count, err := userCollection.CountDocuments(ctx, bson.M{"email": seedData.User.Email})
		if err != nil {
			return fmt.Errorf("failed to check existing users: %v", err)
		}

		if count > 0 {
			fmt.Printf("%s already exists, skipping\n", seedData.User.Email)
			continue
		}

		// Hash password
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(seedData.Password), bcrypt.DefaultCost)
		if err != nil {
			return fmt.Errorf("failed to hash password: %v", err)
		}

		seedData.User.Password = string(hashedPassword)
		seedData.User.CreatedAt = time.Now()
		seedData.User.UpdatedAt = time.Now()
		seedData.User.IsEmailVerified = true
		seedData.User.Permissions = seedData.Permissions

		_, err = userCollection.InsertOne(ctx, seedData.User)
		if err != nil {
			return fmt.Errorf("failed to insert %s: %v", seedData.User.Email, err)
		}

		merchantInfo := "no merchant"
		if seedData.User.MerchantID != nil {
			merchantInfo = seedData.User.MerchantID.Hex()
		}

		fmt.Printf("%s (%s) user seeded with %d permissions [Merchant: %s]\n",
			seedData.User.Email,
			seedData.User.Role,
			len(seedData.Permissions),
			merchantInfo)
	}

	return nil
}
