package router

import (
	"pkm/handler"
	"pkm/middleware"
	"pkm/utils"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func SetupRoutes(h *handler.Handler, db *gorm.DB) *echo.Echo {
	e := echo.New()
	e.Validator = utils.NewValidator()

	v := e.Group("/v1", middleware.Authenticate(db))

	// User
	user := v.Group("/user")
	user.GET("", h.GetUser)
	user.POST("/login", h.Login)
	user.POST("/register", h.Register)
	user.POST("/update/:id", h.UpdateUser)
	user.DELETE("/delete/:id", h.DeleteUser)

	// Admin
	admin := v.Group("/admin")
	admin.GET("", h.GetAdmin)
	admin.GET("/admins", h.GetAllAdmins)
	admin.POST("", h.CreateAdmin)
	admin.POST("/update/:id", h.UpdateAdmin)
	admin.DELETE("/delete/:id", h.DeleteAdmin)

	// Card
	card := v.Group("/card")
	card.POST("/search", h.SearchCard)
	card.GET("/detail", h.CardDetail)
	card.POST("/update", h.CardUpdate)
	card.POST("/update/setNumber", h.CardUpdateSetNumber)

	// provider
	provider := v.Group("/provider")

	psa := provider.Group("/psa")
	psa.POST("", h.InsertPSA)

	tag := provider.Group("/tag")
	tag.POST("", h.InsertTAG)

	bgs := provider.Group("/bgs")
	bgs.POST("", h.InsertBGS)

	cgc := provider.Group("/cgc")
	cgc.POST("", h.InsertCGC)

	priceCharting := provider.Group("/pricecharting")
	priceCharting.POST("", h.PriceChartCards)

	// set
	set := v.Group("/set")
	set.POST("/create", h.CreateSet)
	set.GET("/list", h.SetList)
	set.POST("/update_card_setname", h.UpdateEachCardBasedOnSetName)

	return e
}
