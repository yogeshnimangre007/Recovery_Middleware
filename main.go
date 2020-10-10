package main

import (
	"fmt"
	"log"
	"net/http"
	"runtime/debug"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/panic/", panicDemo)
	mux.HandleFunc("/panic-after/", panicAfterDemo)
	mux.HandleFunc("/", hello)
	log.Fatal(http.ListenAndServe(":8000", recoverMw(mux, true)))
}

//recovery middleware
func recoverMw(app http.Handler, dev bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Println(err)
				stack := debug.Stack()
				fmt.Println(string(stack))
				if !dev {
					http.Error(w, "somethings went wrong :( ", http.StatusInternalServerError)
					return
				}
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprintf(w, "<h1>Panic : %v </h1><pre> %s</pre>", err, string(stack))

			}
		}()
		nw := &responseWriter{ResponseWriter: w}
		app.ServeHTTP(nw, r)
		nw.flush()
	}
}

//wrapper arounf ResponceWriter
type responseWriter struct {
	http.ResponseWriter
	writes [][]byte
	status int
}

//implementing functions for Wrapper
func (rw *responseWriter) Write(b []byte) (int, error) {
	rw.writes = append(rw.writes, b)
	return len(b), nil
}

func (rw *responseWriter) WriteHeader(statusCode int) {
	rw.status = statusCode
}

//this add on function flush responceWriter into ResponceWriter
func (rw *responseWriter) flush() error {
	if rw.status != 0 {
		rw.ResponseWriter.WriteHeader(rw.status)
	}

	for _, write := range rw.writes {
		_, err := rw.ResponseWriter.Write(write)
		if err != nil {
			return err
		}
	}
	return nil
}

//handlers
func panicDemo(w http.ResponseWriter, r *http.Request) {
	funcThatPanics()
}

func panicAfterDemo(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "<h1>Hello!</h1>")
	funcThatPanics()
}

func funcThatPanics() {
	panic("Oh no!")
}

func hello(w http.ResponseWriter, r *http.Request) {

	fmt.Fprintln(w, "<h1>Hello!</h1>")
}
