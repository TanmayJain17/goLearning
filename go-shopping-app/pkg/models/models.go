package models

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"go-fruit-cart/pkg/config"
	"go-fruit-cart/pkg/utils"
	"log"

	/* "strings" */
	"go-fruit-cart/pkg/apperrors"

	_ "github.com/lib/pq"
	uuid "github.com/satori/go.uuid"
)

var db *sql.DB

func init() {
	config.Connect()
	db = config.GetDb()
	/* fmt.Printf("%T\n", db) */
}

type User struct {
	Firstname string `json:"firstname" binding:"required,alpha"`
	Lastname  string `json:"lastname" binding:"required,alpha"`
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required,min=3"`
	/* Cartid    string `json:"cartid"` */
}
type Userdetails struct {
	Firstname string `json:"firstname" binding:"required,alpha"`
	Lastname  string `json:"lastname" binding:"required,alpha"`
	Email     string `json:"email" binding:"required,email"`
	Cartid    string `json:"cartid" binding:"required"`
	/* Cartid    string `json:"cartid"` */
}

type Product struct {
	Name        string  `json:"name" binding:"required"`
	Price       float64 `json:"price" binding:"required,numeric,gte=2"`
	Description string  `json:"description" binding:"required"`
}

type CartProducts struct {
	Productid string  `json:"productid" binding:"required"`
	Name      string  `json:"name" binding:"required"`
	Price     float64 `json:"price" binding:"required,numeric"`
	Quantity  int     `json:"quantity" binding:"required,numeric"`
}

type TokenRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

func InsertNewUser(newUser *User) (int, error) {
	sqlSt := fmt.Sprintf("SELECT * FROM users WHERE email='%v'", newUser.Email)
	val, err := db.Exec(sqlSt)
	if err != nil {
		return 502, apperrors.ErrUnExpectedError
	}
	value, _ := val.RowsAffected()
	if value == 0 {
		sqlStatement := `INSERT INTO users (firstname, lastname, email,password)
		VALUES ($1, $2, $3, $4)`
		_, err := db.Exec(sqlStatement, newUser.Firstname, newUser.Lastname, newUser.Email, newUser.Password)
		if err != nil {
			return 502, apperrors.ErrUnExpectedError
		}
		sqlState := fmt.Sprintf("SELECT cartid, isadmin FROM users WHERE email='%v'", newUser.Email)
		vall, nerr := db.Query(sqlState)
		if nerr != nil {
			return 502, apperrors.ErrUnExpectedError
		}

		for vall.Next() {
			var cartid string
			var isadmin bool
			err := vall.Scan(&cartid, &isadmin)
			if err != nil {
				panic(apperrors.ErrUnExpectedError)
			}
			/* fmt.Println("\n", cartid, isadmin) */
			theUuid, theErr := uuid.FromString(cartid)
			if theErr != nil {
				return 502, apperrors.ErrUnExpectedError
			}
			if !isadmin {
				insertIntoCartSqlSt := fmt.Sprintf("INSERT INTO carts (id) SELECT cartid FROM users WHERE cartid='%v'", theUuid)
				_, err := db.Exec(insertIntoCartSqlSt)
				if err != nil {
					return 502, apperrors.ErrUnExpectedError
				}
				/* fmt.Println("newVal->", newVal) */
			}

		}

		/* fmt.Println("val-->", val) */

		return 200, nil
	} else {
		return 401, apperrors.ErrUserAlreadyExists
	}
}

func FindTheUser(newUser *TokenRequest) (string, string, bool, int, error) {
	sqlStatement := fmt.Sprintf("SELECT email,firstname,lastname,password,isadmin FROM users WHERE email='%v'", newUser.Email)
	vall, err := db.Exec(sqlStatement)
	if err != nil {
		return "", "", false, 502, apperrors.ErrUnExpectedError
	}
	x, err := vall.RowsAffected()
	if err != nil {
		return "", "", false, 502, apperrors.ErrUnExpectedError
	}
	if x == 0 {
		return "", "", false, 404, apperrors.ErrUserDoesNotExist
	} else {
		val, err := db.Query(sqlStatement)
		if err != nil {

			return "", "", false, 502, apperrors.ErrUnExpectedError
		}

		defer val.Close()
		var email string
		var hashedPassword string
		var firstname string
		var lastname string
		var isadmin bool
		for val.Next() {

			err := val.Scan(&email, &firstname, &lastname, &hashedPassword, &isadmin)
			if err != nil {

				return "", "", false, 502, apperrors.ErrUnExpectedError
			}
			passwordCorrect := utils.CheckPasswordHash(newUser.Password, hashedPassword)

			if !passwordCorrect {
				return "", "", false, 401, apperrors.ErrUnauthorized
			}
			/* return hashedPassword, nil */
		}

		return firstname, lastname, isadmin, 200, nil
	}
}

