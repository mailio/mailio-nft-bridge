package router

import (
	"github.com/chryscloud/go-microkit-plugins/auth"
	"github.com/chryscloud/go-microkit-plugins/endpoints"
	jwtModels "github.com/chryscloud/go-microkit-plugins/models/jwt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/mailio/mailio-nft-server/api"
	"github.com/mailio/mailio-nft-server/config"
	"github.com/mailio/mailio-nft-server/model"
	"github.com/mailio/mailio-nft-server/service"
)

func ConfigAPI(router *gin.Engine, env *model.Environment, conf *config.Config) *gin.Engine {

	// initialize services
	nftCatalogService := service.NewNftCatalog(env)
	userService := service.NewUserService(env)
	nftClaimService := service.NewNftClaimService(env)
	nftImageService := service.NewNftImagesService(env)

	// intialize API endpoints
	nftCatalogApi := api.NewNftCatalogAPI(nftCatalogService)
	userApi := api.NewUserAPI(userService)
	nftImageApi := api.NewNftImagesAPI(nftImageService)
	claimApi := api.NewClaimAPI(nftClaimService, nftCatalogService)

	// enable cors
	router.Use(cors.New(cors.Config{
		AllowCredentials: true,
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"*"},
		AllowHeaders:     []string{"*"},
	}))

	// the default ping endpoint (sometimes needed for healthchecks in kubernetes environments)
	root := router.Group("/")
	{
		root.GET("/", endpoints.PingEndpoint)
	}

	// public APIs
	public := router.Group("/api/v1")
	{
		public.GET("/catalog/:id", nftCatalogApi.GetCatalog)
		public.GET("/catalog", nftCatalogApi.ListCatalogs)
		public.POST("/login", userApi.Login)
		public.GET("/claim/:address/payload/:catalogId", claimApi.SigningPayload)
		public.POST("/claim", claimApi.MintClaim)
		public.GET("/user/claims/:walletaddress", claimApi.ListClaimsByUser)
	}

	// init JWT Authentication Middleware for private endpoints
	keys := func(token *jwt.Token) (interface{}, error) {
		return []byte(conf.JWTToken.SecretKey), nil
	}
	auhtMiddleware := auth.JwtMiddleware(&conf.YamlConfig, &jwtModels.UserClaim{}, jwt.SigningMethodHS256, keys)

	// Methods accessible only to a registered user
	private := router.Group("/api/v1", auhtMiddleware)
	{
		private.POST("/catalog", nftCatalogApi.PutCatalog)
		private.PUT("/catalog", nftCatalogApi.PutCatalog)
		private.GET("/bridge/balance", claimApi.GetBridgeBalance)
		private.POST("/nftimage/upload", nftImageApi.Upload)
		private.GET("/nftimage/list", nftImageApi.List)
		private.DELETE("/nftimage/:hash", nftImageApi.RemovePin)
		private.GET("/claim", claimApi.ListClaims)
	}
	return router
}
