package user

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var (
	host     string
	port     int
	user     string
	password string
	dbname   string
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	host = os.Getenv("DB_HOST")
	port, _ = strconv.Atoi(os.Getenv("DB_PORT"))
	user = os.Getenv("DB_USER")
	password = os.Getenv("DB_PASSWORD")
	dbname = os.Getenv("DB_NAME")
}

func connectDB() *sql.DB {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	fmt.Println("Successfully connected to database!")

	return db
}

type User struct {
	ID        string `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
}

func GetUsers(c *gin.Context) {
	db := connectDB()
	defer db.Close()

	stmt, err := db.Prepare("SELECT id, first_name, last_name, email FROM users")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to prepare statement"})
		return
	}
	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve users from database"})
		return
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		err := rows.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan user from database"})
			return
		}
		users = append(users, user)
	}

	c.JSON(http.StatusOK, gin.H{"users": users})
}

func GetUser(c *gin.Context) {
	db := connectDB()
	defer db.Close()

	id := c.Param("id")

	stmt, err := db.Prepare(`SELECT id, first_name, last_name, email FROM users WHERE id = $1`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to prepare statement"})
		return
	}
	defer stmt.Close()

	row := stmt.QueryRow(id)
	var user User

	err = row.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan user from database"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": user,
	})
}

func CreateUser(c *gin.Context) {
	db := connectDB()
	defer db.Close()

	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	stmt, err := db.Prepare("INSERT INTO users (id, first_name, last_name, email) VALUES ($1, $2, $3, $4)")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to prepare statement"})
		return
	}
	defer stmt.Close()

	user.ID = uuid.New().String()

	_, err = stmt.Exec(user.ID, user.FirstName, user.LastName, user.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to execute statement"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "User created successfully",
		"result":  user,
	})
}

func UpdateUser(c *gin.Context) {
	db := connectDB()
	defer db.Close()

	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	stmt, err := db.Prepare("UPDATE users SET first_name = $1, last_name = $2, email = $3 WHERE id = $4")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to prepare statement"})
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(user.FirstName, user.LastName, user.Email, user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to execute statement"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User updated successfully",
		"result":  user,
	})
}

func DeleteUser(c *gin.Context) {
	db := connectDB()
	defer db.Close()

	id := c.Param("id")

	stmt, err := db.Prepare("DELETE FROM users WHERE id = $1")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to prepare statement"})
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to execute statement"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("User deleted successfully"),
	})
}
