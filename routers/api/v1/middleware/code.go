package middleware

type ResCode int

const (
	CodeSuccess      ResCode = 200
	CodeServeBusy    ResCode = 500
	CodeBadRequest   ResCode = 400
	CodeUnauthorized ResCode = 401
	CodeForbidden    ResCode = 403

	CodePasswordError        ResCode = 1000
	CodeLoginError           ResCode = 1001
	CodeUserNotExist         ResCode = 1002
	CodeRegisterRequestError ResCode = 1003
	CodeRegisterError        ResCode = 1004
	CodeInternalCreateError  ResCode = 1010
	CodeInstanceDeleteError  ResCode = 1011
	CodeInstanceQueryError   ResCode = 1012
	CodeInternalPatchError   ResCode = 1013
	CodeInstanceListError    ResCode = 1014
	CodeInstanceDetailError  ResCode = 1015
	CodeInstanceStartError   ResCode = 1016
	CodeInstancePauseError   ResCode = 1017
	CodeInstanceStopError    ResCode = 1018
	CodeInstanceRestartError ResCode = 1019
	CodeInstanceStatusError  ResCode = 1020
	CodeInstanceSaveError    ResCode = 1021

	CodeServerQueryError     ResCode = 1031
	CodeServerSaveError      ResCode = 1032
	CodeResourceInsufficient ResCode = 1040

	CodeAdminServerAddEditError  ResCode = 2001
	CodeAdminServerListError     ResCode = 2002
	CodeAdminServerDetailError   ResCode = 2003
	CodeAdminServerDeleteError   ResCode = 2004
	CodeAdminServerInstanceError ResCode = 2005
	CodeAdminUserQueryError      ResCode = 2010
	CodeAdminUserListError       ResCode = 2011
	CodeAdminUserDetailError     ResCode = 2012
	CodeAdminUserDeleteError     ResCode = 2013
	CodeAdminUserInstanceNoEmpty ResCode = 2014
	CodeAdminUserModifyError     ResCode = 2015
)

var codeMsgMap = map[ResCode]string{
	CodeSuccess:      "success",
	CodeServeBusy:    "server busy",
	CodeUnauthorized: "unauthorized",
	CodeBadRequest:   "bad request",
	CodeForbidden:    "forbidden",

	CodePasswordError:        "password error",
	CodeLoginError:           "login error",
	CodeUserNotExist:         "user not exist",
	CodeRegisterRequestError: "register request error",
	CodeRegisterError:        "username or email exist",
	CodeInternalCreateError:  "create error",
	CodeInstanceDeleteError:  "delete instance error",
	CodeInstanceStatusError:  "instance status error",
	CodeInstanceQueryError:   "query instance error",
	CodeInstanceListError:    "list instance error",
	CodeInstanceDetailError:  "detail instance error",
	CodeInstanceStartError:   "start instance error",
	CodeInstancePauseError:   "pause instance error",
	CodeInstanceStopError:    "stop instance error",
	CodeInstanceRestartError: "restart instance error",
	CodeResourceInsufficient: "resource insufficient",
	CodeInternalPatchError:   "patch error",
	CodeInstanceSaveError:    "save instance error",

	CodeServerQueryError: "query server error",
	CodeServerSaveError:  "save server error",

	CodeAdminServerAddEditError:  "add server error",
	CodeAdminServerListError:     "list server error",
	CodeAdminServerDetailError:   "detail server error",
	CodeAdminServerDeleteError:   "delete server error",
	CodeAdminServerInstanceError: "server instances not empty",
	CodeAdminUserQueryError:      "query user error",
	CodeAdminUserListError:       "list user error",
	CodeAdminUserDetailError:     "detail user error",
	CodeAdminUserDeleteError:     "delete user error",
	CodeAdminUserModifyError:     "modify user error",
	CodeAdminUserInstanceNoEmpty: "user instances not empty",
}
