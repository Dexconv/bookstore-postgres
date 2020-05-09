package main

import (
	"github.com/Dexconv/postgres/controller"
	"net/http"
)





func main(){
	http.HandleFunc("/", index)
	http.HandleFunc("/books", controller.Bookindex)
	http.HandleFunc("/books/show", controller.Onebook)
	http.HandleFunc("/books/create", controller.BookCreateForm)
	http.HandleFunc("/books/create/process", controller.BookCreateProcess)
	http.HandleFunc("/books/update", controller.BookUpdateForm)
	http.HandleFunc("/books/update/process", controller.BookUpdateProcess)
	http.HandleFunc("/books/delete/process", controller.BookDeleteProcess)

	http.ListenAndServe(":8080", nil)

}
func index(w http.ResponseWriter , r *http.Request){
	http.Redirect(w,r,"/books", http.StatusSeeOther)
}

