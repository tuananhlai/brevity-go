package repository

import (
	"context"
	"fmt"
)

func (p *Postgres) ListAllDigitalAuthors(ctx context.Context) ([]*DigitalAuthor, error) {
	rows, err := p.qb.
		Select("da.id",
			"da.display_name",
			"da.system_prompt").
		From("digital_authors da").
		OrderBy("da.created_at DESC").
		RunWith(p.db).QueryContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting list of authors: %v", err)
	}

	items := make([]*DigitalAuthor, 0)
	for rows.Next() {
		var item DigitalAuthor
		err = rows.Scan(&item.ID, &item.DisplayName, &item.SystemPrompt)
		if err != nil {
			return nil, fmt.Errorf("error scanning item: %v", err)
		}
		items = append(items, &item)
	}

	return items, nil
}
