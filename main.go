package main

import (
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
)

func rainbow(w http.ResponseWriter, r *http.Request) {

	fmt.Println("method:", r.Method) //get request method

	// get RAW POST parameters
	rawBody, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(string(rawBody))
	}

	// use strings.split to split the param string on "&"
	params := strings.Split(string(rawBody), "&")
	fmt.Println(params)

	// create an empty map
	m := make(map[string]string)

	// now split each of the array element on "="
	for _, param := range params {
		if param != "" {
			pair := strings.Split(param, "=")
			m[pair[0]], err = url.QueryUnescape(pair[1])
			if err != nil {
				fmt.Println(err)
			}
		}
	}

	fmt.Println(m)

	coloredMessage := colorize(m["submit"])
	fmt.Println(coloredMessage)

	fmt.Fprintf(w, m["cmd"]+" "+coloredMessage)

}

func main() {
	http.HandleFunc("/", rainbow)
	http.ListenAndServe(":8080", nil)

	/*
		buf := bufio.NewReader(os.Stdin)
		fmt.Print("> ")
		sentence, err := buf.ReadBytes('\n')
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(colorize(string(sentence)))
		}
	*/
}

/* function that colorizes the chat in Xonotic */
func colorize(text string) string {
	rainbow := [...]string{"^xF87", "^xF97", "^xFA7", "^xFB7", "^xFC7", "^xFD7", "^xFE7", "^xFF7", "^xEF7", "^xDF7", "^xCF7", "^xBF7",
		"^xAF7", "^x9F7", "^x8F7", "^x7F7", "^x7F8", "^x7F9", "^x7FA", "^x7FB", "^x7FC", "^x7FD", "^x7FE", "^x7FF",
		"^x7EF", "^x7DF", "^x7CF", "^x7BF", "^x7AF", "^x79F", "^x78F", "^x77F", "^x87F", "^x97F", "^xA7F", "^xB7F",
		"^xC7F", "^xD7F", "^xE7F", "^xF7F", "^xF7E", "^xF7D", "^xF7C", "^xF7B", "^xF7A", "^xF79", "^xF78", "^xF77"}
	rainbowPos := rand.Int() % len(rainbow)
	colorLength := 0
	color := ""
	escape := ""
	result := ""
	for _, rune := range text {
		// ';' is a special character in the URL, so we need to escape it
		if rune == ';' {
			escape = ";"
			result += rainbow[rainbowPos] + escape
			rainbowPos = (rainbowPos + 1) % len(rainbow)
		} else if rune == '%' {
			escape = "%%"
			result += rainbow[rainbowPos] + escape
			rainbowPos = (rainbowPos + 1) % len(rainbow)
		} else if rune == '$' {
			escape = "$$"
			result += rainbow[rainbowPos] + escape
			rainbowPos = (rainbowPos + 1) % len(rainbow)
		} else {
			if colorLength == 0 {
				if rune == '^' {
					colorLength = 1
					color = "^"
				} else {
					result += rainbow[rainbowPos] + string(rune)
					rainbowPos = (rainbowPos + 1) % len(rainbow)
					//result += rainbow[rainbowPos]
				}
			} else {
				if colorLength == 1 {
					if rune == 'x' {
						colorLength = 5
						color += string(rune)
					} else if (rune >= '0' && rune <= '9') || (rune >= 'a' && rune <= 'f') || (rune >= 'A' && rune <= 'F') {
						colorLength = 0
						color = ""
					} else {
						colorLength = 0
						color = ""
						result += rainbow[rainbowPos] + string('^')
						rainbowPos = (rainbowPos + 1) % len(rainbow)
						result += rainbow[rainbowPos] + string(rune)
						rainbowPos = (rainbowPos + 1) % len(rainbow)

					}
				} else {
					color += string(rune)
					if len(color) == 5 {
						colorLength = 0
						color = ""
					}
				}
			}
		}
	}

	return result
}