package handlers

import (
	"fmt"
	"forum/utils"
	"io"
	"net/http"
	"os"
	"strings"
)

const MaxFileSize = 20 * 1024 * 1024

func uploadHandler(w http.ResponseWriter, r *http.Request, media_post *string) string {

	r.ParseMultipartForm(MaxFileSize + 1)
	file, header, err := r.FormFile("media_post")
	if err != nil {
		// err := map[string]interface{}{"Error": "Error on copy file", "User": user}
		// utils.RenderTemplate(w, "createPost", err)
		return "no file"
	}
	defer file.Close()

	contentType := header.Header.Get("Content-Type")
	if !isValidImageType(contentType) {

		fmt.Println("Invalid image format")
		err := map[string]interface{}{"Error": "Invalid image format", "User": user}
		utils.RenderTemplate(w, "createPost", err)

		return "err"
	}
	if header.Size > MaxFileSize {
		fmt.Println("Image is more than 20mo")
		err := map[string]interface{}{"Error": "Image is more than 20mo", "User": user}
		utils.RenderTemplate(w, "createPost", err)

		return "err"
	}

	*media_post = header.Filename
	newFile, err := os.Create("./assets/" + header.Filename)
	if err != nil {
		// http.Error(w, err.Error(), http.StatusInternalServerError)
		return "err"
	}
	defer newFile.Close()

	_, err = io.Copy(newFile, file)
	if err != nil {
		// http.Error(w, err.Error(), http.StatusInternalServerError)
		return "err"
	}

	fmt.Println("Image upload")
	return "ok"

}

func isValidImageType(contentType string) bool {
	contentType = strings.ToLower(contentType)
	return contentType == "image/jpeg" ||
		contentType == "image/jpg" ||
		contentType == "image/png" ||
		contentType == "image/gif" ||
		contentType == "image/bmp" ||
		contentType == "image/webp"
}
