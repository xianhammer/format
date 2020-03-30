package oxmsg

type Range struct {
	Min   uint16
	Max   uint16
	short string
	descr string
}

func (r *Range) Contains(value uint16) bool {
	return r.Min <= value && value <= r.Max
}

func (r *Range) Description() string {
	return r.descr
}

func (r *Range) String() string {
	return r.short
}

var (
	RangeMessageEnvelope              = Range{0x0001, 0x0BFF, "Nessage envelope", "Message object envelope property; reserved"}
	RangeRecipient                    = Range{0x0C00, 0x0DFF, "Recipient property", "Recipient property; reserved"}
	RangeMessageNontransmittable      = Range{0x0E00, 0x0FFF, "Message property, non-transmittable", "Non-transmittable Message property; reserved"}
	RangeMessageContent               = Range{0x1000, 0x2FFF, "Message content", "Message content property; reserved"}
	RangeMultipurpose                 = Range{0x3000, 0x33FF, "Multi-purpose", "Multi-purpose property that can appear on all or most objects; reserved"}
	RangeMessageStore                 = Range{0x3400, 0x35FF, "Message store", "Message store property; reserved"}
	RangeContainer                    = Range{0x3600, 0x36FF, "Folder and address book", "Folder and address book container property; reserved"}
	RangeAttachment                   = Range{0x3700, 0x38FF, "Attachment property", "Attachment property; reserved"}
	RangeAddressBook                  = Range{0x3900, 0x39FF, "Address Book", "Address Book object property; reserved"}
	RangeMailUser                     = Range{0x3A00, 0x3BFF, "Mail user property", "Mail user object property; reserved"}
	RangeDistributionList             = Range{0x3C00, 0x3CFF, "Distribution list property", "Distribution list property; reserved"}
	RangeProfile                      = Range{0x3D00, 0x3DFF, "Profile property", "Profile section property; reserved"}
	RangeStatus                       = Range{0x3E00, 0x3EFF, "Status property", "Status object property; reserved"}
	RangeTransportEnvelope            = Range{0x4000, 0x57FF, "Transport-defined envelope property", "Transport-defined envelope property"}
	RangeTransportRecipient           = Range{0x5800, 0x5FFF, "Transport-defined recipient property", "Transport-defined recipient property"}
	RangeUserNontransmittable         = Range{0x6000, 0x65FF, "User-defined property", "User-defined non-transmittable property"}
	RangeProviderNontransmittable     = Range{0x6600, 0x67FF, "Provider-defined property", "Provider-defined internal non-transmittable property"}
	RangeMessageClassContent          = Range{0x6800, 0x7BFF, "Message class-defined content property", "Message class-defined content property"}
	RangeMessageClassNontransmittable = Range{0x7C00, 0x7FFF, "Message class-defined content property, non-transmittable", "Message class-defined non-transmittable property"}
	RangeReserved                     = Range{0x8000, 0xFFFF, "Reserved", "Reserved for mapping to named properties."}
	RangeUnknown                      = Range{0, 0xFFFF, "Unknown", "Reserved for mapping to named properties."}

	rangeAll = []*Range{
		&RangeMessageEnvelope,
		&RangeRecipient,
		&RangeMessageNontransmittable,
		&RangeMessageContent,
		&RangeMultipurpose,
		&RangeMessageStore,
		&RangeContainer,
		&RangeAttachment,
		&RangeAddressBook,
		&RangeMailUser,
		&RangeDistributionList,
		&RangeProfile,
		&RangeStatus,
		&RangeTransportEnvelope,
		&RangeTransportRecipient,
		&RangeUserNontransmittable,
		&RangeProviderNontransmittable,
		&RangeMessageClassContent,
		&RangeMessageClassNontransmittable,
		&RangeReserved,
		&RangeUnknown,
	}
)

func FindRange(propID PropertyID) *Range {
	for _, r := range rangeAll {
		if r.Contains(uint16(propID)) {
			return r
		}
	}
	return &RangeUnknown
}
