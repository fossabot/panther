package email

/**
 * Copyright 2020 Panther Labs Inc
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *    http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

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
