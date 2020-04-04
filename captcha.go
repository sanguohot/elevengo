package elevengo

import (
	"errors"
	"fmt"
	"github.com/deadblue/elevengo/core"
	"github.com/deadblue/elevengo/internal"
	"math/rand"
	"time"
)

const (
	apiCaptcha       = "https://captchaapi.115.com/"
	apiCaptchaSubmit = "https://webapi.115.com/user/captcha"
)

// CaptchaSession holds the information during a CAPTCHA process.
type CaptchaSession struct {
	// The callback function name, hide for caller.
	callback string

	// The CAPTCHA image data.
	// There are 4 Chinese characters on this image, you need call
	// "CaptchaKeysImage" to get a image with 10 Chinese characters and
	// find 4 characters on it which matches with this image, the indexes
	// of the 4 characters is the CAPTCHA code.
	Image []byte
}

// Start a CAPTCHA session.
func (a *Agent) CaptchaStart() (session *CaptchaSession, err error) {
	// Fetch captcha page
	callback := fmt.Sprintf("Close911_%d", time.Now().UnixNano())
	qs := core.NewQueryString().
		WithString("ac", "security_code").
		WithString("type", "web").
		WithString("cb", callback)
	if _, err = a.hc.Get(apiCaptcha, qs); err != nil {
		return
	}
	// Fetch captcha value image
	qs = core.NewQueryString().
		WithString("ct", "index").
		WithString("ac", "code").
		WithInt64("_t", time.Now().Unix())
	image, err := a.hc.Get(apiCaptcha, qs)
	if err != nil {
		return
	}
	// Build session
	session = &CaptchaSession{
		callback: callback,
		Image:    image,
	}
	return
}

// Get CAPTCHA keys image.
// There are 10 Chinese characters on the image, 5 in column and 2 in row.
// You can call this method multiple times, everytime you call it, you will
// get the same 10 characters but in different font.
func (a *Agent) CaptchaKeysImage(session *CaptchaSession) ([]byte, error) {
	qs := core.NewQueryString().
		WithString("ct", "index").
		WithString("ac", "code").
		WithString("t", "all").
		WithInt64("_t", time.Now().Unix())
	return a.hc.Get(apiCaptcha, qs)
}

// Get one CAPTCHA key image.
// You can call this method multiple times, everytime you call it, you will
// get the same character but in different font.
func (a *Agent) CaptchaKeyImage(session *CaptchaSession, index int) ([]byte, error) {
	if index < 0 {
		index = 0
	} else if index > 9 {
		index = 9
	}
	qs := core.NewQueryString().
		WithString("ct", "index").
		WithString("ac", "code").
		WithString("t", "single").
		WithInt("id", index).
		WithInt64("_t", time.Now().Unix())
	return a.hc.Get(apiCaptcha, qs)
}

// Submit the CAPTCHA code.
func (a *Agent) CaptchaSubmit(session *CaptchaSession, code string) (err error) {
	// Get captcha sign
	cb := fmt.Sprintf("jQuery%d_%d", rand.Uint64(), time.Now().UnixNano())
	qs := core.NewQueryString().
		WithString("ac", "code").
		WithString("t", "sign").
		WithString("callback", cb).
		WithInt64("_", time.Now().Unix())
	signResult := &internal.CaptchaSignResult{}
	if err = a.hc.JsonpApi(apiCaptcha, qs, signResult); err != nil {
		return
	}
	// Submit captcha code
	form := core.NewForm().
		WithString("ac", "security_code").
		WithString("type", "web").
		WithString("sign", signResult.Sign).
		WithString("code", code).
		WithString("cb", session.callback)
	submitResult := &internal.CaptchaSubmitResult{}
	err = a.hc.JsonApi(apiCaptchaSubmit, nil, form, submitResult)
	if err == nil && !submitResult.State {
		// TODO: handle error code.
		err = errors.New("submit failed")
	}
	return
}
