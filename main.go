package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"os"
	"sync"

	"github.com/chromedp/chromedp"
)

// set colors
var (
	Red   = Color("\033[1;31m%s\033[0m")
	Green = Color("\033[1;32m%s\033[0m")
)

func Color(colorString string) func(...interface{}) string {
	sprint := func(args ...interface{}) string {
		return fmt.Sprintf(colorString,
			fmt.Sprint(args...))
	}
	return sprint
}

func main() {

	// concurrency flag
	var concurrency int
	flag.IntVar(&concurrency, "c", 20, "")

	// javascript flag
	var userJS string
	flag.StringVar(&userJS, "js", "", "")
	flag.StringVar(&userJS, "javascript", "", "")

	flag.Parse()

	scanner := bufio.NewScanner(os.Stdin)

	var wg sync.WaitGroup

	urls := make(chan string)
	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func() {
			for url := range urls {
				// create context
				ctx, cancel := chromedp.NewContext(context.Background())

				// run task list
				var res string
				err := chromedp.Run(ctx,
					chromedp.Navigate(url),
					chromedp.Evaluate(userJS, &res),
				)

				cancel()

				if err != nil {
					fmt.Printf("%s %s\n", Red("[ERRO]"), url)
					continue
				}

				fmt.Printf("%s %s\n", Green("[VULN]"), url)
			}
			wg.Done()
		}()
	}

	for scanner.Scan() {
		u := scanner.Text()
		urls <- u
	}

	close(urls)
	wg.Wait()
}

func init() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Fuzz for parameter pollution\n\n")
		fmt.Fprintf(os.Stderr, "Usage:\n")
		fmt.Fprintf(os.Stderr, "cat urls | jsfuzz [options]\n\n")
		fmt.Fprintf(os.Stderr, "Options:\n")
		fmt.Fprintf(os.Stderr, "  -c <int>        			set the concurrency level (default 20)\n")
		fmt.Fprintf(os.Stderr, "  -js, --javascript <str>		the JS to run on each page\n")
	}
}
