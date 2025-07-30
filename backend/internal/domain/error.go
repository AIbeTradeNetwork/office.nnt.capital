package domain

import (
	"fmt"
	"log/slog"
	"strings"
)

type Error struct {
	Code   string  `json:"code"`
	Source string  `json:"source"`
	Errors []Error `json:"errors"`
	Native []error `json:"native"`
}

func ErrorIs(err error, target string) bool {
	if err == nil {
		return false
	}
	cErr, ok1 := err.(*Error)
	if ok1 {
		return cErr.Code == target
	}
	dErr, ok2 := err.(Error)
	if ok2 {
		return dErr.Code == target
	}
	return err.Error() == target
}

func AsError(err error) *Error {
	cErr, ok1 := err.(*Error)
	if ok1 {
		return cErr
	}
	dErr, ok2 := err.(Error)
	if ok2 {
		return &dErr
	}
	return nil
}

func NewError(source string) *Error {
	return &Error{Source: source}
}

func (e Error) Error() string {
	err := e.String()
	if len(e.Native) > 0 {
		err = fmt.Sprintf("%s: %s", err, e.GetNative())
	}
	return fmt.Sprintf("%s: %s", err, e.GetErrors())
}

func (e *Error) GetErrors() string {
	var out []string
	for _, err := range e.Errors {
		out = append(out, err.GetErrors())
		if len(err.Native) > 0 {
			out = append(out, err.GetNative())
		}
	}
	return strings.Join(out, " / ")
}

func (e *Error) GetNative() string {
	var out []string
	for _, nErr := range e.Native {
		out = append(out, nErr.Error())
	}
	return strings.Join(out, " + ")
}

func (e *Error) String() string {
	return fmt.Sprintf("%s:%s", e.Source, e.Code)
}

func (e *Error) SetCode(code string) *Error {
	e.Code = code
	return e
}

func (e *Error) SetSource(source string) *Error {
	e.Source = source
	return e
}

func (e *Error) Add(err error) *Error {
	if err == nil {
		return e
	}
	cErr, ok1 := err.(*Error)
	dErr, ok2 := err.(Error)
	if ok1 {
		e.Errors = append(e.Errors, *cErr)
	} else if ok2 {
		e.Errors = append(e.Errors, dErr)
	} else {
		e.Native = append(e.Native, err)
	}
	return e
}

func (e *Error) Log() *Error {
	slog.Error(fmt.Sprintf("%s:%s", e.Source, e.Code), "errors", e.Errors, "native", e.Native)
	return e
}

