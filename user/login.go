package user

import (
	"os"
	"database/sql"
	s "github.com/Fifanon/online_library/session"
	"github.com/gorilla/securecookie"
	vars "github.com/Fifanon/online_library/varsAndFuncs"
	stct "github.com/Fifanon/online_library/structs"
	"net/http"
	dbconfig "github.com/Fifanon/online_library/config"
	"golang.org/x/crypto/bcrypt"
)
var cookieHandler = securecookie.New(
    securecookie.GenerateRandomKey(64),
	securecookie.GenerateRandomKey(32))

var status string
//LoginProcessor **
func LoginProcessor (w http.ResponseWriter, r *http.Request) {

	r.ParseForm()
	//call on dbconfig.GetMySQLDb for connection to the database
	db, err := dbconfig.GetMySQLDb()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	if r.Form.Get("status") == "Member" {
		qResult := db.QueryRow(`select m_email,m_password from members where m_email = $1`, r.Form.Get("email"))
		err = qResult.Scan(&stct.User.Email, &stct.User.Password)
		if err != nil{
		}
	} else {
		qResult := db.QueryRow(`select l_email,l_password from librarian where l_email = $1`, r.Form.Get("email"))
		err = qResult.Scan(&stct.User.Email, &stct.User.Password)
		if err != nil{
		}
	}
	db.Close()
	if err != nil {
		if err == sql.ErrNoRows {
			stct.Msg.EmailExistNot = "user does not exist"
			http.Redirect(w, r, "/home", 303)
			return
		}
	}
	plainPwd := []byte(r.Form.Get("pwd"))
	byteHash := []byte(stct.User.Password)
	err = bcrypt.CompareHashAndPassword(byteHash, plainPwd)
	if err != nil {
		stct.Msg.WrongPwd = "Wrong password"
		http.Redirect(w, r, "/home", http.StatusSeeOther)
		return
	}

	os.Setenv("EMAIL",stct.User.Email)

	status = r.Form.Get("status")
	   //if user authenticated, set session
	var redirectTarget string
	s.SetSession(stct.User.Email, w)
	r.ParseForm()

    if status == "Member" {
		redirectTarget ="/sci-library/welcome"
	}else {
		redirectTarget ="/sci-library/librarian/operations"
	} 
	http.Redirect(w, r, redirectTarget, http.StatusSeeOther)	
}

//MemberLoggedIn **
func MemberLoggedIn(w http.ResponseWriter, r *http.Request) {
	if	validated := s.GetSession(r);!validated{
		http.Redirect(w, r, "/home", http.StatusSeeOther)
		return
	}
		vars.Tpl.ExecuteTemplate(w, "loggedIn.html", stct.Msg)
		stct.Msg.BookExistsNot = ""
		return
}
