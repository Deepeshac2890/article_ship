package Database

import (
	"article_ship/Models"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

var openDBConnection *sql.DB

func InitConnection() {
	// Open up our database connection.
	// I've set up a database on my local machine using phpmyadmin.
	// The database is called testDb
	dataSourceName := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", userName, password, hostName, port, dbName)
	db, err := sql.Open("mysql", dataSourceName)
	if err != nil {
		panic(err.Error())
		return
	}
	openDBConnection = db
}

func closeConnection() {
	err := openDBConnection.Close()
	if err != nil {
		panic(err.Error())
		return
	}
}

func InsertArticle(article Models.Article) {
	if openDBConnection != nil {
		// TODO: Add a check to see if id already present
		sqlQuery := fmt.Sprintf("INSERT INTO Articles (Title, Description, category, Content, id) VALUES ('%s', '%s', '%s', '%s', '%d')", article.Title, article.Desc, article.Category, article.Content, article.Id)
		insert, err := openDBConnection.Query(sqlQuery)
		if err != nil {
			panic(err.Error())
		}
		// be careful deferring Queries if you are using transactions
		defer func(insert *sql.Rows) {
			err := insert.Close()
			if err != nil {

			}
		}(insert)
	} else {
		fmt.Println("Db Connection is close")
	}
}

// TODO: Add Get, delete and update db operations as well
