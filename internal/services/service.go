package services

import (
	"fmt"
	jwt "github.com/golang-jwt/jwt/v4"
	"github.com/pkg/errors"
	"github.com/xuri/excelize/v2"
	"golang.org/x/crypto/bcrypt"
	"moneytracker/internal/api/messages"
	"moneytracker/internal/models"
	"moneytracker/internal/repository"
	"moneytracker/pkg/logging"
	"strconv"
	"strings"
	"time"
)

type Service struct {
	Repository *repository.Repository
}

// NewService конструктор структуры
func NewService(rep *repository.Repository) *Service {
	return &Service{Repository: rep}
}

//type ErrResponse struct {
//	Code int
//	Description string
//}

var log = logging.GetLogger()

// ValidateToken валидация токена
func (s *Service) ValidateToken(tokenString string) (int, error) {

	// парсинг токена
	claims := jwt.MapClaims{}
	_, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return MySigningKey, errors.New("token parse error")
	})
	if err == nil {
		log.Errorln("token parse error")
		return 0, errors.New("token parse error")
	}

	for key, val := range claims {
		switch key {
		case "expire":
			value := fmt.Sprintf("%v", val)
			expireParse, err := time.Parse(time.RFC3339, value)
			if err != nil {
				log.Errorln(err)
				return 0, err
			}
			if expireParse.Before(time.Now()) {
				return 0, messages.ErrExpiredToken
			}
		}
	}

	userId, err := s.Repository.ValidateToken(tokenString)
	if err != nil {
		log.Errorln("token is not available", messages.ErrInvalidToken)
		return 0, errors.WithStack(err)
	}

	return userId, nil
}

// ValidateUser валидация введенных логина и пароля
func (s *Service) ValidateUser(user *models.User) error {
	if len(user.Name) > 20 || len(user.Name) < 3 {
		log.Println(messages.ErrInvalidData)
		return messages.ErrInvalidData
	}
	if len(user.Login) > 20 || len(user.Login) < 3 {
		log.Println(messages.ErrInvalidData)
		return messages.ErrInvalidData
	}
	if len(user.Password) > 20 || len(user.Password) < 6 {
		log.Println(messages.ErrInvalidData)
		return messages.ErrInvalidData
	}
	if strings.Contains(user.Password, "_") || strings.Contains(user.Password, "-") {
		log.Println(messages.ErrInvalidData)
		return messages.ErrInvalidData
	}
	if strings.Contains(user.Password, "@") || strings.Contains(user.Password, "#") {
		log.Println(messages.ErrInvalidData)
		return messages.ErrInvalidData
	}
	if strings.Contains(user.Password, "$") || strings.Contains(user.Password, "%") {
		log.Println(messages.ErrInvalidData)
		return messages.ErrInvalidData
	}
	if strings.Contains(user.Password, "&") || strings.Contains(user.Password, "*") {
		log.Println(messages.ErrInvalidData)
		return messages.ErrInvalidData
	}
	if strings.Contains(user.Password, "(") || strings.Contains(user.Password, ")") {
		log.Println(messages.ErrInvalidData)
		return messages.ErrInvalidData
	}
	if strings.Contains(user.Password, ":") || strings.Contains(user.Password, ".") {
		log.Println(messages.ErrInvalidData)
		return messages.ErrInvalidData
	}
	if strings.Contains(user.Password, "/") || strings.Contains(user.Password, `\`) {
		log.Println(messages.ErrInvalidData)
		return messages.ErrInvalidData
	}
	if strings.Contains(user.Password, ",") || strings.Contains(user.Password, ";") {
		log.Println(messages.ErrInvalidData)
		return messages.ErrInvalidData
	}
	if strings.Contains(user.Password, "?") || strings.Contains(user.Password, `"`) {
		log.Println(messages.ErrInvalidData)
		return messages.ErrInvalidData
	}
	if strings.Contains(user.Password, "!") || strings.Contains(user.Password, "~") {
		log.Println(messages.ErrInvalidData)
		return messages.ErrInvalidData
	}
	return nil
}

