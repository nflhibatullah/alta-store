package helper

func CheckProductStatus(stock int) string {
	
	if stock == 0 {
		return "not available"
	}

	return "available"
}

func CheckAvailableQuantity(quantity, stock int) int {
	if quantity >= stock {
		return stock
	}

	return quantity
}