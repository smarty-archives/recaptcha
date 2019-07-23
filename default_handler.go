package recaptcha

import "net/http"

type DefaultHandler struct {
	inner          http.Handler
	verifier       TokenVerifier
	token          func(*http.Request) string
	clientIP       func(*http.Request) string
	rejectedStatus int
	errorStatus    int
}

func NewHandler(verifier TokenVerifier, options ...HandlerOption) *DefaultHandler {
	this := &DefaultHandler{verifier: verifier}

	WithTokenReader(defaultTokenReader)(this)
	WithClientIPReader(defaultClientIPReader)(this)
	WithRejectedStatus(defaultRejectedStatus)(this)
	WithErrorStatus(defaultErrorStatus)(this)

	for _, option := range options {
		option(this)
	}

	return this
}

func (this *DefaultHandler) Install(inner http.Handler) {
	this.inner = inner
}

func (this *DefaultHandler) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	result, err := this.verify(request)

	if result || err == ErrLookupFailure {
		this.inner.ServeHTTP(response, request)
	} else if err != nil {
		writeResponse(response, this.errorStatus)
	} else {
		writeResponse(response, this.rejectedStatus)
	}
}
func (this *DefaultHandler) verify(request *http.Request) (bool, error) {
	token := this.token(request)
	clientIP := this.clientIP(request)
	return this.verifier.Verify(token, clientIP)
}
func writeResponse(response http.ResponseWriter, statusCode int) {
	http.Error(response, http.StatusText(statusCode), statusCode)
}

/* ------------------------------------------------------------------------------------------------------------------ */

type HandlerOption func(*DefaultHandler)

func WithTokenReader(callback func(*http.Request) string) HandlerOption {
	return func(this *DefaultHandler) { this.token = callback }
}
func WithClientIPReader(callback func(*http.Request) string) HandlerOption {
	return func(this *DefaultHandler) { this.clientIP = callback }
}
func WithRejectedStatus(value int) HandlerOption {
	return func(this *DefaultHandler) { this.rejectedStatus = value }
}
func WithErrorStatus(value int) HandlerOption {
	return func(this *DefaultHandler) { this.errorStatus = value }
}
func WithInnerHandler(value http.Handler) HandlerOption {
	return func(this *DefaultHandler) { this.inner = value }
}

func defaultTokenReader(request *http.Request) string {
	return request.URL.Query().Get(DefaultFormTokenName)
}
func defaultClientIPReader(request *http.Request) string {
	return request.RemoteAddr
}

/* ------------------------------------------------------------------------------------------------------------------ */

const (
	DefaultFormTokenName  = "g-recaptcha-response"
	defaultRejectedStatus = http.StatusForbidden
	defaultErrorStatus    = http.StatusInternalServerError
)
