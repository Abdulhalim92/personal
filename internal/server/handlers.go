package server

import (
	"encoding/json"
	"fmt"
	"io"
	"moneytracker/internal/api/messages"
	"moneytracker/internal/models"
	"moneytracker/internal/server/helpers"
	"net/http"
	"time"
)

// Registration регистрация пользователя
func (s *Server) Registration(w http.ResponseWriter, r *http.Request) {
	var user models.User

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		log.Errorln(err)
		helpers.NewErrorResponse(w, 400, err.Error())
		return
	}
	// валидация введенных логина и пароля
	err = s.Service.ValidateUser(&user)
	if err != nil {
		log.Errorln(err)
		helpers.NewErrorResponse(w, 400, err.Error())
		return
	}

	// проверка логина на дубликат
	isLoginUsed, err := s.Service.IsLoginUsed(user.Login)
	if err != nil {
		log.Error(err)
		helpers.NewErrorResponse(w, 400, err.Error())
		return
	}
	if isLoginUsed {
		log.Println(messages.ErrLoginUsed)
		helpers.NewErrorResponse(w, 400, err.Error())
		return
	}

	// регистрация пользователя
	err = s.Service.Registration(&user)
	if err != nil {
		log.Errorln(err)
		helpers.NewErrorResponse(w, 500, err.Error())
		return
	}

	msg := "Registration is successful"
	helpers.NewSuccessResponse(w, 200, msg)

}

// Login вход пользователя
func (s *Server) Login(w http.ResponseWriter, r *http.Request) {
	var userInfo *models.User

	err := json.NewDecoder(r.Body).Decode(&userInfo)
	if err != nil {
		log.Errorln(err)
		helpers.NewErrorResponse(w, 400, err.Error())
		return
	}

	token, err := s.Service.Login(userInfo)
	if err != nil {
		log.Errorln("Non-existent or expired token", err)
		helpers.NewErrorResponse(w, 400, "Non-existent or expired token")
	}
	if token == "" {
		log.Println("token not received")
		helpers.NewErrorResponse(w, 400, "token not received")
		return
	}

	newToken := models.Token{
		UserID: userInfo.ID,
		Token:  token,
	}

	err = s.Service.SetToken(&newToken)
	if err != nil {
		log.Errorln(err)
		helpers.NewErrorResponse(w, 500, err.Error())
		return
	}

	helpers.NewSuccessResponse(w, 200, token)
}

// CreateAccount регистраций нового счета пользователя
func (s *Server) CreateAccount(w http.ResponseWriter, r *http.Request) {
	var newAccount *models.Account

	err := json.NewDecoder(r.Body).Decode(&newAccount)
	if err != nil {
		log.Errorln(err)
		helpers.NewErrorResponse(w, 400, err.Error())
		return
	}

	ctx := r.Context()
	newAccount.UserId = ctx.Value(UserId).(int)

	// проверка счета пользователя на дубликат
	isAccountUsed, err := s.Service.IsAccountUsed(newAccount.Name)
	if err != nil {
		log.Errorln(err)
		helpers.NewErrorResponse(w, 400, err.Error())
		return
	}
	if isAccountUsed {
		log.Println(messages.ErrAccountUsed)
		helpers.NewErrorResponse(w, 400, messages.ErrAccountUsed.Error())
		return
	}

	// регистрация нового счета пользователя
	err = s.Service.CreateAccount(newAccount)
	if err != nil {
		log.Errorln(err)
		helpers.NewErrorResponse(w, 500, err.Error())
		return
	}

	msg := "Adding new account was successful"
	helpers.NewSuccessResponse(w, 200, msg)
}

// UpdateAccount редактирование счета пользователя
func (s *Server) UpdateAccount(w http.ResponseWriter, r *http.Request) {
	var editAccount *models.EditDelAccount

	err := json.NewDecoder(r.Body).Decode(&editAccount)
	if err != nil {
		log.Errorln(err)
		helpers.NewErrorResponse(w, 400, err.Error())
		return
	}

	accountID := editAccount.AccountID
	newName := editAccount.Name

	ctx := r.Context()
	userId := ctx.Value(UserId).(int)

	isUsersAccount, err := s.Service.IsUsersAccount(userId, accountID)
	if err != nil {
		log.Errorln(err)
		helpers.NewErrorResponse(w, 400, err.Error())
		return
	}
	if !isUsersAccount {
		log.Println(messages.ErrAccountBelongUser)
		helpers.NewErrorResponse(w, 400, messages.ErrAccountBelongUser.Error())
	}

	err = s.Service.UpdateAccount(userId, accountID, newName)
	if err != nil {
		log.Errorln(err)
		helpers.NewErrorResponse(w, 500, err.Error())
		return
	}

	msg := "Editing completed successfully"
	helpers.NewSuccessResponse(w, 200, msg)
}

