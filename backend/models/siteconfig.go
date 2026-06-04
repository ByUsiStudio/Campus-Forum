package models

type SiteConfig struct {
	ID                   uint   `gorm:"primarykey" json:"id"`
	SiteTitle            string `gorm:"type:varchar(255);default:'校园论坛 - 分享与交流'" json:"site_title"`
	ICPNumber            string `gorm:"type:varchar(255)" json:"icp_number"`
	PublicSecurityNumber string `gorm:"type:varchar(255)" json:"public_security_number"`
	SMTPHost             string `gorm:"type:varchar(255)" json:"smtp_host"`
	SMTPPort             int    `gorm:"default:587" json:"smtp_port"`
	SMTPUsername         string `gorm:"type:varchar(255)" json:"smtp_username"`
	SMTPPassword         string `gorm:"type:varchar(255)" json:"smtp_password"`
	SMTPFrom             string `gorm:"type:varchar(255)" json:"smtp_from"`
	SMTPFromName         string `gorm:"type:varchar(255)" json:"smtp_from_name"`
}
