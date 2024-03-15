package repository

import (
	"context"
	"fmt"

	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

const (
	keySetOtp              = "otp:set:id:%v"
	keySetPasswordRecovery = "recovery:set:email:%v"
)

type IRedis interface {
	SetOTP(ctx context.Context, userId uuid.UUID, otpString string) error
	SetEmailVerHash(ctx context.Context, email string, emailVerHash string) error
	GetOTP(ctx context.Context, userId uuid.UUID) (string, error)
	GetEmailVerHash(ctx context.Context, email string) (string, error)
}

type Redis struct {
	r *redis.Client
}

func NewRedis(r *redis.Client) IRedis {
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

func (r *Redis) SetEmailVerHash(ctx context.Context, email string, emailVerHash string) error {
	key := fmt.Sprintf(keySetPasswordRecovery, email)

	// byteVerPassHash, err := json.Marshal(emailVerHash)
	// if err != nil {
	// 	return err
	// }

	err := r.r.SetEx(ctx, key, string(emailVerHash), 2*time.Minute).Err()
	if err != nil {
		return err
	}

	return nil
}

func (r *Redis) GetEmailVerHash(ctx context.Context, email string) (string, error) {
	key := fmt.Sprintf(keySetPasswordRecovery, email)

	verPassHash, err := r.r.Get(ctx, key).Result()
	fmt.Println(verPassHash)
	if err != nil {
		return verPassHash, err
	}

	return verPassHash, nil
}

// func (r *MerchantRepository) GetOTP(ctx context.Context, userId int) (string, error) {
// 	key := fmt.Sprintf(keySetOtp, userId)

// 	stringOTP, err := r.r.Get(ctx, key).Result()
// 	if err != nil {
// 		return stringOTP, err
// 	}

// 	return stringOTP, nil
// }
