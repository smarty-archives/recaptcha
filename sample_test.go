package recaptcha

import (
	"testing"

	"github.com/smartystreets/gunit"
)

func TestSampleFixture(t *testing.T) {
	gunit.Run(new(SampleFixture), t)
}

type SampleFixture struct {
	*gunit.Fixture
}

func (this *SampleFixture) Setup() {
}

func (this *SampleFixture) Test() {
}
