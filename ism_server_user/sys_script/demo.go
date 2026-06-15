

//go:build ignore
// +build ignore


var http = import("net/http")
var json = import("encoding/json")
var fmt = import("fmt")

http.HandleFunc("/DiyChartData", func(w, r) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "content-type")
	w.Header().Set("Access-Control-Expose-Headers", "Access-Control-Allow-Headers, Token")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	

	responseStu = make(struct {
		A int64,
		B float64
	})
	responseStu.A = GetDeviceData("ww->register0")
	var responseData,_ = json.Marshal(responseStu)
	w.Write(responseData)
})

http.HandleFunc("/DiyChartDataFromDb", func(w, r) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "content-type")
	w.Header().Set("Access-Control-Expose-Headers", "Access-Control-Allow-Headers, Token")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	y = GetRemoteDbData("mysql","ismctl","h1ofTL636j","sql.s1267.vhostgo.com",3306,"ismctl","SELECT * FROM devices_alarm_list")
	var responseData,_ = json.Marshal(y)
	w.Write(responseData)
})

http.ListenAndServe(":8099", nil)



