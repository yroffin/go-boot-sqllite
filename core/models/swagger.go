// Package models for all models
// MIT License
//
// Copyright (c) 2017 yroffin
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.
package models

// SwaggerModel the root model of swagger
type SwaggerModel struct {
	Swagger             string                                `json:"swagger"`
	Info                SwaggerInfo                           `json:"info"`
	Host                string                                `json:"host"`
	BasePath            string                                `json:"basePath"`
	Tags                []SwaggerTags                         `json:"tags"`
	Schemes             []string                              `json:"schemes"`
	Paths               map[string]SwaggerRoute               `json:"paths"`
	SecurityDefinitions map[string]SwaggerSecurityDefinitions `json:"securityDefinitions"`
	Definitions         map[string]SwaggerDefinitions         `json:"definitions"`
	ExternalDocs        SwaggerExternalDocs                   `json:"externalDocs"`
}

// SwaggerInfo the info block
type SwaggerInfo struct {
	Description    string         `json:"description"`
	Version        string         `json:"version"`
	Title          string         `json:"title"`
	TermsOfService string         `json:"termsOfService"`
	Contact        SwaggerContact `json:"contact"`
	License        SwaggerLicense `json:"license"`
}

// SwaggerLicense the contact block
type SwaggerContact struct {
	Email string `json:"email"`
}

// SwaggerLicense the license block
type SwaggerLicense struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

// SwaggerTags the tags block
type SwaggerTags struct {
	Name         string              `json:"name"`
	Description  string              `json:"description"`
	ExternalDocs SwaggerExternalDocs `json:"externalDocs"`
}

// SwaggerExternalDocs the docs block
type SwaggerExternalDocs struct {
	Description string `json:"description"`
	URL         string `json:"url"`
}

// SwaggerRoute the route block
type SwaggerRoute map[string]SwaggerMethodBody

// SwaggerMethodBody the method detail block
type SwaggerMethodBody struct {
	Tags        []string                     `json:"tags"`
	Summary     string                       `json:"summary"`
	Description string                       `json:"description"`
	Consumes    []string                     `json:"consumes"`
	Produces    []string                     `json:"produces"`
	Parameters  []SwaggerMethodParamBody     `json:"parameters"`
	Responses   map[string]SwaggerMethodResp `json:"responses"`
}

// SwaggerMethodParamBody the parameter block
type SwaggerMethodParamBody struct {
	In          string `json:"in"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Required    bool   `json:"required"`
	Type        string `json:"type"`
}

// SwaggerMethodResp the docs block
type SwaggerMethodResp struct {
	Description string            `json:"description"`
	Schema      map[string]string `json:"schema"`
}

// SwaggerSecurity the method block
type SwaggerSecurity map[string][]string

// SwaggerSecurityDefinitions the security definition block
type SwaggerSecurityDefinitions struct {
	Type             string            `json:"type"`
	AuthorizationURL string            `json:"AuthorizationUrl"`
	Flow             string            `json:"flow"`
	Scopes           map[string]string `json:"scopes"`
	Name             string            `json:"name"`
	In               string            `json:"in"`
}

// SwaggerDefinitions the definition block
type SwaggerDefinitions struct {
	Type       string                   `json:"type"`
	Properties map[string]SwaggerFormat `json:"properties"`
}

// SwaggerFormat the definition block
type SwaggerFormat struct {
	Type   string `json:"type"`
	Format string `json:"format"`
}

// SwaggerXML the definition block
type SwaggerXML struct {
	Name string `json:"name"`
}
