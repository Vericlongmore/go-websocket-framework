package even

import (
	"application/controllers"
	"application/core/websocket_gateway"
	"testing"
)

func TestEven(t *testing.T) {

	client := &websocket_gateway.Client{}
	message := `{"uid":579501,"platform":"iOS","code":277,"ssid":"","version":"1.30","param":{"ack":"c121c230-eba5-11e7-b63c-00163e0ffda7","idx":"3","type":"人员"},"idfa":"2f3c365736ae9efce5ab2974806a8c3999daded3"} `
	controllers.HandShake(client, []byte(message))

	t.Log(" 7 is not even!")
	t.Fail()

}
