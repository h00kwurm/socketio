package socketio

import (
	"testing"
)

func TestBuildUrl(t *testing.T) {
	const secureTest = "https://something.another.com:8088"
	const secureTestEndpoint = "12345"
	const secureExpected = "wss://something.another.com:8088/socket.io/1/websocket/12345"

	const unsecureTest = "http://127.0.0.1:8080"
	const unsecureTestEndpoint = "593"
	const unsecureExpected = "ws://127.0.0.1:8080/socket.io/1/websocket/593"

	const bustTest = "this.is.totally.busted"
	const bustTestEndpoint = ""

	if secure := buildUrl(secureTest, secureTestEndpoint); secure != secureExpected {
		t.Errorf("[test:buildurl] bad url. ex: %v, ac: %v\n", secureExpected, secure)
	}

	if unsecure := buildUrl(unsecureTest, unsecureTestEndpoint); unsecure != unsecureExpected {
		t.Errorf("[test:buildurl] bad url. ex: %v, ac: %v\n", unsecureExpected, unsecure)
	}

	if bust := buildUrl(bustTest, bustTestEndpoint); bust != bustTest {
		t.Errorf("[test:buildurl] bad url. ex: %v, ac: %v\n", bustTest, bust)
	}

}
