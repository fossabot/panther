package email

import (
	"os"
	"strconv"
	"time"

	"github.com/matcornic/hermes"
)

var (
	// The logo is fetched from panther-public cloudfront CDN
	pantherEmailLogo = "https://d14d54mfia7r7w.cloudfront.net/panther-email-logo-white.png"
	appDomainURL     = os.Getenv("APP_DOMAIN_URL")
	// PantherEmailTemplate is used as a boilerplate for Panther themed email
	PantherEmailTemplate = hermes.Hermes{
		Theme: new(hermes.Flat),
		Product: hermes.Product{
			// Appears in header & footer of e-mails
			Name:      "Panther",
			Link:      appDomainURL,
			Copyright: "Copyright © " + strconv.Itoa(time.Now().Year()) + " Panther Labs Inc. All rights reserved.",
			Logo:      pantherEmailLogo,
		},
	}
)
