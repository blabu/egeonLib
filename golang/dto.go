package golang

import "time"

//ServerInfo - статическая информация о сервисе. Зависит от сборки системы
type ServerInfo struct {
	Name       string                 `json:"name"`
	Version    string                 `json:"version"`
	About      string                 `json:"about"`
	Maintainer string                 `json:"maintainer"`
	Routes     map[string]string      `json:"routes"`
	BaseTypes  map[string]interface{} `json:"types"`
}

// DBStats contains database statistics.
type DBStats struct {
	MaxOpenConnections int `json:"maxOpenCon"` // Maximum number of open connections to the database.

	// Pool Status
	OpenConnections int `json:"establishCon"` // The number of established connections both in use and idle.
	InUse           int `json:"inUseCon"`     // The number of connections currently in use.
	Idle            int `json:"idleCon"`      // The number of idle connections.

	// Counters
	WaitCount         int64         `json:"waitCon"`           // The total number of connections waited for.
	WaitDuration      time.Duration `json:"waitDuration"`      // The total time blocked waiting for a new connection.
	MaxIdleClosed     int64         `json:"maxIdleClosed"`     // The total number of connections closed due to SetMaxIdleConns.
	MaxIdleTimeClosed int64         `json:"maxIdleTimeClosed"` // The total number of connections closed due to SetConnMaxIdleTime.
	MaxLifetimeClosed int64         `json:"maxLifetimeClosed"` // The total number of connections closed due to SetConnMaxLifetime.
}

//ServerStatus - полная инфомация о сервисе в системе Егеон
type ServerStatus struct {
	Info          ServerInfo    `json:"info"`
	StartDate     time.Time     `json:"startDate"`
	UpTime        time.Duration `json:"upTime"`
	UpTimeStr     string        `json:"upTimeStr"`
	SuccesReqCnt  uint64        `json:"succesReqCnt"`
	FaileReqCnt   uint64        `json:"faileReqCnt"`
	FaileGetCnt   uint64        `json:"faileGetCnt"`
	FailePostCnt  uint64        `json:"failePostCnt"`
	FailePutCnt   uint64        `json:"failePutCnt"`
	FaileDelCnt   uint64        `json:"faileDelCnt"`
	MiddleReqTime int64         `json:"middleReqTime"`
}

/*
Каждая из структур реализует отдельную сущность системы
И содержит свое отображение в базе данных
Каждая из структур имеет свою модель, которая предоставляет доступ к данным,
наполняя данными структуры
*/

