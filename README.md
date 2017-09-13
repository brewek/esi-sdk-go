# Introduction
This is a simple GoLang SDK for ESI ([EVE Swagger Interface](https://esi.tech.ccp.is)) for convenient access to their data.
I implemented mostly the methods that I needed for my project, but any pull requests / tips / advices on code is highly appreciated.

# Installation
TODO

# Usage
I import it using
```
esi "lossprevented.pl/esi-sdk-go"
```

And than get client by
```
tb := auth.TokenBody{
	Code: ...,
	RefreshToken: ...,
	GrantType: auth.GrantType...
}
client, _ := auth.GetClient(tb, eveAuthClientID, eveAuthSecretKey)
```

Now it's all relatively easy:
```
verifyResponse = auth.VerifyBearerToken(client.AccessToken)
runnerInfo := client.GetCharacter(verifyResponse.CharacterID)
currentLocation := client.GetCharacterLocation(verifyResponse.CharacterID)
// etc...
```

# Authentication
Good documentation about details of the process can be found [here](https://eveonline-third-party-documentation.readthedocs.io/en/latest/sso/authentication.html).
After you get the Bearer token you can confirm it using auth.VerifyBearerToken(string).

### Example
TODO: provide an example from other project
