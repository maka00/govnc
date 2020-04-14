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
	"fmt"
	"github.com/spf13/cobra"
	"image/jpeg"
	"log"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"time"
)

// The main handler to deliver the image (can be called directly
func myHandler(w http.ResponseWriter, req *http.Request) {
	fmt.Print("Hello, got a request\n")
	mimeWriter := multipart.NewWriter(w)
	mimeWriter.SetBoundary("--boundary")
	contentType := fmt.Sprintf("multipart/x-mixed-replace;boundary=%s", mimeWriter.Boundary())
	w.Header().Add("Cache-Control", "no-store, no-cache, must-revalidate, pre-check=0, post-check=0, max-age=0")
	w.Header().Add("Content-Type", contentType)
	w.Header().Add("Pragma", "no-cache")
	w.Header().Add("Connection", "close")
	s := time.Now()
	for {
		partHeader := make(textproto.MIMEHeader)
		partHeader.Add("Content-Type", "image/jpeg")
		partHeader.Add("X-StartTime", fmt.Sprintf("%v", s.Unix()))
		partHeader.Add("X-Timestamp", fmt.Sprintf("%v", s.Unix()))
		partWriter, _ := mimeWriter.CreatePart(partHeader)
		snapshot, _ := takeShot()
		buf := new(bytes.Buffer)
		jpeg.Encode(buf, snapshot, nil)
		//storeImage(snapshot, "test.png")
		partWriter.Write(buf.Bytes())
	}

}

// small start function for the server (can be refactored in the main function)
func startServer() {
	s := &http.Server{
		Addr:              ":8080",
		Handler:           http.HandlerFunc(myHandler),
		TLSConfig:         nil,
		ReadTimeout:       10 * time.Hour, //.Second,
		ReadHeaderTimeout: 0,
		WriteTimeout:      10 * time.Hour, //Second,
		IdleTimeout:       0,
		MaxHeaderBytes:    1 << 20,
		TLSNextProto:      nil,
		ConnState:         nil,
		ErrorLog:          nil,
		BaseContext:       nil,
		ConnContext:       nil,
	}
	log.Fatal(s.ListenAndServe())
}

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "serves the screen on a http stream",
	Long:  `serves the screen via a multipart mjpeg.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("serve called")
		startServer()
		fmt.Print("finished")
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// serveCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// serveCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