// DeleteAccount удаление счета пользователя
func (s *Server) DeleteAccount(w http.ResponseWriter, r *http.Request) {
	var delAccount *models.EditDelAccount

	err := json.NewDecoder(r.Body).Decode(&delAccount)
	if err != nil {
		log.Errorln(err)
		helpers.NewErrorResponse(w, 400, err.Error())
		return
	}

	ctx := r.Context()
	userId := ctx.Value(UserId).(int)

	accountID := delAccount.AccountID

	isUsersAccount, err := s.Service.IsUsersAccount(userId, accountID)
	if err != nil {
		log.Errorln(err)
		helpers.NewErrorResponse(w, 400, err.Error())
		return
	}
	if !isUsersAccount {
		log.Println(messages.ErrAccountBelongUser)
		helpers.NewErrorResponse(w, 400, messages.ErrAccountBelongUser.Error())
	}

	if delAccount.Name != "" {
		log.Println("When deleting an account is enough account_id")
		helpers.NewErrorResponse(w, 400, "When deleting an account is enough account_id")
		return
	}

	err = s.Service.DeleteAccount(userId, accountID)
	if err != nil {
		log.Errorln(err)
		helpers.NewErrorResponse(w, 500, err.Error())
	}

	msg := "Removal was successful"
	helpers.NewSuccessResponse(w, 200, msg)
}

// CreateOperation добавление операций
func (s *Server) CreateOperation(w http.ResponseWriter, r *http.Request) {
	var operation *models.Operation

	err := json.NewDecoder(r.Body).Decode(&operation)
	if err != nil {
		log.Errorln(err)
		helpers.NewErrorResponse(w, 400, err.Error())
		return
	}

	ctx := r.Context()
	operation.UserID = ctx.Value(UserId).(int)

	// проверка на корректность счета пользователя
	if operation.AccountID <= 0 || operation.AccountIDTo < 0 {
		log.Println("indicated the wrong account")
		helpers.NewErrorResponse(w, 400, "indicated the wrong account")
		return
	}
	isUsersAccount, err := s.Service.IsUsersAccount(operation.UserID, operation.AccountID)
	if err != nil {
		log.Errorln(err)
		helpers.NewErrorResponse(w, 400, err.Error())
		return
	}
	if isUsersAccount == false {
		log.Println(messages.ErrAccountBelongUser)
		helpers.NewErrorResponse(w, 400, messages.ErrAccountBelongUser.Error())
		return
	}

	// проверка на соответствие типа операции и категории
	isTypesCategory, err := s.Service.IsTypesCategory(operation.TypeID, operation.CategoryID)
	if err != nil {
		log.Errorln(err)
		helpers.NewErrorResponse(w, 400, err.Error())
		return
	}
	if isTypesCategory == false {
		log.Println("inconsistency between the type of operations and its category")
		helpers.NewErrorResponse(w, 400, "inconsistency between the type of operations and its category")
		return
	}

	switch operation.TypeID {
	case 1:
		if operation.AccountIDTo != 0 {
			log.Println("income and expense transactions do not have an invoice to send")
			helpers.NewErrorResponse(w, 400, "income and expense transactions do not have an invoice to send")
			return
		}
		err := s.Service.CreateIncome(operation)
		if err != nil {
			log.Errorln(err)
			helpers.NewErrorResponse(w, 500, err.Error())
			return
		}
	case 2:
		if operation.AccountIDTo != 0 {
			log.Println("income and expense transactions do not have an invoice to send")
			helpers.NewErrorResponse(w, 400, "income and expense transactions do not have an invoice to send")
			return
		}
		err := s.Service.CreateExpenditure(operation)
		if err != nil {
			log.Errorln(err)
			helpers.NewErrorResponse(w, 500, err.Error())
			return
		}
	case 3:
		if operation.CategoryID != 0 {
			log.Println("transfer operation has no category")
			helpers.NewErrorResponse(w, 400, "transfer operation has no category")
			return
		}
		isUsersAccountTo, err := s.Service.IsUsersAccount(operation.UserID, operation.AccountIDTo)
		if err != nil {
			log.Errorln(err)
			helpers.NewErrorResponse(w, 400, err.Error())
			return
		}
		if isUsersAccountTo == false {
			log.Println(messages.ErrAccountBelongUser)
			helpers.NewErrorResponse(w, 400, messages.ErrAccountBelongUser.Error())
			return
		}
		err = s.Service.CreateTransfer(operation)
		if err != nil {
			log.Errorln(err)
			helpers.NewErrorResponse(w, 500, err.Error())
			return
		}
	default:
		log.Println("such operation is not registered")
		helpers.NewErrorResponse(w, 400, "such operation is not registered")
		return
	}

	msg := "operation is successful"
	helpers.NewSuccessResponse(w, 200, msg)
}

