package views

import (
	"context"
	"go-demo/internal/app"
)

func HasAuthenticatedUser(ctx context.Context) bool {
	user, ok := ctx.Value(app.ContextUserKey).(*app.AuthenticatedUser)
	if !ok {
		return false
	}

	return user != nil
}

func GetFlash(ctx context.Context) *[]app.FlashMessage {
	flash, ok := ctx.Value(app.ContextFlashKey).(*[]app.FlashMessage)
	if !ok {
		return nil
	}

	return flash
}
