package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"intern-bcc/domain"
	"log"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

var RDB *redis.Client

func ConnectToRedis() {
	opts, err := redis.ParseURL(os.Getenv("REDIS_URL"))
	if err != nil {
		log.Fatal("Failed to Connect to Redis")
	}

	redisClient := redis.NewClient(opts)
	RDB = redisClient
}

const (
	keySetOtp                = "otp:set:id:%v"
	keySetPasswordRecovery   = "recovery:set:name:%v"
	keySetInformationNmentor = "get:all:%v"
	keySetProducts           = "get:all:product:page:%v:name:%v:province:%v:university:%v:category:%v"
)

type IRedis interface {
	SetOTP(ctx context.Context, userId uuid.UUID, otpString string) error
	SetEmailVerHash(ctx context.Context, email string, emailVerHash string) error
	GetOTP(ctx context.Context, userId uuid.UUID) (string, error)
	GetEmailVerHash(ctx context.Context, name string) (string, error)
	SetInformationNmentor(ctx context.Context, key string, anyData any) error
	GetArticles(ctx context.Context, key string) ([]domain.Articles, error)
	GetWebinarNCompetition(ctx context.Context, key string) ([]domain.Information, error)
	GetMentors(ctx context.Context, key string) ([]domain.Mentors, error)
	SetAllProducts(ctx context.Context, key domain.ProductParam, data []domain.Products) error
	GetAllProdoucts(ctx context.Context, key domain.ProductParam) ([]domain.Products, error)
	// SetWebinarNCompetition(ctx context.Context, articles []domain.WebinarNCompetition) error
	// GetWebinarNCompetition(ctx context.Context) ([]domain.Information, error)
}

type Redis struct {
	r *redis.Client
}

func RedisInit(r *redis.Client) IRedis {
	return &Redis{r}
}

func (r *Redis) SetOTP(ctx context.Context, userId uuid.UUID, otpString string) error {
	key := fmt.Sprintf(keySetOtp, userId)

	err := r.r.SetEx(ctx, key, otpString, 2*time.Minute).Err()
	if err != nil {
		return err
	}

	return nil
}

func (r *Redis) GetOTP(ctx context.Context, userId uuid.UUID) (string, error) {
	key := fmt.Sprintf(keySetOtp, userId)

	stringOTP, err := r.r.Get(ctx, key).Result()
	if err != nil {
		return stringOTP, err
	}

	return stringOTP, nil
}

func (r *Redis) SetEmailVerHash(ctx context.Context, name string, emailVerHash string) error {
	key := fmt.Sprintf(keySetPasswordRecovery, name)

	err := r.r.SetEx(ctx, key, string(emailVerHash), 2*time.Minute).Err()
	if err != nil {
		return err
	}

	return nil
}

func (r *Redis) GetEmailVerHash(ctx context.Context, name string) (string, error) {
	key := fmt.Sprintf(keySetPasswordRecovery, name)

	verPassHash, err := r.r.Get(ctx, key).Result()
	fmt.Println(verPassHash)
	if err != nil {
		return verPassHash, err
	}

	return verPassHash, nil
}

func (r *Redis) SetInformationNmentor(ctx context.Context, key string, anyData any) error {
	var data any
	var redisKey string

	fmt.Println("\n\n\nTalking to database")
	switch key {
	case "Articles":
		data = anyData.([]domain.Articles)
		redisKey = fmt.Sprintf(keySetInformationNmentor, key)
	case "WebinarNCompetition":
		data = anyData.([]domain.Information)
		redisKey = fmt.Sprintf(keySetInformationNmentor, key)
	case "Mentors":
		data = anyData.([]domain.Mentors)
		fmt.Println(data)
		redisKey = fmt.Sprintf(keySetInformationNmentor, key)
	}

	byteData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	err = r.r.SetEx(ctx, redisKey, string(byteData), 5*time.Minute).Err()
	if err != nil {
		return err
	}

	return nil
}

func (r *Redis) GetArticles(ctx context.Context, key string) ([]domain.Articles, error) {
	var Articles []domain.Articles
	redisKey := fmt.Sprintf(keySetInformationNmentor, key)

	stringArticles, err := r.r.Get(ctx, redisKey).Result()
	if err != nil {
		return Articles, err
	}

	err = json.Unmarshal([]byte(stringArticles), &Articles)
	if err != nil {
		return Articles, nil
	}

	return Articles, nil
}

