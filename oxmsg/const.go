package oxmsg

import (
	"errors"

	"github.com/xianhammer/format/cbf"
)

type PropertyID uint32
type PropertyType uint16

const (
	PropertyPrefix = "__substg1.0_"

	PtypBinary               PropertyType = 0x0102 // COUNT, 16-bit
	PtypBoolean                           = 0x000B
	PtypCurrency                          = 0x0006
	PtypErrorCode                         = 0x000A
	PtypFloating32                        = 0x0004
	PtypFloating64                        = 0x0005
	PtypFloatingTime                      = 0x0007
	PtypGuid                              = 0x0048
	PtypInteger16                         = 0x0002
	PtypInteger32                         = 0x0003
	PtypInteger64                         = 0x0014
	PtypMultipleBinary                    = 0x1102 // COUNT, 16-bit
	PtypMultipleCurrency                  = 0x1006 // COUNT, 16-bit
	PtypMultipleFloating32                = 0x1004 // COUNT, 16-bit
	PtypMultipleFloating64                = 0x1005 // COUNT, 16-bit
	PtypMultipleFloatingTime              = 0x1007 // COUNT, 16-bit
	PtypMultipleGuid                      = 0x1048 // COUNT, 16-bit
	PtypMultipleInteger16                 = 0x1002 // COUNT, 16-bit
	PtypMultipleInteger32                 = 0x1003 // COUNT, 16-bit
	PtypMultipleInteger64                 = 0x1014 // COUNT, 16-bit
	PtypMultipleString                    = 0x101F // COUNT, 16-bit. UTF16-LE
	PtypMultipleString8                   = 0x101E // COUNT, 16-bit
	PtypMultipleTime                      = 0x1040 // COUNT, 16-bit
	PtypNull                              = 0x0001
	PtypObject                            = 0x000D // A.k.a. PtypEmbeddedTable
	PtypRestriction                       = 0x00FD
	PtypRuleAction                        = 0x00FE
	PtypServerId                          = 0x00FB
	PtypString                            = 0x001F // UTF16-LE
	PtypString8                           = 0x001E
	PtypTime                              = 0x0040
	PtypUnspecified                       = 0x0000

	// The "types" below are NOT from the specification, but added to encompas the "property set"
	// These constants are used instead of ID and does not represent unique ID's but "set (membership) ID".
	PsetAddress              PropertyID = 0xFFFFFFFF // PSETID_Address
	PsetAirSync                         = 0xFFFFFFFE // PSETID_AirSync
	PsetAppointment                     = 0xFFFFFFFD // PSETID_Appointment
	PsetAttachment                      = 0xFFFFFFEC // PSETID_Attachment
	PsetCommon                          = 0xFFFFFFFB // PSETID_Common
	PsetInternetHeaders                 = 0xFFFFFFFA // PS_INTERNET_HEADERS
	PsetLog                             = 0xFFFFFFF9 // PSETID_Log
	PsetMAPI                            = 0xFFFFFFF8 // PS_MAPI
	PsetMeeting                         = 0xFFFFFFF7 // PSETID_Meeting
	PsetMessaging                       = 0xFFFFFFF6 // PSETID_Messaging
	PsetNote                            = 0xFFFFFFF5 // PSETID_Note
	PsetPostRss                         = 0xFFFFFFF4 // PSETID_PostRss
	PsetPublicStrings                   = 0xFFFFFFF3 // PS_PUBLIC_STRINGS
	PsetSharing                         = 0xFFFFFFF2 // PSETID_Sharing
	PsetTask                            = 0xFFFFFFF1 // PSETID_Task
	PsetUnifiedMessaging                = 0xFFFFFFF0 // PSETID_UnifiedMessaging
	PsetXmlExtractedEntities            = 0xFFFFFFEF // PSETID_XmlExtractedEntities
	PsetLAST                            = 0xFFFFFF00 // -- Marker for Pset constants
)

var (
	PS_PUBLIC_STRINGS           = cbf.GUID{0x00020329, 0x0000, 0x0000, [8]byte{0xC0, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x46}}
	PSETID_Common               = cbf.GUID{0x00062008, 0x0000, 0x0000, [8]byte{0xC0, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x46}}
	PSETID_Address              = cbf.GUID{0x00062004, 0x0000, 0x0000, [8]byte{0xC0, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x46}}
	PS_INTERNET_HEADERS         = cbf.GUID{0x00020386, 0x0000, 0x0000, [8]byte{0xC0, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x46}}
	PSETID_Appointment          = cbf.GUID{0x00062002, 0x0000, 0x0000, [8]byte{0xC0, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x46}}
	PSETID_Meeting              = cbf.GUID{0x6ED8DA90, 0x450B, 0x101B, [8]byte{0x98, 0xDA, 0x00, 0xAA, 0x00, 0x3F, 0x13, 0x05}}
	PSETID_Log                  = cbf.GUID{0x0006200A, 0x0000, 0x0000, [8]byte{0xC0, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x46}}
	PSETID_Messaging            = cbf.GUID{0x41F28F13, 0x83F4, 0x4114, [8]byte{0xA5, 0x84, 0xEE, 0xDB, 0x5A, 0x6B, 0x0B, 0xFF}}
	PSETID_Note                 = cbf.GUID{0x0006200E, 0x0000, 0x0000, [8]byte{0xC0, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x46}}
	PSETID_PostRss              = cbf.GUID{0x00062041, 0x0000, 0x0000, [8]byte{0xC0, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x46}}
	PSETID_Task                 = cbf.GUID{0x00062003, 0x0000, 0x0000, [8]byte{0xC0, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x46}}
	PSETID_UnifiedMessaging     = cbf.GUID{0x4442858E, 0xA9E3, 0x4E80, [8]byte{0xB9, 0x00, 0x31, 0x7A, 0x21, 0x0C, 0xC1, 0x5B}}
	PS_MAPI                     = cbf.GUID{0x00020328, 0x0000, 0x0000, [8]byte{0xC0, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x46}}
	PSETID_AirSync              = cbf.GUID{0x71035549, 0x0739, 0x4DCB, [8]byte{0x91, 0x63, 0x00, 0xF0, 0x58, 0x0D, 0xBB, 0xDF}}
	PSETID_Sharing              = cbf.GUID{0x00062040, 0x0000, 0x0000, [8]byte{0xC0, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x46}}
	PSETID_XmlExtractedEntities = cbf.GUID{0x23239608, 0x685D, 0x4732, [8]byte{0x9C, 0x55, 0x4C, 0x95, 0xCB, 0x4E, 0x8E, 0x33}}
	PSETID_Attachment           = cbf.GUID{0x96357F7F, 0x59E1, 0x47D0, [8]byte{0x99, 0xA7, 0x46, 0x51, 0x5C, 0x18, 0x3B, 0x54}}
)

var (
	ErrPropertyParse            = errors.New("Cannot parse as property")
	ErrPropertyID               = errors.New("Not a property ID")
	ErrPropertyNotFound         = errors.New("Property not found")
	ErrPropertyIllegalInstances = errors.New("Property was expected to be defined only once")
)
