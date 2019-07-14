package recaptcha

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strings"
)

type DefaultVerifier struct {
	secret    func() string
	client    httpClient
	endpoint  string
	threshold float32
	hosts     map[string]struct{}
	actions   map[string]struct{}
}

func NewVerifier(options ...VerifierOption) *DefaultVerifier {
	this := &DefaultVerifier{}

	WithSecret(func() string { return "" })(this)
	WithHTTPClient(http.DefaultClient)(this)
	WithEndpoint(defaultEndpoint)(this)
	WithRequiredThreshold(defaultThreshold)(this)
	WithAllowedHosts()(this)
	WithAllowedActions()(this)

	for _, option := range options {
		option(this)
	}

	return this
}

func (this *DefaultVerifier) Verify(token, clientIP string) (bool, error) {
	token = strings.TrimSpace(token)
	if len(token) == 0 {
		return false, nil
	} else {
		return this.verify(token, clientIP)
	}
}
func (this *DefaultVerifier) verify(token, clientIP string) (bool, error) {
	if response, err := this.newRequest(token, clientIP); err != nil {
		return false, ErrLookupFailure
	} else if lookup, err := this.parseLookup(response); err != nil {
		return false, ErrLookupFailure
	} else {
		return lookup.IsValid(this.hosts, this.actions, this.threshold)
	}
}
func (this *DefaultVerifier) newRequest(token, clientIP string) (*http.Response, error) {
	request, _ := http.NewRequest(http.MethodPost, this.endpoint, nil)
	request.Form = url.Values{
		"secret":   []string{this.secret()},
		"response": []string{token},
	}

	if len(clientIP) > 0 {
		request.Form.Set("remoteip", clientIP)
	}

	return this.client.Do(request)
}
func (this *DefaultVerifier) parseLookup(response *http.Response) (lookup defaultLookup, err error) {
	defer func() { _ = response.Body.Close() }()
	err = json.NewDecoder(response.Body).Decode(&lookup)
	return lookup, err
}

/* ------------------------------------------------------------------------------------------------------------------ */

type VerifierOption func(this *DefaultVerifier)

func WithSecret(callback func() string) VerifierOption {
	return func(this *DefaultVerifier) { this.secret = callback }
}
func WithHTTPClient(value httpClient) VerifierOption {
	return func(this *DefaultVerifier) { this.client = value }
}
func WithEndpoint(value string) VerifierOption {
	return func(this *DefaultVerifier) { this.endpoint = value }
}
func WithRequiredThreshold(value float32) VerifierOption {
	return func(this *DefaultVerifier) { this.threshold = value }
}
func WithAllowedHosts(values ...string) VerifierOption {
	return func(this *DefaultVerifier) { this.hosts = createMap(values) }
}
func WithAllowedActions(values ...string) VerifierOption {
	return func(this *DefaultVerifier) { this.actions = createMap(values) }
}
func createMap(values []string) map[string]struct{} {
	allowed := make(map[string]struct{}, len(values))
	for _, value := range values {
		allowed[value] = struct{}{}
	}
	return allowed
}

const (
	defaultEndpoint  = "https://www.google.com/recaptcha/api/siteverify"
	defaultThreshold = 0.3
)

/* ------------------------------------------------------------------------------------------------------------------ */

type httpClient interface {
	Do(*http.Request) (*http.Response, error)
}
