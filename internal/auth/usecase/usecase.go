package usecase

import (
	"context"
	"log"
	"strings"
	"time"

	"errors"

	"github.com/ClearingHouse/config"
	"github.com/ClearingHouse/internal/auth"
	"github.com/ClearingHouse/internal/auth/interfaces"
	authRepo "github.com/ClearingHouse/internal/auth/repository"
	"github.com/ClearingHouse/internal/models"
	orgInterfaces "github.com/ClearingHouse/internal/organizations/interfaces"
	userInterfaces "github.com/ClearingHouse/internal/users/interfaces"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/oauth2"
)

type AuthUsecase struct {
	userRepo      userInterfaces.UsersRepository
	blacklistRepo *authRepo.TokenBlacklistRepository
	orgRepo       orgInterfaces.OrganizationRepository
}

func NewAuthUsecase(userRepo userInterfaces.UsersRepository, blacklistRepo *authRepo.TokenBlacklistRepository, orgRepo orgInterfaces.OrganizationRepository) interfaces.AuthUsecase {
	return &AuthUsecase{
		userRepo:      userRepo,
		blacklistRepo: blacklistRepo,
		orgRepo:       orgRepo,
	}
}

func (u *AuthUsecase) GenerateGoogleLoginURL(state string) string {
	return config.GoogleOauthConfig.AuthCodeURL(state, oauth2.AccessTypeOffline)
}

func (u *AuthUsecase) GenerateGoogleRegisterURL(state string) string {
	return config.GoogleOauthConfig.AuthCodeURL(state, oauth2.AccessTypeOffline)
}

func (u *AuthUsecase) HandleGoogleCallback(code string, c *gin.Context) (*models.User, error) {
	token, err := config.GoogleOauthConfig.Exchange(context.TODO(), code)
	if err != nil {
		return nil, err
	}

	googleUser, err := u.userRepo.GetUserGoogle(token)
	if err != nil {
		return nil, err
	}

	email, emailOk := googleUser["email"].(string)
	firstName, firstNameOk := googleUser["given_name"].(string)
	lastName, lastNameOk := googleUser["family_name"].(string)
	log.Println(email, firstName, lastName)
	if !emailOk || !firstNameOk || !lastNameOk {
		return nil, errors.New("invalid user data from Google")
	}

	// Only find existing users, don't create new ones
	user, err := u.userRepo.GetByEmail(email)
	if err != nil {
		return nil, errors.New("user not registered. Please register first before logging in")
	}

	return user, nil
}

func (u *AuthUsecase) HandleGoogleRegisterCallback(code string, c *gin.Context) (*models.User, error) {
	token, err := config.GoogleOauthConfig.Exchange(context.TODO(), code)
	if err != nil {
		return nil, err
	}

	googleUser, err := u.userRepo.GetUserGoogle(token)
	if err != nil {
		return nil, err
	}

	email, emailOk := googleUser["email"].(string)
	firstName, firstNameOk := googleUser["given_name"].(string)
	lastName, lastNameOk := googleUser["family_name"].(string)
	log.Println(email, firstName, lastName)
	if !emailOk || !firstNameOk || !lastNameOk {
		return nil, errors.New("invalid user data from Google")
	}

	existingUser, err := u.userRepo.GetByEmail(email)
	if existingUser != nil {
		return nil, errors.New("user already registered. Please login instead")
	}

	newUser := &models.User{
		Email:     email,
		FirstName: firstName,
		LastName:  lastName,
	}

	err = u.userRepo.Create(newUser)
	if err != nil {
		return nil, err
	}

	// Extract domain from email
	// emailParts := strings.Split(email, "@")
	// if len(emailParts) == 2 {
	// 	domain := emailParts[1]
	// 	org, err := u.orgRepo.GetOrganizationByDomain(domain)
	// 	if err == nil && org != nil {
	// 		// Add user to organization members
	// 		org.Members = append(org.Members, *newUser)
	// 		err = u.orgRepo.UpdateMembers(org)
	// 		if err != nil {
	// 			log.Printf("Failed to auto-add user to organization: %v", err)
	// 			// Don't fail registration if auto-add fails
	// 		} else {
	// 			log.Printf("User %s auto-added to organization %s based on domain %s", email, org.Name, domain)
	// 		}
	// 	}
	// }

	return newUser, nil
}

