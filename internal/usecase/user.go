package usecase

import (
	"errors"
	"fmt"
	"intern-bcc/domain"
	"intern-bcc/internal/repository"
	"intern-bcc/pkg/gomail"
	"intern-bcc/pkg/jwt"
	"intern-bcc/pkg/response"
	"intern-bcc/pkg/supabase"
	"math/rand"
	"mime/multipart"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/net/context"
)

type IUserUsecase interface {
	GetUser(param domain.UserParam) (domain.Users, error)
	GetLikeProducts(c *gin.Context) ([]domain.ProductResponses, error)
	GetOwnProducts(c *gin.Context) ([]domain.ProductResponses, error)
	GetOwnMentors(c *gin.Context) ([]domain.OwnMentorResponses, error)
	Register(userRequest domain.UserRequest) error
	Login(userLogin domain.UserLogin) (domain.LoginResponse, error)
	UpdateUser(c *gin.Context, userId uuid.UUID, userUpdate domain.UserUpdate) (domain.Users, error)
	UploadUserPhoto(c *gin.Context, userId uuid.UUID, userPhoto *multipart.FileHeader) (domain.Users, error)
	PasswordRecovery(userParam domain.UserParam, ctx context.Context) error
	ChangePassword(ctx context.Context, name string, verPass string, passwordRequest domain.PasswordUpdate) error
	LikeProduct(c *gin.Context, productId uuid.UUID) error
	DeleteLikeProduct(c *gin.Context, productId uuid.UUID) error
}

type UserUsecase struct {
	userRepository    repository.IUserRepository
	productRepository repository.IProductRepository
	jwt               jwt.IJwt
	supabase          supabase.ISupabase
	redis             repository.IRedis
	goMail            gomail.IGoMail
}

func NewUserUsecase(userRepository repository.IUserRepository, productRepository repository.IProductRepository, jwt jwt.IJwt,
	supabase supabase.ISupabase, redis repository.IRedis, goMail gomail.IGoMail) IUserUsecase {
	return &UserUsecase{
		userRepository:    userRepository,
		productRepository: productRepository,
		jwt:               jwt,
		supabase:          supabase,
		redis:             redis,
		goMail:            goMail,
	}
}

func (u *UserUsecase) GetUser(param domain.UserParam) (domain.Users, error) {
	var user domain.Users
	err := u.userRepository.GetUser(&user, param)
	if err != nil {
		return user, response.NewError(http.StatusNotFound, "an error occured when get user", err)
	}

	return user, nil
}

func (u *UserUsecase) GetLikeProducts(c *gin.Context) ([]domain.ProductResponses, error) {
	user, err := u.jwt.GetLoginUser(c)
	if err != nil {
		return []domain.ProductResponses{}, response.NewError(http.StatusInternalServerError, "an error occured when get login user", err)
	}

	err = u.userRepository.GetLikeProducts(&user, user.Id)
	if err != nil {
		return []domain.ProductResponses{}, response.NewError(http.StatusInternalServerError, "an error occured when get liked product", err)
	}

	var productResponses []domain.ProductResponses
	for _, lk := range user.LikeProduct {
		productResponse := domain.ProductResponses{
			Id:           lk.Id,
			Name:         lk.Name,
			MerchantName: lk.Merchant.MerchantName,
			University:   lk.Merchant.University.University,
			Price:        lk.Price,
			ProductPhoto: lk.ProductPhoto,
		}

		productResponses = append(productResponses, productResponse)
	}

	return productResponses, nil
}

func (u *UserUsecase) GetOwnProducts(c *gin.Context) ([]domain.ProductResponses, error) {
	user, err := u.jwt.GetLoginUser(c)
	if err != nil {
		return []domain.ProductResponses{}, response.NewError(http.StatusInternalServerError, "failed to get login user", err)
	}

	err = u.userRepository.GetOwnProducts(&user, user.Id)
	if err != nil {
		return []domain.ProductResponses{}, response.NewError(http.StatusInternalServerError, "an error occured when get own product", err)
	}

	var productResponses []domain.ProductResponses
	for _, p := range user.Merchant.Products {
		productResponse := domain.ProductResponses{
			Id:           p.Id,
			Name:         p.Name,
			MerchantName: user.Merchant.MerchantName,
			University:   user.Merchant.University.University,
			Price:        p.Price,
			ProductPhoto: p.ProductPhoto,
		}

		productResponses = append(productResponses, productResponse)
	}

	return productResponses, nil
}

