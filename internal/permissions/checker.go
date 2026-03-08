package permissions

import (
	"errors"

	"github.com/sup25/gobuy/internal/user/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Checker struct{}

func NewChecker() *Checker {
	return &Checker{}
}

// Role-based default permissions
var rolePermissions = map[models.UserRole][]Permission{
	// System Admin → platform-level only
	models.RoleSystemAdmin: {
		PlatformManageMerchants,
		PlatformViewAnalytics,
		PlatformManageUsers,
		PlatformManageRoles,
	},
	// Merchant Admin → full shop control
	models.RoleMerchantAdmin: {
		ProductCreate, ProductRead, ProductUpdate, ProductDelete, ProductList,
		OrderCreate, OrderRead, OrderUpdate, OrderDelete, OrderList,
		UserCreate, UserRead, UserUpdate, UserDelete, UserList,
		InventoryCreate, InventoryRead, InventoryUpdate, InventoryDelete,
		PaymentProcess, PaymentRefund, PaymentView,
		PermissionManage,
	},
	// Manager → limited shop control
	models.RoleManager: {
		ProductRead, ProductUpdate, ProductList,
		OrderRead, OrderUpdate, OrderList,
		InventoryRead,
	},
	// Staff → minimal shop control
	models.RoleStaff: {
		ProductRead,
		OrderRead,
	},
	// Customer → buyer permissions
	models.RoleCustomer: {
		ProductRead, ProductList,
		OrderCreate, OrderRead,
		PaymentProcess,
	},
}

// HasPermission checks if a user has a specific permission
func (c *Checker) HasPermission(user *models.User, permission Permission) bool {
	if user == nil {
		return false
	}

	// 1) Check role-based permissions
	if perms, ok := rolePermissions[user.Role]; ok {
		for _, p := range perms {
			if p == permission {
				return true
			}
		}
	}

	// 2) Check custom user permissions (only if they're valid for the role)
	permStr := permission.String()
	for _, p := range user.Permissions {
		if p == permStr && c.isPermissionValidForRole(user.Role, permission) {
			return true
		}
	}

	return false
}

// isPermissionValidForRole ensures custom permissions don't exceed role boundaries
func (c *Checker) isPermissionValidForRole(role models.UserRole, permission Permission) bool {
	// System Admin can't be granted more permissions (already has all platform perms)
	if role == models.RoleSystemAdmin {
		return false
	}

	// Merchant roles can't have platform permissions
	if role == models.RoleMerchantAdmin || role == models.RoleManager || role == models.RoleStaff {
		platformPerms := []Permission{
			PlatformManageMerchants,
			PlatformViewAnalytics,
			PlatformManageUsers,
			PlatformManageRoles,
		}
		for _, pp := range platformPerms {
			if permission == pp {
				return false
			}
		}
	}

	// Customer can't have admin/management permissions
	if role == models.RoleCustomer {
		restrictedPerms := []Permission{
			ProductCreate, ProductUpdate, ProductDelete,
			OrderUpdate, OrderDelete,
			InventoryCreate, InventoryUpdate, InventoryDelete,
			PaymentRefund,
			PermissionManage,
		}
		for _, rp := range restrictedPerms {
			if permission == rp {
				return false
			}
		}
	}

	return false
}

// HasPermissionForResource checks permission with merchant scope validation
func (c *Checker) HasPermissionForResource(user *models.User, permission Permission, resourceMerchantID primitive.ObjectID) bool {
	if !c.HasPermission(user, permission) {
		return false
	}

	// System Admin can access all merchants
	if user.Role == models.RoleSystemAdmin {
		return true
	}

	// Customers don't have merchant scope (they buy from any merchant)
	if user.Role == models.RoleCustomer {
		return true
	}

	// For merchant roles, validate merchant scope
	// Check if user has a merchant ID and it matches the resource
	if user.MerchantID != nil && !resourceMerchantID.IsZero() {
		return *user.MerchantID == resourceMerchantID
	}

	// If user has no merchant ID or resource has no merchant, deny access
	return false
}

// CanGrantPermission checks if a user can grant a specific permission to another user
func (c *Checker) CanGrantPermission(granter *models.User, targetUser *models.User, targetRole models.UserRole, permission Permission) bool {
	if granter == nil {
		return false
	}

	// Must have PermissionManage
	if !c.HasPermission(granter, PermissionManage) {
		return false
	}

	// System Admin → can grant platform permissions
	if granter.Role == models.RoleSystemAdmin {
		return true
	}

	// Merchant Admin → can grant permissions only to Manager/Staff inside same merchant
	if granter.Role == models.RoleMerchantAdmin {
		if targetRole != models.RoleManager && targetRole != models.RoleStaff {
			return false
		}

		// Check merchant scope - both must have merchant IDs and they must match
		if targetUser != nil {
			if targetUser.MerchantID == nil || granter.MerchantID == nil {
				return false
			}
			if *targetUser.MerchantID != *granter.MerchantID {
				return false
			}
		}

		// Can only grant permissions they themselves have
		if !c.HasPermission(granter, permission) {
			return false
		}

		// Validate permission is appropriate for target role
		return c.isPermissionValidForRole(targetRole, permission)
	}

	return false
}

// HasAnyPermission checks if user has any of the given permissions
func (c *Checker) HasAnyPermission(user *models.User, permissions ...Permission) bool {
	for _, perm := range permissions {
		if c.HasPermission(user, perm) {
			return true
		}
	}
	return false
}

// HasAllPermissions checks if user has all of the given permissions
func (c *Checker) HasAllPermissions(user *models.User, permissions ...Permission) bool {
	for _, perm := range permissions {
		if !c.HasPermission(user, perm) {
			return false
		}
	}
	return true
}

// GetDefaultPermissions returns default permissions for a role
func GetDefaultPermissions(role models.UserRole) []string {
	perms := rolePermissions[role]
	strPerms := make([]string, len(perms))
	for i, p := range perms {
		strPerms[i] = p.String()
	}
	return strPerms
}

// ValidatePermissions ensures all permissions in the list are valid
func ValidatePermissions(permissions []string) error {
	validPerms := make(map[string]bool)
	for _, p := range AllPermissions() {
		validPerms[p.String()] = true
	}

	for _, p := range permissions {
		if !validPerms[p] {
			return errors.New("invalid permission: " + p)
		}
	}
	return nil
}
