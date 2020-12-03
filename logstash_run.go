package mrlog

import (
	"github.com/labstack/echo"
	"github.com/sirupsen/logrus"
)

//Debugn debug level with only default field
func Debugn(args ...interface{}) {
	logrus.Info(args...)
}

//Infon info level with only default field
func Infon(args ...interface{}) {
	logrus.Info(args...)
}

//Warnn warn level with only default field
func Warnn(args ...interface{}) {
	logrus.Warn(args...)
}

//Errorn error level with only default field
func Errorn(args ...interface{}) {
	logrus.Error(args...)
}

//Fataln fatal level with only default field
func Fataln(args ...interface{}) {
	logrus.Fatal(args...)
}

//Debug debug level with custom field
func Debug(c echo.Context, args ...interface{}) {
	getField(c).Debug(args...)
}

//Info info level with custom field
func Info(c echo.Context, args ...interface{}) {
	getField(c).Info(args...)
}

//Warn warn level with custom field
func Warn(c echo.Context, args ...interface{}) {
	getField(c).Warn(args...)
}

//Error error level with custom field
func Error(c echo.Context, args ...interface{}) {
	getField(c).Error(args...)
}

//Fatal fatal level with custom field
func Fatal(c echo.Context, args ...interface{}) {
	getField(c).Fatal(args...)
}

func getField(c echo.Context) *logrus.Entry {
	defaultFields := c.Get("defaultFields")
	if defaultFields != nil {
		return defaultFields.(*logrus.Entry)
	}

	header := c.Request().Header

	return logrus.WithFields(logrus.Fields{
		"path":     c.Path(),
		"clientid": header.Get("clientid"),
		"uid":      header.Get("uid"),
		"uip":      header.Get("uip"),
		"traceid":  header.Get("traceid"),
		"method":   c.Request().Method,
	})
}
