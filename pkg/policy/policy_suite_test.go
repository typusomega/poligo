package policy_test

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/typusomega/poliGo/pkg/policy"
)

type PolicySuite struct {
	suite.Suite
}

func TestPolicy(t *testing.T) {
	suite.Run(t, new(PolicySuite))
}

func (it *PolicySuite) SetupTest() {

}

func HandleAllBasePolicy() policy.BasePolicy {
	return *policy.DefaultBasePolicy()
}
