package service

import (
	"context"
	"time"

	"github.com/chryscloud/go-microkit-plugins/auth"
	"github.com/chryscloud/go-microkit-plugins/crypto"
	jwtModels "github.com/chryscloud/go-microkit-plugins/models/jwt"
	"github.com/dgrijalva/jwt-go"
	"github.com/ipfs/go-datastore"
	lc "github.com/mailio/mailio-nft-server/config"
	"github.com/mailio/mailio-nft-server/model"
	"github.com/mailio/mailio-nft-server/util"
	"github.com/mitchellh/mapstructure"
)

type UserService struct {
	environment *model.Environment
}

func NewUserService(environment *model.Environment) *UserService {
	return &UserService{
		environment: environment,
	}
}

// Login returns model.ErrUnauthorized if user is not found or password is wrong
func (us *UserService) Login(email string, password string) (*model.JwtTokenOutput, error) {
	user, err := us.getUser(email)
	if err != nil {
		return nil, model.ErrUnauthorized
	}
	// check if passwords match
	ok := crypto.CheckPasswordHash(password, user.Password)
	if !ok {
		lc.Log.Warn("unauthorized access attempt for existing user", email)
		return nil, model.ErrUnauthorized
	}
	userClaim := jwtModels.UserClaim{
		ID: user.ID,
	}
	token, err := auth.NewJWTToken([]byte(lc.Conf.JWTToken.SecretKey), jwt.SigningMethodHS256, userClaim)
	if err != nil {
		lc.Log.Warn("failed to create JWT token", err)
		return nil, err
	}
	jwtToken := model.JwtTokenOutput{
		Token: token,
	}
	return &jwtToken, nil
}

// getUser returns a user from datastore if exists
// otherwise model.ErrNotFound is returned
func (us *UserService) getUser(email string) (*model.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), model.DefaultTimeout)
	defer cancel()
	key := util.CreateKey(model.UserTable, email)
	m, err := us.environment.DB.Get(ctx, key)
	if err != nil {
		if err == datastore.ErrNotFound {
			return nil, model.ErrNotFound
		}
		lc.Log.Error("failed to get catalog", err)
		return nil, err
	}
	userMap, err := util.UnmarshalFromBytes(m)
	if err != nil {
		lc.Log.Error("failed to unmarshal catalog", err)
		return nil, err
	}
	var usr model.User
	err = mapstructure.Decode(userMap, &usr)
	return &usr, err
}

// putUser usperts a new user (insert is exists, or update existing)
func (us *UserService) putUser(user *model.User) (*model.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), model.DefaultTimeout)
	defer cancel()
	id := util.GenerateRandomID()
	if user.ID != "" {
		id = user.ID
	}
	user.ID = id
	user.Modified = time.Now().UnixMilli()
	user.Created = time.Now().UnixMilli()

	m, err := util.MarshalToBytes(user)
	err = us.environment.DB.Put(ctx, util.CreateKey(model.UserTable, id), m)
	if err != nil {
		lc.Log.Error("failed to create new catalog", err)
		return nil, err
	}
	return user, nil
}
