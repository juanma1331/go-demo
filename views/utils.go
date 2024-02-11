package views

import (
	"context"
	"go-demo/internal/shared"
)

func HasAuthenticatedUser(ctx context.Context) bool {
	user, ok := ctx.Value(shared.ContextUserKey).(*shared.AuthenticatedUser)
	if !ok {
		return false
	}

	return user != nil
}

func GetFlash(ctx context.Context) *[]shared.FlashMessage {
	flash, ok := ctx.Value(shared.ContextFlashKey).(*[]shared.FlashMessage)
	if !ok {
		return nil
	}

	return flash
}

func GetAuthenticatedUser(ctx context.Context) *shared.AuthenticatedUser {
	user, ok := ctx.Value(shared.ContextUserKey).(*shared.AuthenticatedUser)
	if !ok {
		return nil
	}

	return user
}
