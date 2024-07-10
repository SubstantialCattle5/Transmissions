package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"os/exec"
	"strconv"
)

// album represents data about a record album.
type album struct {
	ID     string  `json:"id"`
	Title  string  `json:"title"`
	Artist string  `json:"artist"`
	Price  float64 `json:"price"`
	Song   string  `json:"song"`
}

var albums = []album{
	{ID: "1", Title: "Blue Train", Artist: "John Coltrane", Price: 56.99, Song: ""},
	{ID: "2", Title: "Jeru", Artist: "Gerry Mulligan", Price: 17.99, Song: ""},
	{ID: "3", Title: "Sarah Vaughan and Clifford Brown", Artist: "Sarah Vaughan", Price: 39.99, Song: ""},
}

func getAlbums(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, albums)
}

func postAlbums(c *gin.Context) {
	var newAlbum album
	if err := c.BindJSON(&newAlbum); err != nil {
		return
	}

	albums = append(albums, newAlbum)
	c.IndentedJSON(http.StatusCreated, newAlbum)
}

func hideDataHandler(c *gin.Context) {
	// Parse form data
	soundFile, err := c.FormFile("sound_file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error reading sound file"})
		return
	}

	secretFile, err := c.FormFile("secret_file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error reading secret file"})
		return
	}

	numLsb := c.PostForm("num_lsb")
	if numLsb == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Number of LSBs not provided"})
		return
	}

	numLsbInt, err := strconv.Atoi(numLsb)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid number of LSBs"})
		return
	}

	soundPath := "sound.wav"
	secretPath := "secret.txt"
	outputPath := "output.wav"

	// Save the uploaded files
	if err := c.SaveUploadedFile(soundFile, soundPath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error saving sound file"})
		return
	}

	if err := c.SaveUploadedFile(secretFile, secretPath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error saving secret file"})
		return
	}

	// Run the Python script
	cmd := exec.Command("python", "encode.py", soundPath, secretPath, outputPath, strconv.Itoa(numLsbInt))
	err = cmd.Run()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Error running Python script: %v", err)})
		return
	}

	// Send the resulting file back to the user
	c.FileAttachment(outputPath, "output.wav")
}

func main() {
	router := gin.Default()
	router.GET("/albums", getAlbums)
	router.POST("/albums", postAlbums)
	router.POST("/encode", hideDataHandler)

	if err := router.Run("localhost:8080"); err != nil {
		fmt.Println("Failed to run server:", err)
	}
}
