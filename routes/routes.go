package routes

import (
	"net/http"
	"talknity/controllers"

	"github.com/labstack/echo/v4"
)

func Init() *echo.Echo {
	e := echo.New()

	fs := http.FileServer(http.Dir("./images"))
	e.GET("/images/*", echo.WrapHandler(http.StripPrefix("/images/", fs)))
	
	e.GET("/get-posts", controllers.FetchAllPosts)
	e.GET("/get-postshome", controllers.FetchAllPosts)
	e.GET("/get-ownedposts/:user_id", controllers.FetchOwnedPosts)
	e.GET("/search-posts/:search_key", controllers.SearchPosts)

	e.GET("/get-communitycategories", controllers.FetchAllCommunityCategory)
	e.GET("/search-communitycategories/:search_key", controllers.SearchCommunityCategory)
	e.GET("/get-communities", controllers.FetchAllCommunities)
	e.GET("/get-communitieshome", controllers.FetchCommunities)
	e.GET("/get-ownedcommunities/:user_id", controllers.FetchOwnedCommunities)

	e.GET("/get-communitymembers", controllers.FetchAllCommunityMember)
	e.GET("/get-comments", controllers.FetchAllComments)
	
	e.POST("/register-user", controllers.RegisterUser)
	e.POST("/login-user", controllers.CheckLogin)
	
	e.POST("/store-post", controllers.StorePost)
	e.POST("/store-community", controllers.StoreCommunity)
	e.POST("/store-communitycategory", controllers.StoreCommunityCategory)
	
	e.PATCH("/update-post", controllers.UpdatePost)
	
	e.DELETE("/delete-post", controllers.DeletePost)

	return e
}
