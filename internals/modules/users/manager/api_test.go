package manager

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	nbcontext "github.com/netbirdio/netbird/management/server/context"
	"github.com/stretchr/testify/require"

	"github.com/netbirdio/management-refactor/internals/modules/users"
	"github.com/netbirdio/management-refactor/internals/shared/permissions"
	"github.com/netbirdio/management-refactor/internals/shared/permissions/modules"
	"github.com/netbirdio/management-refactor/internals/shared/permissions/operations"
)

func TestGetAllUsersReturnsUsers(t *testing.T) {
	ctrl := gomock.NewController(t)
	managerMock := NewMockManager(ctrl)
	permissionsMock := permissions.NewMockManager(ctrl)
	permissionsMock.EXPECT().WithPermission(modules.Users, operations.Read, gomock.Any()).DoAndReturn(
		func(module modules.Module, operation operations.Operation, handlerFunc func(http.ResponseWriter, *http.Request, *nbcontext.UserAuth)) http.HandlerFunc {
			return func(w http.ResponseWriter, r *http.Request) {
				handlerFunc(w, r, &nbcontext.UserAuth{UserId: "123", AccountId: "abc"})
			}
		}).AnyTimes()

	managerMock.EXPECT().GetAllUsers(gomock.Any(), gomock.Any(), gomock.Any(), "abc").Return([]users.User{{Id: "1"}, {Id: "2"}}, nil).Times(1)

	router := mux.NewRouter()
	RegisterEndpoints(router, permissionsMock, managerMock)

	req := httptest.NewRequest(http.MethodGet, "/users", nil)
	req.Header.Set("Authorization", "Bearer token")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	var response []users.User
	err := json.NewDecoder(w.Body).Decode(&response)
	require.NoError(t, err)
	require.Len(t, response, 2)
	require.Equal(t, "1", response[0].Id)
}
