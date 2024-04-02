package authn

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/26huitailang/yogo/framework/gin"
)

type AuthnApi struct {
}

func Register(r *gin.Engine) error {
	datastore = NewDatasotre()
	webAuthn, _ = NewWebAuthNService()
	api := NewAuthnApi()

	webauthnGroup := r.Group("/webauthn")
	webauthnGroup.GET("/register/begin/:username", api.RegisterBegin)
	webauthnGroup.POST("/register/finish/:username", api.RegisterFinish)
	webauthnGroup.GET("/login/begin/:username", api.LoginBegin)
	webauthnGroup.POST("/login/finish/:username", api.LoginFinish)
	return nil
}

func NewAuthnApi() *AuthnApi {
	return &AuthnApi{}
}

// RegisterBegin godoc
// @Summary 开始注册流程
// @Description
// @Produce  json
// @Tags webauthn
// @Success 200 array *webauthn.protocol.CredentialCreateOptions
// @Router /webauthn/register/begin/:username [get]
func (api *AuthnApi) RegisterBegin(c *gin.Context) {
	username := c.Param("username")
	user, err := datastore.GetUser(username) // Find or create the new user
	if err != nil {
		displayName := strings.Split(username, "@")[0]
		user = NewUser(username, displayName)
		datastore.PutUser(user)
	}
	options, session, err := webAuthn.BeginRegistration(user)
	if err != nil {
		c.AbortWithError(500, err)
		return
	}
	datastore.SaveWebAuthnSession(username, session)
	c.JSON(http.StatusOK, options)
}

// RegisterFinish godoc
// @Summary 完成注册
// @Description
// @Produce  json
// @Tags webauthn
// @Success 200 {array} UserDTO
// @Router /demo/demo2 [get]
func (api *AuthnApi) RegisterFinish(c *gin.Context) {
	logger := c.MustMakeLog()
	username := c.Param("username")
	user, err := datastore.GetUser(username)
	if err != nil {
		c.AbortWithError(500, err)
		return
	}

	// Get the session data stored from the function above
	session, err := datastore.GetWebAuthnSession(username)
	if err != nil {
		c.AbortWithError(500, err)
		return
	}

	credential, err := webAuthn.FinishRegistration(user, *session, c.Request)
	if err != nil {
		logger.Error(c, "finish registration failed", map[string]interface{}{"err": fmt.Sprintf("%+v", err), "session": fmt.Sprintf("%+v", session), "user": fmt.Sprintf("%+v", user)})
		// Handle Error and return.
		c.AbortWithError(500, err)
		return
	}

	logger.Debug(c, "[webauthn] finish registration", map[string]interface{}{"credential": credential})
	user.AddCredential(*credential)

	c.JSON(http.StatusOK, "Registration Success")
}

// LoginBegin godoc
// @Summary 开始登录
// @Description
// @Produce  json
// @Tags webauthn
// @Success 200 array *webauthn.protocol.CredentialCreateOptions
// @Router /webauthn/login/begin/:username [get]
func (api *AuthnApi) LoginBegin(c *gin.Context) {
	username := c.Param("username")
	user, err := datastore.GetUser(username) // Find or create the new user
	if err != nil {
		c.AbortWithError(500, err)
		return
	}
	options, session, err := webAuthn.BeginLogin(user)
	if err != nil {
		c.AbortWithError(500, err)
		return
	}
	datastore.SaveWebAuthnSession(username, session)
	c.JSON(http.StatusOK, options)
}

// LoginFinish godoc
// @Summary 完成登录
// @Description
// @Produce  json
// @Tags webauthn
// @Success 200 {array} UserDTO
// @Router /webauthn/login/finish/:username [post]
func (api *AuthnApi) LoginFinish(c *gin.Context) {
	username := c.Param("username")
	user, err := datastore.GetUser(username)
	if err != nil {
		c.AbortWithError(500, err)
		return
	}

	// Get the session data stored from the function above
	session, err := datastore.GetWebAuthnSession(username)
	if err != nil {
		c.AbortWithError(500, err)
		return
	}

	// in an actual implementation, we should perform additional checks on
	// the returned 'credential', i.e. check 'credential.Authenticator.CloneWarning'
	// and then increment the credentials counter
	_, err = webAuthn.FinishLogin(user, *session, c.Request)
	if err != nil {
		c.AbortWithError(500, err)
		return
	}

	c.JSON(http.StatusOK, "Login Success")
}
