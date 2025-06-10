go:
	go run .

templ:
	~/go/bin/templ generate

dev:
	watchexec -e templ -- make templ &
	watchexec -e go,html,css -- make go
