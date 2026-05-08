package genome

import (
	"strings"
	"testing"

	"github.com/gmlewis/gep/v2/functions"
	"github.com/gmlewis/gep/v2/gene"
	"github.com/gmlewis/gep/v2/grammars"
)

func TestGenerateCodeReportsMissingLinkFunc(t *testing.T) {
	d := &dump{
		gr: &grammars.Grammar{},
		genome: &Genome{
			Genes:    []*gene.Gene{gene.New("d0", functions.Bool)},
			LinkFunc: "MissingLink",
		},
		subs: map[string]string{
			"CHARX": "X",
		},
	}

	_, err := d.generateCode()
	if err == nil || !strings.Contains(err.Error(), "MissingLink") {
		t.Fatalf("generateCode() error = %v, want missing link function error", err)
	}
}
