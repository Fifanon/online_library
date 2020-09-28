package searchbook

import (
	"database/sql"
	"strings"
	"fmt"
	stct "online_library/structs"
	dbconfig "online_library/config"
)


//SearchByTitle **
func SearchByTitle(title string) (books []stct.BookStruct, found bool, errEnc bool) {

	db, err := dbconfig.GetMySQLDb()

	if err != nil {
		panic(err)
	}
	title = strings.ToLower(title)
	qr, err := db.Query(`select book_isbn,book_title,author_name,pages,subject_area,number,b_imagename from book_instances
                      where book_title LIKE ?;`, "%"+title+"%")
	for qr.Next() {
		err = qr.Scan(&stct.Bk.ISBN, &stct.Bk.Title, &stct.Bk.Author, &stct.Bk.Pages, &stct.Bk.Subject, &stct.Bk.Number, &stct.Bk.BookImageName)
		if err != nil {
			if err == sql.ErrNoRows {
			} else {
				errEnc = false
				return
			}
		}
		if stct.Bk.Number != 0 {
			stct.Bk.Availability = "AVAILABLE"
		} else {
			stct.Bk.Availability = "NOT AVAILABLE"
		}
		books = append(books, stct.Bk)
	}

	splited := strings.Fields(title)
  if(len(splited) > 1){
	
	ignoreWords := [8]string{"and","or","from","a","and","the","in", "of"}
	for _, fOrlname := range splited {
		var skip bool = false
		for _,word := range ignoreWords{
			if fOrlname == word{
			   skip = true
			}
		}
		if skip == true{
			continue
		}
		qr, err := db.Query(`select book_isbn,book_title,author_name,pages,subject_area,number,b_imagename from book_instances
                         where book_title LIKE ?;`, "%"+fOrlname+"%")
		for qr.Next() {
			err = qr.Scan(&stct.Bk.ISBN, &stct.Bk.Title, &stct.Bk.Author, &stct.Bk.Pages, &stct.Bk.Subject, &stct.Bk.Number, &stct.Bk.BookImageName)
			if err != nil {
				if err == sql.ErrNoRows {
				} else {
					errEnc = false
					return
				}
			}
			if stct.Bk.Number != 0 {
				stct.Bk.Availability = "AVAILABLE"
			} else {
				stct.Bk.Availability = "NOT AVAILABLE"
			}
			books = append(books, stct.Bk)
		}
	}
  }
  fmt.Println(len(books))
  fmt.Println(found)
	if len(books) == 0 {
		found = false
	} else {
		found = true
	}
	return books, found, errEnc
}
