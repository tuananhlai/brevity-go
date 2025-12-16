package repository

import (
	"context"
	"database/sql"
	"fmt"
	"math"

	"github.com/google/uuid"
)

// ListDigitalAuthorsByUserID returns a list of digital authors belonging to the given user with support
// for pagination. No upper limit is enforced on `page` and `pageSize` parameter.
func (p *Postgres) ListDigitalAuthorsByUserID(ctx context.Context, userID uuid.UUID, page int, pageSize int) (*ListDigitalAuthorsByUserIDResult, error) {
	limit := uint64(pageSize)
	offset := uint64((page - 1) * pageSize)

	tx, err := p.db.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelRepeatableRead,
	})
	if err != nil {
		return nil, fmt.Errorf("error starting transaction: %v", err)
	}

	totalItemsQuery := p.qb.
		Select("COUNT(*)").
		From("digital_authors").
		Where("user_id = ?", userID)

	var totalItems int
	err = totalItemsQuery.RunWith(tx).QueryRowContext(ctx).Scan(&totalItems)
	if err != nil {
		return nil, fmt.Errorf("error getting total number of items: %v", err)
	}

	itemsQuery := p.qb.
		Select("da.id",
			"da.display_name",
			"da.system_prompt",
			"da.api_key_id",
			"lak.encrypted_key AS api_key_encrypted_value",
			"da.created_at").
		From("digital_authors da").
		InnerJoin("llm_api_keys lak ON da.api_key_id = lak.id").
		Where("da.owner_id = ?", userID).
		OrderBy("da.created_at DESC").
		Limit(limit).
		Offset(offset)

	rows, err := itemsQuery.RunWith(tx).QueryContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting list of authors: %v", err)
	}

	items := make([]*DigitalAuthor, 0, pageSize)
	for rows.Next() {
		var item DigitalAuthor
		err = rows.Scan(&item.ID, &item.DisplayName, &item.SystemPrompt, &item.APIKeyID,
			&item.APIKeyEncryptedValue, &item.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("error scanning item: %v", err)
		}
		items = append(items, &item)
	}

	// COMMIT is not actually necessary, since we aren't modifying any data.
	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("error commiting tra")
	}

	totalPage := int(math.Ceil(float64(totalItems) / float64(pageSize)))

	return &ListDigitalAuthorsByUserIDResult{
		Items:      items,
		TotalPage:  totalPage,
		TotalItems: totalItems,
	}, nil
}

type ListDigitalAuthorsByUserIDResult struct {
	Items      []*DigitalAuthor
	TotalPage  int
	TotalItems int
}
