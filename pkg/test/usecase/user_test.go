package test

import (
	"bytes"
	"context"
	"fmt"
	"intern-bcc/domain"
	"intern-bcc/internal/usecase"
	"intern-bcc/pkg/infrastucture"
	"intern-bcc/pkg/jwt"
	testpkg "intern-bcc/pkg/test/pkg"
	test "intern-bcc/pkg/test/repository"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"mime/multipart"
	"os"
	"reflect"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

type mockObjects struct {
	mockUserRepo *test.MockIUserRepository
	mockJWT      *testpkg.MockIJwt
	mockSupabase *testpkg.MockISupabase
	mockGomail   *testpkg.MockIGoMail
	mockProduct *test.MockIProductRepository
}

func TestGetUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := test.NewMockIUserRepository(ctrl)

	mockObject := mockObjects{
		mockUserRepo: mockUserRepo,
	}

	type args struct {
		user      domain.Users
		userParam domain.UserParam
	}

	mockUser := domain.Users{
		Id:         uuid.New(),
		Name:       "TestName",
		Email:      "Test@gmail.com",
		Password:   "TestPassword",
		Gender:     "PEREMPUAN",
		PlaceBirth: "TestPlace",
		DateBirth:  "1 Juni 2005",
	}

	mockUserParam := domain.UserParam{
		Id: mockUser.Id,
	}

	mockArgs := args{
		userParam: mockUserParam,
	}

	tests := []struct {
		name       string
		args       args
		beforeTest func(mockObject mockObjects)
		want       domain.Users
		wantErr    bool
	}{
		{
			name: "Success get user",
			args: mockArgs,
			beforeTest: func(mockObject mockObjects) {
				mockObject.mockUserRepo.EXPECT().GetUser(&mockArgs.user, mockArgs.userParam).SetArg(0, mockUser).Return(nil)
			},
			want:    mockUser,
			wantErr: false,
		},
		{
			name: "failed get user",
			args: mockArgs,
			beforeTest: func(mockObject mockObjects) {
				mockObject.mockUserRepo.EXPECT().GetUser(&mockArgs.user, mockArgs.userParam).Return(assert.AnError)
			},
			want:    mockUser,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := usecase.NewUserUsecase(mockUserRepo, nil, nil, nil, nil)

			tt.beforeTest(mockObject)

			result, err := w.GetUser(tt.args.userParam)

			if err != nil && !tt.wantErr {
				t.Errorf("error = %v", err)
				return
			}

			if !reflect.DeepEqual(result, tt.want) && !tt.wantErr {
				t.Errorf("got=%v, want= %v", result, tt.want)
			}
		})
	}
}

func TestGetLikeProducts(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	productId1 := uuid.New()
	productId2 := uuid.New()
	productId3 := uuid.New()

	mockUserRepo := test.NewMockIUserRepository(ctrl)
	mockJwtPkg := testpkg.NewMockIJwt(ctrl)

	mockObject := mockObjects{
		mockUserRepo: mockUserRepo,
		mockJWT:      mockJwtPkg,
	}

	type args struct {
		user domain.Users
		c    *gin.Context
	}

	mockProducts := []domain.ProductResponses{
		{
			Id:    productId1,
			Name:  "TestProduct1",
			Price: 100000,
		},
		{
			Id:    productId2,
			Name:  "TestProduct2",
			Price: 200000,
		},
		{
			Id:    productId3,
			Name:  "TestProduct3",
			Price: 300000,
		},
	}

	mockUser := domain.Users{
		Id:       uuid.New(),
		Name:     "TestName",
		Email:    "Test@gmail.com",
		Password: "TestPassword",
	}

	mockGetProduct := domain.Users{
		Id:       uuid.New(),
		Name:     "TestName",
		Email:    "Test@gmail.com",
		Password: "TestPassword",
		LikeProduct: []domain.Products{
			{
				Id:    productId1,
				Name:  "TestProduct1",
				Price: 100000,
			},
			{
				Id:    productId2,
				Name:  "TestProduct2",
				Price: 200000,
			},
			{
				Id:    productId3,
				Name:  "TestProduct3",
				Price: 300000,
			},
		},
	}

	var mockContext *gin.Context

	mockArgs := args{
		user: mockUser,
		c:    mockContext,
	}

	tests := []struct {
		name       string
		args       args
		beforeTest func(mockObject mockObjects)
		want       []domain.ProductResponses
		wantErr    bool
	}{
		{
			name: "Success get like products",
			args: mockArgs,
			beforeTest: func(mockObject mockObjects) {
				mockObject.mockJWT.EXPECT().GetLoginUser(mockContext).Return(mockUser, nil)
				mockObject.mockUserRepo.EXPECT().GetLikeProducts(&mockUser, mockUser.Id).SetArg(0, mockGetProduct).Return(nil)
			},
			want:    mockProducts,
			wantErr: false,
		},
		{
			name: "Failed get login user",
			args: mockArgs,
			beforeTest: func(mockObject mockObjects) {
				mockObject.mockJWT.EXPECT().GetLoginUser(mockContext).Return(mockUser, nil)
				mockObject.mockUserRepo.EXPECT().GetLikeProducts(&mockUser, mockUser.Id).SetArg(0, mockGetProduct).Return(nil)
			},
			want:    mockProducts,
			wantErr: true,
		},
		{
			name: "Failed get login user",
			args: mockArgs,
			beforeTest: func(mockObject mockObjects) {
				mockObject.mockJWT.EXPECT().GetLoginUser(mockContext).Return(mockUser, nil)
				mockObject.mockUserRepo.EXPECT().GetLikeProducts(&mockUser, mockUser.Id).SetArg(0, mockGetProduct).Return(nil)
			},
			want:    mockProducts,
			wantErr: true,
		},
		{
			name: "Failed get like products",
			args: mockArgs,
			beforeTest: func(mockObject mockObjects) {
				mockObject.mockJWT.EXPECT().GetLoginUser(mockContext).Return(mockUser, nil)
				mockObject.mockUserRepo.EXPECT().GetLikeProducts(&mockUser, mockUser.Id).SetArg(0, mockGetProduct).Return(assert.AnError)
			},
			want:    mockProducts,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := usecase.NewUserUsecase(mockUserRepo, nil, mockJwtPkg, nil, nil)

			tt.beforeTest(mockObject)

			result, err := w.GetLikeProducts(tt.args.c)

			if tt.wantErr {
				fmt.Println(err)
			}

			if err != nil && !tt.wantErr {
				t.Errorf("error = %v", err)
				return
			}

			if !reflect.DeepEqual(result, tt.want) && !tt.wantErr {
				t.Errorf("got=%v, want= %v", result, tt.want)
			}
		})
	}
}

