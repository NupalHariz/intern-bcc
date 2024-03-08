package repository

import (
	"context"
	"fmt"

	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

const (
	keySetOtp = "otp:set:id:%v"
)

type IMerchantRedis interface {
	SetOTP(ctx context.Context, userId uuid.UUID, otpString string) error
	GetOTP(ctx context.Context, userId uuid.UUID) (string, error)
}

type Redis struct {
	r *redis.Client
}

func NewMerchantRedis(r *redis.Client) IMerchantRedis {
	return &Redis{r}
}

func (r *Redis) SetOTP(ctx context.Context, userId uuid.UUID, otpString string) error {

	key := fmt.Sprintf(keySetOtp, userId)

	// byteOTP, err := json.Marshal(otpString)
	// if err != nil {
	// 	return err
	// }

	// fmt.Println(byteOTP)

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

// func (r *MerchantRepository) GetOTP(ctx context.Context, userId int) (string, error) {
// 	key := fmt.Sprintf(keySetOtp, userId)

// 	stringOTP, err := r.r.Get(ctx, key).Result()
// 	if err != nil {
// 		return stringOTP, err
// 	}

// 	return stringOTP, nil
// }
