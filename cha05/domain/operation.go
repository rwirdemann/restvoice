package domain

type Operation string

func GetOperations(invoice Invoice) []Operation {
	switch invoice.Status {
	case "open":
		return []Operation{"book", "charge", "bookings"}
	case "payment expected":
		return []Operation{"payment", "bookings"}
	case "payed":
		return []Operation{"archive"}
	case "archived":
		return []Operation{"revoke"}
	case "revoked":
		return []Operation{"archive"}
	default:
		return []Operation{}
	}
}
