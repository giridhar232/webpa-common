package fact

import (
	"encoding/base64"
	"github.com/Comcast/webpa-common/canonical"
	"github.com/Comcast/webpa-common/convey"
	"github.com/Comcast/webpa-common/logging"
	"github.com/Comcast/webpa-common/secure"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
	"net/http"
	"os"
	"testing"
)

const (
	conveyPayload string = "eyAicGFyYW1ldGVycyI6IFsgeyAibmFtZSI6ICJEZXZpY2UuRGV2aWNlSW5mby5XZWJwYS5YX0NPTUNBU1QtQ09NX0NJRCIsICJ2YWx1ZSI6ICIwIiwgImRhdGFUeXBlIjogMCB9LCB7ICJuYW1lIjogIkRldmljZS5EZXZpY2VJbmZvLldlYnBhLlhfQ09NQ0FTVC1DT01fQ01DIiwgInZhbHVlIjogIjI2OSIsICJkYXRhVHlwZSI6IDIgfSBdIH0K"
	basicAuth     string = "Basic dXNlcjpwYXNzd29yZA=="
)

func TestLogger(t *testing.T) {
	assert := assert.New(t)
	ctx := context.Background()

	value, ok := Logger(ctx)
	if !assert.Nil(value) || !assert.False(ok) {
		return
	}

	func() {
		defer func() {
			assert.Equal(NoLogger, recover())
		}()

		MustLogger(ctx)
	}()

	logger := &logging.LoggerWriter{os.Stdout}
	ctx = SetLogger(ctx, logger)

	value, ok = Logger(ctx)
	assert.Equal(logger, value)
	assert.True(ok)

	func() {
		defer func() {
			assert.Nil(recover())
		}()

		assert.Equal(logger, MustLogger(ctx))
	}()
}

func TestDeviceId(t *testing.T) {
	assert := assert.New(t)
	ctx := context.Background()

	value, ok := DeviceId(ctx)
	if !assert.Nil(value) || !assert.False(ok) {
		return
	}

	func() {
		defer func() {
			assert.Equal(NoDeviceId, recover())
		}()

		MustDeviceId(ctx)
	}()

	deviceId, err := canonical.ParseId("mac:111122223333")
	if !assert.NotNil(deviceId) || !assert.Nil(err) {
		return
	}

	t.Logf("Parsed device id: %v", deviceId)
	ctx = SetDeviceId(ctx, deviceId)

	value, ok = DeviceId(ctx)
	assert.Equal(deviceId, value)
	assert.True(ok)

	func() {
		defer func() {
			assert.Nil(recover())
		}()

		assert.Equal(deviceId, MustDeviceId(ctx))
	}()
}

func TestConvey(t *testing.T) {
	assert := assert.New(t)
	ctx := context.Background()

	value, ok := Convey(ctx)
	if !assert.Nil(value) || !assert.False(ok) {
		return
	}

	func() {
		defer func() {
			assert.Equal(NoConvey, recover())
		}()

		MustConvey(ctx)
	}()

	payload, err := convey.ParsePayload(base64.StdEncoding, conveyPayload)
	if !assert.NotNil(payload) || !assert.Nil(err) {
		return
	}

	t.Logf("Parsed payload: %v", payload)
	ctx = SetConvey(ctx, payload)

	value, ok = Convey(ctx)
	assert.Equal(value, payload)
	assert.True(ok)

	func() {
		defer func() {
			assert.Nil(recover())
		}()

		assert.Equal(payload, MustConvey(ctx))
	}()
}

func TestToken(t *testing.T) {
	assert := assert.New(t)
	ctx := context.Background()

	value, ok := Token(ctx)
	if !assert.Nil(value) || !assert.False(ok) {
		return
	}

	func() {
		defer func() {
			assert.Equal(NoToken, recover())
		}()

		MustToken(ctx)
	}()

	request, err := http.NewRequest("GET", "", nil)
	if !assert.Nil(err) {
		return
	}

	request.Header.Add("Authorization", basicAuth)
	token, err := secure.NewToken(request)
	if !assert.NotNil(token) || !assert.Nil(err) {
		return
	}

	t.Logf("Parsed token: %v", token)
	ctx = SetToken(ctx, token)

	value, ok = Token(ctx)
	assert.Equal(value, token)
	assert.True(ok)

	func() {
		defer func() {
			assert.Nil(recover())
		}()

		assert.Equal(token, MustToken(ctx))
	}()
}
