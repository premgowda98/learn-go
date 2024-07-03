package maps

import "fmt"

func Maps() {
	// website := []string{"Some", "Nonen", "pen"}
	website := map[string]string{
		"key":  "value",
		"key1": "value2",
	} //like a key value pair, type of key inside [] and outside is for value type

	fmt.Println(website)
	fmt.Println(website["key1"])

}

func LoopMap() {
	website := map[string]string{
		"key":  "value",
		"key1": "value2",
	}

	for key, val := range website {
		fmt.Println(key, val)
	}
}
