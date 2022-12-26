package repository

import (
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"moneytracker/internal/models"
	"moneytracker/pkg/logging"
	"time"
)

type Repository struct {
	Connection *gorm.DB
}

// NewRepository конструктор структуры
func NewRepository(conn *gorm.DB) *Repository {
	return &Repository{Connection: conn}
}

var log = logging.GetLogger()

// ValidateToken валидация токена
func (r Repository) ValidateToken(token string) (int, error) {
	var tokens *models.Token

	//sqlQuery := `select *from tokens where token = ?`
	//
	//if tx := r.Connection.Raw(sqlQuery, token).Scan(&tokens); tx.Error != nil {
	//	log.Errorln(tx.Error)
	//	return 0, tx.Error
	//}
	//if tokens.UserID == 0 {
	//	log.Println("invalid token entered")
	//	return 0, errors.New("invalid token entered")
	//}

	//gorm
	if tx := r.Connection.Where("token = ?", token).Find(&tokens).Scan(&tokens); tx.Error != nil {
		log.Errorln(tx.Error)
		return 0, tx.Error
	}
	if tokens.UserID == 0 {
		log.Println("invalid token entered")
		return 0, errors.New("invalid token entered")
	}

	return tokens.UserID, nil
}

// IsLoginUsed проверка на существование логина
func (r Repository) IsLoginUsed(login string) (bool, error) {
	var user *models.User

	//sqlQuery := `select *from users where login = ?`
	//
	//if tx := r.Connection.Raw(sqlQuery, login).Scan(&user); tx.Error != nil {
	//	log.Errorln("no rows in result set", tx.Error)
	//	return true, tx.Error
	//}
	//if user.Login != "" {
	//	log.Println("login already registered")
	//	return true, errors.New("login already registered")
	//}

	// gorm
	if tx := r.Connection.Where("login = ?", login).Find(&user).Scan(&user); tx.Error != nil {
		log.Errorln("no rows in result set", tx.Error)
		return true, tx.Error
	}
	if user.Login != "" {
		log.Println("login already registered")
		return true, errors.New("login already registered")
	}

	return false, nil
}

// Registration регистрация пользователя
func (r Repository) Registration(user *models.User) error {
	//sqlQuery := `insert into users (name, login, password)
	//			values (?, ?, ?)`
	//if tx := r.Connection.Exec(sqlQuery, user.Name, user.Login, user.Password); tx.Error != nil {
	//	log.Println("failed to add registration data to the database", tx.Error)
	//	return tx.Error
	//}

	// gorm
	tx := r.Connection.Model(&models.User{}).Omit("active", "updated_at", "deleted_at").Create(user)
	if tx.Error != nil {
		log.Errorln("failed to add registration data to the database", tx.Error)
		return tx.Error
	}

	return nil
}

// ValidateLoginAndPassword валидация пользователя
func (r Repository) ValidateLoginAndPassword(login string, password string) (*models.User, error) {
	var user *models.User

	//sqlQuery := `select *from users where login = ?`
	//
	//if tx := r.Connection.Raw(sqlQuery, login).Scan(&user); tx.Error != nil {
	//	log.Println("incorrect login or password", tx.Error)
	//	return &models.User{}, tx.Error
	//}
	//if user == nil {
	//	log.Println("incorrect login or password")
	//	return &models.User{}, errors.New("incorrect login or password")
	//}

	// gorm
	tx := r.Connection.Where("login = ?", login).Find(&user).Scan(&user)
	if tx.Error != nil {
		log.Errorln("incorrect login or password", tx.Error)
		return &models.User{}, tx.Error
	}
	if user == nil {
		log.Println("incorrect login or password")
		return &models.User{}, errors.New("incorrect login or password")
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		log.Errorln("invalid login or password", err)
		return &models.User{}, err
	}

	return user, nil
}

// SetToken установка токена
func (r Repository) SetToken(token *models.Token) error {
	//sqlQuery := `insert into tokens (user_id, token)
	//			values (?, ?)`
	//
	//tx := r.Connection.Exec(sqlQuery, token.UserID, token.Token)
	//if tx.Error != nil {
	//	return tx.Error
	//}

	// gorm
	if tx := r.Connection.Create(token); tx.Error != nil {
		log.Errorln(tx.Error)
		return tx.Error
	}

	return nil
}

// IsAccountUsed проверка на существование счета пользователя
func (r Repository) IsAccountUsed(name string) (bool, error) {
	var account *models.Account

	//sqlQuery := `select name from accounts where name = ?`
	//
	//if tx := r.Connection.Raw(sqlQuery, name).Scan(&account); tx.Error != nil {
	//	log.Println("no rows in result set", tx.Error)
	//	return false, tx.Error
	//}
	//if account.Name != "" {
	//	log.Println("account already registered")
	//	return true, errors.New("account already registered")
	//}

	// gorm
	if tx := r.Connection.Where("name = ?", name).Find(&account).Scan(&account); tx.Error != nil {
		log.Errorln("no rows in result set", tx.Error)
		return false, tx.Error
	}
	if account.Name != "" {
		log.Println("account already registered")
		return true, errors.New("account already registered")
	}

	return false, nil
}

