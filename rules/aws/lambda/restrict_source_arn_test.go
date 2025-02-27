package lambda

import (
	"testing"

	"github.com/aquasecurity/defsec/parsers/types"
	"github.com/aquasecurity/defsec/providers/aws/lambda"
	"github.com/aquasecurity/defsec/rules"
	"github.com/aquasecurity/defsec/state"
	"github.com/stretchr/testify/assert"
)

func TestCheckRestrictSourceArn(t *testing.T) {
	tests := []struct {
		name     string
		input    lambda.Lambda
		expected bool
	}{
		{
			name: "Lambda function permission missing source ARN",
			input: lambda.Lambda{
				Metadata: types.NewTestMetadata(),
				Functions: []lambda.Function{
					{
						Metadata: types.NewTestMetadata(),
						Permissions: []lambda.Permission{
							{
								Metadata:  types.NewTestMetadata(),
								Principal: types.String("sns.amazonaws.com", types.NewTestMetadata()),
								SourceARN: types.String("", types.NewTestMetadata()),
							},
						},
					},
				},
			},
			expected: true,
		},
		{
			name: "Lambda function permission with source ARN",
			input: lambda.Lambda{
				Metadata: types.NewTestMetadata(),
				Functions: []lambda.Function{
					{
						Metadata: types.NewTestMetadata(),
						Permissions: []lambda.Permission{
							{
								Metadata:  types.NewTestMetadata(),
								Principal: types.String("sns.amazonaws.com", types.NewTestMetadata()),
								SourceARN: types.String("source-arn", types.NewTestMetadata()),
							},
						},
					},
				},
			},
			expected: false,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var testState state.State
			testState.AWS.Lambda = test.input
			results := CheckRestrictSourceArn.Evaluate(&testState)
			var found bool
			for _, result := range results {
				if result.Status() == rules.StatusFailed && result.Rule().LongID() == CheckRestrictSourceArn.Rule().LongID() {
					found = true
				}
			}
			if test.expected {
				assert.True(t, found, "Rule should have been found")
			} else {
				assert.False(t, found, "Rule should not have been found")
			}
		})
	}
}
