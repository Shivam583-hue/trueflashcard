package connectapi

import (
	"errors"
	"testing"

	"connectrpc.com/connect"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestToConnectErrorMapsStatusCodes(t *testing.T) {
	cases := []struct {
		name string
		in   codes.Code
		want connect.Code
	}{
		{"not found", codes.NotFound, connect.CodeNotFound},
		{"invalid argument", codes.InvalidArgument, connect.CodeInvalidArgument},
		{"unauthenticated", codes.Unauthenticated, connect.CodeUnauthenticated},
		{"internal", codes.Internal, connect.CodeInternal},
		{"unknown maps to internal", codes.DataLoss, connect.CodeInternal},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := toConnectError(status.Error(tc.in, "boom"))
			if connect.CodeOf(got) != tc.want {
				t.Fatalf("code = %v, want %v", connect.CodeOf(got), tc.want)
			}
		})
	}
}

func TestToConnectErrorNonStatusIsInternal(t *testing.T) {
	got := toConnectError(errors.New("plain error"))
	if connect.CodeOf(got) != connect.CodeInternal {
		t.Fatalf("code = %v, want Internal", connect.CodeOf(got))
	}
}
