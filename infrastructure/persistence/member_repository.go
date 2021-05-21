package persistence

import (
	"context"
	"strings"

	"github.com/majid-cj/go-docker-mongo/domain/entity"
	"github.com/majid-cj/go-docker-mongo/domain/repository"
	"github.com/majid-cj/go-docker-mongo/infrastructure/security"
	"github.com/majid-cj/go-docker-mongo/util"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/bsonx"
	"golang.org/x/crypto/bcrypt"
)

// MemberRepository ...
type MemberRepository struct {
	Ctx context.Context
	DB  *mongo.Collection
}

// NewMemberRepository ...
func NewMemberRepository(db *mongo.Database) *MemberRepository {
	return &MemberRepository{
		Ctx: context.Background(),
		DB:  db.Collection(MEMBER),
	}
}

var _ repository.MemberRepository = &MemberRepository{}

// CreateMember ...
func (m *MemberRepository) CreateMember(member *entity.Member) (*entity.Member, error) {
	if member.Type != 2 {
		return nil, util.GetError("general_error")
	}

	_, err := m.DB.Indexes().CreateOne(
		context.Background(),
		mongo.IndexModel{
			Keys:    bsonx.MDoc{"email": bsonx.Int64(1)},
			Options: options.Index().SetUnique(true),
		},
	)

	if err != nil {
		return nil, util.GetError("general_error")
	}

	_, err = m.DB.InsertOne(m.Ctx, member)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate") {
			return nil, util.GetError("email_taken")
		} else {
			return nil, util.GetError("general_error")
		}
	}
	return member, nil
}

// DeleteMember ...
func (m *MemberRepository) DeleteMember(ID string) error {
	filter := bson.M{"id": ID}
	_, err := m.DB.DeleteOne(m.Ctx, filter)
	if err != nil {
		return util.GetError("general_error")
	}
	return nil
}

// GetMembers ...
func (m *MemberRepository) GetMembers() ([]entity.Member, error) {
	var members entity.Members
	filter := bson.M{}
	cursor, err := m.DB.Find(m.Ctx, filter, nil)
	if err != nil {
		return nil, util.GetError("general_error")
	}

	defer cursor.Close(m.Ctx)

	err = cursor.All(m.Ctx, &members)
	if err != nil {
		return nil, util.GetError("error_retrieve")
	}

	if len(members) == 0 {
		return nil, util.GetError("empty_list")
	}
	return members, nil
}

// GetMember ...
func (m *MemberRepository) GetMember(ID string) (*entity.Member, error) {
	var member entity.Member
	filter := bson.M{"id": ID}
	err := m.DB.FindOne(m.Ctx, filter).Decode(&member)
	if err != nil {
		return nil, util.GetError("member_not_found")
	}
	return &member, nil
}

// GetMembersByType ...
func (m *MemberRepository) GetMembersByType(membertype uint8) ([]entity.Member, error) {
	var members entity.Members

	filter := bson.M{"member_type": membertype}

	cursor, err := m.DB.Find(m.Ctx, filter, nil)
	if err != nil {
		return nil, util.GetError("general_error")
	}

	defer cursor.Close(m.Ctx)

	for cursor.Next(m.Ctx) {
		var member entity.Member
		err := cursor.Decode(member)
		if err != nil {
			return nil, util.GetError("general_error")
		}
		members = append(members, member)
	}

	if len(members) == 0 {
		return nil, util.GetError("empty_list")
	}
	return members, nil
}

// GetMembersBySource ...
func (m *MemberRepository) GetMembersBySource(source uint8) ([]entity.Member, error) {
	var members entity.Members

	filter := bson.M{"source": source}

	cursor, err := m.DB.Find(m.Ctx, filter, nil)
	if err != nil {
		return nil, util.GetError("general_error")
	}

	defer cursor.Close(m.Ctx)

	for cursor.Next(m.Ctx) {
		var member entity.Member
		err := cursor.Decode(member)
		if err != nil {
			return nil, util.GetError("general_error")
		}
		members = append(members, member)
	}

	if len(members) == 0 {
		return nil, util.GetError("empty_list")
	}
	return members, nil
}

// GetMemberByEmailAndPassword ...
func (m *MemberRepository) GetMemberByEmailAndPassword(member *entity.Member) (*entity.Member, error) {
	var getmember entity.Member

	if member.Type != 2 {
		return nil, util.GetError("general_error")
	}

	filter := bson.M{
		"email":       member.Email,
		"member_type": member.Type,
	}

	err := m.DB.FindOne(m.Ctx, filter).Decode(&getmember)

	if err != nil {
		return nil, util.GetError("email_password_wrong")
	}

	err = security.VerifyPassword(getmember.Password, member.Password)

	if err != nil || err == bcrypt.ErrMismatchedHashAndPassword {
		return nil, util.GetError("email_password_wrong")
	}

	return &getmember, nil
}

// GetMemberByEmailAndSource ...
func (m *MemberRepository) GetMemberByEmailAndSource(member *entity.Member) (*entity.Member, uint8, error) {
	var getmember entity.Member

	if member.Type != 2 {
		return nil, 0, util.GetError("general_error")
	}

	filter := bson.M{
		"email":       member.Email,
		"member_type": member.Type,
	}

	err := m.DB.FindOne(m.Ctx, filter).Decode(&getmember)

	if err == nil {
		return &getmember, 0, nil
	}

	member.PrepareSocialMember()
	getnewmember, getError := m.CreateMember(member)

	if getError != nil {
		return &getmember, 0, util.GetError("member_not_found")
	}
	return getnewmember, 1, nil
}

// UpdatePassword ...
func (m *MemberRepository) UpdatePassword(member *entity.Member) error {
	filter := bson.M{"id": member.ID}
	update := bson.M{"$set": member}
	_, err := m.DB.UpdateOne(m.Ctx, filter, update)
	if err != nil {
		return util.GetError("general_error")
	}
	return nil
}
