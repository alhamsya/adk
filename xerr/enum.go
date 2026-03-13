package xerr

import "google.golang.org/grpc/codes"

const (
	TypeUnknown             Type = ""                      // TypeUnknown for handle unknown error type
	TypeOK                  Type = "OK"                    // TypeOK equals to HTTP Status Code [200] or GRPC code [0].
	TypeInvalidParameter    Type = "INVALID_PARAMETER"     // TypeInvalidParameter equals to HTTP Status Code [400] or GRPC code [3].
	TypeUnauthorized        Type = "UNAUTHORIZED"          // TypeUnauthorized equals to HTTP Status Code [401] or GRPC code [16].
	TypeNotFound            Type = "NOT_FOUND"             // TypeNotFound equals to HTTP Status Code [404] or GRPC code [5].
	TypeServiceBusy         Type = "SERVICE_BUSY"          // TypeServiceBusy equals to HTTP Status Code [429] or GRPC code [8].
	TypeSystemError         Type = "SYSTEM_ERROR"          // TypeSystemError equals to HTTP Status Code [500] or GRPC code [13].
	TypeVendorError         Type = "VENDOR_ERROR"          // TypeVendorError equals to HTTP Status Code [500] or GRPC code [13].
	TypeBadGateway          Type = "BAD_GATEWAY"           // TypeBadGateway equals to HTTP Status Code [502] or GRPC code [14].
	TypeMaintenance         Type = "MAINTENANCE"           // TypeMaintenance equals to HTTP Status Code [503] or GRPC code [14].
	TypeGatewayTimeout      Type = "GATEWAY_TIMEOUT"       // TypeGatewayTimeout equals to HTTP Status Code [504] or GRPC code [4].
	TypeNoTradingPermission Type = "NO_TRADING_PERMISSION" // TypeNoTradingPermission equals to HTTP Status Code [403] or GRPC code [7].
	TypeNoSubscription      Type = "NO_SUBSCRIPTION"       // TypeNoSubscription equals to HTTP Status Code [403] or GRPC code [7].
	TypeForbidden           Type = "FORBIDDEN"             // TypeForbidden equals to HTTP Status Code [403] or GRPC code [7].
	TypeDuplicateCall       Type = "DUPLICATE_CALL"        // TypeDuplicateCall equals to HTTP Status Code [412] or GRPC code [9].
	TypeRequestCanceled     Type = "REQUEST_CANCELED"      // TypeRequestCanceled equals to HTTP Status Code [499] or GRPC code [1].
)

var TypeToGRPCCode = map[Type]codes.Code{
	TypeOK:                  codes.OK,
	TypeInvalidParameter:    codes.InvalidArgument,
	TypeUnauthorized:        codes.Unauthenticated,
	TypeNotFound:            codes.NotFound,
	TypeServiceBusy:         codes.ResourceExhausted,
	TypeSystemError:         codes.Internal,
	TypeVendorError:         codes.Internal,
	TypeBadGateway:          codes.Unavailable,
	TypeMaintenance:         codes.Unavailable,
	TypeGatewayTimeout:      codes.DeadlineExceeded,
	TypeNoTradingPermission: codes.PermissionDenied,
	TypeNoSubscription:      codes.PermissionDenied,
	TypeForbidden:           codes.PermissionDenied,
	TypeDuplicateCall:       codes.FailedPrecondition,
	TypeRequestCanceled:     codes.Canceled,
}

var (
	// ErrInvalidParameter equals to HTTP Status Code [400] or GRPC code [3].
	ErrInvalidParameter = &Error{
		Code:    TypeToGRPCCode[TypeInvalidParameter],
		Message: TypeInvalidParameter.String(),
		Type:    TypeInvalidParameter,
	}
	// ErrUnauthorized equals to HTTP Status Code [401] or GRPC code [16].
	ErrUnauthorized = &Error{
		Code:    TypeToGRPCCode[TypeUnauthorized],
		Message: TypeUnauthorized.String(),
		Type:    TypeUnauthorized,
	}
	// ErrNotFound equals to HTTP Status Code [404] or GRPC code [5].
	ErrNotFound = &Error{
		Code:    TypeToGRPCCode[TypeNotFound],
		Message: TypeNotFound.String(),
		Type:    TypeNotFound,
	}
	// ErrServiceBusy equals to HTTP Status Code [429] or GRPC code [8].
	ErrServiceBusy = &Error{
		Code:    TypeToGRPCCode[TypeServiceBusy],
		Message: TypeServiceBusy.String(),
		Type:    TypeServiceBusy,
	}
	// ErrSystemError equals to HTTP Status Code [500] or GRPC code [13].
	ErrSystemError = &Error{
		Code:    TypeToGRPCCode[TypeSystemError],
		Message: TypeSystemError.String(),
		Type:    TypeSystemError,
	}
	// ErrVendorError equals to HTTP Status Code [500] or GRPC code [13].
	ErrVendorError = &Error{
		Code:    TypeToGRPCCode[TypeVendorError],
		Message: TypeVendorError.String(),
		Type:    TypeVendorError,
	}
	// ErrBadGateway equals to HTTP Status Code [502] or GRPC code [14].
	ErrBadGateway = &Error{
		Code:    TypeToGRPCCode[TypeBadGateway],
		Message: TypeBadGateway.String(),
		Type:    TypeBadGateway,
	}
	// ErrMaintenance equals to HTTP Status Code [503] or GRPC code [14].
	ErrMaintenance = &Error{
		Code:    TypeToGRPCCode[TypeMaintenance],
		Message: TypeMaintenance.String(),
		Type:    TypeMaintenance,
	}
	// ErrGatewayTimeout equals to HTTP Status Code [504] or GRPC code [4].
	ErrGatewayTimeout = &Error{
		Code:    TypeToGRPCCode[TypeGatewayTimeout],
		Message: TypeGatewayTimeout.String(),
		Type:    TypeGatewayTimeout,
	}
	// ErrNoTradingPermission equals to HTTP Status Code [403] or GRPC code [7].
	ErrNoTradingPermission = &Error{
		Code:    TypeToGRPCCode[TypeNoTradingPermission],
		Message: TypeNoTradingPermission.String(),
		Type:    TypeNoTradingPermission,
	}
	// ErrNoSubscription equals to HTTP Status Code [403] or GRPC code [7].
	ErrNoSubscription = &Error{
		Code:    TypeToGRPCCode[TypeNoSubscription],
		Message: TypeNoSubscription.String(),
		Type:    TypeNoSubscription,
	}
	// ErrForbidden equals to HTTP Status Code [403] or GRPC code [7].
	ErrForbidden = &Error{
		Code:    TypeToGRPCCode[TypeForbidden],
		Message: TypeForbidden.String(),
		Type:    TypeForbidden,
	}
	// ErrDuplicateCall equals to HTTP Status Code [412] or GRPC code [9].
	ErrDuplicateCall = &Error{
		Code:    TypeToGRPCCode[TypeDuplicateCall],
		Message: TypeDuplicateCall.String(),
		Type:    TypeDuplicateCall,
	}
	// ErrRequestCanceled equals to HTTP Status Code [499] or GRPC code [1].
	ErrRequestCanceled = &Error{
		Code:    TypeToGRPCCode[TypeRequestCanceled],
		Message: TypeRequestCanceled.String(),
		Type:    TypeRequestCanceled,
	}
)
