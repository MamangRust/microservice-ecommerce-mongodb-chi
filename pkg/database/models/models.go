package models

import "time"

// ---------------------------------------------------------------------------
// Core
// ---------------------------------------------------------------------------

type User struct {
	UserID           int32      `bson:"_id" json:"user_id"`
	Firstname        string     `bson:"firstname" json:"firstname"`
	Lastname         string     `bson:"lastname" json:"lastname"`
	Email            string     `bson:"email" json:"email"`
	Password         string     `bson:"password" json:"password"`
	VerificationCode string     `bson:"verification_code" json:"verification_code"`
	IsVerified       *bool      `bson:"is_verified" json:"is_verified"`
	CreatedAt        *time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt        *time.Time `bson:"updated_at" json:"updated_at"`
	DeletedAt        *time.Time `bson:"deleted_at,omitempty" json:"deleted_at"`
}

type Role struct {
	RoleID    int32      `bson:"_id" json:"role_id"`
	RoleName  string     `bson:"role_name" json:"role_name"`
	CreatedAt *time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt *time.Time `bson:"updated_at" json:"updated_at"`
	DeletedAt *time.Time `bson:"deleted_at,omitempty" json:"deleted_at"`
}

type UserRole struct {
	UserRoleID int32      `bson:"_id" json:"user_role_id"`
	UserID     int32      `bson:"user_id" json:"user_id"`
	RoleID     int32      `bson:"role_id" json:"role_id"`
	CreatedAt  *time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt  *time.Time `bson:"updated_at" json:"updated_at"`
	DeletedAt  *time.Time `bson:"deleted_at,omitempty" json:"deleted_at"`
}

// ---------------------------------------------------------------------------
// Merchant
// ---------------------------------------------------------------------------

type Merchant struct {
	MerchantID   int32      `bson:"_id" json:"merchant_id"`
	UserID       int32      `bson:"user_id" json:"user_id"`
	Name         string     `bson:"name" json:"name"`
	Description  *string    `bson:"description" json:"description"`
	Address      *string    `bson:"address" json:"address"`
	ContactEmail *string    `bson:"contact_email" json:"contact_email"`
	ContactPhone *string    `bson:"contact_phone" json:"contact_phone"`
	Status       string     `bson:"status" json:"status"`
	CreatedAt    *time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt    *time.Time `bson:"updated_at" json:"updated_at"`
	DeletedAt    *time.Time `bson:"deleted_at,omitempty" json:"deleted_at"`
}

type MerchantDocument struct {
	DocumentID   int32      `bson:"_id" json:"document_id"`
	MerchantID   int32      `bson:"merchant_id" json:"merchant_id"`
	DocumentType string     `bson:"document_type" json:"document_type"`
	DocumentUrl  string     `bson:"document_url" json:"document_url"`
	Status       string     `bson:"status" json:"status"`
	Note         *string    `bson:"note" json:"note"`
	UploadedAt   *time.Time `bson:"uploaded_at" json:"uploaded_at"`
	CreatedAt    *time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt    *time.Time `bson:"updated_at" json:"updated_at"`
	DeletedAt    *time.Time `bson:"deleted_at,omitempty" json:"deleted_at"`
}

type MerchantDetail struct {
	MerchantDetailID int32      `bson:"_id" json:"merchant_detail_id"`
	MerchantID       int32      `bson:"merchant_id" json:"merchant_id"`
	DisplayName      *string    `bson:"display_name" json:"display_name"`
	CoverImageUrl    *string    `bson:"cover_image_url" json:"cover_image_url"`
	LogoUrl          *string    `bson:"logo_url" json:"logo_url"`
	ShortDescription *string    `bson:"short_description" json:"short_description"`
	WebsiteUrl       *string    `bson:"website_url" json:"website_url"`
	CreatedAt        *time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt        *time.Time `bson:"updated_at" json:"updated_at"`
	DeletedAt        *time.Time `bson:"deleted_at,omitempty" json:"deleted_at"`
}

type MerchantSocialMediaLink struct {
	MerchantSocialID int32      `bson:"_id" json:"merchant_social_id"`
	MerchantDetailID int32      `bson:"merchant_detail_id" json:"merchant_detail_id"`
	Platform         string     `bson:"platform" json:"platform"`
	Url              string     `bson:"url" json:"url"`
	CreatedAt        *time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt        *time.Time `bson:"updated_at" json:"updated_at"`
	DeletedAt        *time.Time `bson:"deleted_at,omitempty" json:"deleted_at"`
}