func TestGetOwnProducts(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := test.NewMockIUserRepository(ctrl)
	mockJWTRepo := testpkg.NewMockIJwt(ctrl)

	productId1 := uuid.New()
	productId2 := uuid.New()
	productId3 := uuid.New()

	mockObject := mockObjects{
		mockUserRepo: mockUserRepo,
		mockJWT:      mockJWTRepo,
	}

	type args struct {
		c    *gin.Context
		user domain.Users
	}

	var c *gin.Context

	mockUserLogin := domain.Users{
		Id:       uuid.New(),
		Name:     "TestName",
		Email:    "Test@gmail.com",
		Password: "TestPassword",
	}

	mockArgs := args{
		c:    c,
		user: mockUserLogin,
	}

	mockUserWithProduct := domain.Users{
		Id:       uuid.New(),
		Name:     mockUserLogin.Name,
		Email:    mockUserLogin.Email,
		Password: mockUserLogin.Password,
		Merchant: domain.Merchants{
			MerchantName: "TestMerchantName",
			Products: []domain.Products{
				{
					Id:           productId1,
					Name:         "TestProduct1",
					Price:        100000,
					ProductPhoto: "TestPhoto1",
				},
				{
					Id:           productId2,
					Name:         "TestProduct2",
					Price:        200000,
					ProductPhoto: "TestPhoto2",
				},
				{
					Id:           productId3,
					Name:         "TestProduct3",
					Price:        300000,
					ProductPhoto: "TestPhoto3",
				},
			},
			University: domain.Universities{
				University: "Test University",
			},
		},
	}

	mockProductResponses := []domain.ProductResponses{
		{
			Id:           productId1,
			Name:         "TestProduct1",
			Price:        100000,
			MerchantName: mockUserWithProduct.Merchant.MerchantName,
			University:   mockUserWithProduct.Merchant.University.University,
			ProductPhoto: "TestPhoto1",
		},
		{
			Id:           productId2,
			Name:         "TestProduct2",
			Price:        200000,
			MerchantName: mockUserWithProduct.Merchant.MerchantName,
			University:   mockUserWithProduct.Merchant.University.University,
			ProductPhoto: "TestPhoto2",
		},
		{
			Id:           productId3,
			Name:         "TestProduct3",
			Price:        300000,
			MerchantName: mockUserWithProduct.Merchant.MerchantName,
			University:   mockUserWithProduct.Merchant.University.University,
			ProductPhoto: "TestPhoto3",
		},
	}

	tests := []struct {
		name       string
		args       args
		beforeTest func(mockObject mockObjects)
		want       []domain.ProductResponses
		wantErr    bool
	}{
		{
			name: "Success get own products",
			args: mockArgs,
			beforeTest: func(mockObject mockObjects) {
				mockObject.mockJWT.EXPECT().GetLoginUser(mockArgs.c).Return(mockUserLogin, nil)
				mockObject.mockUserRepo.EXPECT().GetOwnProducts(&mockArgs.user, mockArgs.user.Id).SetArg(0, mockUserWithProduct).Return(nil)
			},
			want:    mockProductResponses,
			wantErr: false,
		},
		{
			name: "Failed get own products",
			args: mockArgs,
			beforeTest: func(mockObject mockObjects) {
				mockObject.mockJWT.EXPECT().GetLoginUser(mockArgs.c).Return(mockUserLogin, nil)
				mockObject.mockUserRepo.EXPECT().GetOwnProducts(&mockUserLogin, mockArgs.user.Id).Return(assert.AnError)
			},
			want:    mockProductResponses,
			wantErr: true,
		},
		{
			name: "Failed get login user",
			args: mockArgs,
			beforeTest: func(mockObject mockObjects) {
				mockObject.mockJWT.EXPECT().GetLoginUser(mockArgs.c).Return(mockUserLogin, assert.AnError)
			},
			want:    mockProductResponses,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := usecase.NewUserUsecase(mockUserRepo, nil, mockJWTRepo, nil, nil)

			tt.beforeTest(mockObject)

			result, err := w.GetOwnProducts(tt.args.c)

			if err != nil && !tt.wantErr {
				t.Errorf("error = %v", err)
				return
			}

			if !reflect.DeepEqual(result, tt.want) && !tt.wantErr {
				t.Errorf("got=%v, want= %v", result, tt.want)
			}
		})
	}
}

