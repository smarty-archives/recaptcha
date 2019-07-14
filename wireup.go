package recaptcha

import "errors"

func New(options ...interface{}) *DefaultHandler {
	var handlerOptions []HandlerOption
	var verifierOptions []VerifierOption

	for _, option := range options {
		if handlerOption, ok := option.(HandlerOption); ok {
			handlerOptions = append(handlerOptions, handlerOption)
		} else if verifierOption, ok := option.(VerifierOption); ok {
			verifierOptions = append(verifierOptions, verifierOption)
		} else {
			panic(errBadOptionProvided)
		}
	}

	verifier := NewVerifier(verifierOptions...)
	return NewHandler(verifier, handlerOptions...)
}

var errBadOptionProvided = errors.New("configuration option provided was not recognized")
