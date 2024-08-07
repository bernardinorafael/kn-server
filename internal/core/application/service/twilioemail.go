package service

import (
	"errors"
	"fmt"

	"github.com/bernardinorafael/kn-server/internal/core/application/contract"
	"github.com/bernardinorafael/kn-server/internal/core/domain/valueobj/email"
	"github.com/bernardinorafael/kn-server/pkg/logger"
	"github.com/twilio/twilio-go"
	verify "github.com/twilio/twilio-go/rest/verify/v2"
)

type twilioEmailService struct {
	log       logger.Logger
	serviceID string
	client    *twilio.RestClient
}

func NewTwilioEmailService(log logger.Logger, serviceID string, client *twilio.RestClient) contract.EmailVerifier {
	return &twilioEmailService{log, serviceID, client}
}

func (svc twilioEmailService) NotifyEmail(to string) error {
	addr, err := email.New(to)
	if err != nil {
		svc.log.Error("email validation error", "error", err.Error())
		return err
	}

	params := verify.CreateVerificationParams{}

	params.SetTo(string(addr.Email()))
	params.SetChannel("email")

	_, err = svc.client.VerifyV2.CreateVerification(svc.serviceID, &params)
	if err != nil {
		svc.log.Error("cannot send email verify", "error", err.Error())
		return fmt.Errorf("error sending email to %s", to)
	}

	return nil
}

func (svc twilioEmailService) ConfirmEmail(code string, sent string) error {
	if len(code) != smsCodeLength {
		return errors.New("invalid code format")
	}

	addr, err := email.New(sent)
	if err != nil {
		svc.log.Error("email validation error", "error", err.Error())
		return err
	}

	params := verify.CreateVerificationCheckParams{}

	params.SetTo(string(addr.Email()))
	params.SetCode(code)

	res, err := svc.client.VerifyV2.CreateVerificationCheck(svc.serviceID, &params)
	if err != nil {
		svc.log.Error("cannot verify code", "error", err.Error())
		return fmt.Errorf("error verifying code to %s", sent)
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