// GetOwnMentors(c *gin.Context) ([]domain.OwnMentorResponses, error)
func TestGetOwnMentors(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mentorId1 := uuid.New()
	mentorId2 := uuid.New()
	mentorId3 := uuid.New()

	mockUserRepo := test.NewMockIUserRepository(ctrl)
	mockJwtPkg := testpkg.NewMockIJwt(ctrl)

	mockObject := mockObjects{
		mockUserRepo: mockUserRepo,
		mockJWT:      mockJwtPkg,
	}

	var mockContext *gin.Context

	mockUserLogin := domain.Users{
		Id:       uuid.New(),
		Name:     "TestName",
		Email:    "Test@gmail.com",
		Password: "TestPassword",
	}

	type args struct {
		c    *gin.Context
		user domain.Users
	}

	mockArgs := args{
		c:    mockContext,
		user: mockUserLogin,
	}

	mockUserWithMentor := domain.Users{
		Id:       uuid.New(),
		Name:     "TestName",
		Email:    "Test@gmail.com",
		Password: "TestPassword",
		HasMentors: []domain.Mentors{
			{
				Id:            mentorId1,
				Name:          "TestMentor1",
				CurrentJob:    "TestCurrentJob1",
				MentorPicture: "TestMentorPicture1",
			},
			{
				Id:            mentorId2,
				Name:          "TestMentor2",
				CurrentJob:    "TestCurrentJob2",
				MentorPicture: "TestMentorPicture2",
			},
			{
				Id:            mentorId3,
				Name:          "TestMentor3",
				CurrentJob:    "TestCurrentJob3",
				MentorPicture: "TestMentorPicture3",
			},
		},
	}

	mockOwnMentorResponse := []domain.OwnMentorResponses{
		{
			Id:            mentorId1,
			Name:          "TestMentor1",
			CurrentJob:    "TestCurrentJob1",
			MentorPicture: "TestMentorPicture1",
		},
		{
			Id:            mentorId2,
			Name:          "TestMentor2",
			CurrentJob:    "TestCurrentJob2",
			MentorPicture: "TestMentorPicture2",
		},
		{
			Id:            mentorId3,
			Name:          "TestMentor3",
			CurrentJob:    "TestCurrentJob3",
			MentorPicture: "TestMentorPicture3",
		},
	}

	tests := []struct {
		name       string
		args       args
		beforeTest func(mockObject mockObjects)
		want       []domain.OwnMentorResponses
		wantErr    bool
	}{
		{
			name: "Success get own mentors",
			args: mockArgs,
			beforeTest: func(mockObject mockObjects) {
				mockObject.mockJWT.EXPECT().GetLoginUser(mockArgs.c).Return(mockUserLogin, nil)
				mockObject.mockUserRepo.EXPECT().GetOwnMentors(&mockArgs.user, mockArgs.user.Id).SetArg(0, mockUserWithMentor).Return(nil)
			},
			want:    mockOwnMentorResponse,
			wantErr: false,
		},
		{
			name: "failed to get login user",
			args: mockArgs,
			beforeTest: func(mockObject mockObjects) {
				mockObject.mockJWT.EXPECT().GetLoginUser(mockArgs.c).Return(mockUserLogin, assert.AnError)
			},
			want:    mockOwnMentorResponse,
			wantErr: true,
		},
		{
			name: "failed to get own mentor",
			args: mockArgs,
			beforeTest: func(mockObject mockObjects) {
				mockObject.mockJWT.EXPECT().GetLoginUser(mockArgs.c).Return(mockUserLogin, nil)
				mockObject.mockUserRepo.EXPECT().GetOwnMentors(&mockArgs.user, mockArgs.user.Id).SetArg(0, mockUserWithMentor).Return(assert.AnError)
			},
			want:    mockOwnMentorResponse,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := usecase.NewUserUsecase(mockObject.mockUserRepo, nil, mockObject.mockJWT, nil, nil)

			tt.beforeTest(mockObject)

			result, err := w.GetOwnMentors(tt.args.c)

			if err != nil && !tt.wantErr {
				t.Errorf("error = %v", err)
				return
			}

			if !reflect.DeepEqual(result, tt.want) && !tt.wantErr {
				t.Errorf("got = %v\nwant = %v", result, tt.want)
			}
		})
	}
}

//	Register(userRequest domain.UserRequest) error
//
// unit test failed because uuid
func TestRegister(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := test.NewMockIUserRepository(ctrl)

	mockObject := mockObjects{
		mockUserRepo: mockUserRepo,
	}

	type args struct {
		userRequest domain.UserRequest
		newUser     domain.Users
	}

	mockUserRequest := domain.UserRequest{
		Name:     "TestName",
		Email:    "testemail@gmail.com",
		Password: "TestPassword",
	}

	mockNewUser := domain.Users{
		Name:     "TestName",
		Email:    "testemail@gmail.com",
		Password: "TestPassword",
	}

	mockArgs := args{
		userRequest: mockUserRequest,
		newUser:     mockNewUser,
	}

	tests := []struct {
		name       string
		args       args
		beforeTest func(mockObject mockObjects)
		wantErr    bool
	}{
		{
			name: "Test success register",
			args: mockArgs,
			beforeTest: func(mockObject mockObjects) {
				mockObject.mockUserRepo.EXPECT().Register(&mockNewUser).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "Failed register",
			args: mockArgs,
			beforeTest: func(mockObject mockObjects) {
				mockObject.mockUserRepo.EXPECT().Register(&mockNewUser).Return(assert.AnError)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := usecase.NewUserUsecase(mockObject.mockUserRepo, nil, nil, nil, nil)

			tt.beforeTest(mockObject)

			err := w.Register(tt.args.userRequest)

			if err != nil && !tt.wantErr {
				t.Errorf("error = %v", err)
			}
		})
	}
}

// Login(userLogin domain.UserLogin) (domain.LoginResponse, error)
func TestLogin(t *testing.T) {
	infrastucture.LoadEnv()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := test.NewMockIUserRepository(ctrl)
	mockJwtPkg := testpkg.NewMockIJwt(ctrl)

	mockObject := mockObjects{
		mockUserRepo: mockUserRepo,
		mockJWT:      mockJwtPkg,
	}

	type args struct {
		userLogin domain.UserLogin
		user      domain.Users
		userParam domain.UserParam
	}

	mockUserLogin := domain.UserLogin{
		Email:    "test@email.com",
		Password: "testPassword",
	}

	hashPassword, _ := bcrypt.GenerateFromPassword([]byte(mockUserLogin.Password), 10)

	mockUserParam := domain.UserParam{
		Email: "test@email.com",
	}

	mockUser := domain.Users{
		Id:       uuid.New(),
		Email:    "test@email.com",
		Password: string(hashPassword),
	}

	mockArgs := args{
		userLogin: mockUserLogin,
		userParam: mockUserParam,
	}

	mockJwtToken, _ := jwt.JwtInit().GenerateToken(mockUser.Id)

	mockLoginResponse := domain.LoginResponse{
		JWT: mockJwtToken,
	}

	tests := []struct {
		name       string
		args       args
		beforeTest func(mockObject mockObjects)
		want       domain.LoginResponse
		wantErr    bool
	}{
		{
			name: "Success login",
			args: mockArgs,
			beforeTest: func(mockObject mockObjects) {
				mockObject.mockUserRepo.EXPECT().GetUser(&mockArgs.user, mockArgs.userParam).SetArg(0, mockUser).Return(nil)
				mockObject.mockJWT.EXPECT().GenerateToken(mockUser.Id).Return(mockJwtToken, nil)
			},
			want:    mockLoginResponse,
			wantErr: false,
		},
		{
			name: "Failed get user",
			args: mockArgs,
			beforeTest: func(mockObject mockObjects) {
				mockObject.mockUserRepo.EXPECT().GetUser(&mockArgs.user, mockArgs.userParam).SetArg(0, mockUser).Return(assert.AnError)
			},
			want:    mockLoginResponse,
			wantErr: true,
		},
		{
			name: "Failed to generate token",
			args: mockArgs,
			beforeTest: func(mockObject mockObjects) {
				mockObject.mockUserRepo.EXPECT().GetUser(&mockArgs.user, mockArgs.userParam).SetArg(0, mockUser).Return(nil)
				mockObject.mockJWT.EXPECT().GenerateToken(mockUser.Id).Return(mockJwtToken, assert.AnError)
			},
			want:    mockLoginResponse,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := usecase.NewUserUsecase(mockObject.mockUserRepo, nil, mockObject.mockJWT, nil, nil)

			tt.beforeTest(mockObject)

			result, err := w.Login(tt.args.userLogin)

			if err != nil && !tt.wantErr {
				t.Errorf("error = %v", err)
				return
			}

			if !reflect.DeepEqual(result, tt.want) && !tt.wantErr {
				t.Errorf("got = %v\nwant = %v", result, tt.want)
			}
		})
	}
}

