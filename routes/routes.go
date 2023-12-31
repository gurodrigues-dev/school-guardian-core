package routes

import (
	"gin/config"
	"gin/controllers"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Claims struct {
	Cpf string `json:"cpf"`
	jwt.StandardClaims
}

func HandleRequests() {

	config.LoadEnvironmentVariables()

	var secretKey = []byte(config.GetSecretKeyApi())

	r := gin.Default()

	r.Use(func(c *gin.Context) {

		requestID := uuid.New()

		c.Writer.Header().Set("X-Request-ID", requestID.String())

		c.Set("RequestID", requestID)

		c.Next()

	})

	authMiddleware := func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token não fornecido"})
			c.Abort()
			return
		}

		token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
			return secretKey, nil
		})

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token inválido"})
			c.Abort()
			return
		}

		claims, ok := token.Claims.(*Claims)
		if !ok || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token inválido"})
			c.Abort()
			return
		}

		c.Set("cpf", claims.Cpf)
		c.Set("isAuthenticated", true)
		c.Next()
	}

	// health

	r.GET("api/v1/health", controllers.Health)

	// usuarios

	r.POST("api/v1/users", controllers.CreateUser)

	r.GET("api/v1/users/:cpf", authMiddleware, controllers.GetUser)

	r.PUT("api/v1/users/:cpf", authMiddleware, controllers.UpdateUser)

	r.DELETE("api/v1/users/:cpf", authMiddleware, controllers.DeleteUser)

	r.POST("api/v1/users/login", controllers.AuthenticateUser)

	// password

	r.POST("api/v1/password/recovery", controllers.RecoveryPassword)

	r.POST("api/v1/password/verify", controllers.VerifyIdentityToChangePassword)

	r.PUT("api/v1/password/change", controllers.ChangePassword)

	r.Run()

}
