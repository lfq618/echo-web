package www

import (
	"net/http"

	"echo-web/model"
	"echo-web/module/auth"
	"echo-web/module/log"
)

type LoginForm struct {
	Nickname string `form:"nickname" binding:"required"`
	Password string `form:"password" binding:"required"`
}

func LoginHandler(c *Context) error {
	redirect := c.QueryParam(auth.RedirectParam)

	a := c.Auth()
	if a.User.IsAuthenticated() {
		if redirect == "" {
			redirect = "/"
		}
		c.Redirect(http.StatusMovedPermanently, redirect)
		return nil
	}

	c.Set("tmpl", "www/login")
	c.Set("data", map[string]interface{}{
		"title":         "Login",
		"redirectParam": auth.RedirectParam,
		"redirect":      redirect,
	})

	return nil
}

func LoginPostHandler(c *Context) error {
	redirect := c.QueryParam(auth.RedirectParam)
	if redirect == "" {
		redirect = "/"
	}

	a := c.Auth()
	if a.User.IsAuthenticated() {
		c.Redirect(http.StatusMovedPermanently, redirect)
		return nil
	}

	loginURL := c.Request().RequestURI

	var form LoginForm
	if err := c.Bind(&form); err == nil {
		var User model.User
		u := User.GetUserByNicknamePwd(form.Nickname, form.Password)

		if u != nil {
			session := c.Session()
			err := auth.AuthenticateSession(session, u)
			if err != nil {
				c.JSON(http.StatusBadRequest, err)
			}
			c.Redirect(http.StatusMovedPermanently, redirect)
			return nil
		} else {
			c.Redirect(http.StatusMovedPermanently, loginURL)
			return nil
		}
	} else {
		params, _ := c.FormParams()
		log.DebugPrint("Login form params: %v", params)
		log.DebugPrint("Login form bind Error: %v", err)
		c.Redirect(http.StatusMovedPermanently, loginURL)
		return nil
	}

	return nil
}

func LogoutHandler(c *Context) error {
	session := c.Session()
	a := c.Auth()
	auth.Logout(session, a.User)

	redirect := c.QueryParam(auth.RedirectParam)
	if redirect == "" {
		redirect = "/"
	}

	c.Redirect(http.StatusMovedPermanently, redirect)

	return nil
}

func RegisterHandler(c *Context) error {
	redirect := c.QueryParam(auth.RedirectParam)

	a := c.Auth()
	if a.User.IsAuthenticated() {
		if redirect == "" {
			redirect = "/"
		}
		c.Redirect(http.StatusMovedPermanently, redirect)
		return nil
	}

	c.Set("tmpl", "www/register")
	c.Set("data", map[string]interface{}{
		"title":         "Register",
		"redirectParam": auth.RedirectParam,
		"redirect":      redirect,
	})

	return nil
}

func RegisterPostHandler(c *Context) error {
	redirect := c.QueryParam(auth.RedirectParam)
	if redirect == "" {
		redirect = "/"
	}

	a := c.Auth()
	if a.User.IsAuthenticated() {
		c.Redirect(http.StatusMovedPermanently, redirect)
		return nil
	}

	var form LoginForm
	if err := c.Bind(&form); err == nil {
		var User model.User
		u := User.AddUserWithNicknamePwd(form.Nickname, form.Password)
		if u != nil {
			session := c.Session()
			err := auth.AuthenticateSession(session, u)
			if err != nil {
				c.JSON(http.StatusBadRequest, err)
			}
			c.Redirect(http.StatusMovedPermanently, redirect)
			return nil
		} else {
			log.DebugPrint("Register user add error")

			s := c.Session()
			s.AddFlash("Register user add error", "_error")

			// registerURL := c.Request().URI()
			// c.Redirect(http.StatusMovedPermanently, registerURL)
			c.Set("tmpl", "www/register")
			c.Set("data", map[string]interface{}{
				"title":         "Register",
				"redirectParam": auth.RedirectParam,
				"redirect":      redirect,
			})
			return nil
		}
	} else {
		log.DebugPrint("Register form bind Error: %v", err)

		s := c.Session()
		s.AddFlash("Register form bind Error:"+err.Error(), "_error")

		// registerURL := c.Request().URI()
		// c.Redirect(http.StatusMovedPermanently, registerURL)
		c.Set("tmpl", "www/register")
		c.Set("data", map[string]interface{}{
			"title":         "Register",
			"redirectParam": auth.RedirectParam,
			"redirect":      redirect,
		})
		return nil
	}

	return nil
}
