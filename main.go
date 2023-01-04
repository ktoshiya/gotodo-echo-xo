package main

import (
	"context"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/ktoshiya/todo/models"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"golang.org/x/crypto/bcrypt"

	_ "github.com/go-sql-driver/mysql"
)

type UserID int64
type TodoID int64
type TodoStatus string

const (
	TodoStatusTodo  TodoStatus = "todo"
	TodoStatusDoing TodoStatus = "doing"
	TodoStatusDone  TodoStatus = "done"
)

type (
	user struct {
		ID       UserID    `json:"id" db:"id"`
		Name     string    `json:"name" db:"name"`
		Password string    `json:"password" db:"password"`
		Role     string    `json:"role" db:"role"`
		Created  time.Time `json:"created" db:"created"`
		Modified time.Time `json:"modified" db:"modified"`
	}

	todo struct {
		ID       TodoID     `json:"id" db:"id"`
		UserID   UserID     `json:"user_id" db:"user_id"`
		Title    string     `json:"title" db:"title"`
		Created  time.Time  `json:"created" db:"created"`
	}
)

var (
	users = map[int]*user{}
	seq   = user{ID: 1}
	lock  = sync.Mutex{}
)

func (u *user) ComparePassword(pw string) error {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(pw))
}

//----------
// Handlers
//----------

func createUser(ctx context.Context, db *sqlx.DB) echo.HandlerFunc {
	return func(c echo.Context) error {

		u := new(models.User)
		if err := c.Bind(u); err != nil {
			return err
		}

		pw, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		u.Password = string(pw)
		u.Created = time.Now().UTC()

		err = u.Insert(ctx, db)
		if err != nil {
			return err
		}

		return c.JSON(http.StatusCreated, u)
	}
}

func getTodos(ctx context.Context, db *sqlx.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		todos, err := models.TodosByUserID(ctx, db, 1)

		if err != nil {
			return err
		}

		return c.JSON(http.StatusCreated, todos)
	}
}

func getUser(c echo.Context) error {
	lock.Lock()
	defer lock.Unlock()
	id, _ := strconv.Atoi(c.Param("id"))
	return c.JSON(http.StatusOK, users[id])
}

func updateUser(c echo.Context) error {
	lock.Lock()
	defer lock.Unlock()
	u := new(user)
	if err := c.Bind(u); err != nil {
		return err
	}
	id, _ := strconv.Atoi(c.Param("id"))
	users[id].Name = u.Name
	return c.JSON(http.StatusOK, users[id])
}

func deleteUser(c echo.Context) error {
	lock.Lock()
	defer lock.Unlock()
	id, _ := strconv.Atoi(c.Param("id"))
	delete(users, id)
	return c.NoContent(http.StatusNoContent)
}

func getAllUsers(c echo.Context) error {
	lock.Lock()
	defer lock.Unlock()
	return c.JSON(http.StatusOK, users)
}

func main() {
	db, err := sqlx.Connect("mysql", "todo:todo@(127.0.0.1:3306)/todo?parseTime=true")
	if err != nil {
		log.Fatalln(err)
	}

	// SQLのタイムアウト設定
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Routes
	e.GET("/users", getAllUsers)
	e.GET("/todos", getTodos(ctx, db))
	e.POST("/users", createUser(ctx, db))
	e.GET("/users/:id", getUser)
	e.PUT("/users/:id", updateUser)
	e.DELETE("/users/:id", deleteUser)

	// Start server
	e.Logger.Fatal(e.Start(":1323"))
}