// UpdateUser(c *gin.Context, userId uuid.UUID, userUpdate domain.UserUpdate) (domain.UserResponse, error)
func TestUpdateUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := test.NewMockIUserRepository(ctrl)
	mockJwtPkg := testpkg.NewMockIJwt(ctrl)

	mockObject := mockObjects{
		mockUserRepo: mockUserRepo,
		mockJWT:      mockJwtPkg,
	}

	type args struct {
		c           *gin.Context
		updateUser  domain.UserUpdate
		userLogin   domain.Users
		updatedUser domain.Users
	}

	var mockContext *gin.Context

	mockUpdateUser := domain.UserUpdate{
		Name:           "testupdatename",
		Gender:         "PEREMPUAN",
		PlaceBirth:     "testPlace",
		DateBirth:      "1 juni 2005",
		ProfilePicture: "TestPicture",
	}

	mockUserLogin := domain.Users{
		Id:    uuid.New(),
		Name:  "testname",
		Email: "test@email.com",
	}

	mockUpdatedUser := domain.Users{
		Id:             mockUserLogin.Id,
		Name:           "testupdatename",
		Email:          "test@email.com",
		Gender:         "PEREMPUAN",
		PlaceBirth:     "testPlace",
		DateBirth:      "1 juni 2005",
		ProfilePicture: "TestPicture",
	}

	mockArgs := args{
		c:           mockContext,
		updatedUser: mockUpdatedUser,
		userLogin:   mockUserLogin,
		updateUser:  mockUpdateUser,
	}

	mockUpdatedUserResponses := domain.UserResponse{
		Id:             mockUserLogin.Id,
		Name:           "testupdatename",
		Email:          "test@email.com",
		Gender:         "PEREMPUAN",
		PlaceBirth:     "testPlace",
		DateBirth:      "1 juni 2005",
		ProfilePicture: "TestPicture",
	}

	tests := []struct {
		name       string
		args       args
		beforeTest func(mockObject mockObjects)
		want       domain.UserResponse
		wantErr    bool
	}{
		{
			name: "Success update user",
			args: mockArgs,
			beforeTest: func(mockObject mockObjects) {
				mockObject.mockJWT.EXPECT().GetLoginUser(mockArgs.c).Return(mockArgs.userLogin, nil)
				mockObject.mockUserRepo.EXPECT().UpdateUser(&mockArgs.updateUser, mockArgs.userLogin.Id).Return(nil)
				mockObject.mockUserRepo.EXPECT().GetUser(&domain.Users{}, domain.UserParam{Id: mockUserLogin.Id}).SetArg(0, mockArgs.updatedUser).Return(nil)
			},
			want:    mockUpdatedUserResponses,
			wantErr: false,
		},
		{
			name: "Failed to get login user",
			args: mockArgs,
			beforeTest: func(mockObject mockObjects) {
				mockObject.mockJWT.EXPECT().GetLoginUser(mockArgs.c).Return(mockArgs.userLogin, assert.AnError)
			},
			want:    mockUpdatedUserResponses,
			wantErr: true,
		},
		{
			name: "Failed update userr",
			args: mockArgs,
			beforeTest: func(mockObject mockObjects) {
				mockObject.mockJWT.EXPECT().GetLoginUser(mockArgs.c).Return(mockArgs.userLogin, nil)
				mockObject.mockUserRepo.EXPECT().UpdateUser(&mockArgs.updateUser, mockArgs.userLogin.Id).Return(assert.AnError)
			},
			want:    mockUpdatedUserResponses,
			wantErr: true,
		},
		{
			name: "Failed get user",
			args: mockArgs,
			beforeTest: func(mockObject mockObjects) {
				mockObject.mockJWT.EXPECT().GetLoginUser(mockArgs.c).Return(mockArgs.userLogin, nil)
				mockObject.mockUserRepo.EXPECT().UpdateUser(&mockArgs.updateUser, mockArgs.userLogin.Id).Return(nil)
				mockObject.mockUserRepo.EXPECT().GetUser(&domain.Users{}, domain.UserParam{Id: mockUserLogin.Id}).SetArg(0, mockArgs.updatedUser).Return(assert.AnError)
			},
			want:    mockUpdatedUserResponses,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := usecase.NewUserUsecase(mockObject.mockUserRepo, nil, mockObject.mockJWT, nil, nil)

			tt.beforeTest(mockObject)

			result, err := w.UpdateUser(tt.args.c, tt.args.userLogin.Id, tt.args.updateUser)
			if err != nil && !tt.wantErr {
				t.Errorf("error = %v", err)
				return
			}

			if !reflect.DeepEqual(result, tt.want) && !tt.wantErr {
				t.Errorf("got = %v\nwant = %v", result, tt.want)
			}
		})
	}
}

