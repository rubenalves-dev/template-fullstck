package cms

import (
	"encoding/json"
	"log"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nats-io/nats.go"
	"github.com/rubenalves-dev/template-fullstack/server/internal/cms/delivery/events"
	"github.com/rubenalves-dev/template-fullstack/server/internal/cms/delivery/http"
	"github.com/rubenalves-dev/template-fullstack/server/internal/cms/domain"
	"github.com/rubenalves-dev/template-fullstack/server/internal/cms/repositories"
	"github.com/rubenalves-dev/template-fullstack/server/internal/cms/services"
	menuDomain "github.com/rubenalves-dev/template-fullstack/server/internal/platform/menu"
	globalEvents "github.com/rubenalves-dev/template-fullstack/server/pkg/events"
)

type CmsModule struct {
	Service domain.Service
}

func NewModule(pool *pgxpool.Pool, nc *nats.Conn) *CmsModule {
	repo := repositories.NewPgxRepository(pool)
	svc := services.NewService(repo, nc)

	events.RegisterListeners(nc, svc)

	// Register Permissions
	go func() {
		// Small delay to ensure auth module is listening if they start simultaneously,
		// though in main.go auth is started first.
		// But in NATS, listeners might take a ms to register.
		// However, "main.go" does "auth.NewModule" (registers listeners) THEN "cms.NewModule".
		// So it should be fine synchronously, but NATS publish is async anyway.
		payload := globalEvents.SystemPermissionsRegisteredData{
			Module:      "cms",
			Permissions: domain.GetAvailablePermission(),
		}
		data, _ := json.Marshal(payload)
		err := nc.Publish(globalEvents.SystemPermissionsRegister, data)
		if err != nil {
			log.Printf("[ERROR] Failed to publish permissions for CMS module: %v", err)
		}

		menuPayload := globalEvents.SystemMenusRegisteredData{
			Domain:  "cms",
			Version: 1,
			Menu:    []menuDomain.MenuDefinition{MenuDefinition},
		}
		menuData, _ := json.Marshal(menuPayload)
		err = nc.Publish(globalEvents.SystemMenusRegister, menuData)
		if err != nil {
			log.Printf("[ERROR] Failed to publish menus for CMS module: %v", err)
		}
	}()

	return &CmsModule{Service: svc}
}

func (m *CmsModule) RegisterRoutes(r chi.Router) {
	http.RegisterHTTPHandlers(r, m.Service)
}
