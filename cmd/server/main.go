package main

import (
	"context"
	"log"
	"modular_monolith/config"
	"modular_monolith/helper"
	"modular_monolith/internal/blog"
	"modular_monolith/internal/cart"
	"modular_monolith/internal/category"
	"modular_monolith/internal/coupon"
	"modular_monolith/internal/order"
	"modular_monolith/internal/payment"
	"modular_monolith/internal/product"
	"modular_monolith/internal/profile"
	review "modular_monolith/internal/reviews"
	"modular_monolith/internal/user"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func main() {
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

	cld, err := helper.NewCloudinaryUploader(cfg.Clouldinary)
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

	r := gin.Default()
	r.Use(func(c *gin.Context) {
		origin := c.GetHeader("Origin")
		log.Printf("Request from origin: '%s'", origin)
		log.Printf("Request method: %s", c.Request.Method)
		log.Printf("Request path: %s", c.Request.URL.Path)
		c.Next()
	})

	// Cấu hình CORS đơn giản - cho phép tất cả để test
	r.Use(func(c *gin.Context) {
		
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS, PATCH")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization, Accept, X-Requested-With")
		c.Header("Access-Control-Max-Age", "86400")

		if c.Request.Method == "OPTIONS" {
			log.Printf("Handling OPTIONS preflight request")
			c.AbortWithStatus(200)
			return
		}

		c.Next()
	})

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

	reviews := mongoClient.Database(cfg.MongoDB).Collection("reviews")
	reviewsRepository := review.NewReviewRepository(reviews)
	reviewsService := review.NewReviewService(reviewsRepository, userRepository)
	reviewsHandler := review.NewReviewHandler(reviewsService)

	products := mongoClient.Database(cfg.MongoDB).Collection("products")
	productsRepository := product.NewProductRepository(products)
	productsService := product.NewProductService(productsRepository, cld, reviewsService)
	productsHandler := product.NewProductHandler(productsService)

	carts := mongoClient.Database(cfg.MongoDB).Collection("carts")
	cartsRepository := cart.NewCartRepository(carts)
	cartsService := cart.NewCartService(cartsRepository, productsRepository)
	cartsHandler := cart.NewCartHandler(cartsService)

	orders := mongoClient.Database(cfg.MongoDB).Collection("orders")
	ordersRepository := order.NewOrderRepository(orders)
	ordersService := order.NewOrderService(ordersRepository, cartsService)
	ordersHandler := order.NewOrderHandler(ordersService)

	payments := mongoClient.Database(cfg.MongoDB).Collection("payments")
	paymentsRepository := payment.NewPaymentRepository(payments)
	paymentsService := payment.NewPaymentService(paymentsRepository, ordersRepository)
	paymentsHandler := payment.NewPaymentHandler(paymentsService)

	coupons := mongoClient.Database(cfg.MongoDB).Collection("coupons")
	couponsRepository := coupon.NewCouponRepository(coupons)
	couponsService := coupon.NewCouponService(couponsRepository, userService)
	couponsHandler := coupon.NewCouponHandler(couponsService)

	blogs := mongoClient.Database(cfg.MongoDB).Collection("blogs")
	blogsRepository := blog.NewBlogRepository(blogs)
	blogsService := blog.NewBlogService(blogsRepository, cld)
	blogsHandler := blog.NewBlogHandler(blogsService)

	blog.RegisterRoutes(r, blogsHandler)
	coupon.RegisterRoutes(r, couponsHandler)
	payment.RegisterRoutes(r, paymentsHandler)
	order.RegisterRoutes(r, ordersHandler)
	review.RegisterRoutes(r, reviewsHandler)
	user.RegisterRoutes(r, userHandler)
	profile.RegisterRoutes(r, profileHandler)
	category.RegisterRoutes(r, categoryHandler)
	product.RegisterRoutes(r, productsHandler)
	cart.RegisterRoutes(r, cartsHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8003"
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