// ManualRegister creates a user without OAuth for testing purposes
func (u *AuthUsecase) ManualRegister(email, firstName, lastName string) (*models.User, error) {
	existingUser, err := u.userRepo.GetByEmail(email)
	if existingUser != nil {
		return nil, errors.New("user already registered. Please login instead")
	}

	newUser := &models.User{
		Email:     email,
		FirstName: firstName,
		LastName:  lastName,
	}

	err = u.userRepo.Create(newUser)
	if err != nil {
		return nil, err
	}

	// Extract domain from email
	emailParts := strings.Split(email, "@")
	if len(emailParts) == 2 {
		domain := emailParts[1]
		org, err := u.orgRepo.GetOrganizationByDomain(domain)
		if err == nil && org != nil {
			// Add user to organization members
			org.Members = append(org.Members, *newUser)
			err = u.orgRepo.UpdateMembers(org)
			if err != nil {
				log.Printf("Failed to auto-add user to organization: %v", err)
				// Don't fail registration if auto-add fails
			} else {
				log.Printf("User %s auto-added to organization %s based on domain %s", email, org.Name, domain)
			}
		}
	}

	return newUser, nil
}

func (u *AuthUsecase) GenerateTokens(user *models.User) (accessToken string, refreshToken string, err error) {
	privateKey, err := auth.LoadPrivateKey("token_private.pem")
	if err != nil {
		return "", "", err
	}

	accessToken, err = auth.GenerateAccessToken(user.ID.String(), user.Email, user.FirstName, user.LastName, user.IsSuperAdmin, privateKey)
	if err != nil {
		return "", "", err
	}

	refreshToken, err = auth.GenerateRefreshToken(user.ID.String())
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

// Refresh endpoint handler
func (u *AuthUsecase) RefreshAccessToken(refreshToken string) (string, error) {
	// Check if token is blacklisted
	isBlacklisted, err := u.blacklistRepo.IsBlacklisted(refreshToken)
	if err != nil {
		return "", err
	}
	if isBlacklisted {
		return "", errors.New("token has been revoked")
	}

	claims, err := auth.VerifyRefreshToken(refreshToken)
	if err != nil {
		return "", err
	}

	user, err := u.userRepo.GetByID(uuid.MustParse(claims.UserID))
	if err != nil {
		return "", err
	}

	privateKey, err := auth.LoadPrivateKey("token_private.pem")
	if err != nil {
		return "", err
	}

	return auth.GenerateAccessToken(user.ID.String(), user.Email, user.FirstName, user.LastName, user.IsSuperAdmin, privateKey)
}

func (u *AuthUsecase) GetUserByID(userID uuid.UUID) (*models.User, error) {
	return u.userRepo.GetByID(userID)
}

func (u *AuthUsecase) BlacklistToken(token string) error {
	// Verify the token to get its expiration time
	claims, err := auth.VerifyRefreshToken(token)
	if err != nil {
		// If token is invalid/expired, no need to blacklist
		return nil
	}

	// Add to blacklist with expiration time
	expiresAt := time.Unix(claims.ExpiresAt.Unix(), 0)
	return u.blacklistRepo.AddToBlacklist(token, expiresAt)
}

// func (u *AuthUsecase) GenerateToken(email string) (string, error) {
// 	// Generate a JWT token for the user
// 	token, err := createJWT(email)
// 	if err != nil {
// 		return "", err
// 	}
// 	return token, nil
// }
