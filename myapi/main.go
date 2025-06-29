package main

import (
	"database/sql"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	_ "github.com/lib/pq"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "myapi/docs" // Replace with your module name
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
var db *sql.DB

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

func main() {
	var err error
	db, err = sql.Open("postgres", "host=localhost port=5432 user=aravi password=blueberry dbname=yourdb sslmode=disable")
	if err != nil {
		panic("Failed to connect to DB: " + err.Error())
	}
	if err = db.Ping(); err != nil {
		panic("DB unreachable: " + err.Error())
	}

	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	r.POST("/login", LoginHandler)
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

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

// ==================== AUTH ========================

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

// ==================== PRODUCT ========================

func GetProductsHandler(c *gin.Context) {
	rows, err := db.Query("SELECT id, name, description, price, in_stock FROM products")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch products"})
		return
	}
	defer rows.Close()

	var products []Product
	for rows.Next() {
		var p Product
		if err := rows.Scan(&p.ID, &p.Name, &p.Description, &p.Price, &p.InStock); err != nil {
			continue
		}
		products = append(products, p)
	}
	c.JSON(http.StatusOK, products)
}

func AddProductHandler(c *gin.Context) {
	var p Product
	if err := c.ShouldBindJSON(&p); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}
	err := db.QueryRow(
		"INSERT INTO products (name, description, price, in_stock) VALUES ($1, $2, $3, $4) RETURNING id",
		p.Name, p.Description, p.Price, p.InStock,
	).Scan(&p.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Insert failed"})
		return
	}
	c.JSON(http.StatusCreated, p)
}

func UpdateProductHandler(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	var p Product
	if err := c.ShouldBindJSON(&p); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}
	p.ID = id
	_, err = db.Exec(
		"UPDATE products SET name=$1, description=$2, price=$3, in_stock=$4 WHERE id=$5",
		p.Name, p.Description, p.Price, p.InStock, p.ID,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Update failed"})
		return
	}
	c.JSON(http.StatusOK, p)
}

func DeleteProductHandler(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	res, err := db.Exec("DELETE FROM products WHERE id=$1", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Delete failed"})
		return
	}
	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Product deleted"})
}

// ==================== ORDER ========================

func CreateOrderHandler(c *gin.Context) {
	var o Order
	if err := c.ShouldBindJSON(&o); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}
	if o.OrderDate.IsZero() {
		o.OrderDate = time.Now()
	}
	err := db.QueryRow(
		"INSERT INTO orders (product_id, quantity, order_date) VALUES ($1, $2, $3) RETURNING order_id",
		o.ProductID, o.Quantity, o.OrderDate,
	).Scan(&o.OrderID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Order creation failed"})
		return
	}
	c.JSON(http.StatusCreated, o)
}

func GetOrdersHandler(c *gin.Context) {
	rows, err := db.Query("SELECT order_id, product_id, quantity, order_date FROM orders")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch orders"})
		return
	}
	defer rows.Close()

	var orders []Order
	for rows.Next() {
		var o Order
		if err := rows.Scan(&o.OrderID, &o.ProductID, &o.Quantity, &o.OrderDate); err != nil {
			continue
		}
		orders = append(orders, o)
	}
	c.JSON(http.StatusOK, orders)
}
