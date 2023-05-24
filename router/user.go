package router

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"traQuest/model"

	"github.com/google/uuid"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	traqoauth2 "github.com/ras0q/traq-oauth2"
	traq "github.com/traPtitech/go-traq"
	"golang.org/x/oauth2"
)

type User struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

var (
	clientID    = os.Getenv("CLIENT_ID")
	redirectURL = os.Getenv("REDIRECT_URL")
	conf        = traqoauth2.NewConfig(clientID, redirectURL)
)

func authorizeHandler(c echo.Context) error {
	codeVerifier, err := traqoauth2.GenerateCodeVerifier()
	if err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError, "something wrong in generating code_verifier")
	}

	sess, err := session.Get("sessions", c)
	if err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError, "something wrong in getting session")
	}
	sess.Values["code_verifier"] = codeVerifier
	err = sess.Save(c.Request(), c.Response())
	if err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError, "something wrong in saving session")
	}

	codeChallengeMethod, ok := traqoauth2.CodeChallengeMethodFromStr(c.QueryParam("method"))
	if !ok {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusBadRequest, "invalid code_challenge_method")
	}

	codeChallenge, err := traqoauth2.GenerateCodeChallenge(codeVerifier, codeChallengeMethod)
	if err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError, "something wrong in generating code_challenge")
	}

	authCodeURL := conf.AuthCodeURL(
		c.QueryParam("state"),
		traqoauth2.WithCodeChallenge(codeChallenge),
		traqoauth2.WithCodeChallengeMethod(codeChallengeMethod),
	)

	return c.Redirect(http.StatusFound, authCodeURL)
}

func callbackHandler(c echo.Context) error {
	sess, err := session.Get("sessions", c)
	if err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError, "something wrong in getting session")
	}

	codeVerifier, ok := sess.Values["code_verifier"].(string)
	if !ok {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusUnauthorized, "code_verifier is not set")
	}

	code := c.QueryParam("code")
	if code == "" {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusBadRequest, "code is empty")
	}

	ctx := c.Request().Context()
	token, err := conf.Exchange(
		ctx,
		code,
		traqoauth2.WithCodeVerifier(codeVerifier),
	)
	if err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	sess.Values["access_token"] = token
	err = sess.Save(c.Request(), c.Response())
	if err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError, "something wrong in saving session")
	}

	user, err := getMeTraq(c)
	if err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	// まだなければ作成
	userID, err := uuid.Parse(user.Id)
	if err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	isUserExist, err := model.IsUserExist(ctx, userID)
	if err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	if !isUserExist {
		ctx = c.Request().Context()
		createdUser, err := model.CreateUser(ctx, userID, user.Name)
		if err != nil {
			c.Logger().Error(err)
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}
		return echo.NewHTTPError(http.StatusCreated, &createdUser)
	}

	return echo.NewHTTPError(http.StatusOK, "ok")
}

func getMe(c echo.Context) error {
	user, err := getMeTraq(c)
	if err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	userID, err := uuid.Parse(user.Id)
	if err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	ctx := c.Request().Context()
	me, err := model.GetUser(ctx, userID)
	if err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return echo.NewHTTPError(http.StatusOK, &me)
}

func getMeTraq(c echo.Context) (*traq.MyUserDetail, error) {
	ctx := c.Request().Context()

	sess, err := session.Get("sessions", c)
	sess.Options.SameSite = http.SameSiteNoneMode
	sess.Options.Secure = true
	if err != nil {
		return nil, err
	}
	fmt.Println(sess)
	fmt.Println(sess.Values)
	token, ok := sess.Values["access_token"].(*oauth2.Token)
	if !ok {
		return nil, errors.New("access_token is not set")
	}

	traqconf := traq.NewConfiguration()
	traqconf.HTTPClient = conf.Client(ctx, token)
	client := traq.NewAPIClient(traqconf)

	user, _, err := client.MeApi.GetMe(ctx).Execute()
	if err != nil {
		return nil, err
	}
	return user, nil
}
