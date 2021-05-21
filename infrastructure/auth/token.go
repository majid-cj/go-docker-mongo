package auth

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/majid-cj/go-docker-mongo/util"

	"github.com/dgrijalva/jwt-go"
)

// TokenInterface ...
type TokenInterface interface {
	CreateJWTToken(string, string) (*TokenDetail, error)
	ExtractJWTTokenMetadata(*http.Request) (*AccessDetail, error)
}

// Token ...
type Token struct{}

var _ TokenInterface = &Token{}

// NewToken ...
func NewToken() *Token {
	return &Token{}
}

// CreateJWTToken ...
func (token *Token) CreateJWTToken(userid, usertype string) (*TokenDetail, error) {
	tokenDetail := &TokenDetail{}
	tokenDetail.AccessTokenExpire = time.Now().Add(time.Hour * 24).Unix()
	tokenDetail.TokenUUID = util.UUID()

	tokenDetail.RefreshTokenExpire = time.Now().Add(time.Hour * 24 * 7).Unix()
	tokenDetail.RefreshUUID = fmt.Sprintf("%s++%s", tokenDetail.TokenUUID, userid)

	var err error
	accessTokenClaim := jwt.MapClaims{}
	accessTokenClaim["authorization"] = true
	accessTokenClaim["access_uuid"] = tokenDetail.TokenUUID
	accessTokenClaim["user_id"] = userid
	accessTokenClaim["user_type"] = usertype
	accessTokenClaim["exp"] = tokenDetail.AccessTokenExpire
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessTokenClaim)

	tokenDetail.AccessToken, err = accessToken.SignedString([]byte(os.Getenv("ACCESS_SECRET")))
	if err != nil {
		return nil, errors.New("general_error")
	}

	refreshTokenClaim := jwt.MapClaims{}
	refreshTokenClaim["refresh_uuid"] = tokenDetail.RefreshUUID
	refreshTokenClaim["user_id"] = userid
	refreshTokenClaim["user_type"] = usertype
	refreshTokenClaim["exp"] = tokenDetail.RefreshTokenExpire
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshTokenClaim)

	tokenDetail.RefreshToken, err = refreshToken.SignedString([]byte(os.Getenv("REFRESH_SECRET")))
	if err != nil {
		return nil, errors.New("general_error")
	}
	return tokenDetail, nil
}

// ExtractJWTTokenMetadata ...
func (token *Token) ExtractJWTTokenMetadata(request *http.Request) (*AccessDetail, error) {
	_token, err := VerifyToken(request)
	if err != nil {
		return nil, errors.New("general_error")
	}
	claims, ok := _token.Claims.(jwt.MapClaims)
	if ok && _token.Valid {
		accessUUID, ok := claims["access_uuid"].(string)
		if !ok {
			return nil, errors.New("general_error")
		}
		userID, ok := claims["user_id"].(string)
		if !ok {
			return nil, errors.New("general_error")
		}
		return &AccessDetail{
			TokenUUID: accessUUID,
			UserID:    userID,
		}, nil
	}
	return nil, errors.New("general_error")
}

// ExtractMemberType ...
func ExtractMemberType(request *http.Request) (string, error) {
	_token, err := VerifyToken(request)
	if err != nil {
		return "", err
	}
	claims, ok := _token.Claims.(jwt.MapClaims)
	if ok && _token.Valid {
		memberType, ok := claims["user_type"].(string)
		if !ok {
			// error getting token claims (member type)
			return "", errors.New("error_parsing_data")
		}
		return memberType, nil
	}
	return "", err
}

// TokenValid ...
func TokenValid(request *http.Request) error {
	token, err := VerifyToken(request)
	if err != nil {
		return err
	}
	if _, ok := token.Claims.(jwt.Claims); !ok && !token.Valid {
		// invalid token
		return errors.New("general_error")
	}

	return nil
}

// VerifyToken ...
func VerifyToken(request *http.Request) (*jwt.Token, error) {
	_token := ExtractToken(request)
	token, err := jwt.Parse(_token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method %v", token.Header["alg"])
		}
		return []byte(os.Getenv("ACCESS_SECRET")), nil
	})
	if err != nil {
		// error parsing token
		return nil, errors.New("general_error")
	}
	return token, nil
}

// ExtractToken ...
func ExtractToken(request *http.Request) string {
	bearer := request.Header.Get("Authorization")
	token := strings.Split(bearer, " ")
	if token[0] != "Bearer" {
		return ""
	}
	if len(token) == 2 {
		return token[1]
	}
	return ""
}