// CreateAccount добавление счета
func (r Repository) CreateAccount(newAccount *models.Account) error {
	//sqlQuery := `insert into accounts (name, user_id, balance)
	//			values (?, ?, ?)`
	//
	//if tx := r.Connection.Exec(
	//	sqlQuery,
	//	newAccount.Name,
	//	newAccount.UserId,
	//	newAccount.Balance); tx.Error != nil {
	//	log.Println("failed to add account data to the database", tx.Error)
	//	return tx.Error
	//}

	// gorm
	if tx := r.Connection.Omit("active", "created_at", "updated_at", "deleted_at").Create(&newAccount); tx.Error != nil {
		log.Errorln("failed to add account data to the database", tx.Error)
		return tx.Error
	}

	return nil
}

// UpdateAccount редактирование аккаунта
func (r Repository) UpdateAccount(userId int, accountID int, newName string) error {
	//sqlQuery := `update accounts set name = ? where id = ? and user_id = ?`
	//
	//if tx := r.Connection.Exec(sqlQuery, newName, accountID, userId); tx.Error != nil {
	//	log.Errorln(tx.Error)
	//	return tx.Error
	//}

	// gorm
	if tx := r.Connection.Model(&models.Account{}).Where("id = ? and user_id = ?", accountID, userId).Update("name", newName); tx.Error != nil {
		log.Errorln(tx.Error)
		return tx.Error
	}

	return nil
}

// DeleteAccount удаление (отключение) аккаунта
func (r Repository) DeleteAccount(userId int, accountID int) error {
	//sqlQuery := `update accounts set active = false where id = ? and user_id = ?`
	//
	//if tx := r.Connection.Exec(sqlQuery, accountID, userId); tx.Error != nil {
	//	log.Errorln(tx.Error)
	//	return tx.Error
	//}
	//
	// gorm
	if tx := r.Connection.Model(&models.Account{}).Where("id = ? and user_id = ?", accountID, userId).Update("active", false); tx.Error != nil {
		log.Errorln(tx.Error)
		return tx.Error
	}

	return nil
}

// CreateIncome добавление дохода
func (r Repository) CreateIncome(income *models.Operation) error {
	var sqlQuery string

	// добавление данных в таблице operations
	sqlQuery = `insert into operations (category_id, account_id, amount)
				values (?, ?, ?)`

	if tx := r.Connection.Exec(
		sqlQuery,
		income.CategoryID,
		income.AccountID,
		income.Amount); tx.Error != nil {
		log.Errorln("failed to add operation data to the database(operations)", tx.Error)
		return tx.Error // todo
	}

	// добавление данных в таблицу accounts
	sqlQuery = `update accounts set balance = balance + ?, updated_at = current_timestamp where user_id = ? and id = ?`
	if tx := r.Connection.Exec(
		sqlQuery,
		income.Amount,
		income.UserID,
		income.AccountID); tx.Error != nil {
		log.Errorln("failed to add operation data to the database(accounts)", tx.Error)
		return tx.Error // todo
	}

	return nil
}

// CreateExpenditure добавление расхода
func (r Repository) CreateExpenditure(expenditure *models.Operation) error {

	// добавление данных в таблице operations
	sqlQuery := `insert into operations (category_id, account_id, amount)
				values (?, ?, ?)`

	if tx := r.Connection.Exec(
		sqlQuery,
		expenditure.CategoryID,
		expenditure.AccountID,
		expenditure.Amount); tx.Error != nil {
		log.Errorln("failed to add operation data to the database(operations)", tx.Error)
		return tx.Error // todo
	}

	// добавление данных в таблицу accounts
	sqlQuery = `update accounts set balance = balance - ?, updated_at = current_timestamp where user_id = ? and id = ?`
	if tx := r.Connection.Exec(
		sqlQuery,
		expenditure.Amount,
		expenditure.UserID,
		expenditure.AccountID); tx.Error != nil {
		log.Errorln("failed to add operation data to the database(accounts)", tx.Error)
		return tx.Error // todo
	}

	return nil
}

