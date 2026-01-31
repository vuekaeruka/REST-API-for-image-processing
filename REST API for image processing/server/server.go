package main

import (
	"bytes"
	"fmt"
	"html/template"
	"image"
	"image/jpeg"
	"image/png"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/nfnt/resize"
	"golang.org/x/image/draw"
)

func main() {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	os.MkdirAll("./uploads", os.ModePerm)
	os.MkdirAll("./resized_images", os.ModePerm)
	os.MkdirAll("./converted_images", os.ModePerm)
	os.MkdirAll("./cropped_images", os.ModePerm)

	router.Static("/uploads", "./uploads")
	router.Static("/resized_images", "./resized_images")
	router.Static("/converted_images", "./converted_images")
	router.Static("/cropped_images", "./cropped_images")

	router.GET("/", func(c *gin.Context) {
		tmpl, err := template.New("index").Parse(`
		<!DOCTYPE html>
		<html>
		<head>
			<title>Image Processing</title>
		</head>
		<body>
			<h1>Image Processing API</h1>
			<form action="/upload" method="post" enctype="multipart/form-data">
				<h2>Upload Image</h2>
				<input type="file" name="file" required>
				<button type="submit">Upload</button>
			</form>
			{{if .UploadError}}
			<p style="color:red;">{{.UploadError}}</p>
			{{end}}

			<br>

			<form action="/resize" method="post" enctype="multipart/form-data">
				<h2>Resize Image</h2>
				<input type="file" name="file" required><br>
				Width: <input type="text" name="width" required><br>
				Height: <input type="text" name="height" required><br>
				<button type="submit">Resize</button>
			</form>
			{{if .ResizeError}}
			<p style="color:red;">{{.ResizeError}}</p>
			{{end}}

			<br>

			<form action="/convert" method="post" enctype="multipart/form-data">
				<h2>Convert Image Format</h2>
				<input type="file" name="file" required><br>
				Format: <input type="text" name="format" required><br>
				<button type="submit">Convert</button>
			</form>
			{{if .ConvertError}}
			<p style="color:red;">{{.ConvertError}}</p>
			{{end}}

			<br>

			<form action="/crop" method="post" enctype="multipart/form-data">
				<h2>Crop Image</h2>
				<input type="file" name="file" required><br>
				X: <input type="text" name="x" required><br>
				Y: <input type="text" name="y" required><br>
				Width: <input type="text" name="width" required><br>
				Height: <input type="text" name="height" required><br>
				<button type="submit">Crop</button>
			</form>
			{{if .CropError}}
			<p style="color:red;">{{.CropError}}</p>
			{{end}}
		</body>
		</html>
		`)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load template"})
			return
		}
		tmpl.Execute(c.Writer, gin.H{
			"UploadError":  c.DefaultQuery("uploadError", ""),
			"ResizeError":  c.DefaultQuery("resizeError", ""),
			"ConvertError": c.DefaultQuery("convertError", ""),
			"CropError":    c.DefaultQuery("cropError", ""),
		})
	})

	router.POST("/upload", uploadImage)
	router.POST("/resize", resizeImage)
	router.POST("/convert", convertImage)
	router.POST("/crop", cropImage)

	fmt.Println("Сервер запущен на порту 8080")
	if err := router.Run(":8080"); err != nil {
		fmt.Println("Ошибка запуска сервера:", err)
	}
}

func uploadImage(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.Redirect(http.StatusFound, "/?uploadError=Не удалось получить файл")
		return
	}

	outputPath := filepath.Join("./uploads", file.Filename)
	if err := c.SaveUploadedFile(file, outputPath); err != nil {
		c.Redirect(http.StatusFound, "/?uploadError=Не удалось сохранить файл")
		return
	}

	c.Redirect(http.StatusFound, "/")
}

