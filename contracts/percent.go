package contracts

// Percent represents a percentage value as a decimal value.
type Percent float64

// MakePercent creates a Percent from a percentage value (0 to 100).
//
// Parameters:
//   - percentage: The percentage value to convert.
//
// Returns:
//   - A Percent representing the decimal equivalent of the percentage.
func MakePercent(percentage float64) Percent {
	return Percent(percentage / 100.0)
}

// MakePercentFromDecimal creates a Percent from a decimal value (0.0 to 1.0).
//
// Parameters:
//   - decimal: The decimal value to convert.
//
// Returns:
//   - A Percent representing the decimal value.
func MakePercentFromDecimal(decimal float64) Percent {
	return Percent(decimal)
}

// AsDecimal returns the Percent as a decimal value (0.0 to 1.0).
//
// Returns:
//   - The decimal representation of the Percent.
func (this Percent) AsDecimal() float64 {
	return float64(this)
}

// AsPercentage returns the Percent as a percentage value (0 to 100).
//
// Returns:
//   - The percentage representation of the Percent.
func (this Percent) AsPercentage() float64 {
	return float64(this) * 100.0
}