type MerchantBusinessInformation struct {
	MerchantBusinessInfoID int32      `bson:"_id" json:"merchant_business_info_id"`
	MerchantID             int32      `bson:"merchant_id" json:"merchant_id"`
	BusinessType           *string    `bson:"business_type" json:"business_type"`
	TaxID                  *string    `bson:"tax_id" json:"tax_id"`
	EstablishedYear        *int32     `bson:"established_year" json:"established_year"`
	NumberOfEmployees      *int32     `bson:"number_of_employees" json:"number_of_employees"`
	WebsiteUrl             *string    `bson:"website_url" json:"website_url"`
	CreatedAt              *time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt              *time.Time `bson:"updated_at" json:"updated_at"`
	DeletedAt              *time.Time `bson:"deleted_at,omitempty" json:"deleted_at"`
}

type MerchantPolicy struct {
	MerchantPolicyID int32      `bson:"_id" json:"merchant_policy_id"`
	MerchantID       int32      `bson:"merchant_id" json:"merchant_id"`
	PolicyType       string     `bson:"policy_type" json:"policy_type"`
	Title            string     `bson:"title" json:"title"`
	Description      string     `bson:"description" json:"description"`
	CreatedAt        *time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt        *time.Time `bson:"updated_at" json:"updated_at"`
	DeletedAt        *time.Time `bson:"deleted_at,omitempty" json:"deleted_at"`
}

type MerchantCertificationsAndAward struct {
	MerchantCertificationID int32      `bson:"_id" json:"merchant_certification_id"`
	MerchantID              int32      `bson:"merchant_id" json:"merchant_id"`
	Title                   string     `bson:"title" json:"title"`
	Description             *string    `bson:"description" json:"description"`
	IssuedBy                *string    `bson:"issued_by" json:"issued_by"`
	IssueDate               *time.Time `bson:"issue_date" json:"issue_date"`
	ExpiryDate              *time.Time `bson:"expiry_date" json:"expiry_date"`
	CertificateUrl          *string    `bson:"certificate_url" json:"certificate_url"`
	CreatedAt               *time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt               *time.Time `bson:"updated_at" json:"updated_at"`
	DeletedAt               *time.Time `bson:"deleted_at,omitempty" json:"deleted_at"`
}

// ---------------------------------------------------------------------------
// Product
// ---------------------------------------------------------------------------

type Product struct {
	ProductID    int32      `bson:"_id" json:"product_id"`
	MerchantID   int32      `bson:"merchant_id" json:"merchant_id"`
	CategoryID   int32      `bson:"category_id" json:"category_id"`
	Name         string     `bson:"name" json:"name"`
	Description  *string    `bson:"description" json:"description"`
	Price        int32      `bson:"price" json:"price"`
	CountInStock int32      `bson:"count_in_stock" json:"count_in_stock"`
	Brand        *string    `bson:"brand" json:"brand"`
	Weight       *int32     `bson:"weight" json:"weight"`
	Rating       *float64   `bson:"rating" json:"rating"`
	SlugProduct  *string    `bson:"slug_product" json:"slug_product"`
	ImageProduct *string    `bson:"image_product" json:"image_product"`
	CreatedAt    *time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt    *time.Time `bson:"updated_at" json:"updated_at"`
	DeletedAt    *time.Time `bson:"deleted_at,omitempty" json:"deleted_at"`
}

type ProductStockAdjustment struct {
	OperationID string     `bson:"_id" json:"operation_id"`
	ProductID   int32      `bson:"product_id" json:"product_id"`
	Delta       int32      `bson:"delta" json:"delta"`
	CreatedAt   *time.Time `bson:"created_at" json:"created_at"`
}

// ---------------------------------------------------------------------------
// Order
// ---------------------------------------------------------------------------

type Order struct {
	OrderID    int32      `bson:"_id" json:"order_id"`
	UserID     int32      `bson:"user_id" json:"user_id"`
	MerchantID int32      `bson:"merchant_id" json:"merchant_id"`
	TotalPrice int32      `bson:"total_price" json:"total_price"`
	CreatedAt  *time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt  *time.Time `bson:"updated_at" json:"updated_at"`
	DeletedAt  *time.Time `bson:"deleted_at,omitempty" json:"deleted_at"`
}

