package main
import (
    "fmt"
	"code.google.com/p/gorest"
    "net/http"
	"github.com/fzzy/radix/redis"
	"os"
	"time"
	"strconv"
)
func main() {
	//Establish Redis Connection
	c, err := redis.DialTimeout("tcp", "127.0.0.1:6379", time.Duration(10)*time.Second)
	errHndlr(err)
	defer c.Close()

	// select database
	r := c.Cmd("select", 8)
	errHndlr(r.Err)

	r = c.Cmd("flushdb")
	errHndlr(r.Err)

	
	//Establish Listening Web Socket
    gorest.RegisterService(new(ItemService)) //Register our service
    http.Handle("/",gorest.Handle())    
    http.ListenAndServe(":8787",nil)
}

//Service Definition
type ItemService struct {
    gorest.RestService `root:"/items/" consumes:"application/json" produces:"application/json"`
    listItems  gorest.EndPoint `method:"GET" path:"/list/" output:"ItemStore"`
    genItems  gorest.EndPoint `method:"GET" path:"/gen/" output:"string"`
}
func(serv ItemService) ListItems() ItemStore{
    GetItems()
	serv.ResponseBuilder().SetResponseCode(200)
	return itemStore
}
func(serv ItemService) GenItems() string{
    PopulateItems()
	serv.ResponseBuilder().SetResponseCode(200)
	return "OK"
}

type Item struct{
	//Sync Changes in:
	//	readItems
	//	storeItems
    Id              int
    ProductName     string
    Active          bool
}
type ItemStore struct {
    Items []Item
}

var(
    itemStore ItemStore

)

func PopulateItems() {

	o := 1
	fn := "Testing1"
	at := false
	item1 := Item{Id: o, ProductName: fn, Active: at}
	storeItem("Item1", item1)

	o = 2
	fn = "Testing2"
	at = true
	item2 := Item{Id: o, ProductName: fn, Active: at}
	storeItem("Item2", item2)

	o = 3
	fn = "Testing3"
	at = false
	item3 := Item{Id: o, ProductName: fn, Active: at}
	storeItem("Item3", item3)

	itemStore = ItemStore{Items: []Item{item1, item2, item3}}

}

func GetItems() {

	var item1 Item
	readItem("Item1", &item1)

	var item2 Item
	readItem("Item2", &item2)

	var item3 Item
	readItem("Item3", &item3)

	itemStore = ItemStore{Items: []Item{item1, item2, item3}}
}

func readItem(tKey string, stct *Item) {
	//Establish Redis Connection
	c, err := redis.DialTimeout("tcp", "127.0.0.1:6379", time.Duration(10)*time.Second)
	errHndlr(err)
	defer c.Close()

	// select database
	r := c.Cmd("select", 8)
	errHndlr(r.Err)

	s, err := c.Cmd("hgetall", tKey).Hash()
	errHndlr(err)
	stct.Id, err = strconv.Atoi(s["Id"])
	errHndlr(err)
	stct.ProductName = s["ProductName"]
	stct.Active, err = strconv.ParseBool(s["Active"])
	errHndlr(err)
}

func storeItem(tKey string, stct Item) {
	//Establish Redis Connection
	c, err := redis.DialTimeout("tcp", "127.0.0.1:6379", time.Duration(10)*time.Second)
	errHndlr(err)
	defer c.Close()

	// select database
	r := c.Cmd("select", 8)
	errHndlr(r.Err)

//* Strings
	r = c.Cmd("set", "mykey0", "myval0")
	errHndlr(r.Err)

	s, err := c.Cmd("get", "mykey0").Str()
	errHndlr(err)
	fmt.Println("mykey0:", s)


	myhash := map[string]string{
		"Id": strconv.Itoa(stct.Id),
		"ProductName": stct.ProductName,
		"Active": strconv.FormatBool(stct.Active),
		}
	fmt.Println(myhash)
	r = c.Cmd("hmset", tKey, myhash)
	errHndlr(r.Err)
}

func errHndlr(err error) {
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
}