// Registration регистрация пользователя
func (s *Service) Registration(newUser *models.User) error {

	// хеширование пароля
	hash, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Errorln(err)
		return err
	}

	newUser.Password = string(hash)

	// регистрация пользователя
	err = s.Repository.Registration(newUser)
	if err != nil {
		log.Errorln(err)
		return err
	}

	return nil
}

// SetToken Метод для добавления токена
func (s *Service) SetToken(token *models.Token) error {
	err := s.Repository.SetToken(token)
	if err != nil {
		log.Errorln(err)
		return errors.WithStack(err)
	}

	return nil
}

// Login валидация пользователя
func (s *Service) Login(userInfo *models.User) (string, error) {
	userFromDB, err := s.Repository.ValidateLoginAndPassword(userInfo.Login, userInfo.Password)
	if err != nil {
		log.Errorln(err)
		return "", errors.WithStack(err)
	}

	// передача значения ID базы данных к ID user_auth
	userInfo.ID = userFromDB.ID

	err = bcrypt.CompareHashAndPassword([]byte(userFromDB.Password), []byte(userInfo.Password))
	if err != nil {
		log.Errorln(err)
		return "", messages.ErrIncorrectPassword
	}

	// генерация токена с помощью JWT
	token, err := s.GenerateJWT(userInfo.ID)
	if err != nil {
		log.Errorln(err)
		return "", err
	}

	return token, nil
}

var MySigningKey = []byte("secret")

// GenerateJWT генерация токена
func (s *Service) GenerateJWT(userID int) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)

	claims["user_id"] = userID
	claims["authorized"] = true
	claims["created"] = time.Now()
	claims["expire"] = time.Now().Add(24 * time.Hour)

	tokenString, err := token.SignedString(MySigningKey)
	if err != nil {
		log.Errorln(err)
		return "", errors.WithStack(err)
	}

	return tokenString, nil
}

// CreateAccount добавление счета
func (s *Service) CreateAccount(newAccount *models.Account) error {
	err := s.Repository.CreateAccount(newAccount)
	if err != nil {
		err := errors.WithStack(err)
		log.Errorln(err)
		return err
	}

	return nil
}

// UpdateAccount редактирование аккаунта
func (s *Service) UpdateAccount(userId int, accountID int, newName string) error {
	err := s.Repository.UpdateAccount(userId, accountID, newName)
	if err != nil {
		err := errors.WithStack(err)
		log.Errorln(err)
		return err
	}

	return nil
}

// DeleteAccount удаление (отключение) аккаунта
func (s *Service) DeleteAccount(userId int, accountID int) error {
	err := s.Repository.DeleteAccount(userId, accountID)
	if err != nil {
		err := errors.WithStack(err)
		log.Errorln(err)
		return err
	}

	return nil
}

// CreateIncome добавление дохода
func (s *Service) CreateIncome(income *models.Operation) error {
	err := s.Repository.CreateIncome(income)
	if err != nil {
		err := errors.WithStack(err)
		log.Error(err)
		return err
	}

	return nil
}

// CreateExpenditure добавление расхода
func (s *Service) CreateExpenditure(expenditure *models.Operation) error {
	err := s.Repository.CreateExpenditure(expenditure)
	if err != nil {
		err := errors.WithStack(err)
		log.Errorln(err)
		return err
	}

	return nil
}

// CreateTransfer добавление переводов
func (s *Service) CreateTransfer(transfer *models.Operation) error {
	err := s.Repository.CreateTransfer(transfer)
	if err != nil {
		err := errors.WithStack(err)
		log.Errorln(err)
	}

	return nil
}

// GetTypes получение типов операций
func (s *Service) GetTypes() ([]models.Type, error) {

	getTypes, err := s.Repository.GetTypes()
	if err != nil {
		log.Errorln(err)
		return nil, errors.WithStack(err)
	}

	return getTypes, nil
}

// GetCategories получение категорий типов операций
func (s *Service) GetCategories() ([]models.Category, error) {

	getCategories, err := s.Repository.GetCategories()
	if err != nil {
		log.Errorln(err)
		return nil, errors.WithStack(err)
	}

	return getCategories, nil
}

