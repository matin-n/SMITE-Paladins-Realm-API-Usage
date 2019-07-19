package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/tidwall/gjson"
	"io/ioutil"
	"net/http"
	"time"
)

var sessionId = ""
var devId = ""
var authKey = ""

func main() {

	//TODO: Implement Rate Limits
	//concurrent_sessions:  50
	//sessions_per_day: 500
	//session_time_limit:  15 minutes
	//request_day_limit:  7500
	//SMITE:      http://api.smitegame.com/smiteapi.svc

	createSession()

	fmt.Println(getPlayer("MatinMatin"))

}

// Each session only lasts for 15 minutes
func createSession() {
	signature := createSignature(devId, "createsession", authKey, GetCurrentTime())

	resp, err := http.Get("http://api.smitegame.com/smiteapi.svc/createsessionJson/" + devId + "/" + signature + "/" + GetCurrentTime())
	if err != nil {
		// handle error
	}

	defer resp.Body.Close()

	bodyBytes, _ := ioutil.ReadAll(resp.Body)
	bodyString := string(bodyBytes)

	sessionId = gjson.Get(bodyString, "session_id").String()
}

func getPlayer(username string) string {
	signature := createSignature(devId, "getplayer", authKey, GetCurrentTime())

	resp, err := http.Get("http://api.smitegame.com/smiteapi.svc/getplayerjson/" + devId + "/" + signature + "/" + sessionId + "/" + GetCurrentTime() + "/" + username)

	if err != nil {
		// handle error
	}

	defer resp.Body.Close()

	bodyBytes, _ := ioutil.ReadAll(resp.Body)
	bodyString := string(bodyBytes)

	return bodyString
}

func createSignature(devId, functionName, authKey, timestamp string) string {
	return GetMD5Hash(devId + functionName + authKey + timestamp)
}

func GetMD5Hash(text string) string {
	//private static string GetMD5Hash(string input) {
	//	var md5 = new System.Security.Cryptography.MD5CryptoServiceProvider();
	//	var bytes = System.Text.Encoding.UTF8.GetBytes(input);
	//	bytes = md5.ComputeHash(bytes);
	//	var sb = new System.Text.StringBuilder();
	//	foreach (byte b in bytes) {
	//		sb.Append(b.ToString("x2").ToLower());
	//	}
	//	return sb.ToString();
	//}

	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}

func GetCurrentTime() string {
	// string timestamp = DateTime.UtcNow.ToString("yyyyMMddHHmmss");

	return time.Now().UTC().Format("20060102150405") // https://golang.org/src/time/format.go
}
