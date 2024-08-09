package parser

import (
	"context"
	"testing"

	"github.com/environment-toolkit/et/pkg/resolver"
)

func Test_Parse(t *testing.T) {
	ctx := context.TODO()
	environmentVariables := map[string]string{
		"template": "abc",
	}
	resolverSession := resolver.NewSession(ctx, environmentVariables)

	wants := `string {with} abc.`
	input := `string {with} ${{ env:template }}.`
	out, err := New(resolverSession).Parse(input)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if wants != out {
		t.Errorf("unexpected output error: wanted %s ; got %s", wants, out)
		return
	}
}
