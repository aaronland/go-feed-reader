package reader

import (
	"database/sql"
	"encoding/json"
	"github.com/mmcdole/gofeed"
)

func DatabaseRowsToFeeds(rows *sql.Rows) ([]*gofeed.Feed, error) {

	feeds := make([]*gofeed.Feed, 0)

	for rows.Next() {

		var body string
		err := rows.Scan(&body)

		if err != nil {
			return nil, err
		}

		var f gofeed.Feed

		err = json.Unmarshal([]byte(body), &f)

		if err != nil {
			return nil, err
		}

		feeds = append(feeds, &f)
	}

	err := rows.Err()

	if err != nil {
		return nil, err
	}

	return feeds, nil
}

func DatabaseRowToFeed(row *sql.Row) (*gofeed.Feed, error) {

	var body string
	err := row.Scan(&body)

	if err != nil {
		return nil, err
	}

	var i gofeed.Feed

	err = json.Unmarshal([]byte(body), &i)

	if err != nil {
		return nil, err
	}

	return &i, nil
}