func (u *UserUsecase) GetOwnMentors(c *gin.Context) ([]domain.OwnMentorResponses, error) {
	user, err := u.jwt.GetLoginUser(c)
	if err != nil {
		return []domain.OwnMentorResponses{}, response.NewError(http.StatusInternalServerError, "failed to get login user", err)
	}

	err = u.userRepository.GetOwnMentors(&user, user.Id)
	if err != nil {
		return []domain.OwnMentorResponses{}, response.NewError(http.StatusInternalServerError, "an error occured when get own mentors", err)
	}

	var ownMentorResponses []domain.OwnMentorResponses
	for _, m := range user.HasMentors {
		ownMentorResponse := domain.OwnMentorResponses{
			Id:            m.Id,
			Name:          m.Name,
			CurrentJob:    m.CurrentJob,
			MentorPicture: m.MentorPicture,
		}

		ownMentorResponses = append(ownMentorResponses, ownMentorResponse)
	}

	return ownMentorResponses, nil
}

func (u *UserUsecase) Register(userRequest domain.UserRequest) error {
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(userRequest.Password), 10)
	if err != nil {
		return response.NewError(http.StatusInternalServerError, "error when hashing password", err)
	}

	newUser := domain.Users{
		Id:       uuid.New(),
		Name:     userRequest.Name,
		Email:    userRequest.Email,
		Password: string(hashPassword),
	}

	err = u.userRepository.Register(&newUser)
	if err != nil {
		return response.NewError(http.StatusInternalServerError, "error occured when creating user", err)
	}

	return nil
}

