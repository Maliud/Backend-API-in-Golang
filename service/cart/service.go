package cart

import (
	"fmt"

	"github.com/Maliud/Backend-API-in-Golang/types"
)

func getCartItemsIDs(items []types.CartItem) ([]int, error) {
	
	productIds := make([]int, len(items))
	for i, item := range items {
		if item.Quantity <= 0 {
			return nil, fmt.Errorf("Bu Ürün Geçersiz adet isteği %d", item.ProductID)
		}
		productIds[i] = item.ProductID
	}

	return productIds, nil
}

func (h *Handler) createOrder(ps []types.Product, items []types.CartItem, userID int) (int, float64, error) {

	productMap := make(map[int]types.Product)
	for _, product := range ps {
		productMap[product.ID] = product
	}

	// tüm ürünlerin gerçekten stokta olup olmadığını kontrol et
	if err := checkIfCartIsInStock(items, productMap); err != nil {
		return 0, 0, nil
	}

	// toplam fiyatı hesaplama
	totalPrice := calculateTotalPrice(items, productMap)

	// db'mizdeki ürün miktarını azaltılır
	for _, item := range items {
		product := productMap[item.ProductID]
		product.Quantity -= item.Quantity

		h.productStore.UpdateProduct(product)
	}

	// siparişi oluşturun
	orderID, err := h.store.CreateOrder(types.Order{
		UserID: userID,
		Total: totalPrice,
		Status: "Bekleniyor",
		Address: "Adres",
	})
	if err != nil {
		return 0, 0, err
	}

	// sipariş adetlerini olusturma
	for _, item := range items {
		h.store.CreateOrderItem(types.OrderItem{
			OrderID: orderID,
			ProductID: item.ProductID,
			Quantity: item.Quantity,
			Price: productMap[item.ProductID].Price,
		})
	}

	return orderID, totalPrice, nil
}

func checkIfCartIsInStock(cartItems []types.CartItem, products map[int]types.Product) error {
	if len(cartItems) == 0 {
		return fmt.Errorf("Kart Boş!!")
	}

	for _, item := range cartItems {
		product, ok := products[item.ProductID]
		if !ok {
			return fmt.Errorf("Ürün mağazada mevcut değil, lütfen sepetinizi yenileyin", item.ProductID)
		}

		if product.Quantity < item.Quantity {
			return fmt.Errorf("Ürün talep edilen miktarda mevcut değil", product.Name)

		}
	}
	return nil
}

func calculateTotalPrice(cartItems []types.CartItem, products map[int]types.Product) float64 {
	var total float64

	for _, item := range cartItems {
		product := products[item.ProductID]
		total += product.Price * float64(item.Quantity)
	}

	return total
}