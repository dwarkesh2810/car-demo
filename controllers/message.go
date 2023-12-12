package controllers

import (
	"car_demo/conf"
	"car_demo/helper"
	"car_demo/models"
	"encoding/json"
	"net/http"

	beego "github.com/beego/beego/v2/server/web"
	"github.com/twilio/twilio-go"
	twilioApi "github.com/twilio/twilio-go/rest/verify/v2"
)

type MsgController struct {
	beego.Controller
}

var AccountSID string
var AuthToken string
var ServiceID string
var Client *twilio.RestClient

func TwilioSendOTP(phoneNumber string) (string, error) {
	AccountSID = conf.EnvConfig.TwilioAccountSID
	AuthToken = conf.EnvConfig.TwilioAuthToken
	ServiceID = conf.EnvConfig.TwilioServiceSID
	Client = twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: AccountSID,
		Password: AuthToken,
	})
	params := &twilioApi.CreateVerificationParams{}
	params.SetTo(phoneNumber)
	params.SetChannel("sms")

	resp, err := Client.VerifyV2.CreateVerification(ServiceID, params)
	if err != nil {
		return "", err
	}

	return *resp.Sid, nil
}

func TwilioVerifyOTP(phoneNumber string, code string) error {
	params := &twilioApi.CreateVerificationCheckParams{}
	params.SetTo(phoneNumber)
	params.SetCode(code)

	resp, err := Client.VerifyV2.CreateVerificationCheck(ServiceID, params)
	if err != nil {
		return err
	} else if *resp.Status == "approved" {
		return nil
	}

	return nil
}

// SendSMS ...
// @Title SendSMS
// @Description SendSMS
// @Param	body		body 	OTPData	true		"body for Message content"
// @Success 200 {int} string
// @Failure 403 body is empty
// @Failure 400 body is empty
// @router /otp [post]
func (app *MsgController) SendSMS() {
	var v models.OTP
	err := json.Unmarshal(app.Ctx.Input.RequestBody, &v)

	if err != nil {
		return
	}

	_, err = TwilioSendOTP(v.PhoneNumber)
	if err != nil {
		helper.JsonResponse(app.Controller, http.StatusBadRequest, 0, nil, err.Error())
		return
	}

	helper.JsonResponse(app.Controller, http.StatusOK, 1, "OTP send succesfull", "")
}

// VerifySMS ...
// @Title VerifySMS
// @Description VerifySMS
// @Param	body		body 	VerifyData	true		"body for Message content"
// @Success 200 {int} string
// @Failure 403 body is empty
// @Failure 400 body is empty
// @router /verify [post]
func (app *MsgController) VerifySMS() {
	var v models.VerifyData

	err := json.Unmarshal(app.Ctx.Input.RequestBody, &v)

	if err != nil {
		return
	}

	newData := models.VerifyData{
		User: v.User,
		Code: v.Code,
	}

	err = TwilioVerifyOTP(newData.User, newData.Code)
	if err != nil {
		helper.JsonResponse(app.Controller, http.StatusBadRequest, 0, nil, err.Error())
		return
	}

	helper.JsonResponse(app.Controller, http.StatusOK, 1, "OTP verify succesfull", "")
}
