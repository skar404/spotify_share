package spotify_type

type TokenReq struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int64  `json:"expires_in"`
	Scope       string `json:"scope"`
	TokenType   string `json:"token_type"`
}

type TokenOrRefreshReq struct {
	TokenReq
	RefreshToken string `json:"refresh_token"`
}

/*
Test expire token :)
{
  "access_token": "BQDofVcSja2pUfMiZ5_1s5lfdyGsQDMfg7ehUPksEsbujcjK7SwqkdljOGIcr7LP_8dfgHGNhSvBnY6x_gnmVWkhEDyK50sl5vTfkZiR3BqXC0qif68LGZTfttyc6H8CMZRkb3EzMCDcQ4W-zT1iNLpimKFI",
  "expires_in": 3600,
  "scope": "app-remote-control streaming user-read-currently-playing user-read-recently-played",
  "token_type": "Bearer"
}
*/
