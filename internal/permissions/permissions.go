package permissions

type Permission string

// Platform permissions
const (
	PlatformManageMerchants Permission = "platform:manage_merchants"
	PlatformViewAnalytics   Permission = "platform:view_analytics"
	PlatformManageUsers     Permission = "platform:manage_users"
	PlatformManageRoles     Permission = "platform:manage_roles"
)

// Product permissions
const (
	ProductCreate Permission = "product:create"
	ProductRead   Permission = "product:read"
	ProductUpdate Permission = "product:update"
	ProductDelete Permission = "product:delete"
	ProductList   Permission = "product:list"
)

// Order permissions
const (
	OrderCreate Permission = "order:create"
	OrderRead   Permission = "order:read"
	OrderUpdate Permission = "order:update"
	OrderDelete Permission = "order:delete"
	OrderList   Permission = "order:list"
)

// User management permissions
const (
	UserCreate Permission = "user:create"
	UserRead   Permission = "user:read"
	UserUpdate Permission = "user:update"
	UserDelete Permission = "user:delete"
	UserList   Permission = "user:list"
)

// Inventory permissions
const (
	InventoryCreate Permission = "inventory:create"
	InventoryRead   Permission = "inventory:read"
	InventoryUpdate Permission = "inventory:update"
	InventoryDelete Permission = "inventory:delete"
)

// Payment permissions
const (
	PaymentProcess Permission = "payment:process"
	PaymentRefund  Permission = "payment:refund"
	PaymentView    Permission = "payment:view"
)

// Permission management
const (
	PermissionManage Permission = "permission:manage"
)

// AllPermissions returns all available permissions
func AllPermissions() []Permission {
	return []Permission{
		// Products
		ProductCreate, ProductRead, ProductUpdate, ProductDelete, ProductList,
		// Orders
		OrderCreate, OrderRead, OrderUpdate, OrderDelete, OrderList,
		// Users
		UserCreate, UserRead, UserUpdate, UserDelete, UserList,
		// Inventory
		InventoryCreate, InventoryRead, InventoryUpdate, InventoryDelete,
		// Payment
		PaymentProcess, PaymentRefund, PaymentView,
		// Permission management
		PermissionManage,
		// Platform
		PlatformManageMerchants, PlatformViewAnalytics, PlatformManageUsers, PlatformManageRoles,
	}
}

// String converts Permission to string
func (p Permission) String() string {
	return string(p)
}

// IsValid checks if a permission string is valid
func IsValid(perm string) bool {
	for _, p := range AllPermissions() {
		if string(p) == perm {
			return true
		}
	}
	return false
}
