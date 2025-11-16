package main

import "net/http"

func main() {
	server := InitializedServer()

	http.HandleFunc("/api/signup", server.UserController.SignUpAddUser)
	http.HandleFunc("/api/signin", server.UserController.SigninUser)
	http.HandleFunc("/verif/", server.UserController.Verification)

	fs := http.FileServer(http.Dir("./"))
	http.Handle("/", fs)

	http.ListenAndServe(":8081", nil)
}
