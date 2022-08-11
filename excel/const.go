package excel

const (
	Spreadsheet               = "http://schemas.openxmlformats.org/spreadsheetml/2006/main"
	Relationships             = "http://schemas.openxmlformats.org/package/2006/relationships"
	RelationshipsDoc          = "http://schemas.openxmlformats.org/officeDocument/2006/relationships"
	OfficeDocument            = "http://schemas.openxmlformats.org/officeDocument/2006/relationships/officeDocument"
	RelationshipSharedstrings = "http://schemas.openxmlformats.org/officeDocument/2006/relationships/sharedStrings"
	RelationshipStyles        = "http://schemas.openxmlformats.org/officeDocument/2006/relationships/styles"
	RelationshipsWorksheet    = "http://schemas.openxmlformats.org/officeDocument/2006/relationships/worksheet"
	RelationshipsTheme        = "http://schemas.openxmlformats.org/officeDocument/2006/relationships/theme"
	ContentTypes              = "http://schemas.openxmlformats.org/package/2006/content-types"
	CoreProperties            = "http://schemas.openxmlformats.org/package/2006/relationships/metadata/core-properties"
	ExtendedProperties        = "http://schemas.openxmlformats.org/officeDocument/2006/relationships/extended-properties"
	MarkupCompatibility       = "http://schemas.openxmlformats.org/markup-compatibility/2006"
	X14ac                     = "http://schemas.microsoft.com/office/spreadsheetml/2009/9/ac"
	X16r2                     = "http://schemas.microsoft.com/office/spreadsheetml/2015/02/main"
	X15                       = "http://schemas.microsoft.com/office/spreadsheetml/2010/11/main"

	AppVersion     = "16.0300"
	Creator        = "User"  // TODO Set to proper value at startup
	LastModifiedBy = Creator // TODO Set to proper value at startup

	azBase = 'Z' - 'A' + 1
)
