package framework

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"
)

type HandlerFunc func(*Context)

type TemplateData struct {
	Data    map[string]interface{}
	Helpers map[string]HelperFunc
}

type H map[string]interface{}

type HelperFunc func(string) string

type Context struct {
	Writer  http.ResponseWriter
	Request *http.Request

	Params       map[string]string
	Body         map[string]interface{}
	Router       *Router
	TemplateData TemplateData
	Data         H
	SessionData  map[string]interface{}
}

func NewContext(w http.ResponseWriter, r *http.Request, router *Router) *Context {
	return &Context{
		Writer:       w,
		Request:      r,
		Params:       make(map[string]string),
		Router:       router,
		TemplateData: TemplateData{},
		Data:         H{},
	}
}

func (c *Context) SetCookie(name string, value string, maxAge int) {
	http.SetCookie(c.Writer, &http.Cookie{
		Name:   name,
		Value:  value,
		MaxAge: maxAge,
	})
}

func (c *Context) GetCookie(name string) (string, error) {
	cookie, err := c.Request.Cookie(name)
	if err != nil {
		return "", err
	}
	return cookie.Value, nil
}

func (c *Context) GetQuery(name string) (string, error) {
	query := c.Request.URL.Query().Get(name)
	if query == "" {
		error_msg := fmt.Sprintf("Query \"%s\" not found", name)
		return "", fmt.Errorf(error_msg)
	}
	return query, nil
}

func (c *Context) ParseJson() error {
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(body, &c.Body); err != nil {
		return err
	}
	return nil
}

func (c *Context) JSON(code int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		http.Error(c.Writer, err.Error(), http.StatusInternalServerError)
		return
	}
	c.Writer.Header().Set("Content-Type", "application/json")
	c.Writer.WriteHeader(code)
	c.Writer.Write(response)
}

func (c *Context) RenderTemplateFile(filename string) {
	bytes, err := os.ReadFile(c.Router.TemplateDir + "/" + filename)
	if err != nil {
		http.Error(c.Writer, "Error reading template file", http.StatusInternalServerError)
		return
	}

	templateString := string(bytes)
	rendered := c.RenderTemplate(templateString)
	c.Writer.Write([]byte(rendered))
}

func (c *Context) RenderTemplate(templateString string) string {

	for key, value := range c.TemplateData.Data {
		switch v := value.(type) {
		case []string:
			// Handle {{#each key}}...{{/each}} constructs
			startTag := fmt.Sprintf("{{#each %s}}", key)
			endTag := "{{/each}}"
			startIdx := strings.Index(templateString, startTag)
			endIdx := strings.Index(templateString, endTag)
			if startIdx != -1 && endIdx != -1 {
				// Extract the template within the each construct
				inside := templateString[startIdx+len(startTag) : endIdx]
				result := ""
				for _, item := range v {
					result += strings.ReplaceAll(inside, "{{this}}", item)
				}
				// Replace the each construct with the result
				templateString = strings.ReplaceAll(templateString, startTag+inside+endTag, result)
			}
		case bool:
			// Handle {{#if key}}...{{/if}} constructs
			startTag := fmt.Sprintf("{{#if %s}}", key)
			endTag := "{{/if}}"
			startIdx := strings.Index(templateString, startTag)
			endIdx := strings.Index(templateString, endTag)
			if startIdx != -1 && endIdx != -1 {
				// Extract the template within the if construct
				inside := templateString[startIdx+len(startTag) : endIdx]
				result := ""
				if v {
					result = inside
				}
				// Replace the if construct with the result
				templateString = strings.ReplaceAll(templateString, startTag+inside+endTag, result)
			}
		case string:
			// Replace {{key}} with value
			templateString = strings.ReplaceAll(templateString, "{{"+key+"}}", v)
		case int, int32, int64, float32, float64:
			// Replace {{key}} with the string representation of the number
			templateString = strings.ReplaceAll(templateString, "{{"+key+"}}", fmt.Sprintf("%v", v))

		}
	}

	for key, helper := range c.TemplateData.Helpers {
		// Find all occurrences of {{helper key}} in the template
		r := regexp.MustCompile("{{" + key + ` (\w+)}}`)
		matches := r.FindAllStringSubmatch(templateString, -1)

		// For each occurrence, apply the helper function
		for _, match := range matches {
			if len(match) == 2 {
				valueKey := match[1]
				value, exists := c.TemplateData.Data[valueKey]
				if exists {
					// Apply the helper function and replace in the template
					strValue, ok := value.(string)
					if ok {
						templateString = strings.ReplaceAll(templateString, match[0], helper(strValue))
					}
				}
			}
		}
	}
	return templateString
}
