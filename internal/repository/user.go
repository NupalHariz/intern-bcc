package repository

import (
	"intern-bcc/domain"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type IUserRepository interface {
	GetUser(user *domain.Users, param domain.UserParam) error
	GetLikeProduct(likedProduct *domain.LikeProduct, likeProductParam domain.LikeProduct) error
	GetLikeProducts(user *domain.Users, userId uuid.UUID) error
	GetOwnProducts(user *domain.Users, userId uuid.UUID) error
	GetOwnMentors(user *domain.Users, userId uuid.UUID) error
	Register(newUser *domain.Users) error
	UpdateUser(userUpdate *domain.UserUpdate, userId uuid.UUID) error
	LikeProduct(likeProduct *domain.LikeProduct) error
	DeleteLikeProduct(likedProduct *domain.LikeProduct) error
	CreateHasMentor(mentor *domain.HasMentor) error
}

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) IUserRepository {
	return &UserRepository{db}
}

func (r *UserRepository) GetUser(user *domain.Users, param domain.UserParam) error {
	err := r.db.First(user, &param).Error
	if err != nil {
		return err
	}

	return nil
}

func (r *UserRepository) GetLikeProduct(likedProduct *domain.LikeProduct, likeProductParam domain.LikeProduct) error {
	err := r.db.Table("user_like_product").First(likedProduct, likeProductParam).Error
	if err != nil {
		return err
	}

	return nil
}

func (r *UserRepository) GetLikeProducts(user *domain.Users, userId uuid.UUID) error {
	err := r.db.Model(domain.Users{}).Preload("LikeProduct.Merchant.University").Find(user, "id = ?", userId).Error
	if err != nil {
		return err
	}

	return nil
}

func (r *UserRepository) GetOwnProducts(user *domain.Users, userId uuid.UUID) error {
	err := r.db.Preload("Merchant.University").Preload("Merchant.Products").Find(user, "id = ?", userId).Error
	if err != nil {
		return err
	}

	return nil
}

func (r *UserRepository) GetOwnMentors(user *domain.Users, userId uuid.UUID) error {
	err := r.db.Preload("HasMentors").Find(user, "id = ?", userId).Error
	if err != nil {
		return err
	}

	return nil
}

func (r *UserRepository) Register(newUser *domain.Users) error {
	err := r.db.Create(newUser).Error
	if err != nil {
		return err
	}

	return nil
}

func (r *UserRepository) UpdateUser(userUpdate *domain.UserUpdate, userId uuid.UUID) error {
	var user domain.Users
	err := r.db.Model(&user).Where("id = ?", userId).Updates(userUpdate).Error
	if err != nil {
		return err
	}

	return nil
}

func (r *UserRepository) LikeProduct(likeProduct *domain.LikeProduct) error {
	err := r.db.Table("user_like_product").Create(likeProduct).Error
	if err != nil {
		return err
	}

	return nil
}

func (r *UserRepository) DeleteLikeProduct(likedProduct *domain.LikeProduct) error {
	err := r.db.Table("user_like_product").Delete(likedProduct, likedProduct).Error
	if err != nil {
		return err
	}

	return nil
}

func (r *UserRepository) CreateHasMentor(mentor *domain.HasMentor) error {
	err := r.db.Table("has_mentors").Create(mentor).Error
	if err != nil {
		return err
	}

	return nil
}
