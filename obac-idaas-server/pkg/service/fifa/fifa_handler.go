package fifa

import (
	"entrust.com/iat/obac-idaas-server/pkg/models"
	"entrust.com/iat/obac-idaas-server/pkg/repository"
	"github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"
	"time"
)

type FIFAHandlerService interface {
	GetBookingInfo(tokenId float64) (models.BookingInfo, error)
	WriteBookingInfo(tokenId float64, bookingInfo models.BookingInfo)
	ValidateToken(signedToken string) (err error)
	GetTokenClaims(signedToken string) (*jwt.MapClaims, error)
}

func NewFIFAHandlerService(repo repository.Repository) FIFAHandlerService {
	return &fifaHandlerServiceImpl{
		repo: repo,
	}
}

type fifaHandlerServiceImpl struct {
	repo repository.Repository
}

func (h *fifaHandlerServiceImpl) GetBookingInfo(tokenId float64) (models.BookingInfo, error) {
	bookingInfo, ok := h.repo.GetBookingInfo(tokenId)
	if !ok {
		return models.BookingInfo{}, errors.New("booking info not found")
	}
	return bookingInfo, nil
}

func (h *fifaHandlerServiceImpl) WriteBookingInfo(tokenId float64, bookingInfo models.BookingInfo) {
	h.repo.WriteBookingInfo(tokenId, bookingInfo)
}

func (h *fifaHandlerServiceImpl) ValidateToken(signedToken string) (err error) {
	token, err := jwt.ParseWithClaims(
		signedToken,
		jwt.MapClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte("superdopekey"), nil
		},
	)
	if err != nil {
		return
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		err = errors.New("couldn't parse claims")
		return
	}
	if int64(claims["exp"].(float64)) < time.Now().Local().Unix() {
		err = errors.New("expired token")
		return
	}
	return
}

func (h *fifaHandlerServiceImpl) GetTokenClaims(signedToken string) (*jwt.MapClaims, error) {
	var claims jwt.MapClaims
	token, err := jwt.ParseWithClaims(
		signedToken,
		jwt.MapClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte("superdopekey"), nil
		},
	)
	if err != nil {
		return &claims, err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		err = errors.New("couldn't parse claims")
		return &claims, err
	}
	return &claims, nil
}
