package usecase

import (
	"errors"

	"github.com/ClearingHouse/internal/admin/interfaces"
	usersInterfaces "github.com/ClearingHouse/internal/users/interfaces"
	"gorm.io/gorm"
)

type AdminUsecase struct {
	userRepo usersInterfaces.UsersRepository
}

func NewAdminUsecase(userRepo usersInterfaces.UsersRepository) interfaces.AdminUsecase {
	return &AdminUsecase{
		userRepo: userRepo,
	}
}

func (u *AdminUsecase) AssignSuperAdmin(email string) error {
	// Get user by email
	user, err := u.userRepo.GetByEmail(email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("user not found")
		}
		return err
	}

	// Check if already super admin
	if user.IsSuperAdmin {
		return errors.New("user is already a super admin")
	}

	// Update user to super admin
	user.IsSuperAdmin = true
	if err := u.userRepo.Update(user); err != nil {
		return err
	}

	return nil
}

func (u *AdminUsecase) RevokeSuperAdmin(email string) error {
	// Get user by email
	user, err := u.userRepo.GetByEmail(email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("user not found")
		}
		return err
	}

	// Check if user is super admin
	if !user.IsSuperAdmin {
		return errors.New("user is not a super admin")
	}

	// Update user to revoke super admin
	user.IsSuperAdmin = false
	if err := u.userRepo.Update(user); err != nil {
		return err
	}

	return nil
}
