package main

import (
    "context"
    "io"
    "log"
    "net/http"
    "os"
    "strconv"
    "time"

    "cloud.google.com/go/storage"
    "google.golang.org/api/option"
)

func main() {
    log.Println("Starting")

    wwwPort, exist := os.LookupEnv("PORT")
    if !exist {
        wwwPort = "8000"
    }

    if len(os.Getenv("BUCKET")) == 0 {
        panic("No BUCKET environmental variable is set")
    }
    ctx := context.Background()
    client, err := storage.NewClient(
        ctx,
        option.WithCredentialsFile(os.Getenv("SERVICE_ACCOUNT_KEY")),
    )
    if err != nil {
        panic("Unable to create the client")
    }
    bucket := client.Bucket(os.Getenv("BUCKET"))

    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()
	p := r.URL.Path[1:]
        oh := bucket.Object(p)
        objAttrs, err := oh.Attrs(ctx)
        if err != nil {
            if os.Getenv("LOGGING") == "true" {
                elapsed := time.Since(start)
                log.Println("| 404 |", elapsed.String(), r.Host, r.Method, r.URL.Path)
            }
            http.Error(w, "Not found", 404)
            return
        }
        o := oh.ReadCompressed(true)
        rc, err := o.NewReader(ctx)
        if err != nil {
            http.Error(w, "Not found", 404)
            return
        }
        defer rc.Close()

        w.Header().Set("Content-Type", objAttrs.ContentType)
        w.Header().Set("Content-Encoding", objAttrs.ContentEncoding)
        w.Header().Set("Content-Length", strconv.Itoa(int(objAttrs.Size)))
        w.WriteHeader(200)
        if _, err := io.Copy(w, rc); err != nil {
            if os.Getenv("LOGGING") == "true" {
                elapsed := time.Since(start)
                log.Println("| 200 |", elapsed.String(), r.Host, r.Method, r.URL.Path)
            }
            return
        }
        if os.Getenv("LOGGING") == "true" {
            elapsed := time.Since(start)
            log.Println("| 200 |", elapsed.String(), r.Host, r.Method, r.URL.Path)
        }
    })

    http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
        w.Write([]byte("ok"))
    })
    http.HandleFunc("/readiness", func(w http.ResponseWriter, r *http.Request) {
        w.Write([]byte("ok"))
    })

    log.Println("Serving on port", wwwPort)
    http.ListenAndServe(":"+wwwPort, nil)
}
