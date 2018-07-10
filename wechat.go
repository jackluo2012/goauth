package goauth

import (
	"github.com/sanxia/glib"
)

/* ================================================================================
 * Oauth WeChat
 * qq group: 582452342
 * email   : 2091938785@qq.com
 * author  : 美丽的地球啊 - mliu
 * ================================================================================ */

type (
	OauthWeChat struct {
		Oauth
	}

	WeChatAccessTokenResponse struct {
		ErrCode      string `form:"errcode" json:"errcode"`
		ErrMsg       string `form:"errmsg" json:"errmsg"`
		OpenId       string `form:"openid" json:"openid"`
		UnionId      string `form:"unionid" json:"unionid"`
		AccessToken  string `form:"access_token" json:"access_token"`
		RefreshToken string `form:"refresh_token" json:"refresh_token"`
		ExpiresIn    int    `form:"expires_in" json:"expires_in"`
		Scope        string `form:"scope" json:"scope"`
	}

	//微信用户信息响应结构
	WeChatUserInfoResponse struct {
		ErrCode    string   `form:"errcode" json:"errcode"`
		ErrMsg     string   `form:"errmsg" json:"errmsg"`
		OpenId     string   `form:"openid" json:"openid"`         //普通用户的标识，对当前开发者帐号唯一
		UnionId    string   `form:"unionid" json:"unionid"`       //用户统一标识。针对一个微信开放平台帐号下的应用，同一用户的unionid是唯一的。
		Nickname   string   `form:"nickname" json:"nickname"`
		Subscribe  string   `form:"subscribe" json:"subscribe"`   //是否关注
		HeadImgUrl string   `form:"headimgurl" json:"headimgurl"` //用户头像，最后一个数值代表正方形头像大小（有0、46、64、96、132数值可选，0代表640*640正方形头像），用户没有头像时该项为空
		Sex        int      `form:"sex" json:"sex"`               //普通用户性别，1为男性，2为女性
		Country    string   `form:"country" json:"country"`       //国家，如中国为CN
		Province   string   `form:"province" json:"province"`     //普通用户个人资料填写的省份
		City       string   `form:"city" json:"city"`             //普通用户个人资料填写的城市
		Privileges []string `form:"privilege" json:"privilege"`   //用户特权信息，json数组，如微信沃卡用户为（chinaunicom）
	}

	//微信用户信息响应结构
	WeChatSubscribeUserInfoResponse struct {
		Subscribe  int   `form:"subscribe" json:"subscribe"`      //是否关注
		OpenId     string   `form:"openid" json:"openid"`         //普通用户的标识，对当前开发者帐号唯一
		Nickname   string   `form:"nickname" json:"nickname"`
		HeadImgUrl string   `form:"headimgurl" json:"headimgurl"` //用户头像，最后一个数值代表正方形头像大小（有0、46、64、96、132数值可选，0代表640*640正方形头像），用户没有头像时该项为空
		Sex        uint8      `form:"sex" json:"sex"`               //普通用户性别，1为男性，2为女性
		Country    string   `form:"country" json:"country"`       //国家，如中国为CN
		Province   string   `form:"province" json:"province"`     //普通用户个人资料填写的省份
		City       string   `form:"city" json:"city"`             //普通用户个人资料填写的城市
	}
)

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 初始化WeChat授权
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func NewWeChat(clientId, clientSecret, callbackUri string) IOauth {
	oauth := new(OauthWeChat)
	oauth.ClientId = clientId
	oauth.ClientSecret = clientSecret
	oauth.CallbackUri = callbackUri

	oauth.AuthorizeCodeUri = "https://open.weixin.qq.com/connect/oauth2/authorize"
	oauth.AccessTokenUri = "https://api.weixin.qq.com/sns/oauth2/access_token"
	oauth.AccessTokenGlobleUri = "https://api.weixin.qq.com/cgi-bin/token"
	oauth.RefreshTokenUri = "https://api.weixin.qq.com/sns/oauth2/refresh_token"
	oauth.UserInfoUri = "https://api.weixin.qq.com/cgi-bin/user/info"

	return oauth
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 设置Uri
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *OauthWeChat) SetUri(uriType OauthUriType, uri string) {
	switch uriType {
	case AuthorizeCodeUri:
		s.AuthorizeCodeUri = uri
	case AccessTokenUri:
		s.AccessTokenUri = uri
	case RefreshTokenUri:
		s.RefreshTokenUri = uri
	case UserInfoUri:
		s.UserInfoUri = uri
	}
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 获取鉴权地址
 * state, scope
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *OauthWeChat) GetAuthorizeUrl(args ...string) string {
	state, scope := "wechat", "snsapi_login"

	argCount := len(args)
	if argCount > 0 {
		state = args[0]

		if argCount > 1 {
			scope = args[1]
		}
	}

	params := map[string]interface{}{
		"appid":         s.ClientId,
		"redirect_uri":  glib.QueryEncode(s.CallbackUri),
		"scope":         scope,
		"state":         state,
		"response_type": "code",
	}

	queryString := glib.ToQueryString(params)

	return s.AuthorizeCodeUri + "?" + queryString
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 获取AccessToken
 * {
 * "access_token":"ACCESS_TOKEN",
 * "expires_in":7200,
 * "refresh_token":"REFRESH_TOKEN",
 * "openid":"OPENID",
 * "scope":"SCOPE"
 * }
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *OauthWeChat) GetAccessToken(code string) (*OauthToken, error) {
	var oauthToken *OauthToken

	params := map[string]interface{}{
		"appid":      s.ClientId,
		"secret":     s.ClientSecret,
		"code":       code,
		"grant_type": "authorization_code",
	}

	queryString := glib.ToQueryString(params)

	//获取api响应数据
	resp, err := glib.HttpGet(s.AccessTokenUri, queryString)
	if err == nil {
		//解析json数据
		var tokenResponse *WeChatAccessTokenResponse
		glib.FromJson(resp, &tokenResponse)

		if tokenResponse != nil {
			oauthToken = &OauthToken{
				AccessToken:  tokenResponse.AccessToken,
				RefreshToken: tokenResponse.RefreshToken,
				OpenId:       tokenResponse.OpenId,
				UnionId:      tokenResponse.UnionId,
				Scope:        tokenResponse.Scope,
				ExpiresIn:    tokenResponse.ExpiresIn,
			}

			return oauthToken, nil
		}
	}

	return nil, err
}

