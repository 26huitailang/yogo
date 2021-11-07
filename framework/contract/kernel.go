package contract

import "net/http"

const KernelKey = "yogo:kernel"

type Kernel interface {
	HttpEngine() http.Handler
}
