package main

import (
	"fmt"
	"github.com/satori/go.uuid"
	"golangauth/src/components"
	"golangauth/src/config"
	"golangauth/src/modules/user/model"
	"golangauth/src/modules/user/repository"
	"html/template"
	"net/http"
	"os"
)

// INIT TEMPLATE

var tpl *template.Template

func init() {
	tpl = template.Must(template.ParseGlob("public/*"))
}

// END OF INIT TEMPLATE

//GET HANDLER

func homeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		errorHandler(w, r, http.StatusNotFound)
		return
	} else {
		userEmail := components.GetEmailFromCookie(r)

		if userEmail == "" {
			tpl.ExecuteTemplate(w, "home.template", map[string]string{
				"isLoggedIn": "false",
			})
		} else {
			tpl.ExecuteTemplate(w, "home.template", map[string]string{
				"isLoggedIn": "true",
			})
		}
	}
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/register" {
		errorHandler(w, r, http.StatusNotFound)
		return
	} else {
		userEmail := components.GetEmailFromCookie(r)

		if userEmail == "" {

			flashMessage, err := components.GetFlashMessage(w, r, "message")

			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			if components.CheckMessage("full_name_empty", flashMessage) {
				errorMessage := "Nama lengkap masih kosong !"
				tpl.ExecuteTemplate(w, "register.template", map[string]string{
					"Error":   "true",
					"Message": errorMessage,
				})
			} else if components.CheckMessage("email_empty", flashMessage) {
				errorMessage := "Email masih kosong !"
				tpl.ExecuteTemplate(w, "register.template", map[string]string{
					"Error":   "true",
					"Message": errorMessage,
				})
			} else if components.CheckMessage("password_empty", flashMessage) {
				errorMessage := "Password masih kosong !"
				tpl.ExecuteTemplate(w, "register.template", map[string]string{
					"Error":   "true",
					"Message": errorMessage,
				})
			} else if components.CheckMessage("full_name_error_regex", flashMessage) {
				errorMessage := "Nama lengkap hanya boleh berisi huruf dan spasi !"
				tpl.ExecuteTemplate(w, "register.template", map[string]string{
					"Error":   "true",
					"Message": errorMessage,
				})
			} else if components.CheckMessage("email_error_regex", flashMessage) {
				errorMessage := "Format email yang anda masukkan salah !"
				tpl.ExecuteTemplate(w, "register.template", map[string]string{
					"Error":   "true",
					"Message": errorMessage,
				})
			} else if components.CheckMessage("password_error_regex", flashMessage) {
				errorMessage := "Password hanya boleh beirisi huruf dan angka !"
				tpl.ExecuteTemplate(w, "register.template", map[string]string{
					"Error":   "true",
					"Message": errorMessage,
				})
			} else if components.CheckMessage("register_error", flashMessage) {
				errorMessage := "Pendaftaran tidak berhasil silahkan ulangi kembali !"
				tpl.ExecuteTemplate(w, "register.template", map[string]string{
					"Error":   "true",
					"Message": errorMessage,
				})
			} else if components.CheckMessage("register_success", flashMessage) {
				successMessage := "Pendaftaran berhasil !"
				tpl.ExecuteTemplate(w, "register.template", map[string]string{
					"Success": "true",
					"Message": successMessage,
				})
			} else {
				tpl.ExecuteTemplate(w, "register.template", nil)
			}
		} else {
			http.Redirect(w, r, "/", 301)
		}
	}
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/login" {
		errorHandler(w, r, http.StatusNotFound)
		return
	} else {
		userEmail := components.GetEmailFromCookie(r)

		if userEmail == "" {

			flashMessage, err := components.GetFlashMessage(w, r, "message")

			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			if components.CheckMessage("email_empty", flashMessage) {
				errorMessage := "Email masih kosong !"
				tpl.ExecuteTemplate(w, "login.template", map[string]string{
					"Error":   "true",
					"Message": errorMessage,
				})
			} else if components.CheckMessage("password_empty", flashMessage) {
				errorMessage := "Password masih kosong !"
				tpl.ExecuteTemplate(w, "login.template", map[string]string{
					"Error":   "true",
					"Message": errorMessage,
				})
			} else if components.CheckMessage("email_error_regex", flashMessage) {
				errorMessage := "Format email yang anda masukkan salah !"
				tpl.ExecuteTemplate(w, "login.template", map[string]string{
					"Error":   "true",
					"Message": errorMessage,
				})
			} else if components.CheckMessage("password_error_regex", flashMessage) {
				errorMessage := "Password hanya boleh beirisi huruf dan angka !"
				tpl.ExecuteTemplate(w, "login.template", map[string]string{
					"Error":   "true",
					"Message": errorMessage,
				})
			} else if components.CheckMessage("invalid_email_login", flashMessage) {
				errorMessage := "Email yang anda masukkan salah !"
				tpl.ExecuteTemplate(w, "login.template", map[string]string{
					"Error":   "true",
					"Message": errorMessage,
				})
			} else if components.CheckMessage("invalid_password_login", flashMessage) {
				errorMessage := "Password yang anda masukkan salah !"
				tpl.ExecuteTemplate(w, "login.template", map[string]string{
					"Error":   "true",
					"Message": errorMessage,
				})
			} else {
				tpl.ExecuteTemplate(w, "login.template", nil)
			}
		} else {
			http.Redirect(w, r, "/", 301)
		}
	}
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/logout" {
		errorHandler(w, r, http.StatusNotFound)
		return
	} else {
		userEmail := components.GetEmailFromCookie(r)

		if userEmail == "" {
			errorHandler(w, r, http.StatusNotFound)
		} else {

			components.ClearCookie(w)

			http.Redirect(w, r, "/login", 301)
		}
	}
}

