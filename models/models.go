package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Migration represents a schema change applied to the database.
type Migration struct {
	ID        uint      `gorm:"primaryKey;autoIncrement"`
	Name      string    `gorm:"type:text;not null"`
	AppliedAt time.Time `gorm:"not null"`
}

// Admin represents a system administrator.
type Admin struct {
	ID         uint   `gorm:"primaryKey;autoIncrement"`
	Username   string `gorm:"size:100;unique;not null"`
	Password   string `gorm:"size:255;not null"`
	TelegramID int64  `gorm:"not null;unique"`
	Email      string `gorm:"size:100"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

// Server holds connection information to a remote server.
type Server struct {
	ID             uint   `gorm:"primaryKey;autoIncrement"`
	Name           string `gorm:"size:100;not null"`
	UniqueName     string `gorm:"size:100"`
	Version        string `gorm:"size:50"`
	URL            string `gorm:"size:255"`
	Username       string `gorm:"size:100"`
	Password       string `gorm:"size:255"`
	RemovePrefix   string `gorm:"size:100"`
	TotalQuota     int64  `gorm:"not null;default:0"`
	UsedQuota      int64  `gorm:"not null;default:0"`
	RemainingQuota int64  `gorm:"not null;default:0"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

// Config holds configuration links for a server.
type Config struct {
	ID        uint   `gorm:"primaryKey;autoIncrement"`
	ServerID  uint   `gorm:"not null;index"`
	Server    Server `gorm:"foreignKey:ServerID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Link      string `gorm:"type:text;unique;not null"`
	Tag       string `gorm:"size:100"`
	CustomTag string `gorm:"size:100"`
	IsActive  bool   `gorm:"not null;default:true"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Reseller represents a reseller and supports recursive parent-child relationships.
type Reseller struct {
	ID                   uint       `gorm:"primaryKey;autoIncrement"`
	ParentResellerID     *uint      `gorm:"index"`
	ParentReseller       *Reseller  `gorm:"foreignKey:ParentResellerID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
	ChildResellers       []Reseller `gorm:"foreignKey:ParentResellerID"`
	Username             string     `gorm:"size:100;unique;not null"`
	Password             string     `gorm:"size:255;not null"`
	TelegramID           int64      `gorm:"not null;unique"`
	Email                string     `gorm:"size:100"`
	APIToken             string     `gorm:"size:512"`
	CustomName           string     `gorm:"size:100;unique"`
	TotalQuota           int64      `gorm:"not null;default:0"`
	UsedQuota            int64      `gorm:"not null;default:0"`
	RemainingQuota       int64      `gorm:"not null;default:0"`
	Status               string     `gorm:"size:20"`
	CanCreateSubReseller bool       `gorm:"not null;default:false;column:can_create_subreseller"`
	CreatedAt            time.Time
	UpdatedAt            time.Time
}

// User represents an end-user of the service.
type User struct {
	ID             uint      `gorm:"primaryKey;autoIncrement"`
	ResellerID     uint      `gorm:"not null;index"`
	Reseller       Reseller  `gorm:"foreignKey:ResellerID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
	Username       string    `gorm:"size:100;unique;not null"`
	UUID           uuid.UUID `gorm:"type:uuid;unique;not null"`
	TotalQuota     int64     `gorm:"not null;default:0"`
	UsedQuota      int64     `gorm:"not null;default:0"`
	RemainingQuota int64     `gorm:"not null;default:0"`
	ActivationDate time.Time
	Status         string `gorm:"size:20"`
	SubToken       string `gorm:"size:512"`
	TimeLimit      *time.Time
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

// TagSystem groups configurations under a common name.
type TagSystem struct {
	ID        uint              `gorm:"primaryKey;autoIncrement"`
	Name      string            `gorm:"size:100;unique;not null"`
	IsActive  bool              `gorm:"not null;default:true"`
	Configs   []TagSystemConfig `gorm:"foreignKey:TagSystemID"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (TagSystem) TableName() string {
	return "tag_system"
}

// TagSystemConfig connects TagSystems with Configs.
type TagSystemConfig struct {
	ID          uint      `gorm:"primaryKey;autoIncrement"`
	TagSystemID uint      `gorm:"not null;index"`
	TagSystem   TagSystem `gorm:"foreignKey:TagSystemID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	ConfigID    uint      `gorm:"not null;index"`
	Config      Config    `gorm:"foreignKey:ConfigID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func (TagSystemConfig) TableName() string {
	return "tag_system_config"
}

// ResellerTagSystemAllowed defines which resellers can access specific tag systems.
type ResellerTagSystemAllowed struct {
	ID          uint      `gorm:"primaryKey;autoIncrement"`
	TagSystemID uint      `gorm:"not null;index"`
	TagSystem   TagSystem `gorm:"foreignKey:TagSystemID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	ResellerID  uint      `gorm:"not null;index"`
	Reseller    Reseller  `gorm:"foreignKey:ResellerID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

// TagReseller is a label that can be assigned to resellers.
type TagReseller struct {
	ID       uint   `gorm:"primaryKey;autoIncrement"`
	Name     string `gorm:"size:100;unique;not null"`
	IsActive bool   `gorm:"not null;default:true"`
}

// TagResellerConfig associates resellers with configs and optional custom names.
type TagResellerConfig struct {
	ID               uint        `gorm:"primaryKey;autoIncrement"`
	TagResellerID    uint        `gorm:"not null;index"`
	TagReseller      TagReseller `gorm:"foreignKey:TagResellerID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	ResellerID       uint        `gorm:"not null;index"`
	Reseller         Reseller    `gorm:"foreignKey:ResellerID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	ConfigID         uint        `gorm:"not null;index"`
	Config           Config      `gorm:"foreignKey:ConfigID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	ConfigCustomName string      `gorm:"size:100"`
}

// UserTagReseller is the join table between users and tag resellers.
type UserTagReseller struct {
	ID            uint        `gorm:"primaryKey;autoIncrement"`
	UserID        uint        `gorm:"not null;index"`
	User          User        `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	TagResellerID uint        `gorm:"not null;index"`
	TagReseller   TagReseller `gorm:"foreignKey:TagResellerID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

// ResellerConfig provides a many-to-many relationship between resellers and configs.
type ResellerConfig struct {
	ResellerID uint     `gorm:"primaryKey;not null"`
	Reseller   Reseller `gorm:"foreignKey:ResellerID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	ConfigID   uint     `gorm:"primaryKey;not null"`
	Config     Config   `gorm:"foreignKey:ConfigID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	IsActive   bool     `gorm:"not null;default:true"`
	CustomName string   `gorm:"size:100"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

// UserConfig links a user to a configuration.
type UserConfig struct {
	UserID     uint   `gorm:"primaryKey;not null"`
	User       User   `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	ConfigID   uint   `gorm:"primaryKey;not null"`
	Config     Config `gorm:"foreignKey:ConfigID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	AssignedAt time.Time
}

// UsageLog stores bandwidth usage statistics.
type UsageLog struct {
	ID         uint      `gorm:"primaryKey;autoIncrement"`
	UserID     uint      `gorm:"not null;index"`
	User       User      `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	ResellerID uint      `gorm:"not null;index"`
	Reseller   Reseller  `gorm:"foreignKey:ResellerID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	ServerID   uint      `gorm:"not null;index"`
	Server     Server    `gorm:"foreignKey:ServerID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Amount     int64     `gorm:"not null"`
	CreatedAt  time.Time `gorm:"autoCreateTime"`
}

// UserProxy stores generated proxy information for users.
type UserProxy struct {
	ID        uint      `gorm:"primaryKey;autoIncrement"`
	UserID    uint      `gorm:"not null;index"`
	User      User      `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	ServerID  uint      `gorm:"not null;index;uniqueIndex:idx_server_uuid"`
	Server    Server    `gorm:"foreignKey:ServerID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Type      string    `gorm:"size:50"`
	UUID      uuid.UUID `gorm:"type:uuid;index;uniqueIndex:idx_server_uuid"`
	Payload   []byte    `gorm:"type:jsonb"`
	CreatedAt time.Time
}

// Transaction tracks quota changes or other financial operations.
type Transaction struct {
	ID            uint     `gorm:"primaryKey;autoIncrement"`
	ResellerID    uint     `gorm:"not null;index"`
	Reseller      Reseller `gorm:"foreignKey:ResellerID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Type          string   `gorm:"size:20"`
	Amount        int64
	PerformedByID uint     `gorm:"not null;index"`
	PerformedBy   Reseller `gorm:"foreignKey:PerformedByID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Description   string   `gorm:"type:text"`
	ReferenceType string   `gorm:"size:50"`
	ReferenceID   uint
	CreatedAt     time.Time
}

// FilteredInbound represents a blocked inbound tag.
type FilteredInbound struct {
	ID         uint   `gorm:"primaryKey;autoIncrement"`
	InboundTag string `gorm:"size:100;unique;not null"`
}

// FilteredTag represents a blocked tag.
type FilteredTag struct {
	ID  uint   `gorm:"primaryKey;autoIncrement"`
	Tag string `gorm:"size:100;unique;not null"`
}

// ActivityLog tracks actions occurring within the system.
type ActivityLog struct {
	ID          uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	EventType   string         `gorm:"type:varchar(100);not null" json:"event_type"`
	Entity      string         `gorm:"type:varchar(100);not null" json:"entity"`
	EntityID    int            `gorm:"not null" json:"entity_id"`
	RequestID   string         `gorm:"type:varchar(255)" json:"request_id"`
	PerformedBy string         `gorm:"type:varchar(100)" json:"performed_by"`
	IPAddress   string         `gorm:"type:varchar(45)" json:"ip_address"`
	UserAgent   string         `gorm:"type:varchar(255)" json:"user_agent"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

// LogEntry represents a log event that can be stored persistently.
type LogEntry struct {
	ID        uint   `gorm:"primaryKey;autoIncrement"`
	Level     string `gorm:"size:20;index"`
	Message   string `gorm:"type:text"`
	Context   string `gorm:"type:json"`
	CreatedAt time.Time
}
