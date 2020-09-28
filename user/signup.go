package user

import (

	"encoding/json"
	"io/ioutil"
	"io"
	"os"
	"strings"
	vars "github.com/Fifanon/online_library/varsAndFuncs"
	stct "github.com/Fifanon/online_library/structs"
	"net/http"
	dbconfig "github.com/Fifanon/online_library/config"
	"golang.org/x/crypto/bcrypt"
)
//EmailStruct **
type emailStruct struct{
     Email string `json:"email"`
}
//SignupProcessor **
func SignupProcessor(w http.ResponseWriter, r *http.Request) {
	_, err := os.OpenFile("./project_files/public/mphotos/"+vars.PhotoFileName, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		panic(err)
	}
	err = ioutil.WriteFile("./project_files/public/mphotos/"+vars.PhotoFileName, vars.PhotoFilebytes, 0)

	r.ParseMultipartForm(10 << 32)

	//call on dbconfig.GetMySQLDb for connection to the database
	db, err := dbconfig.GetMySQLDb()

	if err != nil {
		panic(err)
	}

	cost := bcrypt.DefaultCost
	hash, err := bcrypt.GenerateFromPassword([]byte(stct.User.Password), cost)
	if err != nil {
		panic(err)
	}
	stct.User.Password = string(hash)

	temp, err := db.Prepare(`insert into temporary_members (mb_firster,mb_laster,mb_email,mb_address,mb_tel,mb_pwd,m_status,m_photo)
            values(?,?,?,?,?,?,?,?);`)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_, err = temp.Exec(&stct.User.FirstName, &stct.User.LastName, &stct.User.Email, &stct.User.Address, &stct.User.PhoneNum, &stct.User.Password, &stct.User.Status, &stct.User.ImageName)
	if err != nil {
		panic(err)
	}
	db.Close()
	vars.Tpl.ExecuteTemplate(w, "signupSucc.html", stct.User)
}

//UploadPhotoFile **
func UploadPhotoFile(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(10 << 20)
	file, handler, err := r.FormFile("imgfile")
	if err != nil {
		panic(err)
	}
	vars.PhotoFileName = handler.Filename
	vars.PhotoFilebytes, err = ioutil.ReadAll(file)
	if err != nil {
		panic(err)
	}
	stct.User.Password = strings.Join(r.Form["password"],"")
	stct.User.ImageName = vars.PhotoFileName
	stct.User.FirstName = strings.Join(r.Form["firster"],"")
	stct.User.LastName = strings.Join(r.Form["laster"],"")
	stct.User.Email = strings.Join(r.Form["email"],"")
	stct.User.Address = strings.Join(r.Form["address"],"")
	stct.User.PhoneNum = strings.Join(r.Form["pnum"],"")
	stct.User.Status = strings.Join(r.Form["status"],"")


	defer file.Close()
	vars.Fileuploadmsg = vars.PhotoFileName
	vars.GotFile = true
	http.Redirect(w, r, "/signup/processing-continue", http.StatusSeeOther)
}

//CheckEmail **
func CheckEmail(w http.ResponseWriter, r *http.Request) {

	 emailjson, err := ioutil.ReadAll(r.Body)
	 if err != nil {
		 panic(err)
	 }
	 bytes := []byte(emailjson)

	 var emailSt emailStruct	
	 
	 json.Unmarshal(bytes, &emailSt)

	 //open database
	 db, err := dbconfig.GetMySQLDb()
	 if err != nil {
		 panic(err)
	 }
	 var email string = ""
	 
    qResult := db.QueryRow(`select m_email from members where m_email = ?;`, emailSt.Email)
     qResult.Scan(&email)
	 db.Close()
     w.Header().Set("Content-Type", "text/plain")
     w.WriteHeader(http.StatusOK)
	 _, err = io.WriteString(w, email)
	 if err != nil{
		 panic(err)
	 }
}