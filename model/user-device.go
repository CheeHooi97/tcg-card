package model

import (
	"pkm/config"
	"pkm/utils"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type UserDevicePlatform string

const (
	UserDevicePlatformAndroid UserDevicePlatform = "ANDROID"
	UserDevicePlatformIOS     UserDevicePlatform = "IOS"
	UserDevicePlatformWeb     UserDevicePlatform = "WEB"
)

type UserDevice struct {
	Id         string             `gorm:"primaryKey" json:"id"`
	UserId     string             `json:"userId"`
	Platform   UserDevicePlatform `json:"platform"`
	DeviceId   string             `json:"deviceId"`
	DeviceInfo string             `json:"deviceInfo"`
	PNSToken   string             `json:"pnsToken"`
	IsHuawei   bool               `json:"isHuawei"`
	BaseModel
}

func (u User) GetAccessToken(device *UserDevice) (string, error) {
	claims := jwt.RegisteredClaims{
		ID:      u.Id,
		Subject: device.Id,
		Audience: jwt.ClaimStrings([]string{
			string(device.Platform),
			device.DeviceId,
		}),
		Issuer:    "doremic",
		NotBefore: jwt.NewNumericDate(time.Now().UTC()),
		IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		// ExpiresAt: jwt.NewNumericDate(time.Now().UTC().AddDate(0, 1, 0)),
	}

	return jwt.NewWithClaims(jwt.SigningMethodRS256, claims).SignedString(config.AuthenticationPrivateKey)
}

func NewUserDevice() *UserDevice {
	now := time.Now().UTC()

	m := new(UserDevice)
	m.Id = utils.UniqueID()
	m.CreatedDateTime = now
	m.UpdatedDateTime = now

	return m
}

func (m *UserDevice) DateTime() {
	m.CreatedDateTime = time.Now().UTC()
	m.UpdatedDateTime = time.Now().UTC()
}

func (m *UserDevice) UpdateDt() {
	m.UpdatedDateTime = time.Now().UTC()
}
