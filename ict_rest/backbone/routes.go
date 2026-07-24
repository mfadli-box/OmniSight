package backbone

import (
	"net/http"
	"time"

	"ict_rest/skeleton/SM01"
	"ict_rest/skeleton/SM02"
	"ict_rest/skeleton/SM03"
	"ict_rest/skeleton/SM04"
	"ict_rest/skeleton/SM05"
	"ict_rest/skeleton/SM06"
	"ict_rest/skeleton/SP00"
	"ict_rest/skeleton/SP01"
	"ict_rest/skeleton/SP02"
	"ict_rest/skeleton/SP03"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetRouter() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	rest := gin.New()
	rest.Use(RequestID())
	rest.Use(CustomRecovery())
	rest.Use(Logger())
	rest.SetTrustedProxies([]string{"localhost", "172.99.66.6"})
	rest.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:36666", "http://172.99.66.6:36666"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length", "X-Request-ID"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	rest.GET("/", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})
	rest.GET("/rest", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "rest"})
	})
	rest.GET("/hook", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "hook"})
	})
	SetDatabase()

	guest := rest.Group("/rest/guest")
	pages := rest.Group("/rest/pages")
	pages.Use(USLoad())
	authu := rest.Group("/rest/pages")
	authu.Use(USAuth())
	admin := rest.Group("/rest/pages")
	admin.Use(USAuth(), USLock())

	guest.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "guest"})
	})
	pages.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "pages"})
	})

	SP00R := SP00.NRepo(PgSQL)
	SP00U := SP00.NCase(SP00R)
	SP00H := SP00.NHand(SP00U)
	guest.GET("/SP00", SP00H.ListHrisCompany)
	guest.POST("/SP00", SP00H.Login)
	pages.DELETE("/SP00", SP00H.Logout)

	SP01R := SP01.NRepo(PgSQL)
	SP01U := SP01.NCase(SP01R)
	SP01H := SP01.NHand(SP01U)
	authu.GET("/SP01/company", SP01H.ListUserCompany)
	authu.GET("/SP01/module", SP01H.ListUserModule)

	SP02R := SP02.NRepo(PgSQL)
	SP02U := SP02.NCase(SP02R)
	SP02H := SP02.NHand(SP02U)
	authu.PUT("/SP02", SP02H.ChangePassword)

	SP03R := SP03.NRepo(PgSQL)
	SP03U := SP03.NCase(SP03R)
	SP03H := SP03.NHand(SP03U)
	authu.GET("/SP03", SP03H.ListActions)

	SM01R := SM01.NRepo(PgSQL)
	SM01U := SM01.NCase(SM01R)
	SM01H := SM01.NHand(SM01U)
	admin.GET("/SM01", SM01H.ListUser)
	admin.POST("/SM01", USLogs("SM01"), SM01H.CreateUser)
	admin.PUT("/SM01/:id", USLogs("SM01"), SM01H.UpdateUser)
	admin.GET("/SM01/hris", SM01H.ListHRISCompany)
	admin.GET("/SM01/company", SM01H.ListAllCompany)
	admin.GET("/SM01/module", SM01H.ListAllModule)
	admin.POST("/SM01/:id/company", USLogs("SM01"), SM01H.AssignCompany)
	admin.GET("/SM01/:id/company", SM01H.ListUserCompany)
	admin.POST("/SM01/:id/company/assign", USLogs("SM01"), SM01H.CreateUserCompany)
	admin.PUT("/SM01/:id/company/:companyId", USLogs("SM01"), SM01H.UpdateUserCompany)
	admin.DELETE("/SM01/:id/company/:companyId", USLogs("SM01"), SM01H.DeleteUserCompany)
	admin.GET("/SM01/:id/privilege", SM01H.ListUserPrivilege)
	admin.POST("/SM01/:id/privilege", USLogs("SM01"), SM01H.CreateUserPrivilege)
	admin.PUT("/SM01/:id/privilege/:privilegeId", USLogs("SM01"), SM01H.UpdateUserPrivilege)
	admin.DELETE("/SM01/:id/privilege/:privilegeId", USLogs("SM01"), SM01H.DeleteUserPrivilege)
	admin.GET("/SM01/type", SM01H.ListLocationTypeByCompany)
	admin.GET("/SM01/:id/location", SM01H.ListUserLocation)
	admin.POST("/SM01/:id/location", USLogs("SM01"), SM01H.CreateUserLocation)
	admin.PUT("/SM01/:id/location/:locationId", USLogs("SM01"), SM01H.UpdateUserLocation)
	admin.DELETE("/SM01/:id/location/:locationId", USLogs("SM01"), SM01H.DeleteUserLocation)

	SM02R := SM02.NRepo(PgSQL)
	SM02U := SM02.NCase(SM02R)
	SM02H := SM02.NHand(SM02U)
	admin.GET("/SM02", SM02H.ListModule)
	admin.POST("/SM02", USLogs("SM02"), SM02H.CreateModule)
	admin.PUT("/SM02/:id", USLogs("SM02"), SM02H.UpdateModule)

	SM03R := SM03.NRepo(PgSQL)
	SM03U := SM03.NCase(SM03R)
	SM03H := SM03.NHand(SM03U)
	admin.GET("/SM03", SM03H.ListCompany)
	admin.POST("/SM03", USLogs("SM03"), SM03H.CreateCompany)
	admin.PUT("/SM03/:id", USLogs("SM03"), SM03H.UpdateCompany)
	admin.GET("/SM03/module", SM03H.ListModule)
	admin.GET("/SM03/:id/module", SM03H.ListCompanyModule)
	admin.POST("/SM03/:id/module", USLogs("SM03"), SM03H.AssignCompanyModule)
	admin.PUT("/SM03/:id/module/:moduleId", USLogs("SM03"), SM03H.UpdateCompanyModule)
	admin.GET("/SM03/:id/type", SM03H.ListLocationType)
	admin.POST("/SM03/:id/type", USLogs("SM03"), SM03H.CreateLocationType)
	admin.PUT("/SM03/:id/type/:typeId", USLogs("SM03"), SM03H.UpdateLocationType)
	admin.DELETE("/SM03/:id/type/:typeId", USLogs("SM03"), SM03H.DeleteLocationType)

	SM04R := SM04.NRepo(PgSQL)
	SM04U := SM04.NCase(SM04R)
	SM04H := SM04.NHand(SM04U)
	admin.GET("/SM04", SM04H.ListSignatureType)
	admin.POST("/SM04", USLogs("SM04"), SM04H.CreateSignatureType)
	admin.PUT("/SM04/:id", USLogs("SM04"), SM04H.UpdateSignatureType)
	admin.GET("/SM04/user", SM04H.ListUser)

	SM05R := SM05.NRepo(PgSQL)
	SM05U := SM05.NCase(SM05R)
	SM05H := SM05.NHand(SM05U)
	admin.GET("/SM05", SM05H.ListSession)
	admin.GET("/SM05/:id", SM05H.GetSessionDetail)
	admin.DELETE("/SM05/:id", USLogs("SM05"), SM05H.RevokeSession)

	SM06R := SM06.NRepo(PgSQL)
	SM06U := SM06.NCase(SM06R)
	SM06H := SM06.NHand(SM06U)
	admin.GET("/SM06/type", SM06H.ListLocationType)
	admin.POST("/SM06/type", USLogs("SM06"), SM06H.CreateLocationType)
	admin.PUT("/SM06/type/:id", USLogs("SM06"), SM06H.UpdateLocationType)
	admin.DELETE("/SM06/type/:id", USLogs("SM06"), SM06H.DeleteLocationType)
	admin.GET("/SM06/type/select", SM06H.ListLocationTypeSelect)
	pages.GET("/SM06", SM06H.ListLocation)
	pages.POST("/SM06", USLogs("SM06"), SM06H.CreateLocation)
	pages.PUT("/SM06/:id", USLogs("SM06"), SM06H.UpdateLocation)
	pages.DELETE("/SM06/:id", USLogs("SM06"), SM06H.DeleteLocation)
	pages.GET("/SM06/select", SM06H.ListLocationSelect)

	return rest
}
