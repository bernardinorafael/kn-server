package twilioclient

import (
	env "github.com/bernardinorafael/kn-server/internal/config"
	"github.com/twilio/twilio-go"
)

func New(env *env.Env) *twilio.RestClient {
	return twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: env.TwilioAccountID,
		Password: env.TwilioAuthToken,
	})
}