// GetTypes получение типов операций
func (s *Server) GetTypes(w http.ResponseWriter, r *http.Request) {

	getTypes, err := s.Service.GetTypes()
	if err != nil {
		log.Errorln(err)
		helpers.NewErrorResponse(w, 500, err.Error())
		return
	}

	err = json.NewEncoder(w).Encode(getTypes)
	if err != nil {
		log.Errorln(err)
		helpers.NewErrorResponse(w, 500, err.Error())
		return
	}
}

// GetCategories получение категорий типов операций
func (s *Server) GetCategories(w http.ResponseWriter, r *http.Request) {

	getCategories, err := s.Service.GetCategories()
	if err != nil {
		log.Errorln(err)
		helpers.NewErrorResponse(w, 500, err.Error())
		return
	}

	err = json.NewEncoder(w).Encode(getCategories)
	if err != nil {
		log.Errorln(err)
		helpers.NewErrorResponse(w, 500, err.Error())
		return
	}
}

// GetAccounts получение счетов пользователя
func (s *Server) GetAccounts(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	userId := ctx.Value(UserId).(int)

	accounts, err := s.Service.GetAccounts(userId)
	if err != nil {
		log.Errorln(err)
		helpers.NewErrorResponse(w, 500, err.Error())
		return
	}

	err = json.NewEncoder(w).Encode(accounts)
	if err != nil {
		log.Errorln(err)
		helpers.NewErrorResponse(w, 500, err.Error())
		return
	}

}

// GetTotalBalance получение всей суммы на балансе пользователя
func (s *Server) GetTotalBalance(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	userId := ctx.Value(UserId).(int)

	totalBalance, err := s.Service.GetTotalBalance(userId)
	if err != nil {
		log.Errorln(err)
		helpers.NewErrorResponse(w, 500, err.Error())
		return
	}

	total := fmt.Sprintf("%f", totalBalance)
	helpers.NewSuccessResponse(w, 200, total)
}

// GetActiveAccounts получение активных аккаунтов пользователя
func (s *Server) GetActiveAccounts(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	userId := ctx.Value(UserId).(int)

	activeAccounts, err := s.Service.GetActiveAccounts(userId)
	if err != nil {
		log.Errorln(err)
		helpers.NewErrorResponse(w, 500, err.Error())
		return
	}

	err = json.NewEncoder(w).Encode(activeAccounts)
	if err != nil {
		log.Errorln(err)
		helpers.NewErrorResponse(w, 500, err.Error())
		return
	}
}

// GetReports получение операций пользователя
func (s *Server) GetReports(w http.ResponseWriter, r *http.Request) {
	var report models.Report
	var getReports []models.GetReports

	err := json.NewDecoder(r.Body).Decode(&report)
	if err != nil {
		log.Errorln(err)
		helpers.NewErrorResponse(w, 400, err.Error())
		return
	}

	ctx := r.Context()
	userId := ctx.Value(UserId).(int)

	// проверка на валидность счета
	if report.AccountID < 0 {
		log.Println(messages.ErrAccountBelongUser)
		helpers.NewErrorResponse(w, 400, messages.ErrAccountBelongUser.Error())
		return
	}
	if report.AccountID > 0 {
		isUsersAccount, err := s.Service.IsUsersAccount(userId, report.AccountID)
		if err != nil {
			log.Errorln(err)
			helpers.NewErrorResponse(w, 400, err.Error())
			return
		}
		if isUsersAccount == false {
			log.Println(messages.ErrAccountBelongUser)
			helpers.NewErrorResponse(w, 400, messages.ErrAccountBelongUser.Error())
			return
		}
	}

	// проверка на существование типов операций
	if report.TypeID < 0 || report.TypeID > 3 {
		log.Println(messages.ErrExistsType)
		helpers.NewErrorResponse(w, 400, messages.ErrExistsType.Error())
		return
	}

	// проверка на правильность введенного значения количества записей
	if report.Limit <= 0 {
		log.Println("the number of entries must be a natural number")
		helpers.NewErrorResponse(w, 400, "the number of entries must be a natural number")
		return
	}

	// проверка на правильность введенного значения номера страницы
	if report.Page <= 0 {
		log.Println("page number must be a natural number")
		helpers.NewErrorResponse(w, 400, "page number must be a natural number")
		return
	}

	// проверка даты
	var startDate time.Time
	var endDate time.Time

	if report.StartDate == "" {
		startDate = time.Time{}
	} else {
		startDate, err = time.Parse("2006-01-02", report.StartDate)
		if err != nil {
			log.Errorln(err)
			helpers.NewErrorResponse(w, 400, err.Error())
			return
		}
	}
	if report.EndDate == "" {
		endDate = time.Now()
	} else {
		endDate, err = time.Parse("2006-01-02", report.EndDate)
		if err != nil {
			log.Errorln(err)
			helpers.NewErrorResponse(w, 400, err.Error())
			return
		}
	}
	if startDate.After(endDate) {
		helpers.NewErrorResponse(w, 400, "incorrect dates entered")
		return
	}

	getReports, err = s.Service.GetReports(userId, &report, startDate, endDate)
	if err != nil {
		log.Errorln(err)
		helpers.NewErrorResponse(w, 500, err.Error())
		return
	}

	// отправить клиенту JSON файл
	err = json.NewEncoder(w).Encode(getReports)
	if err != nil {
		log.Errorln(err)
		helpers.NewErrorResponse(w, 500, err.Error())
		return
	}

}

