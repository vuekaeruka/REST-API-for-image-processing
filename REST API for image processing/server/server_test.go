package main

import (
	"bytes"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/gin-gonic/gin"
)

func setupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.POST("/upload", uploadImage)  
	router.POST("/resize", resizeImage)  
	router.POST("/convert", convertImage)
	router.POST("/crop", cropImage)  
	return router
}

func TestUploadImage(t *testing.T) {
	router := setupRouter()

	file, err := os.Open("./testdata/test_image.jpg")
	if err != nil {
		t.Fatal("Ошибка открытия файла:", err)
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("file", filepath.Base(file.Name()))
	if err != nil {
		t.Fatal("Ошибка создания файла в запросе:", err)
	}

	_, err = io.Copy(part, file)
	if err != nil {
		t.Fatal("Ошибка копирования содержимого файла в запрос:", err)
	}

	err = writer.Close()
	if err != nil {
		t.Fatal("Ошибка закрытия multipart writer:", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/upload", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code == http.StatusOK {
		t.Log("Тест 'UploadImage' прошел успешно")
	}
}

func TestResizeImage(t *testing.T) {
	router := setupRouter()


	file, err := os.Open("./testdata/test_image.jpg")
	if err != nil {
		t.Fatal("Ошибка открытия файла:", err)
	}
	defer file.Close()


	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("file", filepath.Base(file.Name()))
	if err != nil {
		t.Fatal("Ошибка создания файла в запросе:", err)
	}

	_, err = io.Copy(part, file)
	if err != nil {
		t.Fatal("Ошибка копирования содержимого файла в запрос:", err)
	}

	err = writer.WriteField("width", "200")
	if err != nil {
		t.Fatal("Ошибка добавления ширины:", err)
	}
	err = writer.WriteField("height", "200") 
	if err != nil {
		t.Fatal("Ошибка добавления высоты:", err)
	}

	err = writer.Close()
	if err != nil {
		t.Fatal("Ошибка закрытия multipart writer:", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/resize", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code == http.StatusOK {
		t.Log("Тест 'ResizeImage' прошел успешно")
	}
}

func TestConvertImage(t *testing.T) {
	router := setupRouter()

	file, err := os.Open("./testdata/test_image.jpg")
	if err != nil {
		t.Fatal("Ошибка открытия файла:", err)
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("file", filepath.Base(file.Name()))
	if err != nil {
		t.Fatal("Ошибка создания файла в запросе:", err)
	}

	_, err = io.Copy(part, file)
	if err != nil {
		t.Fatal("Ошибка копирования содержимого файла в запрос:", err)
	}

	err = writer.WriteField("format", "png") // Устанавливаем формат конвертации
	if err != nil {
		t.Fatal("Ошибка добавления формата:", err)
	}

	err = writer.Close()
	if err != nil {
		t.Fatal("Ошибка закрытия multipart writer:", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/convert", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code == http.StatusOK {
		t.Log("Тест 'ConvertImage' прошел успешно")
	}
}

func TestCropImage(t *testing.T) {
	router := setupRouter()

	file, err := os.Open("./testdata/test_image.jpg")
	if err != nil {
		t.Fatal("Ошибка открытия файла:", err)
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("file", filepath.Base(file.Name()))
	if err != nil {
		t.Fatal("Ошибка создания файла в запросе:", err)
	}

	_, err = io.Copy(part, file)
	if err != nil {
		t.Fatal("Ошибка копирования содержимого файла в запрос:", err)
	}

	err = writer.WriteField("x", "50")
	if err != nil {
		t.Fatal("Ошибка добавления x:", err)
	}
	err = writer.WriteField("y", "50")
	if err != nil {
		t.Fatal("Ошибка добавления y:", err)
	}
	err = writer.WriteField("width", "200")
	if err != nil {
		t.Fatal("Ошибка добавления ширины:", err)
	}
	err = writer.WriteField("height", "200")
	if err != nil {
		t.Fatal("Ошибка добавления высоты:", err)
	}

	err = writer.Close()
	if err != nil {
		t.Fatal("Ошибка закрытия multipart writer:", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/crop", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code == http.StatusOK {
		t.Log("Тест 'CropImage' прошел успешно")
	}
}
