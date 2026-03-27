package messages

// --- Success Messages ---

// Use these predefined success types for standard CRUD operations.
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
