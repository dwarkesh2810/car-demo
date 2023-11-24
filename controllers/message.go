package controllers

import (
	"car_demo/conf"
	"car_demo/helper"
	"encoding/json"
	"net/http"

	beego "github.com/beego/beego/v2/server/web"
	"github.com/twilio/twilio-go"
	twilioApi "github.com/twilio/twilio-go/rest/verify/v2"
)

type MsgController struct {
	beego.Controller
}

type OTPData struct {
	PhoneNumber string `json:"phone_number"`
}

type VerifyData struct {
	User string `json:"user,omitempty"`
	Code string `json:"code,omitempty"`
}

func (c *MsgController) URLMapping() {
	c.Mapping("SendSMS", c.SendSMS)
	c.Mapping("VerifySMS", c.VerifySMS)
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

func (app *MsgController) SendSMS() {
	var v OTPData
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

func (app *MsgController) VerifySMS() {
	var v VerifyData

	err := json.Unmarshal(app.Ctx.Input.RequestBody, &v)

	if err != nil {
		return
	}

	newData := VerifyData{
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
