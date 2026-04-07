package handler

import (
	"encoding/json"
	"pkm/errcode"
	"pkm/model"
	"pkm/utils"

	"github.com/labstack/echo/v4"
)

func (h *Handler) Register(c echo.Context) error {
	var i struct {
		Email      string            `json:"email" validate:"required,email"`
		DeviceId   string            `json:"deviceId"`
		DeviceInfo map[string]string `json:"deviceInfo"`
		PNSToken   string            `json:"pnsToken"`
	}

	if msg, err := utils.ValidateRequest(c, &i); err != nil {
		return responseValidationError(c, msg)
	}

	checkUser, err := h.User.GetByEmail(i.Email)
	if err != nil {
		return responseError(c, errcode.InternalServerError)
	} else if checkUser.Id != "" {
		return responseError(c, errcode.RegisteredPhoneNumber)
	}

	user := model.NewUser()
	user.Email = i.Email
	user.Status = true

	if err := h.User.Create(user); err != nil {
		return responseError(c, errcode.InternalServerError)
	}

	deviceInfo, err := json.Marshal(i.DeviceInfo)
	if err != nil {
		return responseError(c, errcode.InternalServerError)
	}

	deviceInfoStr, err := utils.EncryptAES(string(deviceInfo))
	if err != nil {
		return responseError(c, errcode.InternalServerError)
	}

	device := model.NewUserDevice()
	device.UserId = user.Id
	device.DeviceId = i.DeviceId
	device.DeviceInfo = deviceInfoStr
	device.PNSToken = i.PNSToken

	if err := h.UserDevice.Create(device); err != nil {
		return responseError(c, errcode.InternalServerError)
	}

	token, err := user.GetAccessToken(device)
	if err != nil {
		return responseError(c, errcode.InternalServerError)
	}

	tk := model.NewToken()
	tk.ReferenceId = user.Id
	tk.DeviceId = i.DeviceId
	tk.AccessToken = token

	if err := h.Token.Create(tk); err != nil {
		return responseError(c, errcode.InternalServerError)
	}

	res := &model.UserWithToken{
		User:  user,
		Token: tk.AccessToken,
	}

	return responseJSON(c, res)
}
