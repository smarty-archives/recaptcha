package recaptcha

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/smartystreets/assertions/should"
	"github.com/smartystreets/gunit"
)

func TestDefaultHandlerFixture(t *testing.T) {
	gunit.Run(new(DefaultHandlerFixture), t)
}

type DefaultHandlerFixture struct {
	*gunit.Fixture

	request  *http.Request
	response *httptest.ResponseRecorder

	handler *DefaultHandler

	innerRequest  *http.Request
	innerResponse http.ResponseWriter
	innerCalls    int

	verifiedToken    string
	verifiedClientIP string
	verifyResult     bool
	verifyError      error
}

func (this *DefaultHandlerFixture) Setup() {
	this.request, _ = http.NewRequest(http.MethodGet, "/some-path/", nil)
	this.response = httptest.NewRecorder()
	this.handler = NewHandler(this)
	this.handler.Install(this)
	this.verifyResult = true
}

func (this *DefaultHandlerFixture) TestInnerHandlerCalled() {
	this.handler.ServeHTTP(this.response, this.request)

	this.assertInnerCalled()
}

func (this *DefaultHandlerFixture) TestBadTokenRequestRejected() {
	this.verifyResult = false

	this.handler.ServeHTTP(this.response, this.request)

	this.assertInnerNotCalled()
	this.assertResponse(defaultRejectedStatus)
}

func (this *DefaultHandlerFixture) TestConfigurationErrorRequestRejected() {
	this.verifyResult = false
	this.verifyError = ErrServerConfig

	this.handler.ServeHTTP(this.response, this.request)

	this.assertInnerNotCalled()
	this.assertResponse(defaultErrorStatus)
}

func (this *DefaultHandlerFixture) TestLookupFailureRequestAllowed() {
	this.verifyResult = false
	this.verifyError = ErrLookupFailure

	this.handler.ServeHTTP(this.response, this.request)

	this.assertInnerCalled()
}

func (this *DefaultHandlerFixture) TestTokenAndClientIPReadFromPOSTRequest() {
	this.request, _ = http.NewRequest(http.MethodPost, "/", ioutil.NopCloser(strings.NewReader(url.Values{
		DefaultFormTokenName: []string{"my-token"},
	}.Encode())))

	this.request.Header.Set("Content-Type", defaultContentType)
	this.request.RemoteAddr = "1.2.3.4"

	this.handler.ServeHTTP(this.response, this.request)

	this.So(this.verifiedToken, should.Equal, "my-token")
	this.So(this.verifiedClientIP, should.Equal, "1.2.3.4")
}

func (this *DefaultHandlerFixture) TestTokenAndClientIPReadFromGETRequest() {
	this.request, _ = http.NewRequest(http.MethodGet, fmt.Sprintf("/?%s=my-token", DefaultFormTokenName), nil)
	this.request.RemoteAddr = "1.2.3.4"

	this.handler.ServeHTTP(this.response, this.request)

	this.So(this.verifiedToken, should.Equal, "my-token")
	this.So(this.verifiedClientIP, should.Equal, "1.2.3.4")
}

func (this *DefaultHandlerFixture) TestAlternateTokenReader() {
	this.request.Header.Set("read-from-different-location", "my-token")

	WithTokenReader(func(request *http.Request) string {
		return this.request.Header.Get("read-from-different-location")
	})(this.handler)

	this.handler.ServeHTTP(this.response, this.request)

	this.So(this.verifiedToken, should.Equal, "my-token")
}

func (this *DefaultHandlerFixture) TestAlternateClientIPReader() {
	this.request.Header.Set("read-from-different-location", "1.2.3.4")

	WithClientIPReader(func(request *http.Request) string {
		return this.request.Header.Get("read-from-different-location")
	})(this.handler)

	this.handler.ServeHTTP(this.response, this.request)

	this.So(this.verifiedClientIP, should.Equal, "1.2.3.4")
}

func (this *DefaultHandlerFixture) TestAlternateRejectedResponseStatus() {
	this.verifyResult = false
	WithRejectedStatus(http.StatusTooManyRequests)(this.handler)

	this.handler.ServeHTTP(this.response, this.request)

	this.assertResponse(http.StatusTooManyRequests)
}

func (this *DefaultHandlerFixture) TestAlternateErrorResponseStatus() {
	this.verifyResult = false
	this.verifyError = ErrServerConfig
	WithErrorStatus(http.StatusBadGateway)(this.handler)

	this.handler.ServeHTTP(this.response, this.request)

	this.assertResponse(http.StatusBadGateway)
}

/* ------------------------------------------------------------------------------------------------------------------ */

func (this *DefaultHandlerFixture) assertInnerCalled() {
	this.So(this.innerRequest, should.Equal, this.request)
	this.So(this.innerResponse, should.Equal, this.response)
	this.So(this.innerCalls, should.Equal, 1)
}
func (this *DefaultHandlerFixture) assertInnerNotCalled() {
	this.So(this.innerRequest, should.BeNil)
	this.So(this.innerResponse, should.BeNil)
	this.So(this.innerCalls, should.BeZeroValue)
}
func (this *DefaultHandlerFixture) assertResponse(statusCode int) {
	this.So(this.response.Code, should.Equal, statusCode)
	this.So(this.response.Body.String(), should.Equal, http.StatusText(statusCode)+"\n")
}

/* ------------------------------------------------------------------------------------------------------------------ */

func (this *DefaultHandlerFixture) Verify(token, clientIP string) (bool, error) {
	this.verifiedToken = token
	this.verifiedClientIP = clientIP
	return this.verifyResult, this.verifyError
}

func (this *DefaultHandlerFixture) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	this.innerRequest = request
	this.innerResponse = response
	this.innerCalls++
}
