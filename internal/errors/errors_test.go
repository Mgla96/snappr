package errors

import (
	"errors"
	"fmt"
	"testing"
)

func TestNew_Error(t *testing.T) {
	exampleError := New("example error")
	exampleWrappedError := fmt.Errorf("wrapped error: %w", exampleError)

	tests := []struct {
		name          string
		e             error
		wantMsg       string
		wantErrorType error
	}{
		{
			name:          "example error",
			e:             exampleError,
			wantMsg:       "example error",
			wantErrorType: exampleError,
		},
		{
			name:          "example wrapped error",
			e:             exampleWrappedError,
			wantMsg:       "wrapped error: example error",
			wantErrorType: exampleError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.e.Error()
			if got != tt.wantMsg {
				t.Errorf("New.Error() = %v, want %v", got, tt.wantMsg)
			}
			if !errors.Is(tt.e, tt.wantErrorType) {
				t.Errorf("New.Error() = %v, want %v", tt.e, tt.wantErrorType)
			}
		})
	}
}
