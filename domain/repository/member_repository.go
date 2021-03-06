package repository

import "github.com/majid-cj/go-docker-mongo/domain/entity"

// MemberRepository ...
type MemberRepository interface {
	CreateMember(*entity.Member) (*entity.Member, error)
	DeleteMember(string) error
	GetMembers() ([]entity.Member, error)
	GetMember(string) (*entity.Member, error)
	GetMembersByType(uint8) ([]entity.Member, error)
	GetMemberByEmailAndPassword(*entity.Member) (*entity.Member, error)
	UpdatePassword(*entity.Member) error
}