// GetExcelReports получение операций пользователя
func (s *Server) GetExcelReports(w http.ResponseWriter, r *http.Request) {
	var report models.Report
	var getReports []models.GetReports

	err := json.NewDecoder(r.Body).Decode(&report)
	if err != nil {
		log.Errorln(err)
		helpers.NewErrorResponse(w, 400, err.Error())
		return
	}

	ctx := r.Context()
	userId := ctx.Value(UserId).(int)

	// проверка на валидность счета
	if report.AccountID < 0 {
		log.Println(messages.ErrAccountBelongUser)
		helpers.NewErrorResponse(w, 400, messages.ErrAccountBelongUser.Error())
		return
	}
	if report.AccountID > 0 {
		isUsersAccount, err := s.Service.IsUsersAccount(userId, report.AccountID)
		if err != nil {
			log.Errorln(err)
			helpers.NewErrorResponse(w, 400, err.Error())
			return
		}
		if isUsersAccount == false {
			log.Println(messages.ErrAccountBelongUser)
			helpers.NewErrorResponse(w, 400, messages.ErrAccountBelongUser.Error())
			return
		}
	}

	// проверка на существование типов операций
	if report.TypeID < 0 || report.TypeID > 3 {
		log.Println(messages.ErrExistsType)
		helpers.NewErrorResponse(w, 400, messages.ErrExistsType.Error())
		return
	}

	// проверка на правильность введенного значения количества записей
	if report.Limit <= 0 {
		log.Println("the number of entries must be a natural number")
		helpers.NewErrorResponse(w, 400, "the number of entries must be a natural number")
		return
	}

	// проверка на правильность введенного значения номера страницы
	if report.Page <= 0 {
		log.Println("page number must be a natural number")
		helpers.NewErrorResponse(w, 400, "page number must be a natural number")
		return
	}

	// проверка даты
	var startDate time.Time
	var endDate time.Time

	if report.StartDate == "" {
		startDate = time.Time{}
	} else {
		startDate, err = time.Parse("2006-01-02", report.StartDate)
		if err != nil {
			log.Errorln(err)
			helpers.NewErrorResponse(w, 400, err.Error())
			return
		}
	}
	if report.EndDate == "" {
		endDate = time.Now()
	} else {
		endDate, err = time.Parse("2006-01-02", report.EndDate)
		if err != nil {
			log.Errorln(err)
			helpers.NewErrorResponse(w, 400, err.Error())
			return
		}
	}
	if startDate.After(endDate) {
		helpers.NewErrorResponse(w, 400, "incorrect dates entered")
		return
	}

	getReports, err = s.Service.GetReports(userId, &report, startDate, endDate)
	if err != nil {
		log.Errorln(err)
		helpers.NewErrorResponse(w, 500, err.Error())
		return
	}

	// отправить клиенту Excel файл
	excelData, err := s.Service.GetExcelData(userId, &getReports)
	if err != nil {
		log.Errorln(err)
		helpers.NewErrorResponse(w, 500, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/vnd.ms-excel")

	_, err = io.Copy(w, excelData)
	if err != nil {
		log.Errorln(err)
		helpers.NewErrorResponse(w, 500, err.Error())
		return
	}

}
