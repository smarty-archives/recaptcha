package recaptcha

import (
	"bytes"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"testing"

	"github.com/smartystreets/assertions/should"
	"github.com/smartystreets/gunit"
)

func TestDefaultVerifierFixture(t *testing.T) {
	gunit.Run(new(DefaultVerifierFixture), t)
}

type DefaultVerifierFixture struct {
	*gunit.Fixture

	verifier *DefaultVerifier

	clientCalls          int
	clientRequest        *http.Request
	clientResponse       *http.Response
	clientError          error
	clientResponseBuffer *bytes.Buffer
}

func (this *DefaultVerifierFixture) Setup() {
	this.clientResponseBuffer = bytes.NewBuffer(nil)
	this.clientResponse = &http.Response{
		Body: ioutil.NopCloser(this.clientResponseBuffer),
	}

	this.verifier = NewVerifier()
	WithHTTPClient(this)(this.verifier)
}

func (this *DefaultVerifierFixture) TestEmptyTokenIsInvalid() {
	result, err := this.verifier.Verify("", "")

	this.So(result, should.BeFalse)
	this.So(err, should.BeNil)
}
func (this *DefaultVerifierFixture) TestTokenLookupIsPerformed() {
	WithSecret(func() string { return "my-secret" })(this.verifier)

	_, _ = this.verifier.Verify(" token ", "client-ip")

	this.So(this.clientCalls, should.Equal, 1)
	this.So(this.clientRequest.Method, should.Equal, http.MethodPost)
	this.So(this.clientRequest.Header.Get(contentTypeHeader), should.Equal, defaultContentType)
	this.So(this.clientRequest.URL.String(), should.Equal, defaultEndpoint)
	this.So(this.clientRequest.PostForm, should.Resemble, url.Values{
		"secret":   []string{"my-secret"},
		"response": []string{"token"},
		"remoteip": []string{"client-ip"},
	})
}
func (this *DefaultVerifierFixture) TestCustomEndpoint() {
	WithEndpoint("/custom-endpoint")(this.verifier)

	_, _ = this.verifier.Verify("token", "")

	this.So(this.clientCalls, should.Equal, 1)
	this.So(this.clientRequest.URL.String(), should.Equal, "/custom-endpoint")
}
func (this *DefaultVerifierFixture) TestDoNotSendEmptyClientIP() {
	WithSecret(func() string { return "my-secret" })(this.verifier)

	_, _ = this.verifier.Verify("token", "")

	this.So(this.clientRequest.PostForm, should.Resemble, url.Values{
		"secret":   []string{"my-secret"},
		"response": []string{"token"},
	})
}

func (this *DefaultVerifierFixture) TestConnectivityError() {
	this.clientError = errors.New("")

	result, err := this.verifier.Verify("token", "ip")

	this.So(result, should.BeFalse)
	this.So(err, should.Equal, ErrLookupFailure)
}
func (this *DefaultVerifierFixture) TestParsingError() {
	this.writeResponseBody("malformed json")

	result, err := this.verifier.Verify("token", "ip")

	this.So(result, should.BeFalse)
	this.So(err, should.Equal, ErrLookupFailure)
}

func (this *DefaultVerifierFixture) TestValidLookup() {
	this.writeResponseBody(`{"Score":1.0}`)

	result, err := this.verifier.Verify("token", "ip")

	this.So(result, should.BeTrue)
	this.So(err, should.BeNil)
}

func (this *DefaultVerifierFixture) TestRequiredThreshold() {
	this.writeResponseBody(`{}`)

	WithRequiredThreshold(0.1)(this.verifier)

	result, err := this.verifier.Verify("token", "ip")

	this.So(result, should.BeFalse)
	this.So(err, should.BeNil)
}
func (this *DefaultVerifierFixture) TestRequiredHostname() {
	this.writeResponseBody(`{"Score":1.0}`)

	WithAllowedHosts("hostname-required")(this.verifier)

	result, err := this.verifier.Verify("token", "ip")

	this.So(result, should.BeFalse)
	this.So(err, should.BeNil)
}
func (this *DefaultVerifierFixture) TestRequiredAction() {
	this.writeResponseBody(`{"Score":1.0}`)

	WithAllowedActions("action-required")(this.verifier)

	result, err := this.verifier.Verify("token", "ip")

	this.So(result, should.BeFalse)
	this.So(err, should.BeNil)
}

/* ------------------------------------------------------------------------------------------------------------------ */

func (this *DefaultVerifierFixture) Do(request *http.Request) (*http.Response, error) {
	_ = request.ParseForm()
	this.clientCalls++
	this.clientRequest = request
	return this.clientResponse, this.clientError
}
func (this *DefaultVerifierFixture) writeResponseBody(value string) {
	this.clientResponseBuffer.WriteString(value)
}
