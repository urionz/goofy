package web

import (
	"testing"
)

func TestNewRequest(t *testing.T) {
	// jsondata := `{
	// 	"name": "lisi",
	// 	"age": 12,
	// 	"isMale": false,
	// 	"sub": [{
	// 	"name": "sub1"
	// }]
	// }`
	// r := httptest.NewRequest("POST", "http://127.0.0.1?name=1&name=3", strings.NewReader("sex=2"))
	// r.Header.Set("CONTENT-TYPE", "application/x-www-form-urlencoded")
	// r.ParseForm()
	// fmt.Println(r.Form.Encode())
	// fmt.Println(r.PostForm.Encode())
	// a, _ := ioutil.ReadAll(r.Body)
	// fmt.Println(string(a))
	// r.Header.Set("CONTENT-TYPE", "application/x-www-form-urlencoded")
	// rr := NewRequest(r)
	// fmt.Println(rr.input().All())
}
