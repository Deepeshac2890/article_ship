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

func InsertArticle(article Models.Article) bool {
	if openDBConnection != nil {
		_, isDuplicate := GetSingleArticle(article.Id)
		if isDuplicate {
			return false
		} else {
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
			return true
		}
	} else {
		fmt.Println("DB Connection is closed")
	}
	return false
}

func GetAllArticles() Models.Articles {
	var articles Models.Articles
	if openDBConnection != nil {
		sqlQuery := fmt.Sprintf("SELECT * FROM Articles")
		results, err := openDBConnection.Query(sqlQuery)
		if err != nil {
			panic(err.Error())
		}
		for results.Next() {
			var article Models.Article
			err = results.Scan(&article.Title, &article.Desc, &article.Category, &article.Content, &article.Id)
			if err != nil {
				panic(err.Error())
			}
			articles = append(articles, article)
		}
	} else {
		fmt.Println("DB Connection is closed")
	}
	return articles
}

func GetSingleArticle(id int32) (Models.Article, bool) {
	var article Models.Article
	if openDBConnection != nil {
		sqlQuery := fmt.Sprintf("SELECT * FROM Articles where id=%d", id)
		result := openDBConnection.QueryRow(sqlQuery)
		err := result.Scan(&article.Title, &article.Desc, &article.Category, &article.Content, &article.Id)
		if err != nil {
			return Models.Article{}, false
		}
		return article, true
	} else {
		fmt.Println("DB Connection is closed")
		return article, false
	}
}

func DeleteSingleArticle(id int32) bool {
	if openDBConnection != nil {
		sqlQuery := fmt.Sprintf("DELETE * FROM Articles where id=%d", id)
		_, err := openDBConnection.Query(sqlQuery)
		if err != nil {
			panic(err.Error())
		}
		return true
	} else {
		fmt.Println("DB Connection is closed")
	}
	return false
}

func DeleteAllArticles() bool {
	if openDBConnection != nil {
		sqlQuery := fmt.Sprintf("DELETE FROM Articles")
		_, err := openDBConnection.Query(sqlQuery)
		if err != nil {
			panic(err.Error())
		}
		return true
	} else {
		fmt.Println("DB Connection is closed")
	}
	return false
}

func UpdateArticle(article Models.Article, id int32) bool {
	if openDBConnection != nil {
		_, found := GetSingleArticle(id)
		if found {
			sqlQuery := fmt.Sprintf("UPDATE Articles SET (Title = '%s', Description = '%s', category = '%s', Content = '%s') where id=%d", article.Title, article.Desc, article.Category, article.Content, id)
			updated, err := openDBConnection.Query(sqlQuery)
			if err != nil {
				panic(err.Error())
			}
			err = updated.Close()
			if err != nil {
				return false
			}
			return true
		}
	} else {
		fmt.Println("DB Connection is closed")
	}
	return false
}
