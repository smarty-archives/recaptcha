package checkpoint

type defaultLookup struct {
	Score    float32  `json:"score"`
	Action   string   `json:"action"`
	Hostname string   `json:"hostname"`
	Errors   []string `json:"errors"`
}

func (this defaultLookup) IsValid(allowedHosts, allowedActions map[string]struct{}, requiredThreshold float32) (bool, error) {
	result, err := this.tokenExists()
	return result &&
		this.meetsRequiredThreshold(requiredThreshold) &&
		this.hasAllowedHost(allowedHosts) &&
		this.hasAllowedAction(allowedActions), err
}

func (this defaultLookup) tokenExists() (bool, error) {
	for _, item := range this.Errors {
		if item == expiredTokenMessage {
			return false, nil
		} else {
			return false, ErrServerConfig
		}
	}

	return true, nil
}

func (this defaultLookup) meetsRequiredThreshold(threshold float32) bool {
	return this.Score >= threshold
}

func (this defaultLookup) hasAllowedHost(allowed map[string]struct{}) bool {
	return isValueAllowed(this.Hostname, allowed)
}

func (this defaultLookup) hasAllowedAction(allowed map[string]struct{}) bool {
	return isValueAllowed(this.Action, allowed)
}

func isValueAllowed(value string, allowed map[string]struct{}) bool {
	if len(allowed) == 0 {
		return true
	}

	_, found := allowed[value]
	return found
}

// Error Code Reference: https://developers.google.com/recaptcha/docs/verify
const expiredTokenMessage = "timeout-or-duplicate"
