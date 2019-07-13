package checkpoint

import "net/http"

type DefaultHandler struct {
	inner    http.Handler
	verifier TokenVerifier
}

func NewHandler(verifier TokenVerifier) *DefaultHandler {
	return &DefaultHandler{verifier: verifier}
}

func (this *DefaultHandler) Install(inner http.Handler) {
	this.inner = inner
}

func (this *DefaultHandler) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	token := request.Form.Get(defaultFormTokenName)
	clientIP := request.RemoteAddr

	result, err := this.verifier.Verify(token, clientIP)
	if !result && err == nil {
		http.Error(response, http.StatusText(defaultRejectedStatus), defaultRejectedStatus)
	} else if err == ErrServerConfig {
		http.Error(response, http.StatusText(defaultErrorStatus), defaultErrorStatus)
	} else {
		this.inner.ServeHTTP(response, request)
	}
}

const (
	defaultFormTokenName  = "g-recaptcha-response"
	defaultRejectedStatus = http.StatusForbidden
	defaultErrorStatus    = http.StatusInternalServerError
)
