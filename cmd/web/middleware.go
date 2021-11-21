package main

import ( 
	"net/http"
)

func secureHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-XSS-Protection", "1; mode=block") 
		w.Header().Set("X-Frame-Options", "deny")
		next.ServeHTTP(w, r) 
	})
}

func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//NEED CHECK
		/* if !isAuthorized(r) {
			w.WriteHeader(http.StatusForbidden)
			return
		} */
		next.ServeHTTP(w, r) 
	})
}