// UploadUserPhoto(c *gin.Context, userId uuid.UUID, userPhoto *multipart.FileHeader) (domain.UserResponse, error)
func TestUploadUserPhoto(t *testing.T) {
	infrastucture.LoadEnv()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := test.NewMockIUserRepository(ctrl)
	mockJwtPkg := testpkg.NewMockIJwt(ctrl)
	mockSupabasePkg := testpkg.NewMockISupabase(ctrl)

	mockUserId := uuid.New()

	mockObject := mockObjects{
		mockUserRepo: mockUserRepo,
		mockJWT:      mockJwtPkg,
		mockSupabase: mockSupabasePkg,
	}

	type args struct {
		c                  *gin.Context
		userId             uuid.UUID
		userPhoto          *multipart.FileHeader
		userLogin          domain.Users
		userLoginNoProfile domain.Users
	}

	var mockContext *gin.Context

	// Create a temporary file
	tempFile, err := ioutil.TempFile(os.TempDir(), "prefix")
	if err != nil {
		log.Fatal("Cannot create temporary file", err)
	}

	// Remember to clean up the file afterwards
	defer os.Remove(tempFile.Name())

	// Write some text line to the file
	text := []byte("This is a test file.")
	if _, err = tempFile.Write(text); err != nil {
		log.Fatal("Failed to write to temporary file", err)
	}

	// Close the file
	if err := tempFile.Close(); err != nil {
		log.Fatal(err)
	}

	// Create a buffer to write our multipart form to
	body := &bytes.Buffer{}

	// Create a multipart writer
	writer := multipart.NewWriter(body)

	// Create a new form file
	fileWriter, err := writer.CreateFormFile("fileField", tempFile.Name())
	if err != nil {
		log.Fatal(err)
	}

	// Open the file
	file, err := os.Open(tempFile.Name())
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// Write the file to the form
	_, err = io.Copy(fileWriter, file)
	if err != nil {
		log.Fatal(err)
	}

	// Close the multipart writer
	// This step is important because it writes the ending boundary
	err = writer.Close()
	if err != nil {
		log.Fatal(err)
	}

	// Create a new multipart reader
	reader := multipart.NewReader(body, writer.Boundary())

	// Read the multipart form from the reader
	form, err := reader.ReadForm(10 << 20) // 10 MB
	if err != nil {
		log.Fatal(err)
	}

	// Now you can get the *multipart.FileHeader from the form
	mockUserPhoto := form.File["fileField"][0]

	mockNewProfilePicture := "NewProfilePicture"

	loginUser := domain.Users{
		Id:             mockUserId,
		Name:           "TestName",
		Email:          "email@email.com",
		ProfilePicture: mockNewProfilePicture,
	}

	loginUserNoProfile := domain.Users{
		Id:    mockUserId,
		Name:  "TestName",
		Email: "email@email.com",
	}

	mockArgs := args{
		c:                  mockContext,
		userId:             mockUserId,
		userPhoto:          mockUserPhoto,
		userLogin:          loginUser,
		userLoginNoProfile: loginUserNoProfile,
	}

	updatedUser := domain.Users{
		Id:             mockUserId,
		Name:           "TestName",
		Email:          "email@email.com",
		ProfilePicture: mockNewProfilePicture,
	}

	updatedUserResponse := domain.UserResponse{
		Id:             mockUserId,
		Name:           "TestName",
		Email:          "email@email.com",
		ProfilePicture: mockNewProfilePicture,
	}

	tests := []struct {
		name       string
		args       args
		beforeTest func(mockObject mockObjects)
		want       domain.UserResponse
		wantErr    bool
	}{
		{
			name: "Success upload photo with delete",
			args: mockArgs,
			beforeTest: func(mockObject mockObjects) {
				mockObject.mockJWT.EXPECT().GetLoginUser(mockArgs.c).Return(mockArgs.userLogin, nil)
				mockObject.mockSupabase.EXPECT().Delete(mockArgs.userLogin.ProfilePicture).Return(nil)
				mockObject.mockSupabase.EXPECT().Upload(mockArgs.userPhoto).Return(mockNewProfilePicture, nil)
				mockObject.mockUserRepo.EXPECT().UpdateUser(&domain.UserUpdate{ProfilePicture: mockNewProfilePicture}, mockArgs.userLogin.Id).Return(nil)
				mockObject.mockUserRepo.EXPECT().GetUser(&domain.Users{}, domain.UserParam{Id: mockArgs.userId}).SetArg(0, updatedUser).Return(nil)
			},
			want:    updatedUserResponse,
			wantErr: false,
		},
		{
			name: "Failed get login user",
			args: mockArgs,
			beforeTest: func(mockObject mockObjects) {
				mockObject.mockJWT.EXPECT().GetLoginUser(mockArgs.c).Return(mockArgs.userLogin, assert.AnError)
			},
			want:    updatedUserResponse,
			wantErr: true,
		},
		{
			name: "Failed Delete Photo",
			args: mockArgs,
			beforeTest: func(mockObject mockObjects) {
				mockObject.mockJWT.EXPECT().GetLoginUser(mockArgs.c).Return(mockArgs.userLogin, nil)
				mockObject.mockSupabase.EXPECT().Delete(mockArgs.userLogin.ProfilePicture).Return(assert.AnError)
			},
			want:    updatedUserResponse,
			wantErr: true,
		},
		{
			name: "Failed upload photo",
			args: mockArgs,
			beforeTest: func(mockObject mockObjects) {
				mockObject.mockJWT.EXPECT().GetLoginUser(mockArgs.c).Return(mockArgs.userLogin, nil)
				mockObject.mockSupabase.EXPECT().Delete(mockArgs.userLogin.ProfilePicture).Return(nil)
				mockObject.mockSupabase.EXPECT().Upload(mockArgs.userPhoto).Return(mockNewProfilePicture, assert.AnError)
			},
			want:    updatedUserResponse,
			wantErr: true,
		},
		{
			name: "Failed update user",
			args: mockArgs,
			beforeTest: func(mockObject mockObjects) {
				mockObject.mockJWT.EXPECT().GetLoginUser(mockArgs.c).Return(mockArgs.userLogin, nil)
				mockObject.mockSupabase.EXPECT().Delete(mockArgs.userLogin.ProfilePicture).Return(nil)
				mockObject.mockSupabase.EXPECT().Upload(mockArgs.userPhoto).Return(mockNewProfilePicture, nil)
				mockObject.mockUserRepo.EXPECT().UpdateUser(&domain.UserUpdate{ProfilePicture: mockNewProfilePicture}, mockArgs.userLogin.Id).Return(assert.AnError)
			},
			want:    updatedUserResponse,
			wantErr: true,
		},
		{
			name: "Failed get user",
			args: mockArgs,
			beforeTest: func(mockObject mockObjects) {
				mockObject.mockJWT.EXPECT().GetLoginUser(mockArgs.c).Return(mockArgs.userLogin, nil)
				mockObject.mockSupabase.EXPECT().Delete(mockArgs.userLogin.ProfilePicture).Return(nil)
				mockObject.mockSupabase.EXPECT().Upload(mockArgs.userPhoto).Return(mockNewProfilePicture, nil)
				mockObject.mockUserRepo.EXPECT().UpdateUser(&domain.UserUpdate{ProfilePicture: mockNewProfilePicture}, mockArgs.userLogin.Id).Return(nil)
				mockObject.mockUserRepo.EXPECT().GetUser(&domain.Users{}, domain.UserParam{Id: mockArgs.userId}).SetArg(0, updatedUser).Return(assert.AnError)
			},
			want:    updatedUserResponse,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := usecase.NewUserUsecase(mockObject.mockUserRepo, nil, mockObject.mockJWT, mockObject.mockSupabase, nil)

			tt.beforeTest(mockObject)

			result, err := w.UploadUserPhoto(tt.args.c, tt.args.userId, tt.args.userPhoto)
			if err != nil && !tt.wantErr {
				t.Errorf("error = %v", err)
				return
			}

			if !reflect.DeepEqual(result, tt.want) && !tt.wantErr {
				t.Errorf("got = %v\nwant = %v", result, tt.want)
			}
		})
	}
}