type OrderItem struct {
	OrderItemID int32      `bson:"_id" json:"order_item_id"`
	OrderID     int32      `bson:"order_id" json:"order_id"`
	ProductID   int32      `bson:"product_id" json:"product_id"`
	Quantity    int32      `bson:"quantity" json:"quantity"`
	Price       int32      `bson:"price" json:"price"`
	CreatedAt   *time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt   *time.Time `bson:"updated_at" json:"updated_at"`
	DeletedAt   *time.Time `bson:"deleted_at,omitempty" json:"deleted_at"`
}

// TableName returns the physical collection name for OrderItem.
func (o *OrderItem) TableName() string { return "order_items" }

type OrderStockReservation struct {
	ReservationID int32      `bson:"_id" json:"reservation_id"`
	OrderID       int32      `bson:"order_id" json:"order_id"`
	ProductID     int32      `bson:"product_id" json:"product_id"`
	Quantity      int32      `bson:"quantity" json:"quantity"`
	Status        string     `bson:"status" json:"status"`
	CreatedAt     *time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt     *time.Time `bson:"updated_at" json:"updated_at"`
}

// ---------------------------------------------------------------------------
// Auth
// ---------------------------------------------------------------------------

type RefreshToken struct {
	RefreshTokenID int32      `bson:"_id" json:"refresh_token_id"`
	UserID         int32      `bson:"user_id" json:"user_id"`
	Token          string     `bson:"token" json:"token"`
	Expiration     time.Time  `bson:"expiration" json:"expiration"`
	CreatedAt      *time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt      *time.Time `bson:"updated_at" json:"updated_at"`
	DeletedAt      *time.Time `bson:"deleted_at,omitempty" json:"deleted_at"`
}

type ResetToken struct {
	ResetTokenID int32     `bson:"_id" json:"id"`
	UserID       int64     `bson:"user_id" json:"user_id"`
	Token        string    `bson:"token" json:"token"`
	ExpiryDate   time.Time `bson:"expiry_date" json:"expiry_date"`
}

// ---------------------------------------------------------------------------
// Content
// ---------------------------------------------------------------------------

type Banner struct {
	BannerID  int32      `bson:"_id" json:"banner_id"`
	Name      string     `bson:"name" json:"name"`
	StartDate *time.Time `bson:"start_date" json:"start_date"`
	EndDate   *time.Time `bson:"end_date" json:"end_date"`
	StartTime *string    `bson:"start_time" json:"start_time"`
	EndTime   *string    `bson:"end_time" json:"end_time"`
	IsActive  *bool      `bson:"is_active" json:"is_active"`
	CreatedAt *time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt *time.Time `bson:"updated_at" json:"updated_at"`
	DeletedAt *time.Time `bson:"deleted_at,omitempty" json:"deleted_at"`
}

type Slider struct {
	SliderID  int32      `bson:"_id" json:"slider_id"`
	Name      string     `bson:"name" json:"name"`
	Image     string     `bson:"image" json:"image"`
	CreatedAt *time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt *time.Time `bson:"updated_at" json:"updated_at"`
	DeletedAt *time.Time `bson:"deleted_at,omitempty" json:"deleted_at"`
}

type Category struct {
	CategoryID    int32      `bson:"_id" json:"category_id"`
	Name          string     `bson:"name" json:"name"`
	Description   *string    `bson:"description" json:"description"`
	SlugCategory  *string    `bson:"slug_category" json:"slug_category"`
	ImageCategory *string    `bson:"image_category" json:"image_category"`
	CreatedAt     *time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt     *time.Time `bson:"updated_at" json:"updated_at"`
	DeletedAt     *time.Time `bson:"deleted_at,omitempty" json:"deleted_at"`
}

// ---------------------------------------------------------------------------
// User Features
// ---------------------------------------------------------------------------

type ShippingAddress struct {
	ShippingAddressID int32      `bson:"_id" json:"shipping_address_id"`
	OrderID           int32      `bson:"order_id" json:"order_id"`
	Alamat            string     `bson:"alamat" json:"alamat"`
	Provinsi          string     `bson:"provinsi" json:"provinsi"`
	Negara            string     `bson:"negara" json:"negara"`
	Kota              string     `bson:"kota" json:"kota"`
	Courier           string     `bson:"courier" json:"courier"`
	ShippingMethod    string     `bson:"shipping_method" json:"shipping_method"`
	ShippingCost      float64    `bson:"shipping_cost" json:"shipping_cost"`
	CreatedAt         *time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt         *time.Time `bson:"updated_at" json:"updated_at"`
	DeletedAt         *time.Time `bson:"deleted_at,omitempty" json:"deleted_at"`
}

