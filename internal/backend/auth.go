package backend

import (
	"context"

	"github.com/roadmap-thesis/backend/internal/domain"
	"github.com/roadmap-thesis/backend/internal/io"
	"github.com/roadmap-thesis/backend/pkg/apperrors"
	"github.com/roadmap-thesis/backend/pkg/auth"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

func (b *backend) Auth(ctx context.Context, input io.AuthInput) (io.AuthOutput, error) {
	ctx, span := tracer.Start(ctx, "(*backend.Auth)", trace.WithAttributes(attribute.String("email", input.Email)))
	defer span.End()

	var output io.AuthOutput

	result, err := b.registerEmail(ctx, input)
	if err != nil {
		return output, err
	}

	token, err := auth.CreateToken(result.id)
	if err != nil {
		return output, err
	}

	output.Created = result.created
	output.Token = token
	return output, nil
}

type registerEmailOutput struct {
	id      int
	email   string
	created bool
}

func (b *backend) registerEmail(ctx context.Context, input io.AuthInput) (registerEmailOutput, error) {
	ctx, span := tracer.Start(ctx, "(*backend.registerEmail)")
	defer span.End()

	existingAccount, err := b.repository.Account.GetByEmail(ctx, input.Email)
	if err != nil && err != domain.ErrAccountNotFound {
		return registerEmailOutput{}, err
	}

	// sign in if account already exists
	if !existingAccount.IsZero() {
		span.SetAttributes(attribute.Bool("create_account", false))
		matched := existingAccount.CheckPassword(input.Password)

		if !matched {
			return registerEmailOutput{}, apperrors.InvalidCredentials()
		}

		return registerEmailOutput{id: existingAccount.ID, email: existingAccount.Email}, nil
	}

	span.SetAttributes(attribute.Bool("create_account", true))
	profile := domain.NewProfile(input.Name, input.Avatar)
	account, err := domain.NewAccount(input.Email, input.Password, profile)
	if err != nil {
		return registerEmailOutput{}, err
	}

	createdAccount, err := b.repository.Account.Save(ctx, account)
	if err != nil {
		span.RecordError(err)
		return registerEmailOutput{}, err
	}

	return registerEmailOutput{id: createdAccount.ID, email: createdAccount.Email, created: true}, err
}
