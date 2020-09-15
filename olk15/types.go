package olk15

// String = 4 byte (le) length + 2*N bytes UTF-8
type String struct {
	Length uint32
	Data   []byte
}

// type Contact struct {
// 	Marker     uint16   // ? 0200
// 	RecordSize uint16   // (le)
// 	Unknown1   [24]byte //
// 	Email      String
// 	Name       String
// }

// Header struct:
// From Contact
// ReplyTo Contact (???)
// Recv01 Contact
// ...
