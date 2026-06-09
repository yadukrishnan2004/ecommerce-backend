package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/yadukrishnan2004/ecommerce-backend/internal/adapter/handler"
	"github.com/yadukrishnan2004/ecommerce-backend/internal/middleware"
	"gorm.io/gorm"
)

func SetUpUserRouter(api fiber.Router, db *gorm.DB, userH *handler.UserHandler) {

	userRoutes := api.Group("/users")

	// PUBLIC ROUTES
	{
		userRoutes.Post(UserSignup, userH.SignUp)
		userRoutes.Post(UserVerify, userH.OtpVerify)
		userRoutes.Post(UserLogin, userH.SignIn)
		userRoutes.Post(UserForgotPassword, userH.Forgotpassword)
		userRoutes.Get(UserAllProducts, userH.GetAll)
		userRoutes.Get(UserSearchProducts, userH.SearchProducts)
		userRoutes.Get(UserFilterProducts, userH.FilterProducts)
		userRoutes.Get(UserProductDetail, userH.GetProduct) 
	}

	// SPECIAL ROUTES (Reset Token Required)
	{
		userRoutes.Post(UserResetPassword, middleware.ResetMiddleware, userH.Resetpassword)
	}

	// AUTHENTICATED ROUTES (JWT Required)

	protected := userRoutes.Group("/")
	protected.Use(middleware.UserMiddleware(db))

	{
		protected.Post(UserLogout, userH.Logout)
		protected.Put(UserUpdateProfile, userH.UpdateProfile)
		protected.Get(UserGetProfile, userH.GetProfile)
		protected.Get(UserOrders, userH.GetOrder)
		protected.Put(UserOrderCancel, userH.CancelOrder)
		protected.Get(UserOrderDetails, userH.GetOrderProduct)
	}
}