func errorHandler(w http.ResponseWriter, r *http.Request, status int) {
	w.WriteHeader(status)
	if status == http.StatusNotFound {
		tpl.ExecuteTemplate(w, "404.template", nil)
	} else {
		tpl.ExecuteTemplate(w, "404.template", nil)
	}
}

// END OF GET HANDLER

// POST ACTION

func registerAction(w http.ResponseWriter, r *http.Request) {

	if r.Method == "POST" {

		db, err := config.GetMongoDB()

		if err != nil {
			fmt.Println("Gagal menghubungkan ke database!")
			os.Exit(2)
		}

		if r.FormValue("fullname") == "" {

			msg := []byte("full_name_empty")

			components.SetFlashMessage(w, "message", msg)

			http.Redirect(w, r, "/register", 301)

		} else if r.FormValue("email") == "" {
			msg := []byte("email_empty")

			components.SetFlashMessage(w, "message", msg)

			http.Redirect(w, r, "/register", 301)
		} else if r.FormValue("password") == "" {
			msg := []byte("password_empty")

			components.SetFlashMessage(w, "message", msg)

			http.Redirect(w, r, "/register", 301)
		} else if components.ValidateFullName(r.FormValue("fullname")) == false {
			msg := []byte("full_name_error_regex")

			components.SetFlashMessage(w, "message", msg)

			http.Redirect(w, r, "/register", 301)
		} else if components.ValidateEmail(r.FormValue("email")) == false {
			msg := []byte("email_error_regex")

			components.SetFlashMessage(w, "message", msg)

			http.Redirect(w, r, "/register", 301)
		} else if components.ValidatePassword(r.FormValue("password")) == false {
			msg := []byte("password_error_regex")

			components.SetFlashMessage(w, "message", msg)

			http.Redirect(w, r, "/register", 301)
		} else {

			var userRepository repository.UserRepository

			userRepository = repository.NewUserRepositoryMongo(db, "pengguna")

			makeID := uuid.NewV1()

			hashedPassword, _ := components.HashPassword(r.FormValue("password"))

			var userModel model.User

			userModel.ID = makeID.String()

			userModel.FullName = r.FormValue("fullname")

			userModel.Email = r.FormValue("email")

			userModel.Password = hashedPassword

			err = userRepository.Insert(&userModel)

			if err != nil {
				msg := []byte("register_error")

				components.SetFlashMessage(w, "message", msg)

				http.Redirect(w, r, "/register", 301)
			} else {
				msg := []byte("register_success")

				components.SetFlashMessage(w, "message", msg)

				http.Redirect(w, r, "/register", 301)
			}
		}
	} else {
		errorHandler(w, r, http.StatusNotFound)
	}
}