func CheckIfAdmin(userEmail interface{}) (bool, error) {
	sqlStatement := fmt.Sprintf("SELECT isadmin FROM users WHERE email ='%v'", userEmail)
	val, err := db.Query(sqlStatement)
	if err != nil {

		return false, apperrors.ErrUnExpectedError
	}
	defer val.Close()
	var isadmin bool

	for val.Next() {
		err := val.Scan(&isadmin)
		if err != nil {
			return false, apperrors.ErrUnExpectedError
		}
	}
	if isadmin {
		return true, nil
	} else {
		return false, nil
	}
}

func InsertNewProduct(newProduct *Product) error {
	sqlStatement := `INSERT INTO products (name, price, description)
		VALUES ($1, $2, $3)`
	_, err := db.Exec(sqlStatement, newProduct.Name, newProduct.Price, newProduct.Description)
	if err != nil {
		return apperrors.ErrUnExpectedError
	}
	/* fmt.Println("val->", val) */
	return nil
}

func FindAllProducts() ([]Product, error) {
	sqlStatement := "select name,price,description from products"
	val, err := db.Query(sqlStatement)
	var allProducts []Product
	if err != nil {
		return allProducts, apperrors.ErrUnExpectedError
	}
	defer val.Close()
	for val.Next() {
		var product Product
		var name string
		var price float64
		var description string
		err := val.Scan(&name, &price, &description)
		if err != nil {
			return allProducts, apperrors.ErrUnExpectedError
		}
		product.Name = name
		product.Price = price
		product.Description = description
		allProducts = append(allProducts, product)
	}

	return allProducts, nil
}

func DeleteProductById(id string) (int, error) {
	sqlStatement := fmt.Sprintf("DELETE FROM products WHERE id='%v'", id)
	val, err := db.Exec(sqlStatement)
	if err != nil {
		return 500, apperrors.ErrUnExpectedError
	}
	count, err := val.RowsAffected()
	if err != nil {
		return 500, apperrors.ErrUnExpectedError
	}
	if count == 1 {
		return 200, nil
	} else {
		return 404, apperrors.ErrProductNotFound
	}
}

func UpdateProductById(product *Product, id string) (int, error) {
	sqlStatement := fmt.Sprintf("SELECT * FROM products WHERE id='%v'", id)
	val, err := db.Exec(sqlStatement)
	if err != nil {
		return 500, apperrors.ErrUnExpectedError
	}
	count, err := val.RowsAffected()
	if err != nil {
		return 500, apperrors.ErrUnExpectedError
	}
	if count == 1 {
		sqlUpdateStatement := fmt.Sprintf("UPDATE products SET name = '%v',price = %v,description = '%v' WHERE id = '%v'", product.Name, product.Price, product.Description, id)

		newVal, err := db.Exec(sqlUpdateStatement)
		if err != nil {
			return 500, apperrors.ErrUnExpectedError
		}
		counter, err := newVal.RowsAffected()
		if err != nil {
			return 500, apperrors.ErrUnExpectedError
		}
		if counter == 1 {
			return 200, nil
		} else {
			return 500, apperrors.ErrUnExpectedError
		}

	} else {
		return 404, apperrors.ErrProductNotFound
	}
}

func FindProduct(id string) (string, float64, error) {
	sqlStatement := fmt.Sprintf("SELECT name,price FROM products WHERE id='%v'", id)
	val, err := db.Query(sqlStatement)
	if err != nil {
		return "", 0, apperrors.ErrUnExpectedError
	}
	defer val.Close()
	var name string
	var price float64
	for val.Next() {
		err := val.Scan(&name, &price)
		if err != nil {
			return "", 0, apperrors.ErrUnExpectedError
		}

	}
	return name, price, nil
}

func GetProductsInCart(email interface{}) ([]byte, error) {
	sqlGetCartId := fmt.Sprintf("SELECT cartid FROM users WHERE email='%v'", email)
	var data []byte
	cartVal, err := db.Query(sqlGetCartId)
	if err != nil {
		return nil, apperrors.ErrUnExpectedError
	}
	defer cartVal.Close()
	var id string
	for cartVal.Next() {
		err := cartVal.Scan(&id)
		if err != nil {
			return data, apperrors.ErrUnExpectedError
		}

	}
	sqlStatement := fmt.Sprintf("SELECT products,totalitems FROM carts WHERE id='%v'", id)
	val, err := db.Query(sqlStatement)

	var count int
	if err != nil {
		return data, apperrors.ErrUnExpectedError
	}
	defer val.Close()

	for val.Next() {
		err := val.Scan(&data, &count)
		if err != nil {
			return data, apperrors.ErrUnExpectedError
		}
	}
	return data, nil
}

