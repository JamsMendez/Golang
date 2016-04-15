package main

import (
	"Golang/ServerHTTP/models"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/securecookie"
)

var hasKey = securecookie.GenerateRandomKey(64)
var blockKey = securecookie.GenerateRandomKey(32)
var cookieHandler = securecookie.New(hasKey, blockKey)

var model = models.UserModel{}

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/", getViewIndex).Methods("GET")
	router.HandleFunc("/other", getViewOther).Methods("GET")
	router.HandleFunc("/login", getViewLogin).Methods("GET")
	router.HandleFunc("/register", getViewRegister).Methods("GET")
	router.HandleFunc("/logout", logoutUser).Methods("GET")

	router.HandleFunc("/login", loginUser).Methods("POST")
	router.HandleFunc("/register", registerUser).Methods("POST")

	http.Handle("/", router)

	fmt.Println("Server running 3000")

	http.ListenAndServe(":3000", nil)
}

func getViewIndex(w http.ResponseWriter, r *http.Request) {
	if username, ok := isAuth(r); ok {
		indexPage := `<p>Bienvenido ` + username + ` !!!</p><p><a href="/logout">Cerrar sesión</a><p>`
		buffer := []byte(indexPage)

		w.Write(buffer)
	} else {
		http.Redirect(w, r, "/login", http.StatusMovedPermanently)
	}
}

func getViewLogin(w http.ResponseWriter, r *http.Request) {
	if _, ok := isAuth(r); ok {
		http.Redirect(w, r, "/", http.StatusMovedPermanently)
	} else {
		loginPage := `<h1>Login</h1>
                  <form method="POST" action="/login">
                      <label for="name">Usuario</label>
                      <input type="text" name="username">
                      <label for="password">Contraseña</label>
                      <input type="password" name="password">
                      <button type="submit">Login</button>
											<a href="/register">Ir a registro</a>
                  </form>`

		buffer := []byte(loginPage)
		w.Write(buffer)
	}
}

func getViewRegister(w http.ResponseWriter, r *http.Request) {
	if _, ok := isAuth(r); ok {
		http.Redirect(w, r, "/", http.StatusMovedPermanently)
	} else {
		registerPage := `<h1>Registro</h1>
                  <form method="POST" action="/register">
                      <label for="name">Usuario</label>
                      <input type="text" name="username">
                      <label for="password">Contraseña</label>
                      <input type="password" name="password">
											<label for="password">Contraseña(*)</label>
                      <input type="password" name="password_">
                      <button type="submit">Registrar</button>
											<a href="/login">Ir a login</a>
                  </form>`

		buffer := []byte(registerPage)
		w.Write(buffer)
	}
}

func getViewOther(w http.ResponseWriter, r *http.Request) {
	if username, ok := isAuth(r); ok {
		indexPage := `<p>Hola de nuevo ` + username + ` !!! </p><p><a href="/logout">Cerrar sesión</a><p>`
		buffer := []byte(indexPage)

		w.Write(buffer)
	} else {
		http.Redirect(w, r, "/login", http.StatusMovedPermanently)
	}
}

func loginUser(w http.ResponseWriter, r *http.Request) {
	redirectPath := "/"

	if _, ok := isAuth(r); !ok {
		username := r.FormValue("username")
		password := r.FormValue("password")

		if username != "" && password != "" {

			userIn := models.User{Username: username, Password: password}
			if userIn.ComparePassword() {
				user := map[string]string{"username": username}

				if encoded, err := cookieHandler.Encode("session", user); err == nil {
					cookie := &http.Cookie{
						Name:  "session",
						Value: encoded,
						Path:  "/",
					}

					http.SetCookie(w, cookie)
				}
			}

		} else {
			redirectPath = "/login"
		}
	}

	http.Redirect(w, r, redirectPath, http.StatusMovedPermanently)
}

func logoutUser(w http.ResponseWriter, r *http.Request) {
	cookie := &http.Cookie{
		Name:   "session",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	}

	http.SetCookie(w, cookie)

	http.Redirect(w, r, "/login", http.StatusMovedPermanently)
}

func registerUser(w http.ResponseWriter, r *http.Request) {
	redirectPath := "/"
	registerSuccess := false

	if _, ok := isAuth(r); !ok {
		username := r.FormValue("username")
		password := r.FormValue("password")
		passwordAgain := r.FormValue("password_")

		if username != "" && password != "" && passwordAgain != "" {
			if password == passwordAgain {
				json := models.User{Username: username, Password: password}
				_, err := model.Insert(json)

				if err == nil {
					registerSuccess = true
					redirectPath = "/login"
				}
			}
		}

		if !registerSuccess {
			redirectPath = "/register"
		}
	}

	http.Redirect(w, r, redirectPath, http.StatusMovedPermanently)
}

func isAuth(r *http.Request) (string, bool) {
	var username string
	if cookie, err := r.Cookie("session"); err == nil {
		cookieValue := make(map[string]string)
		if err = cookieHandler.Decode("session", cookie.Value, &cookieValue); err == nil {
			username = cookieValue["username"]
			return username, true
		}
	}

	return username, false
}