const (
	ErrServer                         = "serverError"
	ErrValidation                     = "validationError"
	ErrEmailEmpty                     = "emailEmpty"
	ErrEmailWrong                     = "emailWrong"
	ErrNicknameWrong                  = "nicknameWrong"
	ErrPasswordEmpty                  = "passwordEmpty"
	ErrPasswordTooShort               = "passwordTooShort"
	ErrPasswordEncode                 = "passwordEncodeError"
	ErrNoDocuments                    = "recordNotFound"
	ErrCount                          = "countError"
	ErrConfig                         = "configError"
	ErrAlreadyInitialized             = "alreadyInitializedError"
	ErrConnect                        = "connectionError"
	ErrRepoInit                       = "repositoryInitError"
	ErrTransportInit                  = "transportInitError"
	ErrProviderInit                   = "providerInitError"
	ErrPasswordWrong                  = "emailPasswordWrong"
	ErrRePasswordWrong                = "passwordsNotEqual"
	ErrRefUidEmpty                    = "refUidEmpty"
	ErrRefUidNotFound                 = "refUidNotFound"
	ErrRefUserCount                   = "refUserCountError"
	ErrRefUserLimit                   = "refUserLimitError"
	ErrUserLevel                      = "userLevelError"
	ErrUserCreate                     = "userCreateError"
	ErrUserAuthCreate                 = "userAuthCreateError"
	ErrUserNotFound                   = "userNotFound"
	ErrUserAuthNotFound               = "userAuthNotFound"
	ErrTokenGeneration                = "tokenGenerationError"
	ErrRefreshGeneration              = "refreshGenerationError"
	ErrTokenMethodWrong               = "tokenMethodWrong"
	ErrAccessDenied                   = "accessDenied"
	ErrWalletExists                   = "walletExists"
	ErrWalletParse                    = "walletParseError"
	ErrWalletCheck                    = "walletCheckError"
	ErrFind                           = "findError"
	ErrRefSelf                        = "refSelfError"
	ErrRefExists                      = "refExistsError"
	ErrCreate                         = "createError"
	ErrUpdate                         = "updateError"
	ErrUpsert                         = "upsertError"
	ErrDelete                         = "deleteError"
	ErrUserBalanceChange              = "userBalanceChangeError"
	ErrUserBalanceRollback            = "userBalanceRollbackError"
	ErrWorkflow                       = "workflowError"
	ErrWorkerRun                      = "workerRunError"
	ErrSignal                         = "signalError"
	ErrPlanRankEmpty                  = "planRankEmpty"
	ErrPlanRankPeriodEmpty            = "planRankPeriodEmpty"
	ErrPlanRankNotFound               = "planRankNotFound"
	ErrWorkflowInit                   = "workflowInitError"
	ErrBuyWorkflowPlan                = "buyWorkflowPlanError"
	ErrBuyWorkflowPlanAdd             = "buyWorkflowPlanAddError"
	ErrBuyWorkflowPlanEnd             = "buyWorkflowPlanEndError"
	ErrBuyWorkflowRankEnd             = "buyWorkflowRankEndError"
	ErrBuyWorkflowRankAddFromPlan     = "buyWorkflowRankAddFromPlanError"
	ErrBuyWorkflowInit                = "buyWorkflowInitError"
	ErrBuyWorkflowPaid                = "buyWorkflowPaidError"
	ErrBuyWorkflowRefund              = "buyWorkflowRefundError"
	ErrBuyWorkflowApprove             = "buyWorkflowApproveError"
	ErrBuyWorkflowSetRowCol           = "buyWorkflowSetRowColError"
	ErrBuyWorkflowSetClientRowCol     = "buyWorkflowSetClientRowColError"
	ErrBuyWorkflowPlace               = "buyWorkflowPlaceError"
	ErrBuyWorkflowPlacesUp            = "buyWorkflowPlacesUpError"
	ErrBuyWorkflowPlacesRefUp         = "buyWorkflowPlacesRefUpError"
	ErrBuyWorkflowCalculateNextRank   = "buyWorkflowCalculateNextRankError"
	ErrBuyWorkflowCalculateActivity   = "buyWorkflowCalculateActivityError"
	ErrBuyWorkflowCancel              = "buyWorkflowCancelError"
	ErrBuyWorkflowChargeRef           = "buyWorkflowChargeRefError"
	ErrBuyWorkflowChargeBin           = "buyWorkflowChargeBinError"
	ErrBuyWorkflowChargeFirstRank     = "buyWorkflowChargeFirstRankError"
	ErrBuyWorkflowChargeFastStart     = "buyWorkflowChargeFastStartError"
	ErrBuyWorkflowChargeMatch         = "buyWorkflowChargeMatchError"
	ErrBuyWorkflowCharged             = "buyWorkflowChargedError"
	ErrBuyWorkflowExecute             = "buyWorkflowChargeError"
	ErrProductWorkflowExecute         = "productWorkflowExecuteError"
	ErrNotifyWorkflowExecute          = "notifyWorkflowExecuteError"
	ErrDepositWorkflowExecute         = "depositWorkflowExecuteError"
	ErrTeamNoFreePlace                = "teamNoFreePlace"
	ErrUserConfig                     = "userConfigError"
	ErrUserPlaceCreate                = "userPlaceCreateError"
	ErrDistCreate                     = "distCreateError"
	ErrUserActivityCreate             = "userActivityCreateError"
	ErrUserConfigCreate               = "userConfigCreateError"
	ErrUserPlaceExists                = "userPlaceExistsError"
	ErrUserActivityExists             = "userActivityExistsError"
	ErrBuyWithoutPlace                = "buyWithoutPlaceError"
	ErrCurrencyCodeInvalid            = "currencyCodeInvalid"
	ErrBadRequest                     = "badRequest"
	ErrPaymentGateway                 = "paymentGatewayError"
	ErrUnexpectedProductType          = "unexpectedProductTypeError"
	ErrTeamTooSmall                   = "teamTooSmall"
	ErrPayoutAmountMin                = "payoutAmountMin"
	ErrPayoutMethod                   = "payoutMethodNotSupported"
	ErrPayoutCurrency                 = "payoutCurrencyNotSupported"
	ErrPayoutAccountNumber            = "payoutAccountNumberInvalid"
	ErrTelegramAuthFailed             = "telegramAuthFailed"
	ErrTelegramHashCheckFailed        = "telegramHashCheckFailed"
	ErrTelegramDataDecodeFailed       = "telegramDataDecodeFailed"
	ErrTelegramDataUserEmpty          = "telegramDataUserEmpty"
	ErrClaimMinPeriod                 = "claimMinPeriod"
	ErrClaimNotAvailable              = "claimNotAvailable"
	ErrNotifyWorkflowCreate           = "notifyWorkflowCreateError"
	ErrNotifyWorkflowGetAllTgIDs      = "notifyWorkflowGetAllTgIDsError"
	ErrNotifyWorkflowUpdate           = "notifyWorkflowUpdateError"
	ErrNotifyWorkflowTgSend           = "notifyWorkflowTgSendError"
	ErrDepositWorkflowGetLastDeposits = "depositWorkflowGetLastDepositsError"
	ErrDepositCreate                  = "depositCreateError"
	ErrDepositTransactionCreate       = "depositTransactionCreateError"
	ErrDepositChangeBalance           = "depositChangeBalanceError"
	ErrTonRate                        = "tonRateError"
	ErrProductWorkflowInit            = "productWorkflowInitError"
	ErrAutofarmWorkflowSleep          = "autofarmWorkflowSleepError"
	ErrAutofarmWorkflowInit           = "autofarmWorkflowInitError"
	ErrAutofarmWorkflowGetUser        = "autofarmWorkflowGetUserError"
	ErrAutofarmWorkflowExecute        = "autofarmWorkflowExecuteError"
	ErrProductWorkflowProduct         = "productWorkflowProductError"
	ErrProductWorkflowPlanAdd         = "productWorkflowPlanAddError"
	ErrProductWorkflowPlanApply       = "productWorkflowPlanApplyError"
	ErrProductWorkflowRefUsersUp      = "productWorkflowRefUsersUpError"
	ErrProductWorkflowRefUserCharge   = "productWorkflowRefUserChargeError"
	ErrBalanceNotEnough               = "balanceNotEnough"
	ErrPayForOrder                    = "payForOrderError"
	ErrCurrencyNotSupported           = "currencyNotSupported"
	ErrTaskAlreadyCompleted           = "taskAlreadyCompleted"
	ErrTaskNotCompleted               = "taskNotCompleted"
	ErrTaskInactive                   = "taskInactive"
	ErrTaskEnded                      = "taskEnded"
	ErrTaskNotStarted                 = "taskNotStarted"
	ErrTaskCurrencyNotSupported       = "taskCurrencyNotSupported"
	ErrAutofarmStart                  = "autofarmStartError"
	ErrComboNotFound                  = "comboNotFound"
	ErrComboAlreadyUsed               = "comboAlreadyUsed"
	ErrComboLimitReached              = "comboLimitReached"
	ErrProductNotFound                = "productNotFound"
	ErrComboReward                    = "comboRewardError"
	ErrComboIncrement                 = "comboIncrementError"
	ErrTaskNeedAutoApprove            = "taskNeedAutoApprove"
	ErrUserNotHaveTgID                = "userNotHaveTgID"
	ErrTaskApprove                    = "taskApproveError"
	ErrUserSafeNotFound               = "userSafeNotFound"
	ErrUserSafeNotOpenedEnough        = "userSafeNotOpenedEnough"
	ErrUserProductExists              = "userProductExists"
	ErrProductLimitReached            = "productLimitReached"
	ErrProductIncCount                = "productIncCountError"
	ErrSafeNotFound                   = "safeNotFound"
	ErrSafeNotHacked                  = "safeNotHacked"
	ErrSafeVariant                    = "safeVariantError"
	ErrSafeUpdate                     = "safeUpdateError"
	ErrUserSafeUpdate                 = "userSafeUpdateError"
	ErrSafeVariantLimit               = "safeVariantLimitError"
	ErrSafeVariantType                = "safeVariantTypeError"
	ErrSafeVariantTransaction         = "safeVariantTransactionError"
)
