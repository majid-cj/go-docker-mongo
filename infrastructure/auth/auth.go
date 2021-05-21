package auth

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

// AccessDetail ...
type AccessDetail struct {
	TokenUUID string
	UserID    string
}

// TokenDetail ...
type TokenDetail struct {
	AccessToken        string `json:"access"`
	RefreshToken       string `json:"refresh"`
	TokenUUID          string `json:"token_uuid"`
	RefreshUUID        string `json:"refresh_uuid"`
	AccessTokenExpire  int64  `json:"token_expire"`
	RefreshTokenExpire int64  `json:"refresh_expire"`
}

// AccessData ...
type AccessData struct {
	redisDB *redis.Client
}

// AuthenticationInterface ...
type AuthenticationInterface interface {
	CreatToken(string, *TokenDetail) error
	FetchToken(string) (string, error)
	DeleteAccessToken(*AccessDetail) error
	DeleteRefreshToken(string) error
}

// NewAccessData ...
func NewAccessData(redisDB *redis.Client) *AccessData {
	return &AccessData{redisDB: redisDB}
}

var _ AuthenticationInterface = &AccessData{}

var ctx = context.Background()

// CreatToken ...
func (access *AccessData) CreatToken(userid string, tokenDetail *TokenDetail) error {
	accessexpire := time.Unix(tokenDetail.AccessTokenExpire, 0)
	refreshexpire := time.Unix(tokenDetail.RefreshTokenExpire, 0)
	timenow := time.Now()

	accesscreated, err := access.redisDB.Set(ctx, tokenDetail.TokenUUID, userid, accessexpire.Sub(timenow)).Result()
	if err != nil {
		return nil
	}

	refreshcreated, err := access.redisDB.Set(ctx, tokenDetail.RefreshUUID, userid, refreshexpire.Sub(timenow)).Result()
	if err != nil {
		return nil
	}

	if accesscreated == "0" || refreshcreated == "0" {
		return errors.New("general_error")
	}
	return nil
}

// FetchToken ...
func (access *AccessData) FetchToken(tokenUUID string) (string, error) {
	userid, err := access.redisDB.Get(ctx, tokenUUID).Result()
	if err != nil {
		return "", errors.New("general_error")
	}
	return userid, nil
}

// DeleteAccessToken ...
func (access *AccessData) DeleteAccessToken(auth *AccessDetail) error {
	refreshUUID := fmt.Sprintf("%s++%s", auth.TokenUUID, auth.UserID)
	accesstoken, err := access.redisDB.Del(ctx, auth.TokenUUID).Result()

	if err != nil {
		return errors.New("general_error")
	}

	refreshtoken, err := access.redisDB.Del(ctx, refreshUUID).Result()
	if err != nil {
		return nil
	}
	if accesstoken != 1 || refreshtoken != 1 {
		return errors.New("general_error")
	}
	return nil
}

// DeleteRefreshToken ...
func (access *AccessData) DeleteRefreshToken(refreshUUID string) error {
	refrestoken, err := access.redisDB.Del(ctx, refreshUUID).Result()
	if err != nil || refrestoken != 1 {
		return errors.New("general_error")
	}
	return nil
}
