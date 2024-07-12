ClientID: "PwJbHppaqfvkLezH4IjEI1kLvzstaVze"
Domain: "dev-8te4x4rpnu1cpxe6.us.auth0.com"

// An application running locally
if #Meta.Environment.Type == "development" && #Meta.Environment.Cloud == "local" {
	CallbackURL: "http://localhost:5173/callback"
	LogoutURL: "http://localhost:5173/"
}
if !(#Meta.Environment.Type == "development" && #Meta.Environment.Cloud == "local"){
    CallbackURL: "https://please-set-the-url.com/callback"
    LogoutURL: "https://please-set-the-url.com/"
}
