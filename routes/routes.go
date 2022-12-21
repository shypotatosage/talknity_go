package routes

import (
	"talknity/controllers"

	"github.com/labstack/echo/v4"
)

func Init() *echo.Echo {
	e := echo.New()
	
	e.GET("/get-posts", controllers.FetchAllPosts)
	e.GET("/get-communitycategories", controllers.FetchAllCommunityCategory)
	e.GET("/get-communities", controllers.FetchAllCommunities)
	e.GET("/get-communitymembers", controllers.FetchAllCommunityMember)
	
	e.POST("/register-user", controllers.RegisterUser)
	e.POST("/login-user", controllers.CheckLogin)
	
	e.POST("/store-post", controllers.StorePost)
	
	e.PATCH("/update-post", controllers.UpdatePost)
	
	e.DELETE("/delete-post", controllers.DeletePost)

	return e
}
