package mage

import (
	"fmt"
	"io/ioutil"
	"regexp"
	"strings"

	"gopkg.in/yaml.v2"
)

const configFile = "deployments/panther_config.yml"

var (
	userEmailRegex = regexp.MustCompile(`UserEmail:[ \t]*'.*'`)
	userFirstRegex = regexp.MustCompile(`UserGivenName:[ \t]*'.*'`)
	userLastRegex  = regexp.MustCompile(`UserFamilyName:[ \t]*'.*'`)
)

// Open and parse a yaml file.
func loadYamlFile(path string) (map[string]interface{}, error) {
	contents, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open '%s': %s", path, err)
	}

	var result map[string]interface{}
	if err := yaml.Unmarshal(contents, &result); err != nil {
		return nil, fmt.Errorf("failed to parse yaml file '%s': %s", path, err)
	}

	return result, nil
}

// Flatten a map of parameter values from the config file into "key=val" strings.
func flattenParameterValues(params interface{}) []string {
	var result []string
	for key, val := range params.(map[interface{}]interface{}) {
		if val != nil {
			result = append(result, fmt.Sprintf("%s=%v", key, val))
		}
	}
	return result
}

// Prompt the user for required values if they aren't specified in the config file.
//
// Also writes the new values back into the file.
func promptRequiredValues(appParams map[interface{}]interface{}) error {
	var email, firstName, lastName string

	fmt.Println("Configuring the initial Panther admin user...")
	for {
		fmt.Printf("Email: ")
		if _, err := fmt.Scanln(&email); err != nil {
			fmt.Println(err) // empty line, for example
			continue
		}

		email = strings.TrimSpace(email)
		if !validEmail(email) {
			fmt.Printf("'%s' is not a valid email, please try again\n", email)
			continue
		}

		break
	}

	fmt.Printf("First name (optional): ")
	_, _ = fmt.Scanln(&firstName)
	firstName = strings.TrimSpace(firstName)

	fmt.Printf("Last name (optional): ")
	_, _ = fmt.Scanln(&lastName)
	lastName = strings.TrimSpace(lastName)

	// Write config values back to the file. The yaml parser doesn't preserve comments so
	// we do a simple find/replace instead of a yaml dump.
	if err := updateConfigFile(email, firstName, lastName); err != nil {
		return fmt.Errorf("failed to update config file %s: %s", configFile, err)
	}

	appParams["UserEmail"] = email
	appParams["UserGivenName"] = firstName
	appParams["UserFamilyName"] = lastName
	return nil
}

func updateConfigFile(email, firstName, lastName string) error {
	content, err := ioutil.ReadFile(configFile)
	if err != nil {
		return err
	}

	content = userEmailRegex.ReplaceAll(content, []byte("UserEmail: '"+email+"'"))
	content = userFirstRegex.ReplaceAll(content, []byte("UserGivenName: '"+firstName+"'"))
	content = userLastRegex.ReplaceAll(content, []byte("UserFamilyName: '"+lastName+"'"))

	return ioutil.WriteFile(configFile, content, 0644)
}

// Very simple email validation to prevent obvious mistakes.
func validEmail(email string) bool {
	return len(email) >= 5 && strings.Contains(email, "@") && strings.Contains(email, ".")
}
