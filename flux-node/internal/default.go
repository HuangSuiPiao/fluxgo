package internal

import (
	"bytes"
	"fmt"
	"github.com/bytepowered/flux/flux-node"
	fluxpkg "github.com/bytepowered/flux/flux-pkg"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/random"
	"io"
	"io/ioutil"
	"net/url"
)

// 默认对RequestBody的表单数据进行解析
func DefaultRequestBodyResolver(webex flux.ServerWebContext) url.Values {
	return webex.FormVars()
}

func DefaultIdentifier(ctx interface{}) string {
	echoc, ok := ctx.(echo.Context)
	fluxpkg.Assert(ok, "<context> must be echo.context")
	id := echoc.Request().Header.Get(flux.XRequestId)
	if "" != id {
		return id
	}
	echoc.Request().Header.Set("X-RequestId-By", "flux")
	return "fxid_" + random.String(32)
}

// Body缓存，允许通过 GetBody 多次读取Body
func RepeatableReader(next echo.HandlerFunc) echo.HandlerFunc {
	// 包装Http处理错误，统一由HttpErrorHandler处理
	return func(echo echo.Context) error {
		request := echo.Request()
		data, err := ioutil.ReadAll(request.Body)
		if nil != err {
			return &flux.ServeError{
				StatusCode: flux.StatusBadRequest,
				ErrorCode:  flux.ErrorCodeGatewayInternal,
				Message:    flux.ErrorMessageRequestPrepare,
				CauseError: fmt.Errorf("read request body, method: %s, uri:%s, err: %w", request.Method, request.RequestURI, err),
			}
		}
		request.GetBody = func() (io.ReadCloser, error) {
			return ioutil.NopCloser(bytes.NewBuffer(data)), nil
		}
		// 恢复Body，但ParseForm解析后，request.Body无法重读，需要通过GetBody
		request.Body = ioutil.NopCloser(bytes.NewBuffer(data))
		return next(echo)
	}
}
