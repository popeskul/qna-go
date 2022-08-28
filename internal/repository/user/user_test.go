package user

import (
	"github.com/popeskul/qna-go/internal/domain"
	"github.com/popeskul/qna-go/internal/util"
)

//var mockDB *sql.DB
//var mockRepo *RepositoryAuth
//
//func TestMain(m *testing.M) {
//	if err := util.ChangeDir("../../"); err != nil {
//		log.Fatal(err)
//	}
//
//	cfg, err := loadConfig()
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	db, err := newDBConnection(cfg)
//	mockDB = db
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	mockRepo = NewRepoAuth(mockDB)
//
//	os.Exit(m.Run())
//}

//func TestRepositoryAuth_CreateUser(t *testing.T) {
//	ctx := context.Background()
//	u := randomUser()
//
//	if err := mockRepo.CreateUser(ctx, u); err != nil {
//		t.Fatalf("error creating user: %v", err)
//	}
//
//	type fields struct {
//		repo *RepositoryAuth
//	}
//	type args struct {
//		u domain.User
//	}
//	tests := []struct {
//		name    string
//		fields  fields
//		args    args
//		wantErr bool
//	}{
//		{
//			name: "CreateUser",
//			fields: fields{
//				repo: mockRepo,
//			},
//			args: args{
//				u: randomUser(),
//			},
//			wantErr: false,
//		},
//		{
//			name: "CreateUser_DuplicateEmail",
//			fields: fields{
//				repo: mockRepo,
//			},
//			args: args{
//				u: domain.User{
//					Name:     util.RandomString(10),
//					Email:    u.Email,
//					Password: util.RandomString(10),
//				},
//			},
//			wantErr: true,
//		},
//	}
//
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			err := tt.fields.repo.CreateUser(ctx, tt.args.u)
//			if (err != nil) != tt.wantErr {
//				t.Errorf("RepositoryAuth.CreateUser() error = %v, wantErr %v", err, tt.wantErr)
//			}
//
//			t.Cleanup(func() {
//				if err == nil {
//					userID, err := findUserIDByEmail(tt.args.u.Email)
//					if err != nil {
//						t.Error(err)
//					}
//
//					helperDeleteUserByID(t, userID)
//				}
//			})
//		})
//	}
//
//	t.Cleanup(func() {
//		userID, err := findUserIDByEmail(u.Email)
//
//		if err = mockRepo.DeleteUserById(ctx, userID); err != nil {
//			t.Error(err)
//		}
//	})
//}

//func TestRepositoryAuth_GetUser(t *testing.T) {
//	ctx := context.Background()
//	u := randomUser()
//	if err := mockRepo.CreateUser(ctx, u); err != nil {
//		t.Fatalf("error creating user: %v", err)
//	}
//
//	userID, err := findUserIDByEmail(u.Email)
//	if err != nil {
//		t.Fatalf("error finding user: %v", err)
//	}
//
//	type fields struct {
//		repo *RepositoryAuth
//	}
//	type args struct {
//		email    string
//		password string
//	}
//	tests := []struct {
//		name   string
//		fields fields
//		args   args
//		want   int
//	}{
//		{
//			name: "GetUser",
//			fields: fields{
//				repo: mockRepo,
//			},
//			args: args{
//				email:    u.Email,
//				password: u.Password,
//			},
//			want: userID,
//		},
//	}
//
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			got, err := tt.fields.repo.GetUser(ctx, tt.args.email, []byte(tt.args.password))
//			if err != nil {
//				t.Error(err)
//			}
//
//			if got.ID != tt.want {
//				t.Errorf("RepositoryAuth.GetUser() = %v, want %v", got.ID, tt.want)
//			}
//		})
//	}
//
//	t.Cleanup(func() {
//		if err = mockRepo.DeleteUserById(ctx, userID); err != nil {
//			t.Error(err)
//		}
//	})
//}
//
//func TestRepositoryAuth_DeleteUserById(t *testing.T) {
//	ctx := context.Background()
//	user := randomUser()
//
//	if err := mockRepo.CreateUser(ctx, user); err != nil {
//		t.Fatalf("error creating user: %v", err)
//	}
//
//	userID, err := findUserIDByEmail(user.Email)
//	if err != nil {
//		t.Fatalf("error finding user: %v", err)
//	}
//
//	type args struct {
//		repo *RepositoryAuth
//	}
//	type want struct {
//		err error
//	}
//	tests := []struct {
//		name string
//		args args
//		want want
//	}{
//		{
//			name: "Success: DeleteUserById",
//			args: args{
//				repo: mockRepo,
//			},
//			want: want{
//				err: nil,
//			},
//		},
//		{
//			name: "Fail: DeleteUserById",
//			args: args{
//				repo: mockRepo,
//			},
//			want: want{
//				err: ErrDeleteUser,
//			},
//		},
//	}
//
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			if err = tt.args.repo.DeleteUserById(ctx, userID); err != tt.want.err {
//				t.Error(err)
//			}
//		})
//	}
//}

func randomUser() domain.User {
	return domain.User{
		Name:     util.RandomString(10),
		Email:    util.RandomString(10) + "@" + util.RandomString(10) + ".com",
		Password: util.RandomString(10),
	}
}

//func helperDeleteUserByID(t *testing.T, id int) {
//	ctx := context.Background()
//	if err := mockRepo.DeleteUserById(ctx, id); err != nil {
//		t.Error(err)
//	}
//}
//
//func findUserIDByEmail(email string) (int, error) {
//	var userID int
//
//	row := mockDB.QueryRow("SELECT id FROM users WHERE email = $1", email).Scan(&userID)
//	if row != nil {
//		return 0, row
//	}
//
//	return userID, nil
//}

//func newDBConnection(cfg *config.Config) (*sql.DB, error) {
//	return postgres.NewPostgresConnection(db.ConfigDB{
//		Host:     cfg.DB.Host,
//		Port:     cfg.DB.Port,
//		User:     cfg.DB.User,
//		Password: cfg.DB.Password,
//		DBName:   cfg.DB.DBName,
//		SSLMode:  cfg.DB.SSLMode,
//	})
//}
//
//func loadConfig() (*config.Config, error) {
//	err := godotenv.Load(".env")
//	if err != nil {
//		log.Fatalf("Some error occured. Err: %s", err)
//	}
//
//	cfg, err := config.New("configs", "test.config")
//	if err != nil {
//		return nil, err
//	}
//	cfg.DB.Password = os.Getenv("DB_PASSWORD")
//
//	return cfg, nil
//}
