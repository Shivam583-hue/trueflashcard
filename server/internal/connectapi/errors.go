package connectapi

import (
	"context"
	"errors"

	"connectrpc.com/connect"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func invoke[Req, Resp any](
	ctx context.Context,
	req *connect.Request[Req],
	fn func(context.Context, *Req) (*Resp, error),
) (*connect.Response[Resp], error) {
	resp, err := fn(ctx, req.Msg)
	if err != nil {
		return nil, toConnectError(err)
	}
	return connect.NewResponse(resp), nil
}

func toConnectError(err error) error {
	st, ok := status.FromError(err)
	if !ok {
		return connect.NewError(connect.CodeInternal, err)
	}
	return connect.NewError(connectCode(st.Code()), errors.New(st.Message()))
}

func connectCode(code codes.Code) connect.Code {
	switch code {
	case codes.InvalidArgument:
		return connect.CodeInvalidArgument
	case codes.NotFound:
		return connect.CodeNotFound
	case codes.AlreadyExists:
		return connect.CodeAlreadyExists
	case codes.PermissionDenied:
		return connect.CodePermissionDenied
	case codes.Unauthenticated:
		return connect.CodeUnauthenticated
	case codes.FailedPrecondition:
		return connect.CodeFailedPrecondition
	default:
		return connect.CodeInternal
	}
}