func AddProductToCart(email interface{}, name string, price float64, productid string) ([]byte, error) {
	sqlGetCartId := fmt.Sprintf("SELECT cartid FROM users WHERE email='%v'", email)
	var data []byte
	cartVal, err := db.Query(sqlGetCartId)
	if err != nil {
		return nil, apperrors.ErrUnExpectedError
	}
	defer cartVal.Close()
	var id string
	for cartVal.Next() {
		err := cartVal.Scan(&id)
		if err != nil {
			return data, apperrors.ErrUnExpectedError
		}

	}
	sqlStatement := fmt.Sprintf("SELECT products,totalitems FROM carts WHERE id='%v'", id)
	val, err := db.Query(sqlStatement)

	if err != nil {
		return data, apperrors.ErrUnExpectedError
	}
	defer val.Close()
	var productArray []CartProducts
	flag := 0
	var count int
	for val.Next() {
		err := val.Scan(&data, &count)
		if err != nil {
			return data, apperrors.ErrUnExpectedError
		}
		if count > 0 {
			err = json.Unmarshal(data, &productArray)
			if err != nil {
				return data, apperrors.ErrUnExpectedError
			}
			/* fmt.Println("productData->", productArray, len(productArray))
			fmt.Println("productData->", string(data)) */
		} else {
			theProduct := CartProducts{
				Productid: productid,
				Name:      name,
				Price:     price,
				Quantity:  1,
			}
			productArray = append(productArray, theProduct)
			count = 1
			flag = 2
		}

	}
	if flag != 2 {
		for i := 0; i < len(productArray); i++ {
			if productArray[i].Productid == productid {
				flag = 1
				productArray[i].Quantity += 1
				count += 1
				break
			}
		}
	}
	if flag != 1 && flag != 2 {

		newProduct := CartProducts{
			Productid: productid,
			Name:      name,
			Price:     price,
			Quantity:  1,
		}
		/* fmt.Println("inside other also but why?") */

		productArray = append(productArray, newProduct)
		/* fmt.Println("new Array", productArray) */
		count += 1
	}

	updatedCart, err := json.Marshal(productArray)
	if err != nil {
		return updatedCart, apperrors.ErrUnExpectedError
	}

	sqlUpdateSt := fmt.Sprintf("UPDATE  carts SET products='%v',totalitems=%v WHERE id='%v'", string(updatedCart), count, id)

	_, err = db.Exec(sqlUpdateSt)
	if err != nil {
		return updatedCart, apperrors.ErrUnExpectedError
	}
	/* fmt.Println(newval.RowsAffected()) */
	return updatedCart, nil
}

func RemoveProductFromCart(email interface{}, name string, price float64, productid string) ([]byte, error) {
	sqlGetCartId := fmt.Sprintf("SELECT cartid FROM users WHERE email='%v'", email)
	var data []byte
	cartVal, err := db.Query(sqlGetCartId)
	if err != nil {
		return nil, apperrors.ErrUnExpectedError
	}
	defer cartVal.Close()
	var id string
	for cartVal.Next() {
		err := cartVal.Scan(&id)
		if err != nil {
			return data, apperrors.ErrUnExpectedError
		}

	}
	sqlStatement := fmt.Sprintf("SELECT products,totalitems FROM carts WHERE id='%v'", id)
	val, err := db.Query(sqlStatement)

	if err != nil {
		return data, apperrors.ErrUnExpectedError
	}
	defer val.Close()
	var productArray []CartProducts
	flag := 0
	var count int
	for val.Next() {
		err := val.Scan(&data, &count)
		if err != nil {
			return data, apperrors.ErrUnExpectedError
		}
		err = json.Unmarshal(data, &productArray)
		if err != nil {
			return data, apperrors.ErrUnExpectedError
		}
		if count == 1 {
			flag = 1
		}
	}
	var arrIndex int
	for i := 0; i < len(productArray); i++ {
		if productArray[i].Productid == productid {
			if productArray[i].Quantity == 1 {
				if count == 1 {
					flag = 3
					break
				} else {
					flag = 2
					arrIndex = i
					/* productArray[i].Quantity = 0 */
					count -= 1
					break
				}

			} else {
				/* flag = 1 */
				/* arrIndex=i */
				productArray[i].Quantity -= 1
				count -= 1
				break
			}

		}
	}
	if flag == 1 || flag == 3 {
		var newArray []CartProducts
		updatedCart, err := json.Marshal(newArray)
		if err != nil {
			return updatedCart, apperrors.ErrUnExpectedError
		}
		sqlUpdateSt := fmt.Sprintf("UPDATE  carts SET products='%v',totalitems=0 WHERE id='%v'", string(updatedCart), id)

		_, err = db.Exec(sqlUpdateSt)
		if err != nil {
			return updatedCart, apperrors.ErrUnExpectedError
		}
		/* fmt.Println(newval.RowsAffected()) */
		return updatedCart, nil
	} else if flag == 2 {
		productArray = append(productArray[:arrIndex], productArray[arrIndex+1:]...)
	}
	updatedCart, err := json.Marshal(productArray)
	if err != nil {
		return updatedCart, apperrors.ErrUnExpectedError
	}
	sqlUpdateSt := fmt.Sprintf("UPDATE  carts SET products='%v',totalitems=%v WHERE id='%v'", string(updatedCart), count, id)

	_, err = db.Exec(sqlUpdateSt)
	if err != nil {
		return updatedCart, apperrors.ErrUnExpectedError
	}
	/* fmt.Println(newval.RowsAffected()) */
	return updatedCart, nil
}

