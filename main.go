package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

type Todo struct {
	Id        int    `json:"id"`
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}

var db *sql.DB
var err error

func main() {
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	connStr := "user=" + dbUser + " password=" + dbPassword + " sslmode=disable dbname=" + dbName
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"PUT", "PATCH", "GET", "POST", "DELETE"},
		AllowHeaders:     []string{"Origin"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))
	router.GET("/api/", getAllTodos)
	router.POST("/api/new", addTodo)
	router.PUT("/api/update/:id", updateTodo)
	router.DELETE("/api/delete/:id", deleteTodo)
	router.Run("localhost:8080")
}

// ? Delete Todo
// @param id
// @return status, message
func deleteTodo(c *gin.Context) {
	id := c.Params.ByName("id")
	_, err := db.Exec("DELETE FROM todo WHERE id = $1", id)
	if err != nil {
		log.Fatal(err)
	}
	c.IndentedJSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "Todo deleted successfully!",
	})
}

// ? Add Todo
// @body title, completed
// @return status, message, title, completed
func addTodo(c *gin.Context) {
	var newTodo Todo
	if err := c.BindJSON(&newTodo); err != nil {
		return
	}
	_, err := db.Exec("INSERT INTO todo (title, completed) VALUES ($1, $2)", newTodo.Title, newTodo.Completed)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("New todo added!")
	c.IndentedJSON(http.StatusCreated, gin.H{
		"status":    http.StatusCreated,
		"message":   "Todo created successfully!",
		"title":     newTodo.Title,
		"completed": newTodo.Completed,
	})

}

// ? Update Todo
// @param id
// @body title, completed
// @return status, message, title, completed
func updateTodo(c *gin.Context) {
	var todo Todo
	if err := c.BindJSON(&todo); err != nil {
		return
	}
	id := c.Params.ByName("id")
	_, err := db.Exec("UPDATE todo SET title = $1, completed = $2 WHERE id = $3", todo.Title, todo.Completed, id)
	if err != nil {
		log.Fatal(err)
	}
	c.IndentedJSON(http.StatusOK, gin.H{
		"status":    http.StatusOK,
		"message":   "Todo updated successfully!",
		"title":     todo.Title,
		"completed": todo.Completed,
	})
}

// ? Get Todo
// @return todos[]
func getAllTodos(c *gin.Context) {
	todos, err := db.Query("SELECT * FROM todo")
	if err != nil {
		log.Fatal(err)
	}
	defer todos.Close()
	var todoList []Todo
	for todos.Next() {
		var t Todo
		if err := todos.Scan(&t.Id, &t.Title, &t.Completed); err != nil {
			log.Fatal(err)
		}
		todoList = append(todoList, t)
	}
	c.JSON(http.StatusOK, todoList)
}
