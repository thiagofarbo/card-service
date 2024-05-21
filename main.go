package main

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"jwt-authentication/db"
	"jwt-authentication/helpers"
	"jwt-authentication/models"
	"net/http"
)

var dbModel = make(map[string]string)

func setupRouter() *gin.Engine {
	r := gin.Default()

	r.POST("/login", func(c *gin.Context) {
		var userRequest models.UserRequest
		if err := c.BindJSON(&userRequest); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		userModel := models.User{
			Username: userRequest.Username,
			Password: userRequest.Password,
		}

		user, err := db.UserMatchPassword(userModel.Username, userModel.Password)

		if user == nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		session := sessions.Default(c)

		if session.Get("hello") != "world" {
			session.Set("hello", "world")
			session.Save()
		}

		c.JSON(http.StatusOK, gin.H{"user": user})
	})

	r.POST("/logout", func(c *gin.Context) {
		//Delete the session
		session := sessions.Default(c)
		session.Clear()
		session.Save()

		var userRequest models.UserRequest
		if err := c.BindJSON(&userRequest); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		userModel := models.User{
			Username: userRequest.Username,
			Password: userRequest.Password,
		}

		user, err := db.UserMatchPassword(userModel.Username, userModel.Password)

		if user == nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		session = sessions.Default(c)

		if session.Get("hello") != "world" {
			session.Set("hello", "world")
			session.Save()
		}

		c.JSON(http.StatusOK, gin.H{"user": user})
	})

	r.POST("/users", func(c *gin.Context) {
		var userRequest models.UserRequest
		if err := c.BindJSON(&userRequest); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		userModel := models.User{
			Username: userRequest.Username,
			Password: userRequest.Password,
		}

		user, err := db.CreateUser(&userModel)

		if user == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		helpers.SetSession(c, user)
		c.JSON(http.StatusCreated, gin.H{"user": user})

		c.Redirect(http.StatusTemporaryRedirect, "/cards")
	})

	r.POST("/notes", func(c *gin.Context) {
		var noteRequest models.NoteRequest
		if err := c.BindJSON(&noteRequest); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		noteModel := models.Note{
			Name:    noteRequest.Name,
			Content: noteRequest.Content,
			UserID:  noteRequest.UserID,
		}

		note, err := db.CreateNote(&noteModel)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, gin.H{"note": note})
	})

	r.POST("/cards", func(c *gin.Context) {
		var cardRequest models.CardRequest
		if err := c.BindJSON(&cardRequest); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		cardModel := models.Card{
			Number: cardRequest.Number,
			UserID: cardRequest.UserID,
		}

		card, err := db.CreateCard(&cardModel)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, gin.H{"card": card})
	})

	// Ping test
	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	//// Get user value
	//r.GET("/user/:name", func(c *gin.Context) {
	//	user := c.Params.ByName("name")
	//	value, ok := dbModel[user]
	//	if ok {
	//		c.JSON(http.StatusOK, gin.H{"user": user, "value": value})
	//	} else {
	//		c.JSON(http.StatusOK, gin.H{"user": user, "status": "no value"})
	//	}
	//})

	// Authorized group (uses gin.BasicAuth() middleware)
	// Same than:
	// authorized := r.Group("/")
	// authorized.Use(gin.BasicAuth(gin.Credentials{
	//	  "foo":  "bar",
	//	  "manu": "123",
	//}))
	authorized := r.Group("/", gin.BasicAuth(gin.Accounts{
		"foo":  "bar", // user:foo password:bar
		"manu": "123", // user:manu password:123
	}))

	/* example curl for /admin with basicauth header
	   Zm9vOmJhcg== is base64("foo:bar")

		curl -X POST \
	  	http://localhost:8080/admin \
	  	-H 'authorization: Basic Zm9vOmJhcg==' \
	  	-H 'content-type: application/json' \
	  	-d '{"value":"bar"}'
	*/
	authorized.POST("admin", func(c *gin.Context) {
		user := c.MustGet(gin.AuthUserKey).(string)

		// Parse JSON
		var json struct {
			Value string `json:"value" binding:"required"`
		}

		if c.Bind(&json) == nil {
			dbModel[user] = json.Value
			c.JSON(http.StatusOK, gin.H{"status": "ok"})
		}
	})

	return r
}
func main() {
	db.ConnectDatabase()
	db.DBMigrate()
	//
	r := gin.Default()
	r.Use(gin.Logger())
	//
	//r.Static("vendor", "./static/vendor")
	//r.LoadHTMLGlob("templates/**/*")

	r = setupRouter()
	// Listen and Server in 0.0.0.0:8080
	r.Run(":8080")
}