//Status - Содержит информацию о статусе счетчика, корректора, объекта, узла учета
type Status struct {
	ID          uint16 `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

//Comment - определяет коментарии в системе
type Comment struct {
	ID           uint64    `json:"id"`
	Head         string    `json:"head"`
	Data         string    `json:"data"`
	ContentType  string    `json:"type"`
	UserAddID    uint32    `json:"userAddId"`
	AddedDate    time.Time `json:"addedDate"`
	ModifiedDate time.Time `json:"modifiedDate"`
	Status       Status    `json:"status"`
}

//Address - сущность определяющая адресс
type Address struct {
	ID            uint64    `json:"id"`
	Country       string    `json:"country"`
	Region        string    `json:"region"`
	City          string    `json:"city"`
	District      string    `json:"district"`
	MicroDistrict string    `json:"microDistrict"`
	Street        string    `json:"street"`
	Build         string    `json:"build"`
	Apartment     string    `json:"apartment"`
	Lat           float64   `json:"lat"`
	Lng           float64   `json:"lng"`
	PostCode      string    `json:"postCode"`
	Comment       Comment   `json:"comment"`
	AddedDate     time.Time `json:"addedDate"`
	ModifiedDate  time.Time `json:"modifDate"`
	AddedByID     uint32    `json:"addedById"`
	Status        Status    `json:"status"`
	FullName      string    `json:"fullName"`
}

//Group - функционал по группировке точек учета
// одни и те же точки учета могут быть в разных группах (связь много ко многим)
type Group struct {
	ID           uint64    `json:"id"`
	Name         string    `json:"name"`
	Description  string    `json:"description"`
	OwnerID      uint32    `json:"ownerId"`   // Создатель группы может быть NULL
	CompanyID    uint32    `json:"comapnyId"` // Если группа создана для компании
	ModifiedBy   uint32    `json:"modifiedBy"`
	AddedDate    time.Time `json:"addedDate"`
	ModifiedDate time.Time `json:"modifDate"`
	Logo         string    `json:"logo"`
	Status       Status    `json:"status"`
}

// Company - объект описывает компании в системе
type Company struct {
	ID           uint32    `json:"id"`
	Name         string    `json:"name"`
	Description  string    `json:"description"`
	Code         uint64    `json:"code"`
	Addr         Address   `json:"address"`
	UserAddID    uint32    `json:"userAddId"`
	IsControl    bool      `json:"isControl"`
	Logo         string    `json:"logo"`
	AddedDate    time.Time `json:"addedDate"`
	ModifiedDate time.Time `json:"modifDate"`
	Status       Status    `json:"status"`
}

// Role - роли
type Role struct {
	ID           uint32    `json:"id"`
	Name         string    `json:"name"`
	URL          string    `json:"url"`
	Method       string    `json:"method"`
	Description  string    `json:"description"`
	AddedDate    time.Time `json:"addedDate"`
	ModifiedDate time.Time `json:"modifDate"`
	UserAddID    uint32    `json:"userAddId"`
	Status       Status    `json:"status"`
}

//RoleSets - набор ролей, определяет наборы предустановленных ролей для определенный вариантов пользователей.
// Создан для упрощение назначения ролей пользователю
// В базе обязательно должен содержатся набор ролей с именем default - это базовые роли пользователя,
// которые будут ему назначены при регистрации в системе. Нужны для удобства в админке
type RoleSets struct {
	Name  string `json:"name"`
	Roles []Role `json:"roles"`
}

// SessionKey - Создание нового типа для ключа сессии позволяет осуществлять switch case по типу
type SessionKey string

// UsersGroup - TODO create user groups interface in the futer (пока не реализовано)
type UsersGroup struct {
	UserID       uint32    `json:"userId"`
	Group        Group     `json:"group"`
	IsUpdate     bool      `json:"isUpdate"`
	IsCreate     bool      `json:"isCreate"`
	IsDelete     bool      `json:"IsDelete"`
	UserAddID    uint32    `json:"userAddId"`
	AddedDate    time.Time `json:"addedDate"`
	ModifiedDate time.Time `json:"modifiedDate"`
}

// User - пользователи системы
type User struct {
	ID                uint32       `json:"id"`
	Profile           UserProfile  `json:"profile"`
	Email             string       `json:"email"`
	IsEmailConfirm    bool         `json:"isEmailConfirm"`
	PassHash          string       `json:"passHash"`
	Salt              string       `json:"salt"`
	Phone             string       `json:"phone"`
	Company           Company      `json:"company"`
	IsComapanyConfirm bool         `json:"isCompanyConfirm"`
	AccessFailedCount int          `json:"accesFailedCnt"`
	RestorePassword   bool         `json:"restorePassword"`
	LastActivity      time.Time    `json:"lastActivity"`
	AddedDate         time.Time    `json:"addedDate"`
	ModifiedDate      time.Time    `json:"modifDate"`
	Status            Status       `json:"status"`
	Roles             []Role       `json:"roles"`
	UsersGroups       []UsersGroup `json:"usersGroups"`
	SessionKey        SessionKey   `json:"sessionKey"`
	ExpireDate        time.Time    `json:"expired"`
}

//UserProfile - хранит данные профиля
type UserProfile struct {
	ID             uint32    `json:"id"`
	LastName       string    `json:"lastName"`
	Name           string    `json:"name"`
	PatronymicName string    `json:"nick"`
	Info           string    `json:"info"`
	Country        string    `json:"country"`
	Region         string    `json:"region"`
	City           string    `json:"city"`
	Lat            float64   `json:"lat"`
	Lng            float64   `json:"lng"`
	AvatarB64      string    `json:"avatarB64"`
	AddedDate      time.Time `json:"addedDate"`
	ModifiedDate   time.Time `json:"modifDate"`
}

// APIToken - структура определяет ключ доступа к API сервера
// Для системы ты будешь пользователем с ограниченным набором ролей
// Позволяет передать часть прав пользователя третьим лицам
// Строго на определенное время
type APIToken struct {
	ID          uint64    `json:"id"`
	OwnerID     uint32    `json:"ownerId"`
	Token       string    `json:"key"`
	Description string    `json:"description"`
	Roles       []Role    `json:"roles"`
	AddedDate   time.Time `json:"addedDate"`
	ExpireDate  time.Time `json:"expired"`
}

//UserLog - структура для хранения активностей пользователя
type UserLog struct {
	UserID       uint32        `json:"userId"`
	SessionKey   SessionKey    `json:"sessionKey"`
	IP           string        `json:"ip"`
	URL          string        `json:"url"`
	Method       string        `json:"method"`
	IsAborted    bool          `json:"isAborted"`
	ResponceTime time.Duration `json:"responceTime"`
	AddeDate     time.Time     `json:"addedDate"`
}
