package main

import (
	"database/sql"
	"net/http"
	"problem1/configs"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	"github.com/labstack/echo/v4"
)

func main() {
	conf := configs.Get()

	db, err := sql.Open(conf.DB.Driver, conf.DB.DataSource)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	e := echo.New()

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "minimal_sns_app")
	})

	// フレンドの情報を格納する構造体
	type Friend struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	}

	// ユーザーのフレンドリストを取得するエンドポイント
	e.GET("/get_friend_list", func(c echo.Context) error {
		// ユーザーIDを取得
		userID := c.QueryParam("ID")
		if userID == "" {
			return c.JSON(http.StatusBadRequest, "ID is required")
		}

		query := `
			SELECT u.id, u.name
			FROM users u
			JOIN friend_link f ON (u.id = IF(f.user1_id = ?, f.user2_id, f.user1_id))
			WHERE f.user1_id = ? OR f.user2_id = ?
		`
		rows, err := db.Query(query, userID)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, "Internal server error")
		}
		defer rows.Close()

		var friends []Friend

		for rows.Next() {
			var friend Friend
			if err := rows.Scan(&friend.ID, &friend.Name); err != nil {
				return c.JSON(http.StatusInternalServerError, "Internal server error")
			}
			friends = append(friends, friend)
		}

		// クエリの実行中にエラーが発生した場合はエラーを処理
		if err := rows.Err(); err != nil {
			return c.JSON(http.StatusInternalServerError, "Internal server error")
		}

		// フレンドリストを返す
		return c.JSON(http.StatusOK, friends)
	})

	e.GET("/get_friend_of_friend_list", func(c echo.Context) error {
		// FIXME
		return nil
	})

	e.GET("/get_friend_of_friend_list_paging", func(c echo.Context) error {
		// FIXME
		return nil
	})

	e.Logger.Fatal(e.Start(":" + strconv.Itoa(conf.Server.Port)))
}
