package main

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"repository/assets"
	"repository/database"
	transporthttp "repository/transport/http"
	"sync"
	"syscall"
	"time"
)

func main() {
	if err := database.Open(); err != nil {
		panic(err)
	}

	assetsService := assets.Service{}
	assetsHandler := assets.Handler{Service: assetsService}

	mux := http.NewServeMux()

	mux.Handle("GET /api/v1/assets/search", transporthttp.AppHandler(assetsHandler.Search))
	mux.Handle("GET /api/v1/assets/{path...}", transporthttp.AppHandler(assetsHandler.Get))
	mux.Handle("POST /api/v1/assets", transporthttp.AppHandler(assetsHandler.Upload))
	mux.Handle("POST /api/v1/assets/{id}/embeddings", transporthttp.AppHandler(assetsHandler.Embed))

	server := &http.Server{
		Addr:    ":8080",
		Handler: loggerMiddleware(mux),
	}

	go func() {
		slog.Info("Starting server on port 8080.")
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("Server forced to shutdown due to error.", "error", err)
		}
	}()
	go func() {
		// time.Sleep(5 * time.Second)
		// ingestInitialData()
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	sig := <-quit
	fmt.Println()
	slog.Info("Shutdown signal caught.", "signal", sig.String())

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		slog.Error("Server forced to shutdown prematurely.", "error", err)
	} else {
		slog.Info("Server stopped completely and gracefully.")
	}

	if err := database.Close(); err != nil {
		slog.Error("Error during database close.", "error", err)
	} else {
		slog.Info("Database closed successfully.")
	}

	slog.Info("Application fully terminated.")
}

type respWriter struct {
	http.ResponseWriter
	statusCode int
	written    bool
}

func (rw *respWriter) WriteHeader(code int) {
	if rw.written {
		return // ignore — header already sent, avoid the superfluous-call log
	}
	rw.written = true
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func loggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		wrapped := &respWriter{ResponseWriter: w, statusCode: http.StatusOK}
		next.ServeHTTP(wrapped, r)
		slog.Info("HTTP Request", "method", r.Method, "path", r.URL.Path, "status", wrapped.statusCode)
	})
}

func ingestInitialData() {
	slog.Info("Ingest Initial Data: Starting...")

	wd, err := os.Getwd()
	if err != nil {
		slog.Error("Ingest Initial Data: unable to start", "error", err)
		return
	}

	base := filepath.Join(filepath.Dir(wd), "resources", "image")

	pahts := make([]string, 0)
	filepath.WalkDir(base, func(path string, d fs.DirEntry, err error) error {
		if err == nil && !d.IsDir() {
			pahts = append(pahts, path)
		}
		return nil
	})

	service := assets.Service{}
	var wg sync.WaitGroup

	for i, p := range pahts {
		file, err := os.Open(p)
		if err != nil {
			slog.Error("Ingest Initial Data: unable to read file", "path", p, "error", err)
			continue
		}

		wg.Go(func() {
			if err := service.Ingest(context.Background(), file, "image", database.JinaClipV2); err != nil {
				slog.Error("Ingest Initial Data: failed to ingest asset", "path", p, "error", err)
			}
		})

		if i%10 == 0 {
			wg.Wait()
		}
	}

	wg.Wait()
	slog.Info("Ingest Initial Data: Finished")

}
