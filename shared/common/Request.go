package common

// QueryParam is type for query parameters
type QueryParam string

// QueryDefault is type for default query values
type QueryDefault string

// Default query values
const (
	DefaultEmpty QueryDefault = ""

	DefaultPage                 QueryDefault = "1"
	DefaultSortDirAsc           QueryDefault = "asc"
	DefaultSortDirDesc          QueryDefault = "desc"
	DefaultSize10               QueryDefault = "10"
	DefaultSize20               QueryDefault = "20"
	DefaultSize30               QueryDefault = "30"
	DefaultOrderType            QueryDefault = "HMD"
	DefaultPartnerType          QueryDefault = "gofood"
	DefaultMenuGroupCode        QueryDefault = "G42"
	DefaultMenuItemCode         QueryDefault = "4839"
	DefaultLat                  QueryDefault = "0"
	DefaultLong                 QueryDefault = "0"
	DefaultTypeOnline           QueryDefault = "online"
	DefaultTypeOutletMap        QueryDefault = "outletMap"
	DefaultTypeOutletGroup      QueryDefault = "outletGroup"
	DefaultSortByCancelTime     QueryDefault = "cancel_time"
	DefaultCategoryOutlet       QueryDefault = "outlet"
	DefaultCategoryMonthly      QueryDefault = "monthly"
	DefaultPaymentMethodAll     QueryDefault = "all"
	DefaultOrderStatusCompleted QueryDefault = "COMPLETED"
	DefaultStatusPending        QueryDefault = "pending"
	PaymentMethodGopayRequest   QueryDefault = "gopay"
	DefaultPartnerTax           float64      = 10
	DefaultCdPrice              float64      = 35000
	DefaultCdTax                float64      = 3500
	DefaultAmount100            float64      = 100
	DefaultAmount110            float64      = 110
)

// Query parameters
const (
	QPage QueryParam = "page"
	QSize QueryParam = "size"

	SortBy  QueryParam = "sortBy"
	SortDir QueryParam = "sortDir"
)
