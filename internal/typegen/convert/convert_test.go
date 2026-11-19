package convert

import (
	"errors"
	"testing"
)

type fakeConverter struct {
	out string
	err error
}

func (f fakeConverter) Convert(string) (string, error) {
	return f.out, f.err
}

func TestConvertService_Convert(t *testing.T) {
	errBoom := errors.New("boom")
	svc := NewConvertService(map[Route]Converter{
		{From: LangGo, To: LangTS}: fakeConverter{out: "export interface A {}"},
		{From: LangTS, To: LangGo}: fakeConverter{err: errBoom},
	})

	tests := []struct {
		name     string
		from     Lang
		to       Lang
		expected string
		wantErr  error
	}{
		{
			name:     "known route returns converter output",
			from:     LangGo,
			to:       LangTS,
			expected: "export interface A {}",
		},
		{
			name:    "converter error is wrapped",
			from:    LangTS,
			to:      LangGo,
			wantErr: errBoom,
		},
		{
			name:    "unknown route",
			from:    LangGo,
			to:      LangGo,
			wantErr: ErrUnsupportedRoute,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := svc.Convert(tt.from, tt.to, "src")
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("Convert() error = %v, want %v", err, tt.wantErr)
			}

			if got != tt.expected {
				t.Errorf("Convert() = %q, want %q", got, tt.expected)
			}
		})
	}
}
