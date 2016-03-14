package sebarmod

func NewServer(host string) *Server{
    s := new(Server)
    s.Host = host
    return s
}