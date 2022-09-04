package services

import (
	"time"

	"github.com/pkg/errors"
)

const (
	NoneDuration = time.Duration(-1)
)

var ErrKeyNotFound = errors.New("key not found")
var ErrTypeNotOk = errors.New("val type not ok")
