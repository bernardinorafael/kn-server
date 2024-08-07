package service

import (
	"errors"
	"fmt"

	"github.com/bernardinorafael/kn-server/internal/core/application/contract"
	"github.com/bernardinorafael/kn-server/internal/core/domain/valueobj/phone"
	"github.com/bernardinorafael/kn-server/pkg/logger"
	"github.com/twilio/twilio-go"
	verify "github.com/twilio/twilio-go/rest/verify/v2"
)

const (
	smsCodeLength = 6
)

type twilioSMSService struct {
	log       logger.Logger
	serviceID string
	client    *twilio.RestClient
}

func NewTwilioSMSService(log logger.Logger, serviceID string, client *twilio.RestClient) contract.SMSVerifier {
	return &twilioSMSService{log, serviceID, client}
}

func (svc twilioSMSService) NotifySMS(to string) error {
	p, err := phone.New(to)
	if err != nil {
		svc.log.Error("phone validation error", "error", err.Error())
		return err
	}

	params := verify.CreateVerificationParams{}

	params.SetTo(fmt.Sprintf("+55%s", p.Phone()))
	params.SetChannel("sms")

	_, err = svc.client.VerifyV2.CreateVerification(svc.serviceID, &params)
	if err != nil {
		svc.log.Error("cannot send sms", "error", err.Error())
		return fmt.Errorf("error sending sms to %s", to)
	}

	return nil
}

func (svc twilioSMSService) ConfirmSMS(code string, to string) error {
	if len(code) != smsCodeLength {
		return errors.New("invalid code format")
	}

	p, err := phone.New(to)
	if err != nil {
		return err
	}

	params := verify.CreateVerificationCheckParams{}

	params.SetTo(fmt.Sprintf("+55%s", p.Phone()))
	params.SetCode(code)

	res, err := svc.client.VerifyV2.CreateVerificationCheck(svc.serviceID, &params)
	if err != nil {
		svc.log.Error("cannot verify code", "error", err.Error())
		return fmt.Errorf("error verifying code to %s", p.Phone())
	}
	status := *res.Status

	switch {
	case status == "pending":
		return errors.New("invalid otp code")
	case status == "canceled":
		return errors.New("verify otp operation canceled")
	case status == "max_attempts_reached":
		return errors.New("max attempts reached")
	case status == "expired":
		return errors.New("verification code expired")
	}

	return nil
}