func (s *OauthWeChat) GetGlobalAccessToken() (*AccessToken, error) {
	var accessToken *AccessToken
	params := map[string]interface{}{
		"appid":      s.ClientId,
		"secret":     s.ClientSecret,
		"grant_type": "client_credential",
	}

	queryString := glib.ToQueryString(params)

	//获取api响应数据
	resp, err := glib.HttpGet(s.AccessTokenGlobleUri, queryString)
	if err == nil {
		//解析json数据
		var tokenResponse *WeChatAccessTokenResponse
		glib.FromJson(resp, &tokenResponse)

		if tokenResponse != nil && tokenResponse.AccessToken != "" {
			accessToken = &AccessToken{
				AccessToken:  tokenResponse.AccessToken,
				ExpiresIn:    tokenResponse.ExpiresIn,
			}
			return accessToken, nil
		}
	}

	return nil, err
}


/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 刷新AccessToken
 * {
 * "access_token":"ACCESS_TOKEN",
 * "expires_in":7200,
 * "refresh_token":"REFRESH_TOKEN",
 * "openid":"OPENID",
 * "scope":"SCOPE"
 * }
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *OauthWeChat) RefreshAccessToken(refreshToken string) (*OauthToken, error) {
	var oauthToken *OauthToken

	params := map[string]interface{}{
		"appid":         s.ClientId,
		"refresh_token": refreshToken,
		"grant_type":    "refresh_token",
	}

	queryString := glib.ToQueryString(params)

	//获取api响应数据
	resp, err := glib.HttpGet(s.RefreshTokenUri, queryString)
	if err == nil {
		//解析json数据
		var tokenResponse *WeChatAccessTokenResponse
		glib.FromJson(resp, &tokenResponse)

		if tokenResponse != nil {
			oauthToken = &OauthToken{
				AccessToken:  tokenResponse.AccessToken,
				RefreshToken: tokenResponse.RefreshToken,
				OpenId:       tokenResponse.OpenId,
				Scope:        tokenResponse.Scope,
				ExpiresIn:    tokenResponse.ExpiresIn,
			}
		}
		return oauthToken, nil
	}

	return nil, err
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 获取用户信息
 * {
 * "openid":"OPENID",
 * "nickname":"NICKNAME",
 * "sex":1,
 * "province":"PROVINCE",
 * "city":"CITY",
 * "country":"COUNTRY",
 * "headimgurl": "http://wx.qlogo.cn/mmopen/g3MonUZtNHkdmzicIlibx6iaFqAc56vxLSUfpb6n5WKSYVY0ChQKkiaJSgQ1dZuTOgvLLrhJbERQQ4eMsv84eavHiaiceqxibJxCfHe/0",
 * "privilege":[
 * "PRIVILEGE1",
 * "PRIVILEGE2"
 * ],
 * "unionid": " o6_bmasdasdsad6_2sgVt7hMZOPfL"
 * }
 * {"subscribe":1,"openid":"o7gih0mnpB3aYyEpNTW96xXjCG9o","nickname":"jack","sex":1,"language":"zh_CN","city":"Chengdu","province":"Sichuan","country":"China","headimgurl":"http:\/\/thirdwx.qlogo.cn\/mmopen\/l8o5Bj65aCS7xsznlibNGNZXznW3ibGHrxS5cBHWIxkRQGXMExu2C3xboZic0NxjDEXdoy9iciblpRaYJNYcsticiaS4Iv4whIZBtyc\/132","subscribe_time":1531130190,"remark":"","groupid":0,"tagid_list":[],"subscribe_scene":"ADD_SCENE_QR_CODE","qr_scene":0,"qr_scene_str":""}
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *OauthWeChat) GetUserInfo(accessToken, openId string) (*OauthUser, error) {
	var oauthUser *OauthUser
	params := map[string]interface{}{
		"access_token": accessToken,
		"openid":       openId,
		"lang":         "zh-CN",
	}

	queryString := glib.ToQueryString(params)

	//获取api响应数据
	resp, err := glib.HttpGet(s.UserInfoUri, queryString)

	if err == nil {
		var userInfoResponse WeChatSubscribeUserInfoResponse

		//解析json数据
		glib.FromJson(resp, &userInfoResponse)

		oauthUser = &OauthUser{
			Nickname: userInfoResponse.Nickname,
			Avatar:   userInfoResponse.HeadImgUrl,
			Sex:      userInfoResponse.Sex,
			OpenId:      userInfoResponse.OpenId,
			Subscribe:userInfoResponse.Subscribe,
			City:userInfoResponse.City,
			Province:userInfoResponse.Province,
		}

	}

	return oauthUser, err
}
