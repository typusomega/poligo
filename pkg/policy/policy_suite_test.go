package policy_test

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type PolicySuite struct {
	suite.Suite
}

func TestPolicy(t *testing.T) {
	suite.Run(t, new(PolicySuite))
}

func (it *PolicySuite) SetupTest() {

}
