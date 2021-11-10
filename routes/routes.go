package routes

import (
	"mime/multipart"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
)

type article struct {
	ID    uint   `json:"id"`
	Title string `json:"title"`
	Body  string `json:"body"`
	Image string `json:"image"`
}

func Serve(r *gin.Engine) {

	articles := []article{
		{ID: 1, Title: "Title#1", Body: "Body#1"},
		{ID: 2, Title: "Title#2", Body: "Body#2"},
		{ID: 3, Title: "Title#3", Body: "Body#3"},
		{ID: 4, Title: "Title#4", Body: "Body#4"},
		{ID: 5, Title: "Title#5", Body: "Body#5"},
	}

	type createArticleForm struct {
		Title string                `form:"title" binding:"required"`
		Body  string                `form:"body" binding:"required"`
		Image *multipart.FileHeader `form:"image" binding:"required"`
	}
	articlesGroup := r.Group("/api/v1/articles")

	articlesGroup.GET("", func(c *gin.Context) {

		result := articles
		if limit := c.Query("limit"); limit != "" {
			n, _ := strconv.Atoi(limit)
			result = result[:n]
		}

		c.JSON(http.StatusOK, gin.H{"articles": result})
	})

	articlesGroup.GET("/:id", func(c *gin.Context) {

		id, _ := strconv.Atoi(c.Param("id"))

		for _, item := range articles {
			if item.ID == uint(id) {
				c.JSON(http.StatusOK, gin.H{"articles": item})
				return
			}
		}
		c.JSON(http.StatusNotFound, gin.H{"error": "Article not found"})

	})

	articlesGroup.POST("", func(c *gin.Context) {

		var form createArticleForm

		if err := c.ShouldBind(&form); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		a := article{
			ID:    uint(len(articles) + 1),
			Title: form.Title,
			Body:  form.Body,
		}

		//Get file
		file, _ := c.FormFile("image")

		//Create Path
		// ID => 8, uploads/articles/8/image.png
		path := "uploads/articles/" + strconv.Itoa(int(a.ID))
		os.MkdirAll(path, 0755)

		//Upload File
		filename := path + "/" + file.Filename
		if err := c.SaveUploadedFile(file, filename); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}

		//Attach File to article
		a.Image = os.Getenv("HOST") + "/" + filename

		articles = append(articles, a)
		c.JSON(http.StatusCreated, gin.H{"article": a})
	})

}
