package organizationuser

import "context"

type Repository interface {
	LinkUserToOrganization(ctx context.Context)
}
