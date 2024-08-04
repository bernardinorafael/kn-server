package service

import (
	"fmt"

	"github.com/bernardinorafael/kn-server/internal/core/application/contract"
	"github.com/bernardinorafael/kn-server/pkg/logger"
	"github.com/twilio/twilio-go"
	verify "github.com/twilio/twilio-go/rest/verify/v2"
)

type twilioSMSService struct {
	log       logger.Logger
	serviceID string
	client    *twilio.RestClient
}

func NewTwilioSMSService(log logger.Logger, serviceID string, client *twilio.RestClient) contract.SMSNotifier {
	return &twilioSMSService{log, serviceID, client}
}

func (svc twilioSMSService) Notify(to string) error {
	params := verify.CreateVerificationParams{}

	params.SetTo(fmt.Sprintf("+55%s", to))
	params.SetChannel("sms")

	_, err := svc.client.VerifyV2.CreateVerification(svc.serviceID, &params)
	if err != nil {
		svc.log.Error("cannot send sms", "error", err.Error())
		return fmt.Errorf("error sending sms to %s", to)
	}

	return nil
}

func (svc twilioSMSService) Confirm(code string, phone string) (status string, err error) {
	params := verify.CreateVerificationCheckParams{}

	params.SetTo(fmt.Sprintf("+55%s", phone))
	params.SetCode(code)

	res, err := svc.client.VerifyV2.CreateVerificationCheck(svc.serviceID, &params)
	if err != nil {
		svc.log.Error("cannot verify code", "error", err.Error())
		return "", fmt.Errorf("error verifying code to %s", phone)
	}

	return *res.Status, nil
}