func GrpcFindTheUser(userEmail string) (string, string, string, string, error) {
	sqlStatement := fmt.Sprintf("SELECT email,firstname,lastname,cartid FROM users WHERE email='%v'", userEmail)
	vall, err := db.Exec(sqlStatement)
	if err != nil {
		return "", "", "", "", apperrors.ErrUnExpectedError
	}
	x, err := vall.RowsAffected()
	if err != nil {
		return "", "", "", "", apperrors.ErrUnExpectedError
	}
	if x == 0 {
		return "", "", "", "", apperrors.ErrNoDataFound
	} else {
		val, err := db.Query(sqlStatement)
		if err != nil {

			return "", "", "", "", apperrors.ErrUnExpectedError
		}

		defer val.Close()
		var email string
		/* var hashedPassword string */
		var firstname string
		var lastname string
		var cartid string
		for val.Next() {

			err := val.Scan(&email, &firstname, &lastname, &cartid)
			if err != nil {

				return "", "", "", "", apperrors.ErrUnExpectedError
			}
			/* passwordCorrect := utils.CheckPasswordHash(newUser.Password, hashedPassword) */

			/* return hashedPassword, nil */
		}

		return userEmail, firstname, lastname, cartid, nil
	}
}

func GrpcFindProduct(id string) (string, string, float64, string, error) {
	sqlStatement := fmt.Sprintf("SELECT name,description,price,image FROM products WHERE id='%v'", id)
	val, err := db.Query(sqlStatement)
	if err != nil {
		return "", "", 0, "", apperrors.ErrUnExpectedError
	}
	defer val.Close()
	var name string
	var description string
	var image string
	var price float64

	for val.Next() {
		err := val.Scan(&name, &description, &price, &image)
		if err != nil {
			return "", "", 0, "", apperrors.ErrUnExpectedError
		}

	}
	return name, description, price, image, nil
}

func GrpcGetAllUsers() ([]Userdetails, error) {
	fmt.Println("inside this function")
	sqlStatement := "SELECT firstname,lastname,email,cartid FROM users"
	val, err := db.Query(sqlStatement)
	var usersArray []Userdetails
	if err != nil {
		return usersArray, apperrors.ErrUnExpectedError
	}
	defer val.Close()

	for val.Next() {
		fmt.Println("error here")
		var newuser = Userdetails{}
		err := val.Scan(&newuser.Firstname, &newuser.Lastname, &newuser.Email, &newuser.Cartid)
		if err != nil {
			return usersArray, apperrors.ErrUnExpectedError
		}
		log.Println(newuser)
		usersArray = append(usersArray, newuser)

	}
	return usersArray, nil
}

/* func GrpcGetCartAmount(email string) (int64, error) {
	sqlSt := fmt.Sprintf("(SELECT cartid FROM users WHERE email=%v)", email)
	sqlStatement := fmt.Sprintf("SELECT products FROM carts where id=%v", sqlSt)
	val, err := db.Query(sqlStatement)
	if err != nil {
		return 0, err
	}
	var prods []CartProducts
	for val.Next() {
		err := val.Scan(&prods)
		if err != nil {
			return 0, err
		}
	}
	fmt.Println(prods)
	return 0, nil
} */