func (u *UserUsecase) Login(userLogin domain.UserLogin) (domain.LoginResponse, error) {
	var user domain.Users
	err := u.userRepository.GetUser(&user, domain.UserParam{
		Email: userLogin.Email,
	})
	if err != nil {
		return domain.LoginResponse{}, response.NewError(http.StatusNotFound, "email or password invalid", err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(userLogin.Password))
	if err != nil {
		return domain.LoginResponse{}, response.NewError(http.StatusNotFound, "email or password invalid", err)
	}

	tokenString, err := u.jwt.GenerateToken(user.Id)
	if err != nil {
		return domain.LoginResponse{}, response.NewError(http.StatusInternalServerError, "failed to generate jwt token", err)
	}

	loginUser := domain.LoginResponse{
		JWT: tokenString,
	}

	return loginUser, nil
}

func (u *UserUsecase) UpdateUser(c *gin.Context, userId uuid.UUID, userUpdate domain.UserUpdate) (domain.Users, error) {
	user, err := u.jwt.GetLoginUser(c)
	if err != nil {
		return domain.Users{}, response.NewError(http.StatusNotFound, "an error occured when get login user", err)
	}

	if userId != user.Id {
		return domain.Users{}, response.NewError(http.StatusUnauthorized, "access denied", errors.New("can not change other people profile"))
	}

	err = u.userRepository.UpdateUser(&userUpdate, userId)
	if err != nil {
		return domain.Users{}, response.NewError(http.StatusInternalServerError, "error occured when update user", err)
	}

	var updatedUser domain.Users
	err = u.userRepository.GetUser(&updatedUser, domain.UserParam{Id: userId})
	if err != nil {
		return domain.Users{}, response.NewError(http.StatusInternalServerError, "an error occured when get updated user", err)
	}

	return updatedUser, nil
}

func (u *UserUsecase) UploadUserPhoto(c *gin.Context, userId uuid.UUID, userPhoto *multipart.FileHeader) (domain.Users, error) {
	user, err := u.jwt.GetLoginUser(c)
	if err != nil {
		return domain.Users{}, response.NewError(http.StatusNotFound, "an error occured when get login user", err)
	}

	if user.Id != userId {
		return domain.Users{}, response.NewError(http.StatusUnauthorized, "access denied", errors.New("can not change other people profile"))
	}

	if user.ProfilePicture != "" {
		err = u.supabase.Delete(user.ProfilePicture)
		if err != nil {
			return domain.Users{}, response.NewError(http.StatusInternalServerError, "error occured when deleting old profile picture", err)
		}
	}

	userPhoto.Filename = fmt.Sprintf("%v-%v", time.Now().String(), userPhoto.Filename)
	userPhoto.Filename = strings.Replace(userPhoto.Filename, " ", "-", -1)

	newProfilePicture, err := u.supabase.Upload(userPhoto)
	if err != nil {
		return domain.Users{}, response.NewError(http.StatusInternalServerError, "failed to upload photo", err)
	}

	err = u.userRepository.UpdateUser(&domain.UserUpdate{ProfilePicture: newProfilePicture}, user.Id)
	if err != nil {
		return domain.Users{}, response.NewError(http.StatusInternalServerError, "error occured when update user", err)
	}

	var updatedUser domain.Users
	err = u.userRepository.GetUser(&updatedUser, domain.UserParam{Id: user.Id})
	if err != nil {
		return domain.Users{}, response.NewError(http.StatusInternalServerError, "an error occured when get updated user", err)
	}

	return updatedUser, nil
}

func (u *UserUsecase) PasswordRecovery(userParam domain.UserParam, ctx context.Context) error {
	var user domain.Users
	err := u.userRepository.GetUser(&user, domain.UserParam{Email: userParam.Email})
	if err != nil {
		return response.NewError(http.StatusNotFound, "account not found", err)
	}

	var emailVerPassword string
	emailVerPassword, emailVerHash, err := newEmailVerification()
	if err != nil {
		return response.NewError(http.StatusInternalServerError, "error occured when send email", err)
	}

	userName := user.Name
	if strings.Contains(user.Name, " ") {
		userName = strings.Replace(userName, " ", "-", -1)
	}

	err = u.redis.SetEmailVerHash(ctx, userName, emailVerHash)
	if err != nil {
		return response.NewError(http.StatusInternalServerError, "error occured when send email", err)
	}

	address := os.Getenv("APP_ADDRESS")
	host := os.Getenv("APP_PORT")
	domainName := fmt.Sprintf("%v:%v", address, host)

	subject := "Account Recovery"
	htmlBody := `<html>
	<h1>Click Link to Change Password</h1>
	<h2><a href="http://` + domainName + `/api/v1/` + `recoveryaccount/` + userName + `/` + emailVerPassword + `">click here</a></h2>
	</html>`

	err = u.goMail.SendGoMail(subject, htmlBody, user.Email)
	if err != nil {
		return response.NewError(http.StatusInternalServerError, "an error occured when send email", err)
	}

	return nil
}

func (u *UserUsecase) ChangePassword(ctx context.Context, name string, verPass string, passwordRequest domain.PasswordUpdate) error {
	verPassHash, err := u.redis.GetEmailVerHash(ctx, name)
	if err != nil {
		return response.NewError(http.StatusInternalServerError, "an error occured when get ver password from database", err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(verPassHash), []byte(verPass))
	if err != nil {
		return response.NewError(http.StatusBadRequest, "verification code invalid", err)
	}

	if passwordRequest.Password != passwordRequest.ConfirmPassword {
		return response.NewError(http.StatusBadRequest, "failed to create new password", errors.New("password and confirm password is different"))
	}

	hashedPass, err := bcrypt.GenerateFromPassword([]byte(passwordRequest.Password), 10)
	if err != nil {
		return response.NewError(http.StatusInternalServerError, "an error occured when make new password", err)
	}

	userName := name
	if strings.Contains(name, "-") {
		userName = strings.Replace(userName, "-", " ", -1)
	}

	var user domain.Users
	err = u.userRepository.GetUser(&user, domain.UserParam{Name: userName})
	if err != nil {
		return response.NewError(http.StatusNotFound, "can not to find user", err)
	}

	err = u.userRepository.UpdateUser(&domain.UserUpdate{Password: string(hashedPass)}, user.Id)
	if err != nil {
		return response.NewError(http.StatusInternalServerError, "an error occured when update password", err)
	}

	return nil
}

func newEmailVerification() (string, string, error) {
	var alphaNumRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")
	emailVerRandRune := make([]rune, 64)
	for i := 0; i < 64; i++ {
		emailVerRandRune[i] = alphaNumRunes[rand.Intn(len(alphaNumRunes)-1)]
	}
	emailVerPassword := string(emailVerRandRune)
	var emailVerPWhash []byte

	emailVerPWhash, err := bcrypt.GenerateFromPassword([]byte(emailVerPassword), bcrypt.DefaultCost)
	if err != nil {
		return "", "", err
	}
	emailVerHash := string(emailVerPWhash)

	return emailVerPassword, emailVerHash, nil
}

func (u *UserUsecase) LikeProduct(c *gin.Context, productId uuid.UUID) error {
	user, err := u.jwt.GetLoginUser(c)
	if err != nil {
		return response.NewError(http.StatusNotFound, "an error occured when get login user", err)
	}

	var product domain.Products
	err = u.productRepository.GetProduct(&product, domain.ProductParam{Id: productId})
	if err != nil {
		return response.NewError(http.StatusNotFound, "failed to find product", err)
	}

	likeProduct := domain.LikeProduct{
		UserId:    user.Id,
		ProductId: product.Id,
	}

	err = u.userRepository.LikeProduct(&likeProduct)
	if err != nil {
		return response.NewError(http.StatusInternalServerError, "an error occured when like product", err)
	}

	return nil
}

func (u *UserUsecase) DeleteLikeProduct(c *gin.Context, productId uuid.UUID) error {
	user, err := u.jwt.GetLoginUser(c)
	if err != nil {
		return response.NewError(http.StatusNotFound, "an error occured when get login user", err)
	}

	var likedProduct domain.LikeProduct
	err = u.userRepository.GetLikeProduct(&likedProduct, domain.LikeProduct{UserId: user.Id, ProductId: productId})
	if err != nil {
		return response.NewError(http.StatusNotFound, "an error occured when get liked product", err)
	}

	err = u.userRepository.DeleteLikeProduct(&likedProduct)
	if err != nil {
		return response.NewError(http.StatusInternalServerError, "failed to deleted liked product", err)
	}

	return nil
}
