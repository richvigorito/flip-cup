package api

import (
    "net/http"
    "path/filepath"
    "log"
)

// SPAHandler returns an http.Handler that serves static files from the "public" directory,
// and falls back to serving index.html for SPA client-side routing.
func SPAHandler(staticDir string) http.Handler {
    fs := http.Dir(staticDir)
    fileServer := http.FileServer(fs)

    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        log.Printf("Request for: %s", r.URL.Path)

        if r.URL.Path == "/ws" || r.URL.Path == "/ws-test" {
            http.NotFound(w, r)
            return
        }

        f, err := fs.Open(r.URL.Path)
        if err != nil {
            log.Printf("File not found, serving index.html for %s", r.URL.Path)
            http.ServeFile(w, r, filepath.Join(staticDir, "index.html"))
            return
        }
        defer f.Close()

        fi, err := f.Stat()
        if err != nil {
            log.Printf("Error stating file, serving index.html for %s", r.URL.Path)
            http.ServeFile(w, r, filepath.Join(staticDir, "index.html"))
            return
        }

        if fi.IsDir() {
            // If directory, serve index.html for SPA routing
            log.Printf("Path is a directory, serving index.html for %s", r.URL.Path)
            http.ServeFile(w, r, filepath.Join(staticDir, "index.html"))
            return
        }

        // File exists and is not a directory — serve it
        log.Printf("Serving static file: %s", r.URL.Path)
        fileServer.ServeHTTP(w, r)
    })
}
