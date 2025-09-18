package main

import (
	"fmt"
	"log"
	"net/http"
	"path/filepath"
)

func main() {
	// Настраиваем обработчик для статических файлов
	http.Handle("/", http.FileServer(http.Dir(".")))

	// Настраиваем CORS для WebAssembly
	http.HandleFunc("/main.wasm", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/wasm")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		http.ServeFile(w, r, "main.wasm")
	})

	// Настраиваем CORS для wasm_exec.js
	http.HandleFunc("/wasm_exec.js", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/javascript")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		http.ServeFile(w, r, "wasm_exec.js")
	})

	port := ":8000"
	fmt.Printf("🚀 HTTP сервер запущен на http://localhost%s\n", port)
	fmt.Printf("📁 Обслуживает файлы из: %s\n", filepath.Abs("."))
	fmt.Printf("🌐 Откройте браузер: http://localhost%s\n", port)
	fmt.Printf("💡 Убедитесь, что файлы main.wasm и wasm_exec.js существуют!\n")

	log.Fatal(http.ListenAndServe(port, nil))
}
