package main

import (
	"fmt"
	"net/http"
	"strconv"
)

func home(w http.ResponseWriter, r *http.Request) {
	// only match exactly "/", return 404 otherwise
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	fmt.Fprint(w, "Hello from SnippetBox")
}

func snippetView(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}

	fmt.Fprintf(w, "Display specific snippet with ID %d...", id)
}

func snippetCreate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		// w.Header()["Date"] = nil // Del() method doesn't remove system-gen headers
		w.Header().Set("Allow", http.MethodPost)                         // has to be done BEFORE `w.WriteHeader` and `w.Write`
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed) // shortcut for `w.WriteHeader` and `w.Write`
		return
	}

	fmt.Fprint(w, "Create a new snippet...")
}
