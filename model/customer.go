package model

type Customer struct {
	BaseModel
	Code               string  `gorm:"type:varchar(32);uniqueIndex"`
	Name               string  `gorm:"not null;type:varchar(128)"`
	Pic                string  `gorm:"type:varchar(128)"`
	Npwp               string  `gorm:"type:varchar(32)"`
	Fax                string  `gorm:"type:varchar(32)"`
	Email              string  `gorm:"type:varchar(128)"`
	Phone              string  `gorm:"type:varchar(32)"`
	Address            string  `gorm:"type:varchar(255)"`
	PostalCode         string  `gorm:"type:varchar(16)"`
	CylinderQuotaLimit int     `gorm:"not null;default:0"`
	CreditLimit        float64 `gorm:"type:decimal(15,2);not null;default:0"`
	// TermOfPaymentId opsional ke master TOP; belum dipakai hitung jatuh tempo (Invoice.DueDate pakai PaymentTermDays).
	TermOfPaymentId    *string        `gorm:"type:varchar(36);index"`
	TermOfPayment      *TermOfPayment `gorm:"foreignKey:TermOfPaymentId"`
	SalesPicName       string         `gorm:"type:varchar(128)"`
	SalesPicPhone      string         `gorm:"type:varchar(32)"`
	SalesPicEmail      string         `gorm:"type:varchar(128)"`
	OutstandingCount   int            `gorm:"not null;default:0"`
	IsActive           bool           `gorm:"not null;default:true"`
	PaymentTermDays    int            `gorm:"not null;default:30"`
	OutstandingBalance float64        `gorm:"type:decimal(15,2);not null;default:0"`
}
