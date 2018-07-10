# goauth
Oauth2 QQ and WeChat for Golang

Example:
---------------
import "github.com/jackluo2012/goauth"

```go
//QQ oauth
qqOauth := goauth.NewQq("you app id", "you app secret", "you callback url")
qqToken, err := qqOauth.GetAccessToken(code)
qqUserInfo, err := qqOauth.GetUserInfo(qqToken.AccessToken, qqToken.OpenId)

```

```go

//WeChat Oauth
weChatOauth := goauth.NewWeChat("you app id", "you app secret", "you callback url")
weChatToken, err := weChatOauth.GetAccessToken(code)
weChatUserInfo, err := weChatOauth.GetUserInfo(weChatToken.AccessToken, weChatToken.OpenId)
```

