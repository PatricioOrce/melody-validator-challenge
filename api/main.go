package main

import (
	"melody-validator-challenge/cmd/server"
)

func main() {

	srv := server.New(":8080")

	if err := srv.ListenAndServe(); err != nil {
		panic(err)
	}

}
