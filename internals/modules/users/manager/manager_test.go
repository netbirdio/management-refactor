package manager

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/netbirdio/management-refactor/internals/modules/users"
	"github.com/netbirdio/management-refactor/internals/shared/db"
)

func TestHideUserIssued(t *testing.T) {
	ctrl := gomock.NewController(t)
	repoMock := NewMockRepository(ctrl)
	repoMock.EXPECT().GetUserByID(gomock.Any(), gomock.Any(), "5").Return(&users.User{Id: "5", Issued: "top secret"}, nil).AnyTimes()

	manager := NewManager(repoMock)

	result, err := manager.GetUserByID(context.Background(), nil, db.LockingStrengthShare, "5")
	assert.NoError(t, err)

	assert.Equal(t, "****", result.Issued)

}
