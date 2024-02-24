# wechat-token-hub
![Coverage](https://img.shields.io/badge/Coverage-94.0%25-brightgreen)

 The project is a centralized control server designed to unify the application, caching, and refreshing of access tokens and tickets for the WeChat API. These credentials have limits on request frequency and number, and applying for new credentials will render the old ones invalid. Some credentials also have IP whitelists, making it best practice to use a global service to manage them separately. For more information, please refer to https://developers.weixin.qq.com/doc/offiaccount/Basic_Information/Get_access_token.html.

## Table of Contents

- [wechat-token-hub](#wechat-token-hub)
  - [Table of Contents](#table-of-contents)
  - [Installation](#installation)
  - [Usage](#usage)
  - [API Documentation](#api-documentation)
    - [Authorization Header:](#authorization-header)
    - [Rotate Query:](#rotate-query)
  - [License](#license)

## Installation

Describe how to install the project and its dependencies.

```sh
$ git clone https://github.com/waynecraig/wechat-token-hub.git
$ cd wechat-token-hub
$ go build -o ./bin/wechat-token-hub ./cmd/main.go
```

## Usage

Describe how to use the project.

```sh
$ WECHAT_API_ROOT=https://api.weixin.qq.com APPID={APPID} APPSECRET={SECRET} JWT_KEY_{kid}={KEY} ./bin/wechat-token-hub
```

| Environment Variable | Description |
| --- | --- |
| WECHAT_API_ROOT | The root URL for the WeChat API |
| APPID | The unique identifier for your WeChat Official Account  |
| APPSECRET | The secret key for your WeChat Official Account |
| JWT_KEY_{kid} | The secret key used for JSON Web Token (JWT) authentication |

## API Documentation

Examples:

1. GET /access-token

Request:
```
GET /access-token HTTP/1.1
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJraWQiOiJrZXktMSIsImlhdCI6MTYzMjQ2NzA4OCwiYXVkIjoid2VjaGF0LXRva2VuLWh1YiJ9.l3Uq3jHMZ-XEUUOLlrVh6by7O3cxNgMpy3eueaEz7To
```

Response:
```raw
{{ACCESS_TOKEN_STRING}}
```

2. GET /ticket?type=jsapi

Request:
```
GET /ticket?type=jsapi HTTP/1.1
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJraWQiOiJrZXktMSIsImlhdCI6MTYzMjQ2NzA4OCwiYXVkIjoid2VjaGF0LXRva2VuLWh1YiJ9.l3Uq3jHMZ-XEUUOLlrVh6by7O3cxNgMpy3eueaEz7To
```

Response:
```raw
{{TICKET_STRING}}
```

### Authorization Header:

The Authorization header is a required header for both endpoints. It should contain a JSON Web Token (JWT) that is signed with the secret found in the environment variable JWT_KEY_{kid}. The kid (key ID) header specifies which key to use for verification. The JWT should be generated by the client's authentication system and should contain the necessary user or application credentials. Additionally, the JWT should have an audience equal to "wechat-token-hub" to ensure that it is authorized for use with the WeChat Token Hub.

### Rotate Query:

The rotate_token query parameter is used to force the server to refresh the access token. If the access token has expired, the server will automatically refresh the token, but if the client needs to refresh the token before it expires, it can make a request with the rotate_token query parameter set to the old token.

Similarly, the rotate_ticket query parameter is used to force the server to refresh the ticket. If the ticket has expired, the server will automatically refresh the ticket, but if the client needs to refresh the ticket before it expires, it can make a request with the rotate_ticket query parameter set to the old ticket.

Note that the rotate query parameters are optional, and if not provided, the server will return the current access token or ticket without refreshing it.

## License

This project is licensed under the [MIT License](LICENSE).