// CreateTransfer добавление перевода
func (r Repository) CreateTransfer(transfer *models.Operation) error {
	sqlQuery := `update accounts set balance = balance - ?, updated_at = current_timestamp where id = ?`
	// снятие с одного счета
	if tx := r.Connection.Exec(
		sqlQuery,
		transfer.Amount,
		transfer.AccountID); tx.Error != nil {
		log.Errorln("failed to update data to the database(accounts)", tx.Error)
		return tx.Error // todo
	}
	// добавление на счет
	if tx := r.Connection.Exec(
		sqlQuery,
		(-1)*transfer.Amount,
		transfer.AccountIDTo); tx.Error != nil {
		log.Errorln("failed to update data to the database(accounts)", tx.Error)
		return tx.Error // todo
	}

	return nil
}

// GetTypes получение всех типов операций
func (r Repository) GetTypes() ([]models.Type, error) {
	var types []models.Type

	sqlQuery := `select *from types`

	if tx := r.Connection.Raw(sqlQuery).Scan(&types); tx.Error != nil {
		log.Errorln("failed to get types from database", tx.Error)
		return nil, tx.Error // todo
	}

	return types, nil
}

// GetCategories получение всех категорий типов операций
func (r Repository) GetCategories() ([]models.Category, error) {
	var categories []models.Category

	sqlQuery := `select *from categories`

	if tx := r.Connection.Raw(sqlQuery).Scan(&categories); tx.Error != nil {
		log.Errorln("failed to get categories from database", tx.Error)
		return nil, tx.Error // todo
	}

	return categories, nil
}

// GetAccounts получение всех счетов пользователя
func (r Repository) GetAccounts(userId int) ([]models.Account, error) {
	var accounts []models.Account

	sqlQuery := `select *from accounts where user_id = ?`

	if tx := r.Connection.Raw(sqlQuery, userId).Scan(&accounts); tx.Error != nil {
		log.Errorln("failed to get accounts for user from database", tx.Error)
		return nil, tx.Error // todo
	}

	return accounts, nil
}

// GetActiveAccounts получение активных счетов пользователя
func (r Repository) GetActiveAccounts(userId int) ([]models.Account, error) {
	var accounts []models.Account

	sqlQuery := `select *from accounts where user_id = ? and active = true`

	if tx := r.Connection.Raw(sqlQuery, userId).Scan(&accounts); tx.Error != nil {
		if tx := r.Connection.Raw(sqlQuery, userId).Scan(&accounts); tx.Error != nil {
			log.Errorln("failed to get active accounts for user from database", tx.Error)
			return nil, tx.Error // todo
		}
	}

	return accounts, nil
}

// ExistsAccount проверка на существования аккаунта
func (r Repository) ExistsAccount(accountID int) (bool, error) {
	sqlQuery := `select id from accounts where id = ?`

	if tx := r.Connection.Raw(sqlQuery, accountID); tx.Error != nil {
		log.Errorln(tx.Error)
		return false, tx.Error
	}

	return true, nil
}

// IsUsersAccount проверка на принадлежность счета к пользователю
func (r Repository) IsUsersAccount(userID int, accountID int) (bool, error) {
	var usersAccountID int
	sqlQuery := `select id from accounts where id = ? and user_id = ?`

	if tx := r.Connection.Raw(sqlQuery, accountID, userID).Scan(&usersAccountID); tx.Error != nil {
		log.Errorln(tx.Error) // todo
		return false, tx.Error
	}
	if usersAccountID == 0 {
		log.Println("incorrect account entered or does not belong to the user")
		return false, errors.New("incorrect account entered or does not belong to the user")
	}

	return true, nil
}

// IsTypesCategory проверка на соответствие типа операции и категории
func (r Repository) IsTypesCategory(typeID int, categoryID int) (bool, error) {
	var id int
	sqlQuery := `select id from categories where id = ? and type_id = ?`

	if tx := r.Connection.Raw(sqlQuery, categoryID, typeID).Scan(&id); tx.Error != nil {
		log.Errorln(tx.Error)
		return false, tx.Error
	}
	if typeID == 3 && categoryID == 0 {
		return true, nil
	}
	if id == 0 {
		log.Println("inconsistency between the type of operations and its category")
		return false, errors.New("inconsistency between the type of operations and its category")
	}

	return true, nil
}

// GetAll получение отчетов по всем типам и всем аккаунтам
func (r Repository) GetAll(userId int, report *models.Report, startDate time.Time, endDate time.Time) ([]models.GetReports, error) {
	var reports []models.GetReports

	sqlQuery := `select *from operations 
        					where account_id in (select id from accounts where user_id = ?) and created_at between ? and ? limit ? offset ?`
	if tx := r.Connection.Raw(
		sqlQuery,
		userId,
		startDate,
		endDate,
		report.Limit,
		(report.Page-1)*report.Limit).Scan(&reports); tx.Error != nil {
		log.Errorln(tx.Error)
		return nil, tx.Error
	}

	return reports, nil
}

