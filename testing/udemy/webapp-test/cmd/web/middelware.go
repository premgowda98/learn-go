package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
)

type ContextKey string

const contextUseKey ContextKey = "user_ip"

func (app *application) ipFromConext(ctx context.Context) string {
	return ctx.Value(contextUseKey).(string)
}

func (app *application) addIPToContext(next http.Handler) http.Handler{
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
		var ctx = context.Background()

		ip, err := getIP(r)

		if err !=nil {
			ip = "unknown"
		}

		ctx = context.WithValue(r.Context(),contextUseKey, ip)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func getIP(r *http.Request)(string, error){
	ip, _, err := net.SplitHostPort(r.RemoteAddr)

	if err != nil {
		return "unknown", nil
	}

	userIP := net.ParseIP(ip)
	if userIP == nil{
		return "", fmt.Errorf("not in correct format")
	}

	forward := r.Header.Get("X-Forwarded-For")

	if len(forward) >0 {
		ip = forward
	}

	if len(ip)==0{
		ip = "forward"
	}

	return ip, nil

}
