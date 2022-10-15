package golang

import "time"

//ServerInfo - статическая информация о сервисе. Зависит от сборки системы
type ServerInfo struct {
	Name       string                 `json:"name"`
	Version    string                 `json:"version"`
	About      string                 `json:"about,omitempty"`
	Maintainer string                 `json:"maintainer,omitempty"`
	Routes     map[string]string      `json:"routes,omitempty"`
	BaseTypes  map[string]interface{} `json:"types,omitempty"`
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
	Info          ServerInfo             `json:"info,omitempty"`
	Addition      map[string]interface{} `json:"addition,omitempty"`
	StartDate     time.Time              `json:"startDate,omitempty"`
	UpTime        time.Duration          `json:"upTime,omitempty"`
	UpTimeStr     string                 `json:"upTimeStr,omitempty"`
	SuccesReqCnt  uint64                 `json:"succesReqCnt,omitempty"`
	FaileReqCnt   uint64                 `json:"faileReqCnt,omitempty"`
	FaileGetCnt   uint64                 `json:"faileGetCnt,omitempty"`
	FailePostCnt  uint64                 `json:"failePostCnt,omitempty"`
	FailePutCnt   uint64                 `json:"failePutCnt,omitempty"`
	FaileDelCnt   uint64                 `json:"faileDelCnt,omitempty"`
	MiddleReqTime int64                  `json:"middleReqTime,omitempty"`
}

/*
Каждая из структур реализует отдельную сущность системы
И содержит свое отображение в базе данных
Каждая из структур имеет свою модель, которая предоставляет доступ к данным,
наполняя данными структуры
*/

//Status - Содержит информацию о статусе счетчика, корректора, объекта, узла учета
type Status struct {
	ID          uint16 `json:"id,omitempty"`
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
}

//Comment - определяет коментарии в системе
type Comment struct {
	ID           uint64    `json:"id,omitempty"`
	Head         string    `json:"head,omitempty"`
	Data         string    `json:"data,omitempty"`
	ContentType  string    `json:"type,omitempty"`
	UserAddID    uint32    `json:"userAddId,omitempty"`
	AddedDate    time.Time `json:"addedDate,omitempty"`
	ModifiedDate time.Time `json:"modifiedDate,omitempty"`
	Status       Status    `json:"status,omitempty"`
}

//Address - сущность определяющая адресс
type Address struct {
	ID            uint64    `json:"id"`
	Country       string    `json:"country"`
	Region        string    `json:"region"`
	City          string    `json:"city"`
	District      string    `json:"district,omitempty"`
	MicroDistrict string    `json:"microDistrict,omitempty"`
	Street        string    `json:"street"`
	Build         string    `json:"build"`
	Apartment     string    `json:"apartment,omitempty"`
	Lat           float64   `json:"lat"`
	Lng           float64   `json:"lng"`
	PostCode      string    `json:"postCode,omitempty"`
	Comment       Comment   `json:"comment,omitempty"`
	AddedDate     time.Time `json:"addedDate,omitempty"`
	ModifiedDate  time.Time `json:"modifDate,omitempty"`
	AddedByID     uint32    `json:"addedById,omitempty"`
	Status        Status    `json:"status,omitempty"`
	FullName      string    `json:"fullName,omitempty"`
}

//Group - функционал по группировке точек учета
// одни и те же точки учета могут быть в разных группах (связь много ко многим)
type Group struct {
	ID           uint64    `json:"id"`
	Name         string    `json:"name"`
	Description  string    `json:"description,omitempty"`
	OwnerID      uint32    `json:"ownerId,omitempty"`   // Создатель группы может быть NULL
	CompanyID    uint32    `json:"comapnyId,omitempty"` // Если группа создана для компании
	ModifiedBy   uint32    `json:"modifiedBy,omitempty"`
	AddedDate    time.Time `json:"addedDate,omitempty"`
	ModifiedDate time.Time `json:"modifDate,omitempty"`
	Logo         string    `json:"logo,omitempty"`
	Status       Status    `json:"status,omitempty"`
}

// Company - объект описывает компании в системе
type Company struct {
	ID           uint32    `json:"id"`
	Name         string    `json:"name"`
	Description  string    `json:"description,omitempty"`
	Code         uint64    `json:"code,omitempty"`
	Addr         Address   `json:"address,omitempty"`
	UserAddID    uint32    `json:"userAddId,omitempty"`
	IsControl    bool      `json:"isControl,omitempty"`
	Logo         string    `json:"logo,omitempty"`
	AddedDate    time.Time `json:"addedDate,omitempty"`
	ModifiedDate time.Time `json:"modifDate,omitempty"`
	Status       Status    `json:"status,omitempty"`
}

