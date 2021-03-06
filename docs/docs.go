// Package docs GENERATED BY THE COMMAND ABOVE; DO NOT EDIT
// This file was generated by swaggo/swag
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {},
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/{namespace}/job": {
            "get": {
                "description": "Finds a job that either is not in-use or has been inactive for more than the specified time.",
                "consumes": [
                    "*/*"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "root"
                ],
                "summary": "Get an unused job and lock it as in-use",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Namespace of job(s)",
                        "name": "namespace",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "integer",
                        "description": "Minimum age of job (last_used_on) in minutes before assuming it's no longer in use (optional, defaults to never)",
                        "name": "min_age",
                        "in": "path"
                    },
                    {
                        "type": "string",
                        "description": "API Key",
                        "name": "api_key",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/main.JobResponse"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/main.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/main.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/{namespace}/jobs/": {
            "get": {
                "description": "Gets all jobs",
                "consumes": [
                    "*/*"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "root"
                ],
                "summary": "Get all jobs",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Namespace of job(s)",
                        "name": "namespace",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "API Key",
                        "name": "api_key",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/main.JobsResponse"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/main.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/main.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/{namespace}/jobs/{jobId}": {
            "get": {
                "description": "Gets a job with specified id",
                "consumes": [
                    "*/*"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "root"
                ],
                "summary": "Get a job with specified id",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Namespace of job(s)",
                        "name": "namespace",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Job ID",
                        "name": "jobId",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "API Key",
                        "name": "api_key",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/main.JobResponse"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/main.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/main.ErrorResponse"
                        }
                    }
                }
            },
            "put": {
                "description": "Overwrites job if it exists.",
                "consumes": [
                    "application/json",
                    " */*"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "root"
                ],
                "summary": "Insert/replace job with specified id",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Namespace of job(s)",
                        "name": "namespace",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Job ID",
                        "name": "jobId",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "API Key",
                        "name": "api_key",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Who's using this job (optional) eg. hostname of machine using it.",
                        "name": "in_use_by",
                        "in": "query"
                    },
                    {
                        "description": "Job metadata (optional) - arbitrary json can be stored in {meta: {...}}",
                        "name": "body",
                        "in": "body",
                        "schema": {
                            "$ref": "#/definitions/main.PutJobByIdBody"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/main.SuccessResponse"
                        }
                    }
                }
            },
            "post": {
                "description": "Used together with GET /job's min_age parameter so that inactive jobs can be reused. Fails if the job doesn't exist",
                "consumes": [
                    "*/*"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "root"
                ],
                "summary": "Update a job, marking it as still in use",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Namespace of job(s)",
                        "name": "namespace",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Job ID",
                        "name": "jobId",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "API Key",
                        "name": "api_key",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/main.SuccessResponse"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/main.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/main.ErrorResponse"
                        }
                    }
                }
            },
            "delete": {
                "description": "Deletes a job with specified id",
                "consumes": [
                    "*/*"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "root"
                ],
                "summary": "Delete a job with specified id",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Namespace of job(s)",
                        "name": "namespace",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Job ID",
                        "name": "jobId",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "API Key",
                        "name": "api_key",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/main.SuccessResponse"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/main.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/main.ErrorResponse"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "main.ErrorResponse": {
            "type": "object",
            "properties": {
                "error": {
                    "type": "string"
                }
            }
        },
        "main.Job": {
            "type": "object",
            "properties": {
                "created_on": {
                    "type": "string"
                },
                "in_use": {
                    "type": "boolean"
                },
                "in_use_by": {
                    "type": "string"
                },
                "job_key": {
                    "type": "string"
                },
                "last_used_on": {
                    "type": "string"
                },
                "meta": {}
            }
        },
        "main.JobResponse": {
            "type": "object",
            "properties": {
                "result": {
                    "$ref": "#/definitions/main.Job"
                }
            }
        },
        "main.JobsResponse": {
            "type": "object",
            "properties": {
                "result": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/main.Job"
                    }
                }
            }
        },
        "main.PutJobByIdBody": {
            "type": "object",
            "properties": {
                "meta": {}
            }
        },
        "main.SuccessResponse": {
            "type": "object",
            "properties": {
                "result": {
                    "type": "string"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "",
	Host:             "",
	BasePath:         "",
	Schemes:          []string{},
	Title:            "Job Manager API",
	Description:      "",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
