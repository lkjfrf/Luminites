package main

import (
	hand "Luminites/handlers"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		// 요청에서 파일 업로드를 처리합니다.
		file, handler, err := r.FormFile("uploadfile")
		if err != nil {
			http.Error(w, "파일 업로드 중 오류가 발생했습니다.", http.StatusInternalServerError)
			return
		}
		defer file.Close()

		// 업로드된 파일을 서버의 파일 시스템에 저장합니다.
		uploadedFileName := handler.Filename
		savePath := filepath.Join("./uploads", uploadedFileName)
		out, err := os.Create(savePath)
		if err != nil {
			http.Error(w, "파일 저장 중 오류가 발생했습니다.", http.StatusInternalServerError)
			return
		}
		defer out.Close()

		_, err = io.Copy(out, file)
		if err != nil {
			http.Error(w, "파일 복사 중 오류가 발생했습니다.", http.StatusInternalServerError)
			return
		}

		fmt.Fprintf(w, "파일 %s 업로드 완료\n", uploadedFileName)
	} else {
		// 파일 업로드 페이지를 렌더링합니다.
		html := `
		<!DOCTYPE html>
		<html>
		<head>
			<title>파일 업로드</title>
		</head>
		<body>
			<h1>파일 업로드</h1>
			<form method="POST" action="/upload" enctype="multipart/form-data">
				<input type="file" name="uploadfile">
				<input type="submit" value="업로드">
			</form>
		</body>
		</html>
		`
		fmt.Fprintln(w, html)
	}
}

func main() {
	// 업로드된 파일을 저장할 디렉토리를 생성합니다.
	err := os.MkdirAll("./uploads", os.ModePerm)
	if err != nil {
		fmt.Println("디렉토리 생성 중 오류가 발생했습니다:", err)
		return
	}

	http.HandleFunc("/upload", uploadHandler)

	// 파일 다운로드를 위한 핸들러를 등록합니다.
	http.HandleFunc("/download/", func(w http.ResponseWriter, r *http.Request) {
		// 다운로드할 파일의 경로를 가져옵니다.
		filePath := "./uploads/" + filepath.Base(r.URL.Path)
		http.ServeFile(w, r, filePath)
	})

	http.HandleFunc("/posts", hand.PostHandler)
	http.HandleFunc("/adminAuth", hand.AdminHandler)
	http.HandleFunc("/email", hand.EmailHandler)

	// Server Start
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("웹 서버를 시작하는 동안 오류가 발생했습니다:", err)
	}
}
