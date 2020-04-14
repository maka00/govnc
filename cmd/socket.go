/*
Copyright © 2020 Manfred Kalan <manfred.kalan@gmail.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"bytes"
	"encoding/base64"
	"fmt"
	mux "github.com/gorilla/mux"
	websocket "github.com/gorilla/websocket"
	"github.com/spf13/cobra"
	"image"
	"image/jpeg"
	"log"
	"net/http"
	"time"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  2048,
	WriteBufferSize: 2048,
}

// The go routine with a channel to get the screen
func getImage(c chan *image.RGBA) {
	const timeBetweenFrames = time.Duration(1*time.Second) / 30.0
	avgTimeBetweenFrames := timeBetweenFrames
	for {
		start := time.Now()
		img, err := takeShot()
		if err != nil {
			log.Println("error on snapshot")
			continue
		}
		c <- img
		dTime := time.Since(start)
		waitTime := dTime - avgTimeBetweenFrames
		// slow the server down if it's going too fast
		if waitTime < 0 {
			time.Sleep(-1 * waitTime)
		} else {
			log.Println("fast enough")
		}
	}
}

// The handler for the websocket
func server(w http.ResponseWriter, r *http.Request) {
	log.Println("received a request")
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	// create a channel and start a go routine to fill it
	c := make(chan *image.RGBA)
	defer close(c)
	go getImage(c)

	// read the incomming message (we don't care what the content is)
	messageType, p, err := conn.ReadMessage()
	if err != nil {
		log.Println(err)
		return
	}
	log.Printf("got message %s\n", string(p))

	// constantly send updates to the client from the channel c
	for {
		img := <-c
		var msg string
		buf := new(bytes.Buffer)
		jpeg.Encode(buf, img, nil)
		msg = fmt.Sprintf("data:image/jpeg;base64,  %s", base64.StdEncoding.EncodeToString(buf.Bytes()))
		if err := conn.WriteMessage(messageType, []byte(msg)); err != nil {
			log.Println(err)
			return
		}
	}
}

// socketCmd represents the socket command
var socketCmd = &cobra.Command{
	Use:   "socket",
	Short: "Shares the screen via websockets",
	Long: `Shares the screen via websockets.
Uses a go channel to gather the image and transforms it in parallel to a base64 encoded JPEG`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("socket called")
		r := mux.NewRouter()
		r.PathPrefix("/socket").HandlerFunc(server)
		r.PathPrefix("/").Handler(http.FileServer(http.Dir("./static/")))

		log.Fatal(http.ListenAndServe(":8888", r))
	},
}

func init() {
	rootCmd.AddCommand(socketCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// socketCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// socketCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