type Cart struct {
	CartID    int32      `bson:"_id" json:"cart_id"`
	UserID    int32      `bson:"user_id" json:"user_id"`
	ProductID int32      `bson:"product_id" json:"product_id"`
	Name      string     `bson:"name" json:"name"`
	Price     int32      `bson:"price" json:"price"`
	Image     string     `bson:"image" json:"image"`
	Quantity  int32      `bson:"quantity" json:"quantity"`
	Weight    int32      `bson:"weight" json:"weight"`
	CreatedAt *time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt *time.Time `bson:"updated_at" json:"updated_at"`
	DeletedAt *time.Time `bson:"deleted_at,omitempty" json:"deleted_at"`
}

// ---------------------------------------------------------------------------
// Reviews
// ---------------------------------------------------------------------------

type Review struct {
	ReviewID  int32      `bson:"_id" json:"review_id"`
	UserID    int32      `bson:"user_id" json:"user_id"`
	ProductID int32      `bson:"product_id" json:"product_id"`
	Name      string     `bson:"name" json:"name"`
	Comment   string     `bson:"comment" json:"comment"`
	Rating    int32      `bson:"rating" json:"rating"`
	CreatedAt *time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt *time.Time `bson:"updated_at" json:"updated_at"`
	DeletedAt *time.Time `bson:"deleted_at,omitempty" json:"deleted_at"`
}

type ReviewDetail struct {
	ReviewDetailID int32      `bson:"_id" json:"review_detail_id"`
	ReviewID       int32      `bson:"review_id" json:"review_id"`
	Type           string     `bson:"type" json:"type"`
	Url            string     `bson:"url" json:"url"`
	Caption        *string    `bson:"caption" json:"caption"`
	CreatedAt      *time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt      *time.Time `bson:"updated_at" json:"updated_at"`
	DeletedAt      *time.Time `bson:"deleted_at,omitempty" json:"deleted_at"`
}

// ---------------------------------------------------------------------------
// Transaction
// ---------------------------------------------------------------------------

type Transaction struct {
	TransactionID int32      `bson:"_id" json:"transaction_id"`
	OrderID       int32      `bson:"order_id" json:"order_id"`
	MerchantID    int32      `bson:"merchant_id" json:"merchant_id"`
	PaymentMethod string     `bson:"payment_method" json:"payment_method"`
	Amount        int32      `bson:"amount" json:"amount"`
	PaymentStatus string     `bson:"payment_status" json:"payment_status"`
	CreatedAt     *time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt     *time.Time `bson:"updated_at" json:"updated_at"`
	DeletedAt     *time.Time `bson:"deleted_at,omitempty" json:"deleted_at"`
}

// ---------------------------------------------------------------------------
// Infrastructure
// ---------------------------------------------------------------------------

type OutboxEvent struct {
	OutboxID      int64      `bson:"_id" json:"outbox_id"`
	Topic         string     `bson:"topic" json:"topic"`
	EventKey      string     `bson:"event_key" json:"event_key"`
	Payload       []byte     `bson:"payload" json:"payload"`
	Status        string     `bson:"status" json:"status"`
	Attempts      int32      `bson:"attempts" json:"attempts"`
	NextAttemptAt time.Time  `bson:"next_attempt_at" json:"next_attempt_at"`
	CreatedAt     *time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt     *time.Time `bson:"updated_at" json:"updated_at"`
}

type ConsumerInbox struct {
	ConsumerName       string     `bson:"consumer_name" json:"consumer_name"`
	EventKey           string     `bson:"event_key" json:"event_key"`
	Topic              string     `bson:"topic" json:"topic"`
	PartitionID        int32      `bson:"partition_id" json:"partition_id"`
	MessageOffset      int64      `bson:"message_offset" json:"message_offset"`
	Status             string     `bson:"status" json:"status"`
	Attempts           int32      `bson:"attempts" json:"attempts"`
	LeaseUntil         *time.Time `bson:"lease_until" json:"lease_until"`
	LastError          string     `bson:"last_error" json:"last_error"`
	ProcessedAt        *time.Time `bson:"processed_at" json:"processed_at"`
	ReservationVersion int64      `bson:"reservation_version" json:"reservation_version"`
}
