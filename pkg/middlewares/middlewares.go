package middlewares

type MiddlewareName string

func (m MiddlewareName) String() string {
	return string(m)
}
