package article

import (
	"time"

	"github.com/google/uuid"
)

type Article struct {
	ID          uuid.UUID `db:"id"`
	Slug        string    `db:"slug"`
	Title       string    `db:"title"`
	Description string    `db:"description"`
	TextContent string    `db:"text_content"`
	Content     string    `db:"content"`
	AuthorID    uuid.UUID `db:"author_id"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}
