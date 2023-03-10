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
	e.GET("/get-post/:post_id", controllers.FetchPost)
	e.GET("/get-ownedcomments/:user_id", controllers.FetchOwnedComments)
	e.GET("/get-postshome", controllers.FetchAllPosts)
	e.GET("/get-ownedposts/:user_id", controllers.FetchOwnedPosts)
	e.GET("/search-posts/:search_key", controllers.SearchPosts)

	e.GET("/get-communitycategories", controllers.FetchAllCommunityCategory)
	e.GET("/search-communitycategories/:search_key", controllers.SearchCommunityCategory)

	e.GET("/get-communities", controllers.FetchAllCommunities)
	e.GET("/get-communitieshome", controllers.FetchCommunities)
	e.GET("/get-communitiescategory/:category_id", controllers.FetchCommunitiesCategory)
	e.GET("/search-communitiescategory/:category_id/:search_key", controllers.SearchCommunitiesCategory)
	e.GET("/get-ownedcommunities/:user_id", controllers.FetchOwnedCommunities)

	e.GET("/get-communitymembers/:community_id", controllers.FetchAllCommunityMember)
	
	e.POST("/register-user", controllers.RegisterUser)
	e.POST("/login-user", controllers.CheckLogin)
	e.POST("/user-profile", controllers.UserProfile)
	
	e.POST("/store-post", controllers.StorePost)
	e.POST("/store-comment", controllers.StoreComment)
	e.POST("/store-community", controllers.StoreCommunity)
	e.POST("/join-community", controllers.JoinCommunity)
	e.POST("/store-communitycategory", controllers.StoreCommunityCategory)
	
	e.PATCH("/update-post", controllers.UpdatePost)
	e.PATCH("/update-community", controllers.UpdateCommunity)
	e.PATCH("/update-profile", controllers.UpdateProfile)
	e.PATCH("/update-comment", controllers.UpdateComment)
	
	e.DELETE("/remove-member/:community_member_id", controllers.DeleteMember)
	e.DELETE("/signout-community/:user_id/:community_id", controllers.SignoutCommunity)
	e.DELETE("/delete-post/:post_id", controllers.DeletePost)
	e.DELETE("/delete-comment/:comment_id", controllers.DeleteComment)
	e.DELETE("/delete-community/:community_id", controllers.DeleteCommunity)

	return e
}
