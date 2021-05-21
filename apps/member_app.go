package apps

import (
	"github.com/majid-cj/go-docker-mongo/domain/entity"
	"github.com/majid-cj/go-docker-mongo/domain/repository"
)

// MemberApp ...
type MemberApp struct {
	repositorymember repository.MemberRepository
}

// MemberAppInterface ...
type MemberAppInterface interface {
	CreateMember(*entity.Member) (*entity.Member, error)
	DeleteMember(string) error
	GetMembers() ([]entity.Member, error)
	GetMember(string) (*entity.Member, error)
	GetMembersByType(uint8) ([]entity.Member, error)
	GetMemberByEmailAndPassword(*entity.Member) (*entity.Member, error)
	UpdatePassword(*entity.Member) error
}

var _ MemberAppInterface = &MemberApp{}

// CreateMember ...
func (m *MemberApp) CreateMember(member *entity.Member) (*entity.Member, error) {
	return m.repositorymember.CreateMember(member)
}

// DeleteMember ...
func (m *MemberApp) DeleteMember(ID string) error {
	return m.repositorymember.DeleteMember(ID)
}

// GetMembers ...
func (m *MemberApp) GetMembers() ([]entity.Member, error) {
	return m.repositorymember.GetMembers()
}

// GetMember ...
func (m *MemberApp) GetMember(id string) (*entity.Member, error) {
	return m.repositorymember.GetMember(id)
}

// GetMembersByType ...
func (m *MemberApp) GetMembersByType(membertype uint8) ([]entity.Member, error) {
	return m.repositorymember.GetMembersByType(membertype)
}

// GetMemberByEmailAndPassword ...
func (m *MemberApp) GetMemberByEmailAndPassword(member *entity.Member) (*entity.Member, error) {
	return m.repositorymember.GetMemberByEmailAndPassword(member)
}

// UpdatePassword ...
func (m *MemberApp) UpdatePassword(member *entity.Member) error {
	return m.repositorymember.UpdatePassword(member)
}
