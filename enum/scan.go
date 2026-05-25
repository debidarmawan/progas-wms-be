package enum

import "fmt"

func scanString(dest *string, value interface{}) error {
	if value == nil {
		return nil
	}
	switch v := value.(type) {
	case string:
		*dest = v
	case []byte:
		*dest = string(v)
	default:
		return fmt.Errorf("cannot scan %T into string enum", value)
	}
	return nil
}
