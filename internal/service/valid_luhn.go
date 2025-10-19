package service

import "strconv"

func IsValidLuhn(order string) bool {
	sum := 0
	alternate := false // Флаг для чередования

	for i := len(order) - 1; i >= 0; i-- {
		digit, err := strconv.Atoi(string(order[i]))
		if err != nil {
			// номер содержит нечисловые символы
			return false
		}

		if alternate {
			digit *= 2
			if digit > 9 {
				digit -= 9
			}
		}
		sum += digit
		alternate = !alternate
	}

	return sum%10 == 0
 }