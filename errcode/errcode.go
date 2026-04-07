package errcode

type ErrorCode struct {
	Message string `json:"message"`
}

var (
	InternalServerError               = ErrorCode{Message: "Internal Server Error"}
	RegisteredEmail                   = ErrorCode{Message: "Email has been registered"}
	InvalidRequest                    = ErrorCode{Message: "Invalid Request"}
	ValidationError                   = ErrorCode{Message: "Validation error"}
	InvalidEncryptedText              = ErrorCode{Message: "Invalid encrypted text"}
	UserNotFound                      = ErrorCode{Message: "User not found"}
	FailedGetUser                     = ErrorCode{Message: "Failed to get user"}
	RegisteredUsername                = ErrorCode{Message: "Username has been used"}
	AdminNotFound                     = ErrorCode{Message: "Admin not found"}
	FailedRetrieveMessage             = ErrorCode{Message: "Failed to retrieve messages"}
	MessageNotFound                   = ErrorCode{Message: "Message not found"}
	FriendNotFound                    = ErrorCode{Message: "Friend not found"}
	FriendRequestNotFound             = ErrorCode{Message: "Friend request not found"}
	CompanyNotFound                   = ErrorCode{Message: "Company not found"}
	TicketNotFound                    = ErrorCode{Message: "Ticket not found"}
	TicketAssigned                    = ErrorCode{Message: "Ticket has been assigned"}
	ChatGroupMessageNotFound          = ErrorCode{Message: "Chat group message not found"}
	InvalidOwnerId                    = ErrorCode{Message: "Invalid ownerId of group"}
	CompanyIdAndUserNameFieldRequired = ErrorCode{Message: "CompanyId and Username are required"}
	MediaFieldRequired                = ErrorCode{Message: "Media must be provided"}
	TokenFinished                     = ErrorCode{Message: "Token has all been used"}
	RegisteredPhoneNumber             = ErrorCode{Message: "Phone number has been registered"}
	ActorNotFound                     = ErrorCode{Message: "Actor not found"}
)
