package healthcheck

import (
	"github.com/romberli/das/pkg/message"
	"github.com/romberli/go-util/config"
)

func init() {
	initServiceDebugMessage()
	initServiceInfoMessage()
	initServiceErrorMessage()
}

const (
	// debug
	DebugHealthcheckGetResultByOperationID = 101001
	DebugHealthcheckCheck                  = 101002
	DebugHealthcheckCheckByHostInfo        = 101003
	DebugHealthcheckReviewAccuracy         = 101004
	// info
	InfoHealthcheckGetResultByOperationID = 201001
	InfoHealthcheckCheck                  = 201002
	InfoHealthcheckCheckByHostInfo        = 201003
	InfoHealthcheckReviewAccuracy         = 201004
	// error
	ErrHealthcheckDefaultEngineRun       = 401013
	ErrHealthcheckGetResultByOperationID = 401014
	ErrHealthcheckCheck                  = 401015
	ErrHealthcheckCheckByHostInfo        = 401016
	ErrHealthcheckReviewAccuracy         = 401017
	ErrHealthcheckCloseConnection        = 401018
)

func initServiceDebugMessage() {
	message.Messages[DebugHealthcheckGetResultByOperationID] = config.NewErrMessage(
		message.DefaultMessageHeader, DebugHealthcheckGetResultByOperationID,
		"healthcheck: get result by operation id message: %s")
	message.Messages[DebugHealthcheckCheck] = config.NewErrMessage(
		message.DefaultMessageHeader, DebugHealthcheckCheck,
		"healthcheck: check message: %s")
	message.Messages[DebugHealthcheckCheckByHostInfo] = config.NewErrMessage(
		message.DefaultMessageHeader, DebugHealthcheckCheckByHostInfo,
		"healthcheck: check by host info message: %s")
	message.Messages[DebugHealthcheckReviewAccuracy] = config.NewErrMessage(
		message.DefaultMessageHeader, DebugHealthcheckReviewAccuracy,
		"healthcheck: review accuracy message: %s")
}

func initServiceInfoMessage() {
	message.Messages[InfoHealthcheckGetResultByOperationID] = config.NewErrMessage(
		message.DefaultMessageHeader, InfoHealthcheckGetResultByOperationID,
		"healthcheck: get result by operation id completed. operation_id: %d")
	message.Messages[InfoHealthcheckCheck] = config.NewErrMessage(
		message.DefaultMessageHeader, InfoHealthcheckCheck,
		"healthcheck: check completed. %s")
	message.Messages[InfoHealthcheckCheckByHostInfo] = config.NewErrMessage(
		message.DefaultMessageHeader, InfoHealthcheckCheckByHostInfo,
		"healthcheck: check by host info completed. %s")
	message.Messages[InfoHealthcheckReviewAccuracy] = config.NewErrMessage(
		message.DefaultMessageHeader, InfoHealthcheckReviewAccuracy,
		"healthcheck: review accuracy completed. %s")
}

func initServiceErrorMessage() {
	message.Messages[ErrHealthcheckDefaultEngineRun] = config.NewErrMessage(
		message.DefaultMessageHeader, ErrHealthcheckDefaultEngineRun,
		"default engine run failed.\n%s")
	message.Messages[ErrHealthcheckGetResultByOperationID] = config.NewErrMessage(
		message.DefaultMessageHeader, ErrHealthcheckGetResultByOperationID,
		"healthcheck: get result by operation id failed. operation_id: %d\n%s")
	message.Messages[ErrHealthcheckCheck] = config.NewErrMessage(
		message.DefaultMessageHeader, ErrHealthcheckCheck,
		"healthcheck: check failed.  %s")
	message.Messages[ErrHealthcheckCheckByHostInfo] = config.NewErrMessage(
		message.DefaultMessageHeader, ErrHealthcheckCheckByHostInfo,
		"healthcheck: check by host info failed. host info %s")
	message.Messages[ErrHealthcheckCloseConnection] = config.NewErrMessage(
		message.DefaultMessageHeader, ErrHealthcheckCloseConnection,
		"healthcheck: close middleware connection failed.\n%s")

}
