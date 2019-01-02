// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/line/line-bot-sdk-go/linebot"
)

var bot *linebot.Client
var serverURL string

func main() {
	var err error
	serverURL = os.Getenv("LINECORP_PLATFORM_CHANNEL_SERVERURL")
	bot, err = linebot.New(os.Getenv("LINECORP_PLATFORM_CHANNEL_CHANNELSECRET"), os.Getenv("LINECORP_PLATFORM_CHANNEL_CHANNELTOKEN"))
	log.Println("Bot:", bot, " err:", err)

	//BOT APIs
	http.HandleFunc("/callback", callbackHandler)

	//Web APIs
	http.HandleFunc("/", listCust)
	http.HandleFunc("/link", link)
	http.HandleFunc("/login", login)

	//provide by Heroku
	port := os.Getenv("PORT")
	addr := fmt.Sprintf(":%s", port)
	http.ListenAndServe(addr, nil)
}
