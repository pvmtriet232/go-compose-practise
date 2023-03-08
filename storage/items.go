// storage/items.go
package storage

import (
	"context"
	"fmt"
)

type CreateItemRequest struct {
	Name string
}

type Item struct {
	ID   string
	Name string
}

func (s *Storage) CreateItem(ctx context.Context, i CreateItemRequest) (*Item, error) {
	row := s.conn.QueryRowContext(ctx, "INSERT INTO items(name) VALUES($1) RETURNING id, name", i.Name)
	return ScanItem(row)
}

func (s *Storage) ListItems(ctx context.Context) ([]*Item, error) {
	rows, err := s.conn.QueryContext(ctx, "SELECT id, name FROM items")
	if err != nil {
		return nil, fmt.Errorf("could not retrieve items: %w", err)
	}
	defer rows.Close()

	var items []*Item
	for rows.Next() {
		item, err := ScanItem(rows)
		if err != nil {
			return nil, fmt.Errorf("could not scan item: %w", err)
		}

		items = append(items, item)
	}

	return items, nil
}

func ScanItem(s Scanner) (*Item, error) {
	i := &Item{}
	if err := s.Scan(&i.ID, &i.Name); err != nil {
		return nil, err
	}

	return i, nil
}
