package invoice

import "context"

type Usecase interface {
	BuildInvoice(ctx context.Context, request Content) error
}
