package routers

import "github.com/rs/cors"

// C - Cors initialize variable
var C = cors.New(cors.Options{
	AllowedOrigins:     []string{"*"},
	AllowCredentials:   true,
	AllowedMethods:     []string{"GET", "PUT", "POST", "DELETE", "OPTIONS"},
	AllowedHeaders:     []string{"Content-Type", "Bearer", "Bearer ", "content-type", "Origin", "Accept", "authorization", "Authorization"},
	OptionsPassthrough: true,
})
