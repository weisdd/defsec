package elasticsearch

import (
	"testing"

	"github.com/aquasecurity/defsec/parsers/types"
	"github.com/aquasecurity/defsec/providers/aws/elasticsearch"
	"github.com/aquasecurity/defsec/rules"
	"github.com/aquasecurity/defsec/state"
	"github.com/stretchr/testify/assert"
)

func TestCheckEnforceHttps(t *testing.T) {
	tests := []struct {
		name     string
		input    elasticsearch.Elasticsearch
		expected bool
	}{
		{
			name: "Elasticsearch domain with enforce HTTPS disabled",
			input: elasticsearch.Elasticsearch{
				Metadata: types.NewTestMetadata(),
				Domains: []elasticsearch.Domain{
					{
						Metadata: types.NewTestMetadata(),
						Endpoint: elasticsearch.Endpoint{
							Metadata:     types.NewTestMetadata(),
							EnforceHTTPS: types.Bool(false, types.NewTestMetadata()),
						},
					},
				},
			},
			expected: true,
		},
		{
			name: "Elasticsearch domain with enforce HTTPS enabled",
			input: elasticsearch.Elasticsearch{
				Metadata: types.NewTestMetadata(),
				Domains: []elasticsearch.Domain{
					{
						Metadata: types.NewTestMetadata(),
						Endpoint: elasticsearch.Endpoint{
							Metadata:     types.NewTestMetadata(),
							EnforceHTTPS: types.Bool(true, types.NewTestMetadata()),
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
			testState.AWS.Elasticsearch = test.input
			results := CheckEnforceHttps.Evaluate(&testState)
			var found bool
			for _, result := range results {
				if result.Status() == rules.StatusFailed && result.Rule().LongID() == CheckEnforceHttps.Rule().LongID() {
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
