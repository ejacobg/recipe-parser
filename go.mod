module main

go 1.18

replace parser => ./parser

replace recipe => ./recipe

require (
	golang.org/x/net v0.0.0-20220325170049-de3da57026de
	parser v0.0.0-00010101000000-000000000000
	recipe v0.0.0-00010101000000-000000000000
)

require models v0.0.0-00010101000000-000000000000 // indirect

replace models => ./models
