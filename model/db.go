package model

import (
	"database/sql"
	_ "github.com/lib/pq"
	"fmt"
	"log"
)

type Book struct{
	Isbn    string
	Title   string
	Author  string
	Price   float32
}

var DB *sql.DB
func init(){
	var err error
	DB , err = sql.Open("postgres", "postgres://bond:password@localhost/bookstore?sslmode=disable")
	if err != nil{
		log.Fatalln(err)
	}
	err = DB.Ping()
	if err != nil{
		log.Fatalln(err)
	}
	fmt.Println("database connected successfully")
}

func AllBooks()(*sql.Rows, error){
	rows ,err := DB.Query("SELECT * FROM books")
	return rows, err
}
func OneBooks(isbn string)*sql.Row{
	rows  := DB.QueryRow("SELECT * FROM books where isbn = $1", isbn)
	return rows
}
func InsertBook(bk Book)error{
	_ ,err := DB.Exec("INSERT INTO books (isbn, title, author, price) VALUES ($1, $2, $3, $4)", bk.Isbn,bk.Title,bk.Author,bk.Price)
	return  err
}
func UpdateBook(bk Book)error{
	_ ,err := DB.Exec("UPDATE books SET isbn = $1, title=$2, author=$3, price=$4 WHERE isbn=$1;", bk.Isbn, bk.Title, bk.Author, bk.Price)
	return  err
}
func DeleteBook(isbn string)error{
	_ ,err := DB.Exec("DELETE FROM books WHERE isbn = $1", isbn)
	return  err
}