// GetAllTypesCustomAccount получение ответов по всем типам и выборочного счета
func (r Repository) GetAllTypesCustomAccount(report *models.Report, startDate time.Time, endDate time.Time) ([]models.GetReports, error) {
	var reports []models.GetReports

	sqlQuery := `select *from operations 
        					where account_id = ? and created_at between ? and ? limit ? offset ?`
	if tx := r.Connection.Raw(
		sqlQuery,
		report.AccountID,
		startDate,
		endDate,
		report.Limit,
		(report.Page-1)*report.Limit).Scan(&reports); tx.Error != nil {
		log.Errorln(tx.Error)
		return nil, tx.Error
	}

	return reports, nil
}

// GetCustomTypeCustomAccount получение отчетов по выборочному типу и выборочного счета
func (r Repository) GetCustomTypeCustomAccount(userId int, report *models.Report, startDate time.Time, endDate time.Time) ([]models.GetReports, error) {
	var reports []models.GetReports

	sqlQuery := `select *from operations 
        					where account_id = ? and account_id in (select id from accounts where user_id = ?) and created_at between ? and ? limit ? offset ?`
	if tx := r.Connection.Raw(
		sqlQuery,
		report.AccountID,
		userId,
		startDate,
		endDate,
		report.Limit,
		(report.Page-1)*report.Limit).Scan(&reports); tx.Error != nil {
		log.Errorln(tx.Error)
		return nil, tx.Error
	}

	return reports, nil
}

// GetCustomTypeAllAccounts получение отчетов по выборочному типу и всех счетов
func (r Repository) GetCustomTypeAllAccounts(report *models.Report, startDate time.Time, endDate time.Time) ([]models.GetReports, error) {
	var reports []models.GetReports

	sqlQuery := `select *from operations
        					where category_id in (select id from categories where type_id = ?) and created_at between ? and ? limit ? offset ?`
	if tx := r.Connection.Raw(
		sqlQuery,
		report.TypeID,
		startDate,
		endDate,
		report.Limit,
		(report.Page-1)*report.Limit).Scan(&reports); tx.Error != nil {
		log.Errorln(tx.Error)
		return nil, tx.Error
	}

	return reports, nil
}

// GetInfoByUserId получение пользователя по его идентификатору
func (r Repository) GetInfoByUserId(userID int) (*models.User, error) {
	var user *models.User

	sqlQuery := `select *from users where id = ?`

	if tx := r.Connection.Raw(sqlQuery, userID).Scan(&user); tx.Error != nil {
		log.Errorln(tx.Error)
		return nil, tx.Error
	}
	if user == nil {
		log.Println("user with this id does not exist")
		return nil, errors.New("user with this id does not exist")
	}

	return user, nil
}

// GetInfoByAccountId получение информации по id счета пользователя
func (r Repository) GetInfoByAccountId(accountId int) (*models.Account, error) {
	var account *models.Account

	sqlQuery := `select *from accounts where id = ?`

	if tx := r.Connection.Raw(sqlQuery, accountId).Scan(&account); tx.Error != nil {
		log.Errorln(tx.Error)
		return nil, tx.Error
	}
	if account == nil {
		log.Println("there is no such account with the given identifier")
		return nil, errors.New("there is no such account with the given identifier")
	}

	return account, nil
}

// GetTypeByCategoryId получение типа операции по id категории
func (r Repository) GetTypeByCategoryId(categoryID int) (*models.Type, error) {
	var operationType *models.Type

	sqlQuery := `select *from types where id = (select type_id from categories where id = ?)`

	if tx := r.Connection.Raw(sqlQuery, categoryID).Scan(&operationType); tx.Error != nil {
		log.Errorln(tx.Error)
		return nil, tx.Error
	}
	if operationType == nil {
		log.Println("there is no such type with this category")
		return nil, errors.New("there is no such type with this category")
	}

	return operationType, nil
}

// GetCategoryById получение категории по id
func (r Repository) GetCategoryById(categoryId int) (*models.Category, error) {
	var category *models.Category

	sqlQuery := `select *from categories where id = ?`

	if tx := r.Connection.Raw(sqlQuery, categoryId).Scan(&category); tx.Error != nil {
		log.Errorln(tx.Error)
		return nil, tx.Error
	}
	if category == nil {
		log.Println("category with this id does not exist")
		return nil, errors.New("category with this id does not exist")
	}

	return category, nil
}

// GetTotalBalance получение всей суммы на балансе пользователя
func (r Repository) GetTotalBalance(userId int) (float64, error) {
	var totalBalance float64

	sqlQuery := `select sum(balance) from accounts where user_id = ?`

	if tx := r.Connection.Raw(sqlQuery, userId).Scan(&totalBalance); tx.Error != nil {
		log.Errorln(tx.Error)
		return 0, tx.Error
	}

	return totalBalance, nil
}
