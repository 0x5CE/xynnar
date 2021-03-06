// Package docs GENERATED BY SWAG; DO NOT EDIT
// This file was generated by swaggo/swag
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {
            "name": "Muazzam Ali Kazmi",
            "url": "https://github.com/0x5CE/xynnar",
            "email": "muazzam_ali@live.com"
        },
        "license": {
            "name": "Public Domain"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/api/characters/{movie_id}": {
            "get": {
                "description": "Retreive characters in a particular movie specified by ID. Also outputs their combined height.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Characters"
                ],
                "summary": "Get movie characters",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Movie ID",
                        "name": "movie_id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Sort by name, gender, or height. Place '-' before it for descending order",
                        "name": "sort",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "Filter by gender",
                        "name": "filter",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/main.charactersGETResp"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/main.errorResp"
                        }
                    }
                }
            }
        },
        "/api/comment/": {
            "post": {
                "description": "Commment on a movie (public comment)",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Comments"
                ],
                "summary": "Comment on a movie",
                "parameters": [
                    {
                        "description": "Movie ID",
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/main.commentPOSTBody"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/main.commentsGETResp"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/main.errorResp"
                        }
                    }
                }
            }
        },
        "/api/comments/{movie_id}": {
            "get": {
                "description": "Comments returns either all comments, or comments about a particular movie. For instance, /comments/4",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Comments"
                ],
                "summary": "Get all comments about a movie by ID.",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Movie ID",
                        "name": "movie_id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/main.commentsGETResp"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/main.errorResp"
                        }
                    }
                }
            }
        },
        "/api/films": {
            "get": {
                "description": "Retreive list of all Star War movies along with title, opening crawl, etc. Retreived in chronological order.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Films"
                ],
                "summary": "Get films list",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/main.filmsGETResp"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/main.errorResp"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "main.Character": {
            "type": "object",
            "properties": {
                "birth_year": {
                    "type": "string"
                },
                "gender": {
                    "type": "string"
                },
                "hair_color": {
                    "type": "string"
                },
                "height": {
                    "type": "string"
                },
                "height_ft": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                }
            }
        },
        "main.Comment": {
            "type": "object",
            "properties": {
                "content": {
                    "type": "string"
                },
                "ip": {
                    "type": "string"
                },
                "movie_id": {
                    "type": "string"
                },
                "timestamp": {
                    "type": "string"
                }
            }
        },
        "main.Film": {
            "type": "object",
            "properties": {
                "episode_id": {
                    "type": "integer"
                },
                "opening_crawl": {
                    "type": "string"
                },
                "release_date": {
                    "type": "string"
                },
                "title": {
                    "type": "string"
                },
                "total_comments": {
                    "type": "integer"
                }
            }
        },
        "main.charactersGETResp": {
            "type": "object",
            "properties": {
                "characters": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/main.Character"
                    }
                },
                "count": {
                    "type": "integer",
                    "example": 10
                },
                "totalHeight": {
                    "type": "string"
                },
                "totalHeight_ft": {
                    "type": "string"
                }
            }
        },
        "main.commentPOSTBody": {
            "type": "object",
            "properties": {
                "comment": {
                    "type": "string",
                    "example": "Great movie!"
                },
                "movie_id": {
                    "type": "integer",
                    "example": 5
                }
            }
        },
        "main.commentsGETResp": {
            "type": "object",
            "properties": {
                "comments": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/main.Comment"
                    }
                },
                "count": {
                    "type": "integer",
                    "example": 10
                }
            }
        },
        "main.errorResp": {
            "type": "object",
            "properties": {
                "count": {
                    "type": "integer",
                    "example": 500
                },
                "message": {
                    "type": "string",
                    "example": "Error fetching film"
                }
            }
        },
        "main.filmsGETResp": {
            "type": "object",
            "properties": {
                "count": {
                    "type": "integer",
                    "example": 10
                },
                "results": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/main.Film"
                    }
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "",
	BasePath:         "/",
	Schemes:          []string{},
	Title:            "Xynnar API",
	Description:      "This is an API to get Star Wars movies list and comment on them.",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