func resizeImage(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.Redirect(http.StatusFound, "/?resizeError=Не удалось получить файл")
		return
	}

	width, err := strconv.Atoi(c.PostForm("width"))
	if err != nil {
		c.Redirect(http.StatusFound, "/?resizeError=Некорректная ширина")
		return
	}

	height, err := strconv.Atoi(c.PostForm("height"))
	if err != nil {
		c.Redirect(http.StatusFound, "/?resizeError=Некорректная высота")
		return
	}

	srcFile, err := file.Open()
	if err != nil {
		c.Redirect(http.StatusFound, "/?resizeError=Не удалось открыть файл")
		return
	}
	defer srcFile.Close()

	img, _, err := image.Decode(srcFile)
	if err != nil {
		c.Redirect(http.StatusFound, "/?resizeError=Не удалось декодировать изображение")
		return
	}

	resizedImage := resize.Resize(uint(width), uint(height), img, resize.Lanczos3)

	var buf bytes.Buffer
	if err := jpeg.Encode(&buf, resizedImage, nil); err != nil {
		c.Redirect(http.StatusFound, "/?resizeError=Не удалось закодировать изображение")
		return
	}

	outputFilename := "resized_" + file.Filename
	outputPath := filepath.Join("./resized_images", outputFilename)

	if err := os.WriteFile(outputPath, buf.Bytes(), 0644); err != nil {
		c.Redirect(http.StatusFound, "/?resizeError=Не удалось сохранить изображение")
		return
	}

	c.Redirect(http.StatusFound, "/")
}

func convertImage(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.Redirect(http.StatusFound, "/?convertError=Не удалось получить файл")
		return
	}

	format := c.PostForm("format")
	if format != "png" && format != "jpeg" {
		c.Redirect(http.StatusFound, "/?convertError=Поддерживаются только форматы png и jpeg")
		return
	}

	srcFile, err := file.Open()
	if err != nil {
		c.Redirect(http.StatusFound, "/?convertError=Не удалось открыть файл")
		return
	}
	defer srcFile.Close()

	img, _, err := image.Decode(srcFile)
	if err != nil {
		c.Redirect(http.StatusFound, "/?convertError=Не удалось декодировать изображение")
		return
	}

	var buf bytes.Buffer
	outputFilename := file.Filename + "." + format
	outputPath := filepath.Join("./converted_images", outputFilename)

	switch format {
	case "png":
		if err := png.Encode(&buf, img); err != nil {
			c.Redirect(http.StatusFound, "/?convertError=Не удалось закодировать PNG")
			return
		}
	case "jpeg":
		if err := jpeg.Encode(&buf, img, nil); err != nil {
			c.Redirect(http.StatusFound, "/?convertError=Не удалось закодировать JPEG")
			return
		}
	}

	if err := os.WriteFile(outputPath, buf.Bytes(), 0644); err != nil {
		c.Redirect(http.StatusFound, "/?convertError=Не удалось сохранить изображение")
		return
	}

	c.Redirect(http.StatusFound, "/")
} // ← ЗАКРЫВАЮЩАЯ СКОБКА ДОБАВЛЕНА ЗДЕСЬ

func cropImage(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.Redirect(http.StatusFound, "/?cropError=Не удалось получить файл")
		return
	}

	x, err := strconv.Atoi(c.PostForm("x"))
	if err != nil {
		c.Redirect(http.StatusFound, "/?cropError=Некорректное значение X")
		return
	}

	y, err := strconv.Atoi(c.PostForm("y"))
	if err != nil {
		c.Redirect(http.StatusFound, "/?cropError=Некорректное значение Y")
		return
	}

	width, err := strconv.Atoi(c.PostForm("width"))
	if err != nil {
		c.Redirect(http.StatusFound, "/?cropError=Некорректная ширина")
		return
	}

	height, err := strconv.Atoi(c.PostForm("height"))
	if err != nil {
		c.Redirect(http.StatusFound, "/?cropError=Некорректная высота")
		return
	}

	srcFile, err := file.Open()
	if err != nil {
		c.Redirect(http.StatusFound, "/?cropError=Не удалось открыть файл")
		return
	}
	defer srcFile.Close()

	img, _, err := image.Decode(srcFile)
	if err != nil {
		c.Redirect(http.StatusFound, "/?cropError=Не удалось декодировать изображение")
		return
	}

	croppedImage := image.NewRGBA(image.Rect(0, 0, width, height))
	draw.Draw(croppedImage, croppedImage.Bounds(), img, image.Point{X: x, Y: y}, draw.Src)

	var buf bytes.Buffer
	if err := jpeg.Encode(&buf, croppedImage, nil); err != nil {
		c.Redirect(http.StatusFound, "/?cropError=Не удалось закодировать изображение")
		return
	}

	outputFilename := "cropped_" + file.Filename
	outputPath := filepath.Join("./cropped_images", outputFilename)

	if err := os.WriteFile(outputPath, buf.Bytes(), 0644); err != nil {
		c.Redirect(http.StatusFound, "/?cropError=Не удалось сохранить обрезанное изображение")
		return
	}

	c.Redirect(http.StatusFound, "/")
}