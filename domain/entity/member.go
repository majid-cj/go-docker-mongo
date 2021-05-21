package entity

import (
	"time"

	"github.com/majid-cj/go-docker-mongo/util"
)

// MemberType ...
var MemberType = map[uint8]string{
	1: "Admin",
	2: "End User",
}

// Member ...
type Member struct {
	ID          string    `bson:"id" json:"id"`
	Type        uint8     `bson:"member_type" json:"member_type"`
	DisplayName string    `bson:"display_name" json:"display_name"`
	Email       string    `bson:"email" json:"email"`
	Password    string    `bson:"password" json:"password"`
	Verified    bool      `bson:"verified" json:"verified"`
	Active      bool      `bson:"active" json:"active"`
	CreatedAt   time.Time `bson:"created_at"`
	UpdateAt    time.Time `bson:"update_at"`
}

// MemberSerializer ...
type MemberSerializer struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	Type      uint8  `json:"member_type"`
	Verified  bool   `json:"verified"`
	Active    bool   `json:"active"`
	CreatedAt string `json:"created_at"`
}

// Members ...
type Members []Member

// PrepareMember ...
func (m *Member) PrepareMember(password []byte) {
	m.ID = util.UUID()
	m.Email = util.EscapeString(m.Email)
	m.Password = string(password)
	m.Active = true
	m.Verified = false
	m.CreatedAt = util.GetTimeNow()
	m.UpdateAt = util.GetTimeNow()
}

// PrepareSocialMember ...
func (m *Member) PrepareSocialMember() {
	m.ID = util.UUID()
	m.Email = util.EscapeString(m.Email)
	m.Active = true
	m.Verified = true
	m.CreatedAt = util.GetTimeNow()
	m.UpdateAt = util.GetTimeNow()
}

// GetMemberSerializer ...
func (m Member) GetMemberSerializer() MemberSerializer {
	return MemberSerializer{
		ID:        m.ID,
		Email:     m.Email,
		Type:      m.Type,
		Verified:  m.Verified,
		Active:    m.Active,
		CreatedAt: m.CreatedAt.String(),
	}
}

// GetMembersSerializer ...
func (members Members) GetMembersSerializer() []interface{} {
	results := make([]interface{}, len(members))
	for index, user := range members {
		results[index] = user.GetMemberSerializer()
	}
	return results
}

// ValidateMember ...
func (member *Member) ValidateMember() error {
	if err := util.ValidateDisplayName(member.DisplayName); err != nil {
		return err
	}
	if err := util.ValidateFormat(member.Email); err != nil {
		return err
	}
	if err := util.ValidatePassword(member.Password); err != nil {
		return err
	}
	return nil
}
