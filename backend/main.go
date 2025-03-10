package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	_ "github.com/go-sql-driver/mysql" // Ensure this import is present

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

// JWT Secret Key
var jwtSecret = []byte("your_secret_key")

// Database connection
var db *sql.DB

// Student Struct
type Student struct {
	ID       int    `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// JWT Claims
type Claims struct {
	Email string `json:"email"`
	jwt.RegisteredClaims
}

// Connect to MySQL
func connectDB() {
	dsn := "root:root@tcp(127.0.0.1:3306)/studentdb"
	var err error
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("Database connection failed:", err)
	}
	if err := db.Ping(); err != nil {
		log.Fatal("Database not reachable:", err)
	}
	fmt.Println("Database connected!")
}

// Signup API
func signup(c *gin.Context) {
	var student Student
	if err := c.ShouldBindJSON(&student); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// Hash Password
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(student.Password), bcrypt.DefaultCost)

	// Insert into Database
	_, err := db.Exec("INSERT INTO students (email, password) VALUES (?, ?)", student.Email, string(hashedPassword))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Signup failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Signup successful"})
}

// Login API
func login(c *gin.Context) {
	var student Student
	if err := c.ShouldBindJSON(&student); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// Fetch Student from DB
	var dbPassword string
	err := db.QueryRow("SELECT password FROM students WHERE email = ?", student.Email).Scan(&dbPassword)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Compare Hashed Password
	if bcrypt.CompareHashAndPassword([]byte(dbPassword), []byte(student.Password)) != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Generate JWT Token
	expirationTime := time.Now().Add(1 * time.Hour)
	claims := &Claims{
		Email: student.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString(jwtSecret)

	c.JSON(http.StatusOK, gin.H{"token": tokenString})
}

// Middleware to Protect Routes
func authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing token"})
			c.Abort()
			return
		}

		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		c.Set("email", claims.Email)
		c.Next()
	}
}

// Protected Dashboard API
func dashboard(c *gin.Context) {
	email, _ := c.Get("email")
	c.JSON(http.StatusOK, gin.H{"message": "Welcome to the Dashboard", "email": email})
}

// func main() {
// 	connectDB()
// 	router := gin.Default()

// 	// Routes
// 	router.POST("/signup", signup)
// 	router.POST("/login", login)
// 	router.GET("/dashboard", authMiddleware(), dashboard)

//		router.Run(":8080")
//	}
func main() {
	connectDB()
	router := gin.Default()

	// Enable CORS
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"}, // Allow frontend URL
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
	}))

	// Routes
	router.POST("/signup", signup)
	router.POST("/login", login)
	router.GET("/dashboard", authMiddleware(), dashboard)

	router.Run(":8080")
}
