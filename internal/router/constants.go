package router

// User / Auth Route Constants
const (
	UserSignup         = "/signup"
	UserVerify         = "/verify"
	UserLogin          = "/login"
	UserForgotPassword = "/forgot-password"
	UserAllProducts    = "/allproducts"
	UserSearchProducts = "/search"
	UserFilterProducts = "/filter"
	UserProductDetail  = "/products/:id"
	UserResetPassword  = "/reset-password"
	UserLogout         = "/logout"
	UserUpdateProfile  = "/profile"
	UserGetProfile     = "/profile"
	UserOrders         = "/:id/orders"
	UserOrderCancel    = "/:id/cancel"
	UserOrderDetails   = "/:id/orders/details"
)

// Admin Route Constants
const (
	AdminUpdateUser         = "/users/:id"
	AdminUpdateProdStatus   = "/products/:id/status"
	AdminBlockUser          = "/users/:id/block"
	AdminAddProduct         = "/products"
	AdminUpdateProduct      = "/products/:id"
	AdminGetAllProducts     = "/products"
	AdminGetProductsStatus  = "/products/status/:status"
	AdminSearchProducts     = "/products/search"
	AdminFilterProducts     = "/products/filter"
	AdminGetProduct         = "/products/:id"
	AdminGetAllOrders       = "/orders"
	AdminGetOrderDetails    = "/orders/:id"
	AdminDeleteProduct      = "/products/:id"
	AdminDeleteUser         = "/users/:id"
	AdminUpdateOrderStatus  = "/orders/:id"
	AdminUpdateOrderStatus2 = "/orders/status/:id"
	AdminSearchUsers        = "/users/search"
	AdminGetAllUsers        = "/users"
	AdminDashboardGraphs    = "/dashboard-graphs"
	AdminGetUser            = "/users/:id"
	AdminGetUserCart        = "/users/:id/cart"
	AdminGetUserWishlist    = "/users/:id/wishlist"
	AdminGetUserAddresses   = "/users/:id/addresses"
	AdminCreateCategory     = "/categories"
	AdminGetAllCategories   = "/categories"
	AdminUpdateCategory     = "/categories/:id"
	AdminDeleteCategory     = "/categories/:id"
	AdminLowStock           = "/inventory/low-stock"
	AdminDashboardKpis      = "/dashboard/kpis"
)

// Cart Route Constants
const (
	CartGet    = "/"
	CartAdd    = "/add"
	CartUpdate = "/:id"
	CartRemove = "/:id"
	CartClear  = "/clear"
)

// Wishlist Route Constants
const (
	WishlistGet    = "/"
	WishlistAdd    = "/:id"
	WishlistRemove = "/:id"
	WishlistClear  = "/clear"
)

// Order Route Constants
const (
	OrderGet           = "/"
	OrderGetByID       = "/:id"
	OrderBuyNow        = "/buy-now"
	OrderPlace         = "/"
	OrderVerifyPayment = "/verify-payment"
)

// Address Route Constants
const (
	AddressGet  = "/"
	AddressPost = "/"
)
