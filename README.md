In the recover script I have  been looking at the panic and recover mechanisms in Go and utilizing them to create middleware for an HTTP server.

Given a simple web server (see main.go) that can potentially panic, I have created an http.Handler that wraps the existing mux and will recover from any panics and then does the following:

1. Logs the error, as well as the stack trace.
2. Sets the status code to http.StatusInternalServerError (500) whenever a panic occurs.
3. Write a "Something went wrong" message when a panic occurs.
4. If the environment is set to be development, print the stack trace and the error to the webpage as well as to the   logs. Otherwise default to the "Something went wrong" message described in (3).
