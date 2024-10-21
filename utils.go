package majordomo_ai

import (
	"errors"
	"os"
	"strconv"
)

func setCredentialsFromEnv(c *Credentials) error {

	if c == nil {
		return errors.New("Invalid credentials passed")
	}

	val, isSet := os.LookupEnv("MAJORDOMO_AI_ACCOUNT")
	if !isSet {
		return errors.New("Environment variable MAJORDOMO_AI_ACCOUNT not set")
	}
	var err error
	c.AccountId, err = strconv.Atoi(val)
	if err != nil {
		return errors.New("Invalid format for MAJORDOMO_AI_ACCOUNT")
	}

	val, isSet = os.LookupEnv("MAJORDOMO_AI_WORKSPACE")
	if !isSet {
		return errors.New("Environment variable MAJORDOMO_AI_WORKSPACE not set")
	}
	c.Workspace = val

	val, isSet = os.LookupEnv("MAJORDOMO_AI_API_KEY")
	if !isSet {
		return errors.New("Environment variable MAJORDOMO_AI_API_KEY not set")
	}
	c.MdApiKey = val

	return nil
}
