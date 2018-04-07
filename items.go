package reader

import (
	"database/sql"
	"encoding/json"
	"github.com/mmcdole/gofeed"
)

func DatabaseRowsToFeedItems(rows *sql.Rows) ([]*gofeed.Item, error) {

	items := make([]*gofeed.Item, 0)

	for rows.Next() {

		var body string
		err := rows.Scan(&body)

		if err != nil {
			return nil, err
		}

		var i gofeed.Item

		err = json.Unmarshal([]byte(body), &i)

		if err != nil {
			return nil, err
		}

		items = append(items, &i)
	}

	err := rows.Err()

	if err != nil {
		return nil, err
	}

	return items, nil
}

func DatabaseRowToFeedItem(row *sql.Row) (*gofeed.Item, error) {

	var body string
	err := row.Scan(&body)

	if err != nil {
		return nil, err
	}

	var i gofeed.Item

	err = json.Unmarshal([]byte(body), &i)

	if err != nil {
		return nil, err
	}

	return &i, nil
}
