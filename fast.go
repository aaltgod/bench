package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
)


func FastSearch(out io.Writer) {

	file, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}

	r := regexp.MustCompile("@")
	var seenBrowsers []string
	uniqueBrowsers := 0
	foundUsers := ""
	scanner := bufio.NewScanner(file)
	var users []User

	for scanner.Scan() {
		user := User{}
		err := user.UnmarshalJSON(scanner.Bytes())
		if err != nil {
			panic(err)
		}
		users = append(users, user)

	}

	for i, user := range users {
			isAndroid := false
			isMSIE := false
			browsers := user.Browsers
			for _, browserRaw := range browsers {
				if strings.Contains(browserRaw, "Android") {
					isAndroid = true
					notSeenBefore := true

					for _, item := range seenBrowsers {
						if item == browserRaw {
							notSeenBefore = false
						}
					}
					if notSeenBefore {
						seenBrowsers = append(seenBrowsers, browserRaw)
						uniqueBrowsers++
					}
				}
			}

			for _, browserRaw := range browsers {
				if strings.Contains(browserRaw, "MSIE") {
					isMSIE = true
					notSeenBefore := true

					for _, item := range seenBrowsers {
						if item == browserRaw {
							notSeenBefore = false
						}
					}
					if notSeenBefore {
						seenBrowsers = append(seenBrowsers, browserRaw)
						uniqueBrowsers++
					}
				}
			}

			if !(isAndroid && isMSIE) {
				continue
			}

			email := r.ReplaceAllString(user.Email, " [at] ")
			foundUsers += fmt.Sprintf("[%d] %s <%s>\n", i, user.Name, email)
		}

		fmt.Fprintln(out, "found users:\n"+foundUsers)
		fmt.Fprintln(out, "Total unique browsers", len(seenBrowsers))
	}
