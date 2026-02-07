package domain

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*Page, error)
	GetBySlug(ctx context.Context, slug string) (*Page, error)
	List(ctx context.Context) ([]Page, error)
	Create(ctx context.Context, page *Page) error
	Update(ctx context.Context, page *Page) error
	Delete(ctx context.Context, id uuid.UUID) error

	// Layout Management
	// Use transactions here to ensure all or nothing updates
	SaveLayout(ctx context.Context, pageID uuid.UUID, rows []Row) error
	GetFullLayout(ctx context.Context, pageID uuid.UUID) ([]Row, error)

	// SEO & Status
	UpdateStatus(ctx context.Context, id uuid.UUID, status string) error
}

type Service interface {
	CreateDraft(ctx context.Context, title string) error
	PublishPage(ctx context.Context, id uuid.UUID) error
	ArchivePage(ctx context.Context, id uuid.UUID) error

	// Content Management
	UpdatePageMetadata(ctx context.Context, id uuid.UUID, req PageUpdateRequest) error
	UpdatePageLayout(ctx context.Context, id uuid.UUID, layout []RowRequest) error

	// Public Facing
	GetPageBySlug(ctx context.Context, Slug string) (*Page, error)
}
