package main

import (
	"context"
	"log"
	"modular_monolith/config"
	"modular_monolith/internal/category"
	"modular_monolith/internal/cloudinaryutil"
	"modular_monolith/internal/profile"
	"modular_monolith/internal/user"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func main() {
	// Load .env file only if it exists (for local development)
	if _, err := os.Stat(".env"); err == nil {
		if err := godotenv.Load(); err != nil {
			log.Printf("Warning: Error loading .env file: %v", err)
		} else {
			log.Println("Successfully loaded .env file")
		}
	} else {
		log.Println("No .env file found, using environment variables")
	}

	cfg := config.LoadConfig()

	cld, err := cloudinaryutil.NewCloudinaryUploader(cfg.Clouldinary)
	if err != nil {
		panic(err)
	}

	mongoClient, err := connectToMongoDB(cfg.MongoURI)
	if err != nil {
		panic(err)
	}

	defer func() {
		if err := mongoClient.Disconnect(context.Background()); err != nil {
			panic(err)
		}
	}()

	profilesCollection := mongoClient.Database(cfg.MongoDB).Collection("profiles")
	profileRepository := profile.NewProfileRepository(profilesCollection)
	profileService := profile.NewProfileService(profileRepository, cld)
	profileHandler := profile.NewProfileHandler(profileService)

	usersCollection := mongoClient.Database(cfg.MongoDB).Collection("users")
	userRepository := user.NewUserRepository(usersCollection)
	userService := user.NewUserService(userRepository, profileService)
	userHandler := user.NewUserHandler(userService)

	categories := mongoClient.Database(cfg.MongoDB).Collection("categories")
	categoryRepository := category.NewCategoryRepository(categories)
	categoryService := category.NewCategoryService(categoryRepository)
	categoryHandler := category.NewCategoryHandler(categoryService)

	r := gin.Default()
	r.LoadHTMLGlob("web/*")
	r.GET("/login", func(c *gin.Context) {
		c.HTML(http.StatusOK, "login.html", gin.H{})
	})

	user.RegisterRoutes(r, userHandler)
	profile.RegisterRoutes(r, profileHandler)
	category.RegisterRoutes(r, categoryHandler)

	// Get port from environment variable (Railway sets this automatically)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8003" // fallback port
	}

	log.Printf("Server starting on port %s", port)
	r.Run(":" + port)
}

func connectToMongoDB(uri string) (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		log.Println("Failed to connect to MongoDB")
		return nil, err
	}

	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		log.Println("Failed to ping to MongoDB")
		return nil, err
	}

	log.Println("Successfully connected to MongoDB")
	return client, nil
}