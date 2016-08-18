# proxy_server
execute:

go run proxy_server.go -host=habr.ru -old=python -new=Python

# proxy_server with NewSingleHostReverseProxy  
go run proxy_server2.go -host=habr.ru -old=python -new=Python
