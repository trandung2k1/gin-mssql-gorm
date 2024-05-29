package main

import (
	"gin-gorm/docs"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger" // gin-swagger middleware

	// swagger embed files
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type User struct {
	gorm.Model
	ID int64 `gorm:"primaryKey;autoIncrement:true"`
	// ID       uuid.UUID `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	Name     string `gorm:"type:varchar(60);not null"`
	Email    string `gorm:"index;unique;not null"`
	Password string `gorm:"not null"`
}

// @BasePath /api/v1

// PingExample godoc
// @Summary ping example
// @Schemes
// @Description do ping
// @Tags example
// @Accept json
// @Produce json
// @Success 200 {string} Helloworld
// @Router /example/helloworld [get]
func Helloworld(g *gin.Context) {
	g.JSON(http.StatusOK, "helloworld")
}

func main() {
	gin.SetMode(gin.ReleaseMode)
	docs.SwaggerInfo.BasePath = "/api/v1"
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second,   // Slow SQL threshold
			LogLevel:                  logger.Silent, // Log level
			IgnoreRecordNotFoundError: true,          // Ignore ErrRecordNotFound error for logger
			ParameterizedQueries:      true,          // Don't include params in the SQL log
			Colorful:                  false,         // Disable color
		},
	)
	dsn := os.Getenv("DSN")
	db, err := gorm.Open(sqlserver.Open(dsn), &gorm.Config{Logger: newLogger})
	if err != nil {
		log.Fatal(err)
	}
	// db.AutoMigrate(&User{})
	db.Debug().AutoMigrate(&User{})
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.SetTrustedProxies(nil)

	v1 := r.Group("/api/v1")
	{
		eg := v1.Group("/example")
		{
			eg.GET("/helloworld", Helloworld)
		}

	}
	r.GET("/ping", func(c *gin.Context) {
		// user := User{Name: "Dung", Password: "Tranvandung", Email: "trandungksnb00@gmail.com"}
		// result := db.Create(&user) // pass pointer of data to Creat
		var user []User
		db.Take(&user)
		c.JSON(http.StatusOK, gin.H{
			"data": user,
		})
	})
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	r.Run()

}