// PasswordRecovery(userParam domain.UserParam, ctx context.Context) error
// Gagal karena function newEmailVerificationPassword
func TestPasswordRecovery(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := test.NewMockIUserRepository(ctrl)
	mockGomail := testpkg.NewMockIGoMail(ctrl)

	mockObject := mockObjects{
		mockUserRepo: mockUserRepo,
		mockGomail:   mockGomail,
	}

	type args struct {
		user         domain.Users
		userParam    domain.UserParam
		ctx          context.Context
		subject      string
		htmlBody     string
		emailVerHash string
	}

	mockUser := domain.Users{
		Id:       uuid.New(),
		Name:     "testname",
		Email:    "test@email.com",
		Password: "testpassword",
	}

	mockUserParam := domain.UserParam{
		Email: "test@email.com",
	}

	mockSubject := "test recovery account"

	mockHtmlBody := "test html body"

	//Implementation of newEmailVerification function on usercase/user.go
	var alphaNumRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")
	emailVerRandRune := make([]rune, 64)
	for i := 0; i < 64; i++ {
		emailVerRandRune[i] = alphaNumRunes[rand.Intn(len(alphaNumRunes)-1)]
	}
	emailVerPassword := string(emailVerRandRune)
	var emailVerPWhash []byte

	emailVerPWhash, err := bcrypt.GenerateFromPassword([]byte(emailVerPassword), bcrypt.DefaultCost)
	if err != nil {
		panic("error hashing")
	}
	mockEmailVerHash := string(emailVerPWhash)

	var mockContext context.Context

	mockArgs := args{
		user:         mockUser,
		userParam:    mockUserParam,
		ctx:          mockContext,
		subject:      mockSubject,
		htmlBody:     mockHtmlBody,
		emailVerHash: mockEmailVerHash,
	}

	tests := []struct {
		name       string
		args       args
		beforeTest func(mockObject mockObjects)
		wantErr    bool
	}{
		{
			name: "Success send password recovery",
			args: mockArgs,
			beforeTest: func(mockObject mockObjects) {
				mockObject.mockUserRepo.EXPECT().GetUser(&domain.Users{}, mockArgs.userParam).SetArg(0, mockArgs.user).Return(nil)
				mockObject.mockUserRepo.EXPECT().CreatePasswordVerification(mockArgs.ctx, mockArgs.emailVerHash, mockArgs.user.Name).Return(nil)
				mockObject.mockGomail.EXPECT().SendGoMail(mockArgs.subject, mockArgs.htmlBody, mockArgs.user.Email).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "Success send password recovery",
			args: mockArgs,
			beforeTest: func(mockObject mockObjects) {
				mockObject.mockUserRepo.EXPECT().GetUser(&domain.Users{}, mockArgs.userParam).SetArg(0, mockArgs.user).Return(assert.AnError)
			},
			wantErr: true,
		},
		{
			name: "Success send password recovery",
			args: mockArgs,
			beforeTest: func(mockObject mockObjects) {
				mockObject.mockUserRepo.EXPECT().GetUser(&domain.Users{}, mockArgs.userParam).SetArg(0, mockArgs.user).Return(nil)
				mockObject.mockUserRepo.EXPECT().CreatePasswordVerification(mockArgs.ctx, mockArgs.emailVerHash, mockArgs.user.Name).Return(assert.AnError)
			},
			wantErr: true,
		},
		{
			name: "Success send password recovery",
			args: mockArgs,
			beforeTest: func(mockObject mockObjects) {
				mockObject.mockUserRepo.EXPECT().GetUser(&domain.Users{}, mockArgs.userParam).SetArg(0, mockArgs.user).Return(nil)
				mockObject.mockUserRepo.EXPECT().CreatePasswordVerification(mockArgs.ctx, mockArgs.emailVerHash, mockArgs.user.Name).Return(nil)
				mockObject.mockGomail.EXPECT().SendGoMail(mockArgs.subject, mockArgs.htmlBody, mockArgs.user.Email).Return(assert.AnError)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := usecase.NewUserUsecase(mockObject.mockUserRepo, nil, nil, nil, mockObject.mockGomail)

			tt.beforeTest(mockObject)

			err := w.PasswordRecovery(mockArgs.userParam, mockArgs.ctx)

			if err != nil && !tt.wantErr {
				t.Errorf("error = %v", err)
			}
		})
	}
}

// ChangePassword(ctx context.Context, name string, verPass string, passwordRequest domain.PasswordUpdate) error
//Failed because bycrpt
func TestChangePassword(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := test.NewMockIUserRepository(ctrl)

	mockObject := mockObjects{
		mockUserRepo: mockUserRepo,
	}

	//Implementation of newEmailVerification function on usercase/user.go
	var alphaNumRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")
	emailVerRandRune := make([]rune, 64)
	for i := 0; i < 64; i++ {
		emailVerRandRune[i] = alphaNumRunes[rand.Intn(len(alphaNumRunes)-1)]
	}
	mockEmailVerPassword := string(emailVerRandRune)
	var emailVerPWhash []byte

	emailVerPWhash, err := bcrypt.GenerateFromPassword([]byte(mockEmailVerPassword), bcrypt.DefaultCost)
	if err != nil {
		panic("error hashing")
	}
	mockEmailVerHash := string(emailVerPWhash)

	type args struct {
		ctx             context.Context
		name            string
		verPass         string
		verPassHash     string
		passwordRequest domain.PasswordUpdate
		user            domain.Users
	}

	var mockCtx context.Context

	mockName := "testname"

	mockPasswordRequest := domain.PasswordUpdate{
		Password:        "testpassword",
		ConfirmPassword: "testpassword",
	}

	hashPassword, _ := bcrypt.GenerateFromPassword([]byte(mockPasswordRequest.Password), 10)

	mockUser := domain.Users{
		Id:       uuid.New(),
		Name:     "testname",
		Email:    "test@email.com",
		Password: "testpassword",
	}

	mockArgs := args{
		ctx:             mockCtx,
		name:            mockName,
		verPass:         mockEmailVerPassword,
		verPassHash:     mockEmailVerHash,
		passwordRequest: mockPasswordRequest,
		user:            mockUser,
	}

	tests := []struct {
		name       string
		args       args
		beforeTest func(mockObject mockObjects)
		wantErr    bool
	}{
		{
			name: "Failed get password verification",
			args: mockArgs,
			beforeTest: func(mockObject mockObjects) {
				mockObject.mockUserRepo.EXPECT().GetPasswordVerification(mockArgs.ctx, mockArgs.name).Return(mockArgs.verPassHash, assert.AnError)
			},
			wantErr: true,
		},
		{
			name: "Failed to get user",
			args: mockArgs,
			beforeTest: func(mockObject mockObjects) {
				mockObject.mockUserRepo.EXPECT().GetPasswordVerification(mockArgs.ctx, mockArgs.name).Return(mockArgs.verPassHash, nil)
				mockObject.mockUserRepo.EXPECT().GetUser(&domain.Users{}, domain.UserParam{Name: mockArgs.name}).SetArg(0, mockArgs.user).Return(assert.AnError)
			},
			wantErr: true,
		},
		{
			name: "Failed to update password",
			args: mockArgs,
			beforeTest: func(mockObject mockObjects) {
				mockObject.mockUserRepo.EXPECT().GetPasswordVerification(mockArgs.ctx, mockArgs.name).Return(mockArgs.verPassHash, nil)
				mockObject.mockUserRepo.EXPECT().GetUser(&domain.Users{}, domain.UserParam{Name: mockArgs.name}).SetArg(0, mockArgs.user).Return(nil)
				mockObject.mockUserRepo.EXPECT().UpdateUser(&domain.UserUpdate{Password: string(hashPassword)}, mockArgs.user.Id).Return(assert.AnError)
			},
			wantErr: true,
		},
		{
			name: "Success change password",
			args: mockArgs,
			beforeTest: func(mockObject mockObjects) {
				mockObject.mockUserRepo.EXPECT().GetPasswordVerification(mockArgs.ctx, mockArgs.name).Return(mockArgs.verPassHash, nil)
				mockObject.mockUserRepo.EXPECT().GetUser(&domain.Users{}, domain.UserParam{Name: mockArgs.name}).SetArg(0, mockArgs.user).Return(nil)
				mockObject.mockUserRepo.EXPECT().UpdateUser(&domain.UserUpdate{Password: string(hashPassword)}, mockArgs.user.Id).Return(nil)
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := usecase.NewUserUsecase(mockObject.mockUserRepo, nil, nil, nil, nil)

			tt.beforeTest(mockObject)

			err := w.ChangePassword(tt.args.ctx, tt.args.name, tt.args.verPass, tt.args.passwordRequest)

			if err != nil && !tt.wantErr {
				t.Errorf("error = %v", err)
			}
		})
	}
}

// 	LikeProduct(c *gin.Context, productId uuid.UUID) error
func TestLikeProduct(t *testing.T){
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockJwt := testpkg.NewMockIJwt(ctrl)
	mockUserRepo := test.NewMockIUserRepository(ctrl)
	mockProductRepo := test.NewMockIProductRepository(ctrl)

	mockObject := mockObjects{
		mockUserRepo: mockUserRepo,
		mockJWT: mockJwt,
		mockProduct: mockProductRepo,
	}

	type args struct{
		c *gin.Context
		productId uuid.UUID
		product domain.Products
		loginUser domain.Users
		likeProduct domain.LikeProduct
	}

	var mockContext *gin.Context

	mockProductId := uuid.New()

	mockProduct := domain.Products{
		Id: mockProductId,
		Name: "test product",
		Price: 10000,
		Description: "test decription",
	}

	mockLoginUser := domain.Users{
		Id: uuid.New(),
		Name: "testname",
		Email: "test@email.com",
	}

	mockLikeProduct := domain.LikeProduct{
		UserId: mockLoginUser.Id,
		ProductId: mockProductId,
	}

	mockArgs := args{
		c: mockContext,
		productId: mockProductId,
		product: mockProduct,
		loginUser: mockLoginUser,
		likeProduct: mockLikeProduct,
	}

	tests := []struct{
		name string
		args args
		beforeTest func(mockObject mockObjects)
		wantErr bool
	}{
		{
			name: "Success like product",
			args: mockArgs,
			beforeTest: func(mockObject mockObjects) {
				mockObject.mockJWT.EXPECT().GetLoginUser(mockArgs.c).Return(mockArgs.loginUser, nil)
				mockObject.mockProduct.EXPECT().GetProduct(&domain.Products{}, domain.ProductParam{Id: mockArgs.productId}).SetArg(0, mockArgs.product).Return(nil)
				mockObject.mockUserRepo.EXPECT().LikeProduct(&mockArgs.likeProduct).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "Failed get login user",
			args: mockArgs,
			beforeTest: func(mockObject mockObjects) {
				mockObject.mockJWT.EXPECT().GetLoginUser(mockArgs.c).Return(mockArgs.loginUser, assert.AnError)
			},
			wantErr: true,
		},
		{
			name: "Failed get product",
			args: mockArgs,
			beforeTest: func(mockObject mockObjects) {
				mockObject.mockJWT.EXPECT().GetLoginUser(mockArgs.c).Return(mockArgs.loginUser, nil)
				mockObject.mockProduct.EXPECT().GetProduct(&domain.Products{}, domain.ProductParam{Id: mockArgs.productId}).SetArg(0, mockArgs.product).Return(assert.AnError)
			},
			wantErr: true,
		},
		{
			name: "Failed like product",
			args: mockArgs,
			beforeTest: func(mockObject mockObjects) {
				mockObject.mockJWT.EXPECT().GetLoginUser(mockArgs.c).Return(mockArgs.loginUser, nil)
				mockObject.mockProduct.EXPECT().GetProduct(&domain.Products{}, domain.ProductParam{Id: mockArgs.productId}).SetArg(0, mockArgs.product).Return(nil)
				mockObject.mockUserRepo.EXPECT().LikeProduct(&mockArgs.likeProduct).Return(assert.AnError)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := usecase.NewUserUsecase(mockObject.mockUserRepo, mockObject.mockProduct, mockObject.mockJWT, nil, nil)

			tt.beforeTest(mockObject)

			err := w.LikeProduct(tt.args.c, tt.args.productId)

			if err != nil && !tt.wantErr {
				t.Errorf("error = %v", err)
			}
		})
	}
}

// 	DeleteLikeProduct(c *gin.Context, productId uuid.UUID) error
func TestDeleteLikeProduct(t *testing.T){
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockJwt := testpkg.NewMockIJwt(ctrl)
	mockUserRepo := test.NewMockIUserRepository(ctrl)
	
	mockObject := mockObjects{
		mockUserRepo: mockUserRepo,
		mockJWT: mockJwt,
	}

	type args struct{
		c *gin.Context
		productId uuid.UUID
		loginUser domain.Users
		likeProduct domain.LikeProduct
	}

	var mockContext *gin.Context

	mockProductId := uuid.New()

	mockLoginUser := domain.Users{
		Id: uuid.New(),
		Name: "test name",
		Email: "test@email.com",
		Password: "testpassword",
		LikeProduct: []domain.Products{
			{
				Id: mockProductId,
				Name: "testproduct",
				Price: 100000,
				Description: "testdescription",
			},
		},
	}

	mockLikeProduct := domain.LikeProduct{
		UserId: mockLoginUser.Id,
		ProductId: mockProductId,
	}

	mockArgs := args{
		c: mockContext,
		productId: mockProductId,
		loginUser: mockLoginUser,
		likeProduct: mockLikeProduct,
	}

	fmt.Println("\n\n\\n\n", mockArgs.likeProduct)

	tests := []struct{
		name string
		args args
		beforeTest func(mockObject mockObjects)
		wantErr bool
	}{
		{
			name: "Success delete product",
			args: mockArgs,
			beforeTest: func(mockObject mockObjects) {
				mockObject.mockJWT.EXPECT().GetLoginUser(mockArgs.c).Return(mockArgs.loginUser, nil)
				mockObject.mockUserRepo.EXPECT().GetLikeProduct(&domain.LikeProduct{}, mockArgs.likeProduct).SetArg(0, mockArgs.likeProduct).Return(nil)
				mockObject.mockUserRepo.EXPECT().DeleteLikeProduct(&mockArgs.likeProduct).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "Failed get login user",
			args: mockArgs,
			beforeTest: func(mockObject mockObjects) {
				mockObject.mockJWT.EXPECT().GetLoginUser(mockArgs.c).Return(mockArgs.loginUser, assert.AnError)
			},
			wantErr: true,
		},
		{
			name: "Failed get login product",
			args: mockArgs,
			beforeTest: func(mockObject mockObjects) {
				mockObject.mockJWT.EXPECT().GetLoginUser(mockArgs.c).Return(mockArgs.loginUser, nil)
				mockObject.mockUserRepo.EXPECT().GetLikeProduct(&domain.LikeProduct{}, mockArgs.likeProduct).SetArg(0, mockArgs.likeProduct).Return(assert.AnError)
			},
			wantErr: true,
		},
		{
			name: "Failed delete like product",
			args: mockArgs,
			beforeTest: func(mockObject mockObjects) {
				mockObject.mockJWT.EXPECT().GetLoginUser(mockArgs.c).Return(mockArgs.loginUser, nil)
				mockObject.mockUserRepo.EXPECT().GetLikeProduct(&domain.LikeProduct{}, mockArgs.likeProduct).SetArg(0, mockArgs.likeProduct).Return(nil)
				mockObject.mockUserRepo.EXPECT().DeleteLikeProduct(&mockArgs.likeProduct).Return(assert.AnError)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := usecase.NewUserUsecase(mockObject.mockUserRepo, nil, mockObject.mockJWT, nil, nil)

			tt.beforeTest(mockObject)

			err := w.DeleteLikeProduct(tt.args.c, tt.args.productId)

			if err != nil && !tt.wantErr {
				t.Errorf("error = %v", err)
			}
		})
	}
}
