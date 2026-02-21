package repository

import (
	"context"
	"fmt"

	"github.com/lib/pq"
)

// ListDigitalAuthorsWithArticleSlugs returns a list of digital authors along with the slugs of their existing articles.
// This is used avoid duplication when generating new articles.
func (p *Postgres) ListDigitalAuthorsWithArticleSlugs(ctx context.Context) ([]*DigitalAuthorWithArticleSlugs, error) {
	rows, err := p.qb.
		Select("da.id", "da.system_prompt", "COALESCE(ARRAY_AGG(a.slug) FILTER (WHERE a.slug IS NOT NULL), '{}') AS article_slugs").
		From("digital_authors da").
		LeftJoin("articles a ON a.author_id = da.id").
		GroupBy("da.id", "da.system_prompt").
		RunWith(p.db).QueryContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting list of authors with slugs: %v", err)
	}
	defer rows.Close()

	items := make([]*DigitalAuthorWithArticleSlugs, 0)
	for rows.Next() {
		var item DigitalAuthorWithArticleSlugs
		err = rows.Scan(&item.ID, &item.SystemPrompt, pq.Array(&item.ArticleSlugs))
		if err != nil {
			return nil, fmt.Errorf("error scanning item: %v", err)
		}
		items = append(items, &item)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating authors with slugs: %v", err)
	}

	return items, nil
}
