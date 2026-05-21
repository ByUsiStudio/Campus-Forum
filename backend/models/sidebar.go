package models

type SidebarConfig struct {
    ID      uint   `gorm:"primarykey" json:"id"`
    Items   string `gorm:"type:text" json:"items"` // JSON string
}

type SidebarItem struct {
    Title string `json:"title"`
    Link  string `json:"link"`
    Icon  string `json:"icon"`
}