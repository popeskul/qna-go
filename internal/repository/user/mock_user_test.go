package user

import (
	"context"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/popeskul/qna-go/internal/domain"
	"github.com/popeskul/qna-go/internal/repository/mock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"testing"
)

type RepositoryAuthSuite struct {
	suite.Suite
	*require.Assertions

	ctrl         *gomock.Controller
	mockRepoAuth *mock.MockAuth
}

func TestRepositoryAuthSuiteSuite(t *testing.T) {
	suite.Run(t, new(RepositoryAuthSuite))
}

func (s *RepositoryAuthSuite) SetupTest() {
	s.Assertions = require.New(s.T())

	s.ctrl = gomock.NewController(s.T())
	s.mockRepoAuth = mock.NewMockAuth(s.ctrl)
}

func (s *RepositoryAuthSuite) TearDownTest() {
	s.ctrl.Finish()
}

func (s *RepositoryAuthSuite) TestCreateUser() {
	ctx := context.Background()
	user := randomUser()

	type args struct {
		u   domain.User
		ctx context.Context
	}
	tests := []struct {
		name string
		mock func(ctx context.Context, u domain.User)
		args args
		err  error
	}{
		{
			name: "Success: create user",
			mock: func(ctx context.Context, u domain.User) {
				s.mockRepoAuth.
					EXPECT().
					CreateUser(ctx, u).
					Return(nil)
			},
			args: args{
				u:   user,
				ctx: ctx,
			},
			err: nil,
		},
		{
			name: "Fail: duplicate email",
			mock: func(ctx context.Context, u domain.User) {
				s.mockRepoAuth.
					EXPECT().
					CreateUser(ctx, u).
					Return(ErrCreateUser)
			},
			args: args{
				u:   domain.User{},
				ctx: ctx,
			},
			err: ErrCreateUser,
		},
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			tt.mock(tt.args.ctx, tt.args.u)
			err := s.mockRepoAuth.CreateUser(tt.args.ctx, tt.args.u)
			s.Equal(tt.err, err)
		})
	}
}

func (s *RepositoryAuthSuite) TestGetUser() {
	ctx := context.Background()
	user1 := randomUser()
	user2 := randomUser()

	type args struct {
		ctx  context.Context
		user domain.User
	}
	type want struct {
		err  bool
		user domain.User
	}
	tests := []struct {
		name string
		mock func(ctx context.Context, u domain.User)
		args args
		want want
	}{
		{
			name: "Success: get user",
			args: args{
				ctx:  ctx,
				user: user1,
			},
			mock: func(ctx context.Context, u domain.User) {
				s.mockRepoAuth.
					EXPECT().
					GetUser(ctx, u.Email, []byte(u.Password)).
					Return(u, nil)
			},
			want: want{
				err:  false,
				user: user1,
			},
		},
		{
			name: "Fail: get user",
			args: args{
				ctx:  ctx,
				user: user2,
			},
			mock: func(ctx context.Context, u domain.User) {
				s.mockRepoAuth.
					EXPECT().
					GetUser(ctx, u.Email, []byte(u.Password)).
					Return(domain.User{}, errors.New(""))
			},
			want: want{
				err: true,
			},
		},
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			tt.mock(tt.args.ctx, tt.args.user)
			u, err := s.mockRepoAuth.GetUser(tt.args.ctx, tt.args.user.Email, []byte(tt.args.user.Password))
			s.Equal(tt.want.err, err != nil)
			s.Equal(tt.want.user, u)
		})
	}
}

func (s *RepositoryAuthSuite) TestGetUserByEmail() {
	ctx := context.Background()
	user1 := randomUser()
	user2 := randomUser()

	type args struct {
		ctx  context.Context
		user domain.User
	}
	type want struct {
		err  bool
		user domain.User
	}
	tests := []struct {
		name string
		mock func(ctx context.Context, u domain.User)
		args args
		want want
	}{
		{
			name: "Success: get user",
			args: args{
				ctx:  ctx,
				user: user1,
			},
			mock: func(ctx context.Context, u domain.User) {
				s.mockRepoAuth.
					EXPECT().
					GetUserByEmail(ctx, u.Email).
					Return(u, nil)
			},
			want: want{
				err:  false,
				user: user1,
			},
		},
		{
			name: "Fail: get user",
			args: args{
				ctx:  ctx,
				user: user2,
			},
			mock: func(ctx context.Context, u domain.User) {
				s.mockRepoAuth.
					EXPECT().
					GetUserByEmail(ctx, u.Email).
					Return(domain.User{}, errors.New(""))
			},
			want: want{
				err: true,
			},
		},
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			tt.mock(tt.args.ctx, tt.args.user)
			u, err := s.mockRepoAuth.GetUserByEmail(tt.args.ctx, tt.args.user.Email)
			s.Equal(tt.want.err, err != nil)
			s.Equal(tt.want.user, u)
		})
	}
}

func (s *RepositoryAuthSuite) TestDeleteUserById() {
	ctx := context.Background()
	user := randomUser()

	type args struct {
		ctx  context.Context
		user domain.User
	}
	type want struct {
		err bool
	}
	tests := []struct {
		name string
		mock func(ctx context.Context, userID int)
		args args
		want want
	}{
		{
			name: "Success: delete user",
			args: args{
				ctx:  ctx,
				user: user,
			},
			mock: func(ctx context.Context, userID int) {
				s.mockRepoAuth.
					EXPECT().
					DeleteUserById(ctx, userID).
					Return(nil)
			},
			want: want{
				err: false,
			},
		},
		{
			name: "Fail: delete user",
			args: args{
				ctx:  ctx,
				user: user,
			},
			mock: func(ctx context.Context, userID int) {
				s.mockRepoAuth.
					EXPECT().
					DeleteUserById(ctx, userID).
					Return(errors.New(""))
			},
			want: want{
				err: true,
			},
		},
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			tt.mock(tt.args.ctx, tt.args.user.ID)
			err := s.mockRepoAuth.DeleteUserById(tt.args.ctx, tt.args.user.ID)
			s.Equal(tt.want.err, err != nil)
		})
	}
}
