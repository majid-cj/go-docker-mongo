package routers

import (
	"github.com/majid-cj/go-docker-mongo/apps"
	"github.com/majid-cj/go-docker-mongo/domain/entity"
	"github.com/majid-cj/go-docker-mongo/util"

	"github.com/kataras/iris/v12"
)

// MemberRouters ...
type MemberRouters struct {
	memberapp apps.MemberAppInterface
}

// NewMemberRouters ...
func NewMemberRouters(
	memberapp apps.MemberAppInterface,
) *MemberRouters {
	return &MemberRouters{
		memberapp: memberapp,
	}
}

// GetAllMembers ...
func (router *MemberRouters) GetAllMembers(c iris.Context) {
	var members entity.Members
	members, err := router.memberapp.GetMembers()
	if err != nil {
		util.ResponseT(err, iris.StatusNotFound, c)
		return
	}
	util.Response(members.GetMembersSerializer(), iris.StatusOK, c)
}

// GetMembersByType ...
func (router *MemberRouters) GetMembersByType(c iris.Context) {
	var members entity.Members
	members, err := router.memberapp.GetMembersByType(c.Params().GetUint8Default("type", 1))
	if err != nil {
		util.ResponseT(err, iris.StatusNotFound, c)
		return
	}
	util.Response(members.GetMembersSerializer(), iris.StatusOK, c)
}
