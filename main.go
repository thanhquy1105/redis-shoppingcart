package main

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/redis/go-redis/v9"
)

var rdb *redis.Client
var ctx = context.Background()

func main() {
	rdb = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	seeder()

	unPaid()
	cartsThatTotalQtyIsGreaterThan5()
	numberOfCartsThatProductXAppears("0002")
	calculateAmoutOfACart("1000")
	productsOfACart("1001")
	removeProductInACart("1000", "0005")
	increaseQuantityOfProductInACartBy1("1000", "0001")
	removeAllProductsInACart("1111")
	cartWithGreatestAmount()
	productAppearsInMostCarts()
}

func seeder() {
	rdb.HSet(ctx, "product:0001", "id", "0001", "name", "iphone14", "image", "https://www.apple.com/newsroom/images/product/iphone/standard/Apple-iPhone-14-iPhone-14-Plus-hero-220907_Full-Bleed-Image.jpg.large.jpg", "price", 800)
	rdb.HSet(ctx, "product:0002", "id", "0002", "name", "ao", "image", "https://wetrek.vn/pic/products/Ao-Co-Do-Sao-Vang-Viet-Nam_HasThumb.jpg", "price", 30)
	rdb.HSet(ctx, "product:0003", "id", "0003", "name", "aoquan", "image", "https://nemsport.com/wp-content/uploads/2021/04/quan-ao-da-bong-doi-tuyen-viet-nam-do.jpg", "price", 300)
	rdb.HSet(ctx, "product:0004", "id", "0004", "name", "giay", "image", "https://www.sport9.vn/images/uploaded/kamito/kamito%20qh19/8b3d7b484f59b607ef48.jpg", "price", 500)

	rdb.HSet(ctx, "cart:1000", "product:0001", "2", "product:0002", "3", "product:0004", "1", "paid", "true")
	rdb.HSet(ctx, "cart:1001", "product:0001", "1", "product:0002", "1", "product:0003", "2", "paid", "false")
}

// Done
func unPaid() {
	fmt.Println("Carts are unpaid")
	carts, _ := rdb.Keys(ctx, "cart:*").Result()
	result := []string{}
	for _, cart := range carts {
		v, _ := rdb.HGet(ctx, cart, "paid").Result()

		if v == "false" {
			result = append(result, cart)
		}
	}
	fmt.Println(result)
}

// Done
func cartsThatTotalQtyIsGreaterThan5() {
	fmt.Println("Carts that total quantity is greater than 5")
	carts, _ := rdb.Keys(ctx, "cart:*").Result()
	result := []string{}
	for _, cart := range carts {
		sum := 0
		products, _ := rdb.HGetAll(ctx, cart).Result()
		for productID, qty := range products {
			if strings.Contains(productID, "product:") {
				qtyInt, _ := strconv.Atoi(qty)
				sum += qtyInt
			}
		}
		if sum > 5 {
			result = append(result, cart)
		}
	}
	fmt.Println(result)
}

// Done
func numberOfCartsThatProductXAppears(productID string) {
	fmt.Println("Number of carts that product:ID appears")
	carts, _ := rdb.Keys(ctx, "cart:*").Result()
	result := []string{}
	for _, cart := range carts {
		values, _ := rdb.HGetAll(ctx, cart).Result()
		for key, _ := range values {
			if key == "product:"+productID {
				result = append(result, cart)
			}
		}
	}
	fmt.Println(result)
	fmt.Println(len(result))
}

// Done
func calculateAmoutOfACart(cartID string) float64 {
	fmt.Println("Calculate amout of a cartID")
	cart, _ := rdb.HGetAll(ctx, "cart:"+cartID).Result()
	var total float64 = 0.0
	for productID, qty := range cart {
		price, _ := rdb.HGet(ctx, productID, "price").Result()

		priceFloat, _ := strconv.ParseFloat(price, 64)
		qtyInt, _ := strconv.Atoi(qty)

		total += float64(qtyInt) * priceFloat
	}
	fmt.Println(total)
	return total
}

func productsOfACart(cartID string) {
	fmt.Println("list products of a cart")
	cart, _ := rdb.HGetAll(ctx, "cart:"+cartID).Result()
	products := make(map[string]map[string]string)
	for productID, _ := range cart {
		if strings.Contains(productID, "product:") {
			product, _ := rdb.HGetAll(ctx, productID).Result()
			products[productID] = make(map[string]string)
			products[productID] = product
		}
	}
	fmt.Println(cart)
	fmt.Println(products)
}

func removeProductInACart(cartID string, productID string) {
	fmt.Println("remove productID of a cartID")
	num, _ := rdb.HDel(ctx, "cart:"+cartID, "product:"+productID).Result()
	fmt.Println(num)
}

func increaseQuantityOfProductInACartBy1(cartID string, productID string) {
	fmt.Println("increase the quantity of productID in cartID by 1 unit")
	num, _ := rdb.HIncrBy(ctx, "cart:"+cartID, "product:"+productID, 1).Result()
	fmt.Println(num)
}

func removeAllProductsInACart(cartID string) {
	fmt.Println("remove all products of a cartID")
	num, _ := rdb.Del(ctx, "cart:"+cartID).Result()
	fmt.Println(num)
}

func cartWithGreatestAmount() {
	fmt.Println("cart with greatest amount")
	carts, _ := rdb.Keys(ctx, "cart:*").Result()
	result := ""
	amout := 0.0
	for _, cart := range carts {
		temp := calculateAmoutOfACart(cart[5:])
		if temp > amout {
			amout = temp
			result = cart
		}
	}
	fmt.Println(result)
}

func productAppearsInMostCarts() {
	fmt.Println("product appears in most shopping carts")
	products, _ := rdb.Keys(ctx, "product:*").Result()
	carts, _ := rdb.Keys(ctx, "cart:*").Result()
	result := []string{}
	max := 0
	for _, product := range products {
		num := 0
		for _, cart := range carts {
			exist, _ := rdb.HExists(ctx, cart, product).Result()
			if exist {
				num++
			}
		}
		if num > max {
			result = nil
			result = append(result, product)
			max = num
		} else if num == max {
			result = append(result, product)
		}
	}
	fmt.Println(result)
}