func (r *Redis) GetWebinarNCompetition(ctx context.Context, key string) ([]domain.Information, error) {
	var webinarNCompetition []domain.Information
	redisKey := fmt.Sprintf(keySetInformationNmentor, key)
	stringWebinarNCompetition, err := r.r.Get(ctx, redisKey).Result()
	if err != nil {
		return webinarNCompetition, err
	}

	err = json.Unmarshal([]byte(stringWebinarNCompetition), &webinarNCompetition)
	if err != nil {
		return webinarNCompetition, nil
	}

	return webinarNCompetition, nil
}

func (r *Redis) GetMentors(ctx context.Context, key string) ([]domain.Mentors, error) {
	var mentors []domain.Mentors
	redisKey := fmt.Sprintf(keySetInformationNmentor, key)

	stringMentor, err := r.r.Get(ctx, redisKey).Result()
	if err != nil {
		return mentors, err
	}

	err = json.Unmarshal([]byte(stringMentor), &mentors)
	if err != nil {
		return mentors, nil
	}
	fmt.Println(mentors)

	return mentors, nil
}

func (r *Redis) SetAllProducts(ctx context.Context, key domain.ProductParam, data []domain.Products) error {
	redisKey := fmt.Sprintf(keySetProducts, key.Page, key.Name, key.ProvinceId, key.UniversityId, key.CategoryId)

	byteData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	err = r.r.SetEx(ctx, redisKey, string(byteData), 5*time.Minute).Err()
	if err != nil {
		return err
	}

	fmt.Printf("\n\n\nTalking to database")

	return nil
}

func (r *Redis) GetAllProdoucts(ctx context.Context, key domain.ProductParam) ([]domain.Products, error) {
	var products []domain.Products
	redisKey := fmt.Sprintf(keySetProducts, key.Page, key.Name, key.ProvinceId, key.UniversityId, key.CategoryId)

	stringProducts, err := r.r.Get(ctx, redisKey).Result()
	if err != nil {
		return products, err
	}

	err = json.Unmarshal([]byte(stringProducts), &products)
	if err != nil {
		return products, nil
	}

	return products, nil
}

// func (r *Redis) GetInformationNmentor(ctx context.Context, key string) (any, error) {
// 	// var articles []domain.Articles
// 	var result any
// 	var redisKey string
// 	switch key {
// 	case "Articles":
// 		redisKey = fmt.Sprintf(keySetInformationNmentor, key)
// 		// result = result.([]domain.Articles)
// 	case "WebinarNCompetition":
// 		redisKey = fmt.Sprintf(keySetInformationNmentor, key)
// 		// result = result.([]domain.Information)
// 	case "Mentors":
// 		redisKey = fmt.Sprintf(keySetInformationNmentor, key)
// 		// result = result.([]domain.Mentors)
// 	}

// 	stringArticles, err := r.r.Get(ctx, redisKey).Result()
// 	if err != nil {
// 		return result, err
// 	}

// 	err = json.Unmarshal([]byte(stringArticles), &result)
// 	if err != nil {
// 		return result, nil
// 	}
// 	fmt.Println(result)
// 	switch key {
// 	case "Articles":
// 		// redisKey = fmt.Sprintf(keySetInformationNmentor, key)
// 		result = result.([]domain.Articles)
// 	case "WebinarNCompetition":
// 		// redisKey = fmt.Sprintf(keySetInformationNmentor, key)
// 		result = result.([]domain.Information)
// 	case "Mentors":
// 		// redisKey = fmt.Sprintf(keySetInformationNmentor, key)
// 		result = result.([]domain.Mentors)
// 	}

// 	return result, nil
// }

// func (r *Redis) SetWebinarNCompetition(ctx context.Context, articles []domain.WebinarNCompetition) error {
// 	byteWebinarNCompetition, err := json.Marshal(articles)
// 	if err != nil {
// 		return err
// 	}

// 	err = r.r.SetEx(ctx, keySetArticle, string(byteWebinarNCompetition), 5*time.Minute).Err()
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }
