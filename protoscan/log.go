package protoscan

import (
	"reflect"

	"github.com/Sirupsen/logrus"
)

// -----------------------------------------------------------------------------

type pkg struct{}

var log = logrus.WithField("pkg", reflect.TypeOf(pkg{}).PkgPath())