// GetAccounts получение всех счетов пользователя
func (s *Service) GetAccounts(userId int) ([]models.Account, error) {
	accounts, err := s.Repository.GetAccounts(userId)
	if err != nil {
		log.Errorln(err)
		return nil, errors.WithStack(err)
	}

	return accounts, nil
}

// GetActiveAccounts получение активных счетов пользователя
func (s *Service) GetActiveAccounts(userId int) ([]models.Account, error) {
	activeAccounts, err := s.Repository.GetActiveAccounts(userId)
	if err != nil {
		log.Errorln(err)
		return nil, errors.WithStack(err)
	}

	return activeAccounts, nil
}

// GetReports получение отчетов операций пользователя
func (s *Service) GetReports(userId int, report *models.Report, startDate time.Time, endDate time.Time) ([]models.GetReports, error) {

	switch {
	case report.TypeID == 0 && report.AccountID != 0:
		reports, err := s.Repository.GetAllTypesCustomAccount(report, startDate, endDate)
		if err != nil {
			log.Errorln(err)
			return nil, err
		}

		return reports, nil

	case report.TypeID == 0 && report.AccountID == 0:
		reports, err := s.Repository.GetAll(userId, report, startDate, endDate)
		if err != nil {
			log.Errorln(err)
			return nil, err
		}

		return reports, nil

	case report.TypeID != 0 && report.AccountID != 0:
		reports, err := s.Repository.GetCustomTypeCustomAccount(userId, report, startDate, endDate)
		if err != nil {
			log.Errorln(err)
			return nil, err
		}

		return reports, nil

	case report.TypeID != 0 && report.AccountID == 0:
		reports, err := s.Repository.GetCustomTypeAllAccounts(report, startDate, endDate)
		if err != nil {
			log.Errorln(err)
			return nil, err
		}

		return reports, nil

	default:
		log.Println("such transaction is not registered or wrong user account is entered")
		return nil, errors.New("such transaction is not registered or wrong user account is entered")
	}

}

// IsUsersAccount проверка на принадлежность счета к пользователю
func (s *Service) IsUsersAccount(userID int, accountID int) (bool, error) {

	isUsersAccount, err := s.Repository.IsUsersAccount(userID, accountID)
	if err != nil {
		log.Errorln(err)
		return false, err
	}

	return isUsersAccount, nil
}

// IsTypesCategory проверка на соответствие типа операции и категории
func (s *Service) IsTypesCategory(typeID int, categoryID int) (bool, error) {
	isTypesCategory, err := s.Repository.IsTypesCategory(typeID, categoryID)
	if err != nil {
		log.Errorln(err)
		return false, err
	}

	return isTypesCategory, nil
}

// IsLoginUsed проверка на дубликат логина
func (s *Service) IsLoginUsed(login string) (bool, error) {
	isLoginUsed, err := s.Repository.IsLoginUsed(login)
	if err != nil {
		log.Errorln(err)
		return true, err
	}
	if isLoginUsed {
		log.Println(messages.ErrLoginUsed)
		return true, messages.ErrLoginUsed
	}

	return false, nil
}

// IsAccountUsed проверка на дубликат счета пользователя
func (s *Service) IsAccountUsed(name string) (bool, error) {
	isAccountUsed, err := s.Repository.IsAccountUsed(name)
	if err != nil {
		log.Errorln(err)
		return true, err
	}
	if isAccountUsed {
		log.Println(messages.ErrAccountUsed)
		return true, messages.ErrAccountUsed
	}

	return false, nil
}

