package imapclient_test

import (
	"bufio"
	"net"
	"testing"

	"github.com/emersion/go-imap/v2"
	"github.com/emersion/go-imap/v2/imapclient"
)

func TestID_extra(t *testing.T) {
	clientConn, serverConn := net.Pipe()
	defer serverConn.Close()

	sent := make(chan string, 1)
	go func() {
		serverConn.Write([]byte("* OK [CAPABILITY IMAP4rev1 ID] ready\r\n"))
		line, _ := bufio.NewReader(serverConn).ReadString('\n')
		sent <- line
		serverConn.Write([]byte(`* ID ("name" "srv" "host" "mx1")` + "\r\n"))
		serverConn.Write([]byte("T1 OK ID completed\r\n"))
	}()

	client := imapclient.New(clientConn, nil)
	defer client.Close()

	data, err := client.ID(&imap.IDData{
		Name: "go-imap",
		Raw: map[string]string{
			"x-originating-ip": "192.0.2.1",
			"x-client":         "a",
		},
	}).Wait()
	if err != nil {
		t.Fatalf("ID() = %v", err)
	}

	want := `T1 ID ("name" "go-imap" "x-client" "a" "x-originating-ip" "192.0.2.1")` + "\r\n"
	if got := <-sent; got != want {
		t.Errorf("sent %q, want %q", got, want)
	}
	if data.Name != "srv" || data.Raw["host"] != "mx1" {
		t.Errorf("ID() = %+v, want name srv and host mx1", data)
	}
}
