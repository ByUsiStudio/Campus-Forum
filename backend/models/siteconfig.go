package models

type SiteConfig struct {
	ID        uint   `gorm:"primarykey" json:"id"`
	SiteTitle string `gorm:"type:varchar(255);default:'校园论坛 - 分享与交流'" json:"site_title"`
}
