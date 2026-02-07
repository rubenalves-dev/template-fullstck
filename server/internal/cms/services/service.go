package services

import (
	"context"
	"encoding/json"
	"regexp"
	"strings"

	"github.com/google/uuid"
	"github.com/nats-io/nats.go"
	"github.com/rubenalves-dev/template-fullstack/server/internal/cms/domain"
	"github.com/rubenalves-dev/template-fullstack/server/pkg/events"
)

type service struct {
	repo domain.Repository
	nc   *nats.Conn
}

func NewService(repo domain.Repository, nc *nats.Conn) domain.Service {
	return &service{
		repo: repo,
		nc:   nc,
	}
}

func (s service) CreateDraft(ctx context.Context, title string) error {
	page := &domain.Page{
		ID:     uuid.New(),
		Title:  title,
		Slug:   slugify(title),
		Status: "draft",
	}

	err := s.repo.Create(ctx, page)
	if err != nil {
		return err
	}

	event := events.CmsPageDraftedData{
		PageID: page.ID,
		Title:  page.Title,
		Slug:   page.Slug,
	}
	eventBytes, _ := json.Marshal(event)
	return s.nc.Publish(events.CmsPageDrafted, eventBytes)
}

func (s service) PublishPage(ctx context.Context, id uuid.UUID) error {
	err := s.repo.UpdateStatus(ctx, id, "published")
	if err != nil {
		return err
	}

	page, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	event := events.CmsPagePublishedData{
		PageID: page.ID,
		Title:  page.Title,
		Slug:   page.Slug,
	}
	eventBytes, _ := json.Marshal(event)
	return s.nc.Publish(events.CmsPagePublished, eventBytes)
}

func (s service) ArchivePage(ctx context.Context, id uuid.UUID) error {
	page, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	err = s.repo.UpdateStatus(ctx, page.ID, "archived")
	if err != nil {
		return err
	}

	event := events.CmsPageArchivedData{
		PageID: page.ID,
		Title:  page.Title,
		Slug:   page.Slug,
	}
	eventBytes, _ := json.Marshal(event)
	return s.nc.Publish(events.CmsPageArchived, eventBytes)
}

func (s service) UpdatePageMetadata(ctx context.Context, id uuid.UUID, req domain.PageUpdateRequest) error {
	page, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if req.Title != nil {
		page.Title = *req.Title
	}
	if req.Slug != nil {
		page.Slug = *req.Slug
	}
	if req.SEODescription != nil {
		page.SEODescription = *req.SEODescription
	}
	if req.Keywords != nil {
		page.SEOKeywords = req.Keywords
	}

	return s.repo.Update(ctx, page)
}

func (s service) UpdatePageLayout(ctx context.Context, id uuid.UUID, layout []domain.RowRequest) error {
	domainRows := make([]domain.Row, len(layout))
	for i, rowReq := range layout {
		rowID := uuid.New()
		if rowReq.ID != nil {
			rowID = *rowReq.ID
		}

		domainRows[i] = domain.Row{
			ID:               rowID,
			PageID:           id,
			OrderIndex:       rowReq.SortOrder,
			CSSClass:         rowReq.CSSClass,
			BackgroundConfig: rowReq.BackgroundConfig,
		}

		domainCols := make([]domain.Column, len(rowReq.Columns))
		for j, colReq := range rowReq.Columns {
			colID := uuid.New()
			domainCols[j] = domain.Column{
				ID:         colID,
				RowID:      rowID,
				OrderIndex: j,
				CSSClass:   colReq.CSSClass,
				WidthSM:    colReq.WidthSM,
				WidthMD:    colReq.WidthMD,
				WidthLG:    colReq.WidthLG,
				WidthXL:    colReq.WidthXL,
			}

			domainBlocks := make([]domain.Block, len(colReq.Blocks))
			for k, blockReq := range colReq.Blocks {
				domainBlocks[k] = domain.Block{
					ID:         uuid.New(),
					ColumnID:   colID,
					Type:       blockReq.Type,
					OrderIndex: k,
					Content:    blockReq.Content,
				}
			}
			domainCols[j].Blocks = domainBlocks
		}
		domainRows[i].Columns = domainCols
	}

	err := s.repo.SaveLayout(ctx, id, domainRows)
	if err != nil {
		return err
	}

	event := events.CmsPageLayoutUpdatedData{
		PageID: id,
	}
	eventBytes, _ := json.Marshal(event)
	return s.nc.Publish(events.CmsPageLayoutUpdated, eventBytes)
}

func (s service) GetPageBySlug(ctx context.Context, Slug string) (*domain.Page, error) {
	page, err := s.repo.GetBySlug(ctx, Slug)
	if err != nil {
		return nil, err
	}

	layout, err := s.repo.GetFullLayout(ctx, page.ID)
	if err != nil {
		return nil, err
	}

	page.Rows = layout
	return page, nil
}

func slugify(text string) string {
	var re = regexp.MustCompile("[^a-z0-9]+")
	return strings.Trim(re.ReplaceAllString(strings.ToLower(text), "-"), "-")
}
