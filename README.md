# Player Developer tech assignment

## Package structure

The package is structured as follows:

### `api`

Contains the API definition and the generated client.

#### How to run

```bash
cd api
go run main.go
```

#### How to run test

```bash
cd api
go test
```

### `tool`

Contains the tool that executes the requests for each player.

#### How to run

```bash
cd tool
go run main.go
```

#### How to run test

```bash
cd tool
go test
```


### Example of a .csv file: (see `tool/mac_addresses.csv`)

```
mac_addresses, id1, id2, id3
a1:bb:cc:dd:ee:ff, 1, 2, 3
a2:bb:cc:dd:ee:ff, 1, 2, 3
a3:bb:cc:dd:ee:ff, 1, 2, 3
a4:bb:cc:dd:ee:ff, 1, 2, 3
```

## API Documentation

### Requests

#### Create a new profile

`POST /profiles/{mac_address}`

#### Parameters

mac_address : the mac address of the player

#### Body

```
Headers

Content-Type: application/json
x-client-id: required (see below)
x-authentication-token: required

Body

{
  "profile": {    
    "applications": [
      {
        "applicationId": "music_app"
        "version": "v1.4.10"
      },
      {
        "applicationId": "diagnostic_app",
        "version": "v1.2.6"
      },
      {
        "applicationId": "settings_app",
        "version": "v1.1.5"
      }
    ]
  }
}
```

### Responses

#### 200
```
Headers

Content-Type: application/json

Body

{
  "profile": {    
    "applications": [
      {
        "applicationId": "music_app"
        "version": "v1.4.10"
      },
      {
        "applicationId": "diagnostic_app",
        "version": "v1.2.6"
      },
      {
        "applicationId": "settings_app",
        "version": "v1.1.5"
      }
    ]
  },
  "clientId": "abcd",
  "macAddress": "a1:bb:cc:dd:ee:ff",
}
```

#### 401

```
Headers

Content-Type: application/json

Body

{
  "statusCode": 401,
  "error": "Unauthorized",
  "message": "invalid clientId or token supplied"
}
```

#### 409

```
Headers

Content-Type: application/json

Body

{
  "statusCode": 409,
  "error": "Conflict",
  "message": "client already exists"
}
```

#### 500

```
Headers

Content-Type: application/json

Body

{
  "statusCode": 500,
  "error": "Internal Server Error",
  "message": "An internal server error occurred"
}
```

#### Get a profile

`GET /profiles/{mac_address}`

#### Parameters

mac_address : the mac address of the player

#### Body

```
Headers

Content-Type: application/json
x-client-id: required (see below)
x-authentication-token: required
```

#### Responses

#### 200

```
Headers

Content-Type: application/json

Body

{
  "profile": {    
    "applications": [
      {
        "applicationId": "music_app"
        "version": "v1.4.10"
      },
      {
        "applicationId": "diagnostic_app",
        "version": "v1.2.6"
      },
      {
        "applicationId": "settings_app",
        "version": "v1.1.5"
      }
    ]
  },
  "clientId": "abcd",
  "macAddress": "a1:bb:cc:dd:ee:ff",
}
```

#### 401

```
Headers

Content-Type: application/json

Body

{
  "statusCode": 401,
  "error": "Unauthorized",
  "message": "invalid clientId or token supplied"
}
```

#### 404

```
Headers

Content-Type: application/json

Body

{
  "statusCode": 404,
  "error": "Not Found",
  "message": "profile of client mac_address not found"
}
```


#### 500

```
Headers

Content-Type: application/json

Body

{
  "statusCode": 500,
  "error": "Internal Server Error",
  "message": "An internal server error occurred"
}
```

#### Delete a profile

`DELETE /profiles/{mac_address}`

#### Body

```
Headers

Content-Type: application/json
x-client-id: required (see below)
x-authentication-token: required
```

#### Parameters

mac_address : the mac address of the player

#### Responses

#### 200

```
Headers

Content-Type: application/json

Body

{
  "profile": {    
    "applications": [
      {
        "applicationId": "music_app"
        "version": "v1.4.10"
      },
      {
        "applicationId": "diagnostic_app",
        "version": "v1.2.6"
      },
      {
        "applicationId": "settings_app",
        "version": "v1.1.5"
      }
    ]
  },
  "clientId": "abcd",
  "macAddress": "a1:bb:cc:dd:ee:ff",
}
```

#### 401

```
Headers

Content-Type: application/json

Body

{
  "statusCode": 401,
  "error": "Unauthorized",
  "message": "invalid clientId or token supplied"
}
```

#### 500

```
Headers

Content-Type: application/json

Body

{
  "statusCode": 500,
  "error": "Internal Server Error",
  "message": "An internal server error occurred"
}
```

#### Update a profile

`PUT /profiles/{mac_address}`

#### Parameters

mac_address : the mac address of the player

#### Body

```
Headers

Content-Type: application/json
x-client-id: required (see below)
x-authentication-token: required

Body

{
  "profile": {    
    "applications": [
      {
        "applicationId": "music_app"
        "version": "v1.4.10"
      },
      {
        "applicationId": "diagnostic_app",
        "version": "v1.2.6"
      },
      {
        "applicationId": "settings_app",
        "version": "v1.1.5"
      }
    ]
  }
}
```

### Responses

#### 200
```
Headers

Content-Type: application/json

Body

{
  "profile": {    
    "applications": [
      {
        "applicationId": "music_app"
        "version": "v1.4.10"
      },
      {
        "applicationId": "diagnostic_app",
        "version": "v1.2.6"
      },
      {
        "applicationId": "settings_app",
        "version": "v1.1.5"
      }
    ]
  },
  "clientId": "abcd",
  "macAddress": "a1:bb:cc:dd:ee:ff",
}
```

#### 401

```
Headers

Content-Type: application/json

Body

{
  "statusCode": 401,
  "error": "Unauthorized",
  "message": "invalid clientId or token supplied"
}
```

#### 409

```
Headers

Content-Type: application/json

Body

{
  "statusCode": 409,
  "error": "Conflict",
  "message": "child \"profile\" fails because [child \"applications\" fails because [\"applications\" is required]]"
}
```

#### 500

```
Headers

Content-Type: application/json

Body

{
  "statusCode": 500,
  "error": "Internal Server Error",
  "message": "An internal server error occurred"
}
```
