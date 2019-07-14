package recaptcha

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/smartystreets/assertions/should"
	"github.com/smartystreets/gunit"
)

func TestWireupFixture(t *testing.T) {
	gunit.Run(new(WireupFixture), t)
}

type WireupFixture struct {
	*gunit.Fixture

	resolvedToken string
}

func (this *WireupFixture) TestUnrecognizedOptionShouldPanic() {
	this.So(func() { New(0) }, should.PanicWith, errBadOptionProvided)
}

func (this *WireupFixture) TestWireup() {
	const expectedToken = "expected-token"
	handler := New(
		WithInnerHandler(this),
		WithHTTPClient(this),
		WithTokenReader(func(*http.Request) string { return expectedToken }))

	handler.ServeHTTP(
		httptest.NewRecorder(),
		httptest.NewRequest(http.MethodGet, "/", nil))

	this.So(this.resolvedToken, should.Equal, expectedToken)
}

func (this *WireupFixture) Do(request *http.Request) (*http.Response, error) {
	this.resolvedToken = request.Form.Get("response")
	return nil, errors.New("")
}
func (this *WireupFixture) ServeHTTP(http.ResponseWriter, *http.Request) {}
