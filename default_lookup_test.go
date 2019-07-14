package recaptcha

import (
	"testing"

	"github.com/smartystreets/assertions/should"
	"github.com/smartystreets/gunit"
)

func TestDefaultLookupFixture(t *testing.T) {
	gunit.Run(new(DefaultLookupFixture), t)
}

type DefaultLookupFixture struct {
	*gunit.Fixture
}

func (this *DefaultLookupFixture) TestAcceptedWhenScoreMeetsThreshold() {
	lookup := defaultLookup{Score: 0.5}

	result, err := lookup.IsValid(nil, nil, 0.5)

	this.So(result, should.BeTrue)
	this.So(err, should.BeNil)
}
func (this *DefaultLookupFixture) TestRejectedWhenScoreDoesNotMeetThreshold() {
	lookup := defaultLookup{Score: 0.5}

	result, err := lookup.IsValid(nil, nil, lookup.Score+0.1)

	this.So(result, should.BeFalse)
	this.So(err, should.BeNil)
}

func (this *DefaultLookupFixture) TestRejectWhenRequiredHostMissing() {
	lookup := defaultLookup{}
	allowedHosts := map[string]struct{}{"some-hostname": {}}

	result, err := lookup.IsValid(allowedHosts, nil, 0.0)

	this.So(result, should.BeFalse)
	this.So(err, should.BeNil)
}
func (this *DefaultLookupFixture) TestAcceptedWhenRequiredHostFound() {
	lookup := defaultLookup{Hostname: "some-hostname"}
	allowedHosts := map[string]struct{}{lookup.Hostname: {}}

	result, err := lookup.IsValid(allowedHosts, nil, 0.0)

	this.So(result, should.BeTrue)
	this.So(err, should.BeNil)
}

func (this *DefaultLookupFixture) TestRejectWhenRequiredActionMissing() {
	lookup := defaultLookup{}
	allowedActions := map[string]struct{}{"some-action": {}}

	result, err := lookup.IsValid(nil, allowedActions, 0.0)

	this.So(result, should.BeFalse)
	this.So(err, should.BeNil)
}
func (this *DefaultLookupFixture) TestAcceptedWhenRequiredActionFound() {
	lookup := defaultLookup{Action: "some-action"}
	allowedActions := map[string]struct{}{lookup.Action: {}}

	result, err := lookup.IsValid(nil, allowedActions, 0.0)

	this.So(result, should.BeTrue)
	this.So(err, should.BeNil)
}

func (this *DefaultLookupFixture) TestRejectedWhenTokenExpired() {
	lookup := defaultLookup{Errors: []string{expiredTokenMessage}}

	result, err := lookup.IsValid(nil, nil, 0.0)

	this.So(result, should.BeFalse)
	this.So(err, should.BeNil)
}
func (this *DefaultLookupFixture) TestServerErrors() {
	lookup := defaultLookup{Errors: []string{"other-error"}}

	result, err := lookup.IsValid(nil, nil, 0.0)

	this.So(result, should.BeFalse)
	this.So(err, should.Equal, ErrServerConfig)
}

func (this *DefaultLookupFixture) TestFullValidation() {
	lookup := defaultLookup{
		Score:    0.5,
		Action:   "some-action",
		Hostname: "some-hostname",
	}

	allowedHosts := map[string]struct{}{
		lookup.Hostname: {},
		"another-host":  {},
	}
	allowedActions := map[string]struct{}{
		lookup.Action:    {},
		"another-action": {},
	}

	result, err := lookup.IsValid(allowedHosts, allowedActions, lookup.Score)

	this.So(result, should.BeTrue)
	this.So(err, should.BeNil)
}
