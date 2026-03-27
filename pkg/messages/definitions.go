package messages

// --- Success Messages ---

// Use these predefined success types for standard operations
var (
	SuccessCreate = newSuccess("create", "%s berhasil disimpan", "%s successfully saved")
	SuccessUpdate = newSuccess("update", "%s berhasil diperbarui", "%s successfully updated")
	SuccessDelete = newSuccess("delete", "%s berhasil dihapus", "%s successfully deleted")
	SuccessGet    = newSuccess("get", "%s berhasil diambil", "%s successfully retrieved")
)

// --- Error Definitions ---
// Use these as standard error variables that can be checked with errors.Is()
var (
	// General API Errors
	ErrInternalServer      = NewError("terjadi kesalahan pada server", "Internal server error occurred")
	ErrBadRequest          = NewError("permintaan tidak valid", "Invalid request")
	ErrGeneralInvalidInput = NewError("input tidak valid", "Invalid input")
	ErrUnauthorized        = NewError("tidak ada akses (unauthorized)", "Unauthorized access")
	ErrForbidden           = NewError("akses ditolak (forbidden)", "Access forbidden")
)

var (
	// Organization Errors
	ErrOrganizationNotFound   = NewError("organization tidak ditemukan", "Organization not found")
	ErrOrganizationNameExists = NewError("nama organization sudah digunakan", "Organization name already exists")

	// Standard CRUD Errors
	ErrGeneralCreated         = NewError("data gagal dibuat", "Failed to create data")
	ErrGeneralUpdated         = NewError("data gagal diperbarui", "Failed to update data")
	ErrGeneralDeleted         = NewError("data gagal dihapus", "Failed to delete data")
	ErrGeneralNotFound        = NewError("data tidak ditemukan", "Data not found")
	ErrGeneralNotFoundPartial = NewError("sebagian data tidak ditemukan", "Some data not found")
	ErrGeneralID              = NewError("id tidak valid", "Invalid ID")
)

var (
	// Datatable Errors
	ErrSortColumn = NewError("kolom pengurutan tidak valid", "Invalid sort column")
	ErrSortOrder  = NewError("urutan pengurutan tidak valid", "Invalid sort order")
)

var (
	// JWT / Auth Token Errors
	ErrTokenInvalid          = NewError("token tidak valid", "Token is invalid")
	ErrTokenExpired          = NewError("token sudah kadaluarsa", "Token has expired")
	ErrTokenUnexpectedMethod = NewError("metode penandatanganan token tidak didukung", "Unexpected token signing method")
)

var (
	// Crypto / Hashing Errors
	ErrHashInvalidFormat       = NewError("format hash tidak valid", "Encoded hash is not in the correct format")
	ErrHashIncompatibleVersion = NewError("versi argon2 tidak kompatibel", "Incompatible argon2 version")
)

var (
	// Password Validation Errors
	ErrPasswordTooWeak    = NewError("password terlalu lemah", "Password is too weak")
	ErrRegisterValidation = NewError("password lemah atau email tidak valid", "weak password or invalid email")
)

var (
	// Auth Errors
	ErrEmailAlreadyExists = NewError("email sudah digunakan", "Email already in use")
	ErrTooManyRequests    = NewError("terlalu banyak permintaan", "Too many requests. Please try again later.")
)
