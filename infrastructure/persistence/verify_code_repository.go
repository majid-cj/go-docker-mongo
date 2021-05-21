package persistence

import (
	"context"

	"github.com/majid-cj/go-docker-mongo/domain/entity"
	"github.com/majid-cj/go-docker-mongo/infrastructure/security"
	"github.com/majid-cj/go-docker-mongo/util"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// VerifyCodeRepository ...
type VerifyCodeRepository struct {
	Ctx      context.Context
	Db       *mongo.Collection
	DbMember *mongo.Collection
}

// NewVerifyCodeRepository ...
func NewVerifyCodeRepository(db *mongo.Database) *VerifyCodeRepository {
	return &VerifyCodeRepository{
		Ctx:      context.Background(),
		Db:       db.Collection(VERIFYCODE),
		DbMember: db.Collection(MEMBER),
	}
}

// CreateVerificationCode ...
func (vcr *VerifyCodeRepository) CreateVerificationCode(code *entity.VerificationCode) (*entity.VerificationCode, error) {
	filter := bson.M{"member": code.Member}
	vcr.Db.DeleteMany(vcr.Ctx, filter)
	_, err := vcr.Db.InsertOne(vcr.Ctx, code)
	if err != nil {
		return nil, err
	}
	return code, nil
}

// CreateVerificationCodeFromEmail ...
func (vcr *VerifyCodeRepository) CreateVerificationCodeFromEmail(code *entity.VerificationCode) (*entity.VerificationCode, error) {
	var member entity.Member
	var verifycode entity.VerificationCode
	err := vcr.DbMember.FindOne(vcr.Ctx, bson.M{"email": code.Email, "member_type": 3, "source": 1}).Decode(&member)
	if err != nil {
		return nil, util.GetError("no_email_account")
	}

	vcr.Db.DeleteMany(vcr.Ctx, bson.M{"member": member.ID})
	verifycode.PrepareVerificationCode(member.ID, code.CodeType)
	_, err = vcr.Db.InsertOne(vcr.Ctx, &verifycode)
	if err != nil {
		return nil, util.GetError("general_error")
	}
	return &verifycode, nil
}

// ResetPassword ...
func (vcr *VerifyCodeRepository) ResetPassword(code *entity.VerificationCode) error {
	var verifycode entity.VerificationCode
	var member entity.Member
	filterMember := bson.M{"email": code.Email, "member_type": 3, "source": 1}
	err := vcr.DbMember.FindOne(vcr.Ctx, filterMember).Decode(&member)
	if err != nil {
		return util.GetError("general_error")
	}

	filterCode := bson.M{"member": member.ID, "code": code.Code, "code_type": code.CodeType, "taken": false}
	err = vcr.Db.FindOne(vcr.Ctx, filterCode).Decode(&verifycode)
	if err != nil {
		return util.GetError("general_error")
	}

	if util.GetTimeNow().Unix() > verifycode.ExpiredAt.Unix() {
		return util.GetError("token_expired")
	}

	hashedpassword, err := security.HashPassword(code.Password)
	if err != nil {
		return err
	}
	_, err = vcr.Db.UpdateMany(vcr.Ctx, filterCode, bson.M{"$set": bson.M{"taken": true}})
	if err != nil {
		return util.GetError("general_error")
	}

	member.Password = string(hashedpassword)
	member.UpdateAt = util.GetTimeNow()
	_, err = vcr.DbMember.UpdateOne(vcr.Ctx, bson.M{"id": member.ID}, bson.M{"$set": member})
	if err != nil {
		return util.GetError("general_error")
	}

	return nil
}

// CheckVerificationCode ...
func (vcr *VerifyCodeRepository) CheckVerificationCode(code *entity.VerificationCode) error {
	var verifycode entity.VerificationCode
	filter := bson.M{"member": code.Member, "code": code.Code, "code_type": code.CodeType}

	err := vcr.Db.FindOne(vcr.Ctx, filter).Decode(&verifycode)
	if err != nil {
		return util.GetError("general_error")
	}

	if util.GetTimeNow().Unix() > verifycode.ExpiredAt.Unix() {
		return util.GetError("token_expired")
	}

	update := bson.M{"$set": bson.M{"taken": true}}
	_, err = vcr.Db.UpdateOne(vcr.Ctx, filter, update)
	if err != nil {
		return util.GetError("general_error")
	}

	_, err = vcr.DbMember.UpdateOne(vcr.Ctx, bson.M{"id": code.Member}, bson.M{"$set": bson.M{"verified": true}})
	if err != nil {
		return util.GetError("general_error")
	}
	return nil
}

// RenewVerificationCode ...
func (vcr *VerifyCodeRepository) RenewVerificationCode(code *entity.VerificationCode) (*entity.VerificationCode, error) {
	var verifyCode entity.VerificationCode
	filter := bson.M{"member": code.Member}
	_, err := vcr.Db.DeleteMany(vcr.Ctx, filter)
	if err != nil {
		return nil, util.GetError("general_error")
	}
	verifyCode.PrepareVerificationCode(code.Member, code.CodeType)
	_, err = vcr.Db.InsertOne(vcr.Ctx, verifyCode)
	if err != nil {
		return nil, util.GetError("general_error")
	}
	return &verifyCode, nil
}
