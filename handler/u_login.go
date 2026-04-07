package handler

import (
	"encoding/json"
	"pkm/errcode"
	"pkm/model"
	"pkm/utils"

	"github.com/labstack/echo/v4"
)

func (h *Handler) Login(c echo.Context) error {
	var i struct {
		Email      string            `json:"email" validate:"required"`
		Platform   string            `json:"platform"`
		DeviceId   string            `json:"deviceId"`
		DeviceInfo map[string]string `json:"deviceInfo"`
		PNSToken   string            `json:"pnsToken"`
	}

	if msg, err := utils.ValidateRequest(c, &i); err != nil {
		return responseValidationError(c, msg)
	}

	user, err := h.User.GetByEmail(i.Email)
	if err != nil {
		return responseError(c, errcode.InternalServerError)
	}

	device, err := h.UserDevice.FindByDeviceId(i.DeviceId)
	if err != nil {
		return responseError(c, errcode.InternalServerError)
	}

	if err := h.UserDevice.UpdateByPnsToken(i.PNSToken); err != nil {
		return responseError(c, errcode.InternalServerError)
	}

	if device.Id == "" {
		device = model.NewUserDevice()
	}

	deviceInfo, err := json.Marshal(i.DeviceInfo)
	if err != nil {
		return responseError(c, errcode.InternalServerError)
	}

	deviceInfoStr, err := utils.EncryptAES(string(deviceInfo))
	if err != nil {
		return responseError(c, errcode.InternalServerError)
	}

	device.UserId = user.Id
	device.Platform = model.UserDevicePlatform(i.Platform)
	device.DeviceId = i.DeviceId
	device.DeviceInfo = deviceInfoStr
	device.PNSToken = i.PNSToken
	device.UpdateDt()

	if err := h.UserDevice.Upsert(device); err != nil {
		return responseError(c, errcode.InternalServerError)
	}

	tk, err := h.Token.FindByReferenceIdAndDeviceId(user.Id, device.DeviceId)
	if err != nil {
		return responseError(c, errcode.InternalServerError)
	}

	if tk.Id == "" {
		token, err := user.GetAccessToken(device)
		// token, err := user.GetAccessToken(device, model.RoleTenant)
		if err != nil {
			return responseError(c, errcode.InternalServerError)
		}

		tk = model.NewToken()
		tk.ReferenceId = user.Id
		tk.DeviceId = i.DeviceId
		tk.AccessToken = token

		if err := h.Token.Create(tk); err != nil {
			return responseError(c, errcode.InternalServerError)
		}
	}

	res := &model.UserWithToken{
		User:  user,
		Token: tk.AccessToken,
	}

	return responseJSON(c, res)
}
