package persistence

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/majid-cj/go-docker-mongo/domain/repository"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Repository ...
type Repository struct {
	Member     repository.MemberRepository
	VerifyCode repository.VerificationCodeRepository
	Ctx        context.Context
	Client     *mongo.Client
}

// NewRepository ...
func NewRepository() (*Repository, error) {
	URL := fmt.Sprintf("%s://%s:%s/?connect=direct", os.Getenv("DB_DRIVER"), os.Getenv("DB_HOST"), os.Getenv("DB_PORT"))
	clientOption := options.Client().ApplyURI(URL)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)

	client, err := mongo.Connect(ctx, clientOption)

	db := client.Database(os.Getenv("DB_NAME"))

	defer cancel()

	if err != nil {
		return nil, err
	}

	return &Repository{
		Member:     NewMemberRepository(db),
		VerifyCode: NewVerifyCodeRepository(db),
		Ctx:        ctx,
		Client:     client,
	}, nil
}

const (
	// MEMBER ...
	MEMBER = "member"
	// VERIFYCODE ...
	VERIFYCODE = "verify_code"
)
