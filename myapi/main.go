package main

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"

	_ "myapi/docs" // Replace with your actual module path

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Product and Order API with JWT Auth
// @version 1.0
// @description Swagger-enabled API with Login, Token Auth, Product and Order Management
// @host localhost:8080
// @BasePath /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

var jwtKey = []byte("supersecretkey")

type Credentials struct {
	Username string `json:"username" example:"admin"`
	Password string `json:"password" example:"admin123"`
}

type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

type Product struct {
	ID          int     `json:"id" example:"1"`
	Name        string  `json:"name" example:"Laptop"`
	Description string  `json:"description" example:"Powerful laptop"`
	Price       float64 `json:"price" example:"999.99"`
	InStock     bool    `json:"in_stock" example:"true"`
}

type Order struct {
	OrderID   int       `json:"order_id" example:"101"`
	ProductID int       `json:"product_id" example:"1"`
	Quantity  int       `json:"quantity" example:"2"`
	OrderDate time.Time `json:"order_date" example:"2024-01-01T15:04:05Z"`
}

var (
	products = []Product{
		{ID: 1, Name: "Laptop", Description: "Powerful laptop", Price: 999.99, InStock: true},
		{ID: 2, Name: "Phone", Description: "Android smartphone", Price: 499.99, InStock: true},
	}
	orders = []Order{}
)

func main() {
	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Public routes
	r.POST("/login", LoginHandler)
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Protected routes
	auth := r.Group("/")
	auth.Use(JWTMiddleware())
	{
		auth.GET("/product", GetProductsHandler)
		auth.POST("/product", AddProductHandler)
		auth.PUT("/product/:id", UpdateProductHandler)
		auth.DELETE("/product/:id", DeleteProductHandler)

		auth.GET("/order", GetOrdersHandler)
		auth.POST("/order", CreateOrderHandler)
	}

	r.Run(":8080")
}

// ================== AUTH =====================

// @Summary Login to get token
// @Tags auth
// @Accept json
// @Produce json
// @Param credentials body Credentials true "User credentials"
// @Success 200 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /login [post]
func LoginHandler(c *gin.Context) {
	var creds Credentials
	if err := c.ShouldBindJSON(&creds); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}
	if creds.Username != "admin" || creds.Password != "admin123" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	expiration := time.Now().Add(1 * time.Hour)
	claims := &Claims{
		Username: creds.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiration),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString(jwtKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Token generation failed"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": tokenStr})
}

func JWTMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if len(authHeader) < 8 || authHeader[:7] != "Bearer " {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Missing or invalid token"})
			return
		}
		tokenStr := authHeader[7:]
		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})
		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}
		c.Next()
	}
}

// ================== PRODUCTS =====================

// @Summary Get all products
// @Tags product
// @Security BearerAuth
// @Produce json
// @Success 200 {array} Product
// @Router /product [get]
func GetProductsHandler(c *gin.Context) {
	c.JSON(http.StatusOK, products)
}

// @Summary Add a new product
// @Tags product
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param product body Product true "New product"
// @Success 201 {object} Product
// @Failure 400 {object} map[string]string
// @Router /product [post]
func AddProductHandler(c *gin.Context) {
	var newProduct Product
	if err := c.ShouldBindJSON(&newProduct); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}
	newProduct.ID = getNextProductID()
	products = append(products, newProduct)
	c.JSON(http.StatusCreated, newProduct)
}

// @Summary Update a product by ID
// @Tags product
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "Product ID"
// @Param product body Product true "Updated product"
// @Success 200 {object} Product
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /product/{id} [put]
func UpdateProductHandler(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	var updated Product
	if err := c.ShouldBindJSON(&updated); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}
	for i := range products {
		if products[i].ID == id {
			updated.ID = id
			products[i] = updated
			c.JSON(http.StatusOK, updated)
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
}

// @Summary Delete a product by ID
// @Tags product
// @Security BearerAuth
// @Produce json
// @Param id path int true "Product ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /product/{id} [delete]
func DeleteProductHandler(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	for i, p := range products {
		if p.ID == id {
			products = append(products[:i], products[i+1:]...)
			c.JSON(http.StatusOK, gin.H{"message": "Product deleted"})
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
}

func getNextProductID() int {
	maxID := 0
	for _, p := range products {
		if p.ID > maxID {
			maxID = p.ID
		}
	}
	return maxID + 1
}

// ================== ORDERS =====================

// @Summary Create a new order
// @Tags order
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param order body Order true "Order details"
// @Success 201 {object} Order
// @Failure 400 {object} map[string]string
// @Router /order [post]
func CreateOrderHandler(c *gin.Context) {
	var order Order
	if err := c.ShouldBindJSON(&order); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}
	order.OrderID = getNextOrderID()
	if order.OrderDate.IsZero() {
		order.OrderDate = time.Now()
	}
	orders = append(orders, order)
	c.JSON(http.StatusCreated, order)
}

// @Summary Get all orders
// @Tags order
// @Security BearerAuth
// @Produce json
// @Success 200 {array} Order
// @Router /order [get]
func GetOrdersHandler(c *gin.Context) {
	c.JSON(http.StatusOK, orders)
}

func getNextOrderID() int {
	maxID := 0
	for _, o := range orders {
		if o.OrderID > maxID {
			maxID = o.OrderID
		}
	}
	return maxID + 1
}