func loginAction(w http.ResponseWriter, r *http.Request) {

	if r.Method == "POST" {

		if r.FormValue("email") == "" {
			msg := []byte("email_empty")

			components.SetFlashMessage(w, "message", msg)

			http.Redirect(w, r, "/login", 301)
		} else if r.FormValue("email") != "" && r.FormValue("password") == "" {
			msg := []byte("password_empty")

			components.SetFlashMessage(w, "message", msg)

			http.Redirect(w, r, "/login", 301)
		} else if r.FormValue("email") != "" && r.FormValue("password") != "" && components.ValidateEmail(r.FormValue("email")) == false {

			msg := []byte("email_error_regex")

			components.SetFlashMessage(w, "message", msg)

			http.Redirect(w, r, "/login", 301)
		} else if r.FormValue("email") != "" && r.FormValue("password") != "" && components.ValidateEmail(r.FormValue("email")) == true && components.ValidatePassword(r.FormValue("password")) == false {
			msg := []byte("password_error_regex")

			components.SetFlashMessage(w, "message", msg)

			http.Redirect(w, r, "/login", 301)
		} else if r.FormValue("email") != "" && r.FormValue("password") != "" && components.ValidateEmail(r.FormValue("email")) == true && components.ValidatePassword(r.FormValue("password")) == true {

			statusLogin := checkUserIsRegistered(r.FormValue("email"), r.FormValue("password"))

			msg := []byte(statusLogin)

			components.SetFlashMessage(w, "message", msg)

			if statusLogin == "login_success" {

				components.SetCookie(r.FormValue("email"), w)

				http.Redirect(w, r, "/", 301)

			} else {

				http.Redirect(w, r, "/login", 301)

			}
		}
	} else {
		errorHandler(w, r, http.StatusNotFound)
	}
}

// END OF POST ACTION

// CHECK USER FUNCTION

func checkUserIsRegistered(email string, password string) string {
	db, err := config.GetMongoDB()

	if err != nil {
		fmt.Println("Gagal menghubungkan ke database!")
		os.Exit(2)
	}

	var userRepository repository.UserRepository

	userRepository = repository.NewUserRepositoryMongo(db, "pengguna")

	userData, err1 := userRepository.FindAll()

	if err1 != nil {
		return "invalid_email_login"
	} else {
		for _, user := range userData {
			if email == user.Email {
				if components.CheckPasswordHash(password, user.Password) == true {
					return "login_success"
				} else {
					return "invalid_password_login"
				}
			} else {
				return "invalid_email_login"
			}
		}
	}

	return "invalid_email_login"
}

// END OF CHECK USER FUNCTION

//MAIN FUNCTION

func main() {

	http.HandleFunc("/", homeHandler)

	//Register Page

	http.HandleFunc("/register", registerHandler)

	http.HandleFunc("/registerAction", registerAction)

	//Login Page

	http.HandleFunc("/login", loginHandler)

	http.HandleFunc("/loginAction", loginAction)

	//Logout

	http.HandleFunc("/logout", logoutHandler)

	fmt.Println("web server berjalan akses http://localhost:8045/")

	http.ListenAndServe(":8045", nil)
}

// END OF MAIN FUNCTION