// Role - роли
type Role struct {
	ID           uint32    `json:"id"`
	Name         string    `json:"name"`
	URL          string    `json:"url"`
	Method       string    `json:"method"`
	Description  string    `json:"description,omitempty"`
	AddedDate    time.Time `json:"addedDate,omitempty"`
	ModifiedDate time.Time `json:"modifDate,omitempty"`
	UserAddID    uint32    `json:"userAddId,omitempty"`
	Status       Status    `json:"status,omitempty"`
}

//RoleSets - набор ролей, определяет наборы предустановленных ролей для определенный вариантов пользователей.
// Создан для упрощение назначения ролей пользователю
// В базе обязательно должен содержатся набор ролей с именем default - это базовые роли пользователя,
// которые будут ему назначены при регистрации в системе. Нужны для удобства в админке
type RoleSets struct {
	Name  string `json:"name,omitempty"`
	Roles []Role `json:"roles,omitempty"`
}

// SessionKey - Создание нового типа для ключа сессии позволяет осуществлять switch case по типу
type SessionKey string

// UsersGroup - TODO create user groups interface in the futer (пока не реализовано)
type UsersGroup struct {
	UserID       uint32    `json:"userId"`
	Group        Group     `json:"group,omitempty"`
	IsUpdate     bool      `json:"isUpdate"`
	IsCreate     bool      `json:"isCreate"`
	IsDelete     bool      `json:"isDelete"`
	UserAddID    uint32    `json:"userAddId,omitempty"`
	AddedDate    time.Time `json:"addedDate,omitempty"`
	ModifiedDate time.Time `json:"modifiedDate,omitempty"`
}

// User - пользователи системы
type User struct {
	ID                uint32       `json:"id"`
	Profile           UserProfile  `json:"profile,omitempty"`
	Email             string       `json:"email"`
	IsEmailConfirm    bool         `json:"isEmailConfirm"`
	PassHash          string       `json:"passHash,omitempty"`
	Salt              string       `json:"salt,omitempty"`
	Phone             string       `json:"phone,omitempty"`
	Company           Company      `json:"company,omitempty"`
	IsComapanyConfirm bool         `json:"isCompanyConfirm,omitempty"`
	AccessFailedCount int          `json:"accesFailedCnt,omitempty"`
	RestorePassword   bool         `json:"restorePassword,omitempty"`
	LastActivity      time.Time    `json:"lastActivity,omitempty"`
	AddedDate         time.Time    `json:"addedDate,omitempty"`
	ModifiedDate      time.Time    `json:"modifDate,omitempty"`
	Status            Status       `json:"status,omitempty"`
	Roles             []Role       `json:"roles,omitempty"`
	UsersGroups       []UsersGroup `json:"usersGroups,omitempty"`
	SessionKey        SessionKey   `json:"sessionKey,omitempty"`
	ExpireDate        time.Time    `json:"expired,omitempty"`
}

//UserProfile - хранит данные профиля
type UserProfile struct {
	ID             uint32    `json:"id,omitempty"`
	LastName       string    `json:"lastName,omitempty"`
	Name           string    `json:"name,omitempty"`
	PatronymicName string    `json:"nick,omitempty"`
	Info           string    `json:"info,omitempty"`
	Country        string    `json:"country,omitempty"`
	Region         string    `json:"region,omitempty"`
	City           string    `json:"city,omitempty"`
	Lat            float64   `json:"lat,omitempty"`
	Lng            float64   `json:"lng,omitempty"`
	AvatarB64      string    `json:"avatarB64,omitempty"`
	AddedDate      time.Time `json:"addedDate,omitempty"`
	ModifiedDate   time.Time `json:"modifDate,omitempty"`
}

// APIToken - структура определяет ключ доступа к API сервера
// Для системы ты будешь пользователем с ограниченным набором ролей
// Позволяет передать часть прав пользователя третьим лицам
// Строго на определенное время
type APIToken struct {
	ID          uint64    `json:"id"`
	OwnerID     uint32    `json:"ownerId,omitempty"`
	Token       string    `json:"token"`
	Description string    `json:"description,omitempty"`
	Roles       []Role    `json:"roles,omitempty"`
	AddedDate   time.Time `json:"addedDate"`
	ExpireDate  time.Time `json:"expired"`
}

//UserLog - структура для хранения активностей пользователя
type UserLog struct {
	UserID       uint32        `json:"userId"`
	SessionKey   SessionKey    `json:"sessionKey,omitempty"`
	IP           string        `json:"ip,omitempty"`
	URL          string        `json:"url"`
	Method       string        `json:"method"`
	IsAborted    bool          `json:"isAborted"`
	RequestID    string        `json:"requestId"`
	UserAgent    string        `json:"browser,omitempty"`
	ResponceTime time.Duration `json:"responceTime"`
	AddeDate     time.Time     `json:"addedDate"`
}
