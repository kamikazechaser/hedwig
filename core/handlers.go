// package main

// import (
// 	"net/http"

// 	"github.com/go-chi/render"
// )

// // placeholder
// type Res struct {
// 	Status int     `json:"status"`
// 	Ok     bool    `json:"ok"`
// 	Data   *Config `json:"data"`
// }

// func telegramHandler(w http.ResponseWriter, r *http.Request) {
// 	// placeholder
// 	ctx := r.Context()
// 	config, ok := ctx.Value(configKey).(*Config)
// 	if !ok {
// 		http.Error(w, http.StatusText(422), 422)
// 		return
// 	}

// 	render.JSON(w, r, Res{200, true, config})
// }