// GetExcelData получение данных в виде Excel
func (s *Service) GetExcelData(userId int, reports *[]models.GetReports) (*excelize.File, error) {

	excelFile := excelize.NewFile()

	sheet := excelFile.NewSheet("NewSheet")

	err := excelFile.SetCellValue("NewSheet", "A1", "Имя пользователя")
	if err != nil {
		return nil, err
	}
	err = excelFile.SetCellValue("NewSheet", "B1", "Логин")
	if err != nil {
		return nil, err
	}
	err = excelFile.SetCellValue("NewSheet", "C1", "Тип операции")
	if err != nil {
		return nil, err
	}
	err = excelFile.SetCellValue("NewSheet", "D1", "Категория типа")
	if err != nil {
		return nil, err
	}
	err = excelFile.SetCellValue("NewSheet", "E1", "Счет пользователя")
	if err != nil {
		return nil, err
	}
	err = excelFile.SetCellValue("NewSheet", "G1", "Сумма")
	if err != nil {
		return nil, err
	}
	err = excelFile.SetCellValue("NewSheet", "H1", "Дата операции")
	if err != nil {
		return nil, err
	}

	for i, report := range *reports {
		i += 2

		user, err := s.GetInfoByUserId(userId)
		if err != nil {
			log.Errorln(err)
			return nil, err
		}

		category, err := s.GetCategoryById(report.CategoryID)
		if err != nil {
			log.Errorln(err)
			return nil, err
		}

		typeByCategoryId, err := s.GetTypeByCategoryId(report.CategoryID)
		if err != nil {
			log.Errorln(err)
			return nil, err
		}

		infoByAccountId, err := s.GetInfoByAccountId(report.AccountID)
		if err != nil {
			log.Errorln(err)
			return nil, err
		}

		err = excelFile.SetCellValue("NewSheet", "A"+strconv.Itoa(i), user.Name)
		if err != nil {
			return nil, err
		}
		err = excelFile.SetCellValue("NewSheet", "B"+strconv.Itoa(i), user.Login)
		if err != nil {
			return nil, err
		}
		err = excelFile.SetCellValue("NewSheet", "C"+strconv.Itoa(i), typeByCategoryId.Name)
		if err != nil {
			return nil, err
		}
		err = excelFile.SetCellValue("NewSheet", "D"+strconv.Itoa(i), category.Name)
		if err != nil {
			return nil, err
		}
		err = excelFile.SetCellValue("NewSheet", "E"+strconv.Itoa(i), infoByAccountId.Name)
		if err != nil {
			return nil, err
		}
		err = excelFile.SetCellValue("NewSheet", "G"+strconv.Itoa(i), report.Amount)
		if err != nil {
			return nil, err
		}
		err = excelFile.SetCellValue("NewSheet", "H"+strconv.Itoa(i), report.CreatedAt.Format("02.01.2006 15:04"))
		if err != nil {
			return nil, err
		}
	}

	excelFile.SetActiveSheet(sheet)

	err = excelFile.SaveAs("report.xlsx")
	if err != nil {
		log.Errorln(err)
		return nil, err
	}

	return excelFile, nil
}

// GetInfoByUserId получение пользователя по его идентификатору
func (s *Service) GetInfoByUserId(userID int) (*models.User, error) {
	user, err := s.Repository.GetInfoByUserId(userID)
	if err != nil {
		log.Errorln(err)
		return nil, err
	}

	return user, nil
}

// GetInfoByAccountId получение информации по id счета пользователя
func (s *Service) GetInfoByAccountId(accountId int) (*models.Account, error) {
	account, err := s.Repository.GetInfoByAccountId(accountId)
	if err != nil {
		log.Errorln(err)
		return nil, err
	}

	return account, nil
}

// GetTypeByCategoryId получение типа операции по id категории
func (s *Service) GetTypeByCategoryId(categoryID int) (*models.Type, error) {
	typeByCategoryId, err := s.Repository.GetTypeByCategoryId(categoryID)
	if err != nil {
		log.Errorln(err)
		return nil, err
	}

	return typeByCategoryId, nil
}

// GetCategoryById получение категории по id
func (s *Service) GetCategoryById(categoryId int) (*models.Category, error) {
	category, err := s.Repository.GetCategoryById(categoryId)
	if err != nil {
		log.Errorln(err)
		return nil, err
	}

	return category, nil
}

// GetTotalBalance получение всей суммы на балансе пользователя
func (s *Service) GetTotalBalance(userId int) (float64, error) {
	totalBalance, err := s.Repository.GetTotalBalance(userId)
	if err != nil {
		log.Errorln(err)
		return 0, errors.WithStack(err)
	}

	return totalBalance, nil
}
