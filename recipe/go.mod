module recipe

go 1.18

replace parser => ../parser

require (
	github.com/google/go-cmp v0.5.7
	golang.org/x/net v0.0.0-20220325170049-de3da57026de
	models v0.0.0-00010101000000-000000000000
	parser v0.0.0-00010101000000-000000000000
)

replace models => ../models
