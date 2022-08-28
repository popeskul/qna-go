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
	u := randomUser()

	testCases := []struct {
		name      string
		mock      func(store *mock.MockAuth)
		check     func(t *testing.T, err error)
		wantError error
	}{
		{
			name: "Success: create user",
			mock: func(store *mock.MockAuth) {
				store.EXPECT().CreateUser(ctx, u).Return(nil)
			},
			check: func(t *testing.T, err error) {
				s.NoError(err)
			},
		},
		{
			name: "Error: create user",
			mock: func(store *mock.MockAuth) {
				store.EXPECT().CreateUser(ctx, u).Return(ErrCreateUser)
			},
			check: func(t *testing.T, err error) {
				s.Equal(ErrCreateUser, err)
			},
			wantError: ErrCreateUser,
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			tc.mock(s.mockRepoAuth)
			err := s.mockRepoAuth.CreateUser(ctx, u)
			tc.check(s.T(), err)
		})
	}

	// another way to test

	//type args struct {
	//	ctx  context.Context
	//	user domain.User
	//}
	//type want struct {
	//	err error
	//}
	//
	//tests := []struct {
	//	name string
	//	args args
	//	want want
	//}{
	//	{
	//		name: "Success: create user",
	//		args: args{
	//			ctx:  ctx,
	//			user: u,
	//		},
	//		want: want{
	//			err: nil,
	//		},
	//	},
	//	{
	//		name: "Fail: create user",
	//		args: args{
	//			ctx:  ctx,
	//			user: u,
	//		},
	//		want: want{
	//			err: ErrCreateUser,
	//		},
	//	},
	//}
	//
	//for _, tt := range tests {
	//	s.T().Run(tt.name, func(t *testing.T) {
	//		s.mockRepoAuth.EXPECT().CreateUser(tt.args.ctx, tt.args.user).Return(tt.want.err)
	//		err := s.mockRepoAuth.CreateUser(tt.args.ctx, tt.args.user)
	//		s.Equal(tt.want.err, err)
	//	})
	//}
}

func (s *RepositoryAuthSuite) TestGetUser() {
	ctx := context.Background()
	u := randomUser()
	u2 := randomUser()

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
		mock func()
		args args
		want want
	}{
		{
			name: "Success: get user",
			args: args{
				ctx:  ctx,
				user: u,
			},
			mock: func() {
				s.mockRepoAuth.
					EXPECT().
					GetUser(ctx, u.Email, []byte(u.Password)).
					Return(u, nil)
			},
			want: want{
				err:  false,
				user: u,
			},
		},
		{
			name: "Fail: get user",
			args: args{
				ctx:  ctx,
				user: u2,
			},
			mock: func() {
				s.mockRepoAuth.
					EXPECT().
					GetUser(ctx, u2.Email, []byte(u2.Password)).
					Return(domain.User{}, errors.New("error"))
			},
			want: want{
				err: true,
			},
		},
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			tt.mock()
			got, err := s.mockRepoAuth.GetUser(tt.args.ctx, tt.args.user.Email, []byte(tt.args.user.Password))
			t.Log(got, err)
			s.Equal(tt.want.err, err != nil)
			s.Equal(tt.want.user, got)
		})
	}

	// another way to test

	//tests := []struct {
	//	name    string
	//	mock    func()
	//	wantErr error
	//}{
	//	{
	//		name: "Success: create user",
	//		mock: func() {
	//			s.mockRepoAuth.EXPECT().CreateUser(ctx, u).Return(nil)
	//		},
	//		wantErr: nil,
	//	},
	//	{
	//		name: "Fail: duplicate email",
	//		mock: func() {
	//			s.mockRepoAuth.EXPECT().CreateUser(ctx, u2).Return(ErrCreateUser)
	//		},
	//		wantErr: ErrCreateUser,
	//	},
	//}
	//
	//for _, tt := range tests {
	//	s.T().Run(tt.name, func(t *testing.T) {
	//		tt.mock()
	//
	//		err := s.mockRepoAuth.CreateUser(ctx, u)
	//		s.Equal(tt.wantErr, err)
	//	})
	//}

	//type args struct {
	//	ctx  context.Context
	//	user domain.User
	//}
	//type want struct {
	//	err  error
	//	user domain.User
	//}
	//tests := []struct {
	//	name string
	//	args args
	//	want want
	//}{
	//	{
	//		name: "Success: get user",
	//		args: args{
	//			ctx:  ctx,
	//			user: u,
	//		},
	//		want: want{
	//			err:  nil,
	//			user: u,
	//		},
	//	},
	//	{
	//		name: "Fail: get user",
	//		args: args{
	//			ctx:  ctx,
	//			user: randomUser(),
	//		},
	//		want: want{
	//			err: nil,
	//		},
	//	},
	//}
	//
	//for _, tt := range tests {
	//	s.T().Run(tt.name, func(t *testing.T) {
	//		s.mockRepoAuth.
	//			EXPECT().
	//			GetUser(tt.args.ctx, tt.args.user.Email, []byte(tt.args.user.Password)).
	//			Return(tt.args.user, tt.want.err)
	//
	//		userID, err := s.mockRepoAuth.GetUser(tt.args.ctx, tt.args.user.Email, []byte(tt.args.user.Password))
	//		s.Equal(tt.want.err, err)
	//		s.Equal(tt.want.user, userID)
	//	})
	//}
}
