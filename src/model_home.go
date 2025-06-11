package main

func Home_page(req Req) Res {
	res := Res(`
		<div hx-swap-oob="innerHTML:#ws">
			<h2>Hello {name}!</h2>
		</div>
	`)
	return res
}
