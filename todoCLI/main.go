package main

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
	todo "todoCLI/cmd"
)

const (
	todoFile = ".todos.json"
)

func main() {
	add := flag.Bool("add", false, "add a new todo")
	done := flag.Int("done", 0, "mark a todo as completed")
	undone := flag.Int("undone", 0, "mark a todo as not completed")
	del := flag.Int("del", 0, "delete a todo")
	show := flag.Bool("show", false, "show all todos")
	help := flag.Bool("help", false, "get the list of commands")
	hide := flag.Int("hide", 0, "hide a task from the list")
	unhide := flag.Int("unhide", 0, "unhide a task from the list")
	showall := flag.Bool("showall", false, "show all todos, even hidden")
	cal10 := flag.Bool("cal10", false, "show ten next events from your calendar")
	cal100 := flag.Bool("cal100", false, "show hundred next events from your calendar")
	caltomorrow := flag.Bool("caltomorrow", false, "show events from your calendar for tomorrow")
	cal1000 := flag.Bool("cal1000", false, "show next thousand events")
	calweek := flag.Bool("calweek", false, "show upcoming week")
	calmonth := flag.Bool("calmonth", false, "show upcoming month")
	calyear := flag.Bool("calyear", false, "show upcoming year")
	flag.Parse()

	todos := &todo.Todos{}
	if err := todos.Load(todoFile); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
	todos.CheckForDoneItemsYesterday(todoFile)
	switch {
	case *cal10:
		events, _ := getEvents("10")
		todos.Print(false, events)
	case *cal100:
		events, _ := getEvents("100")
		todos.Print(false, events)
	case *cal1000:
		events, _ := getEvents("1000")
		todos.Print(false, events)
	case *caltomorrow:
		events, _ := getEvents("tomorrow")
		todos.Print(false, events)
	case *calweek:
		events, _ := getEvents("week")
		todos.Print(false, events)
	case *calmonth:
		events, _ := getEvents("month")
		todos.Print(false, events)
	case *calyear:
		events, _ := getEvents("year")
		todos.Print(false, events)
	case *hide > 0:
		err := todos.Hide(*hide)
		err = todos.Store(todoFile)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(1)
		}
		events, _ := getEvents("")
		todos.Print(false, events)
	case *unhide > 0:
		err := todos.Unhide(*unhide)
		err = todos.Store(todoFile)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(1)
		}
		events, _ := getEvents("")
		todos.Print(false, events)
	case *help:
		fmt.Println("Hello there here is quick reference to all the commands that you can use in the CLI. The commands without -cal are working on the tasks, -cal is referencing to the calendar API, to show different time frame for events from calendar, or amount of them.")
		fmt.Println("If you want to do something with the task (hide, unhide, delete, mark as done, or undone) then you have to provide the index from the table.")
		fmt.Println("The links next to google calendar event are provided by Google API, and they direct you to this event in the browser.")
		fmt.Println("All of the tasks you can find in the json file .todos.json")
		fmt.Println("If you want to inc")
		fmt.Println("-add task name")
		fmt.Println("-done index")
		fmt.Println("-undone index")
		fmt.Println("-del index")
		fmt.Println("-show ")
		fmt.Println("-help ")
		fmt.Println("-hide index")
		fmt.Println("-unhide index")
		fmt.Println("-showall ")
		fmt.Println("-cal10 ")
		fmt.Println("-cal100 ")
		fmt.Println("-caltomorrow ")
		fmt.Println("-cal1000 ")
		fmt.Println("-calweek ")
		fmt.Println("-calmonth ")
		fmt.Println("-calyear ")
	case *add:
		task, err := getInput(os.Stdin, flag.Args()...)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(1)
		}
		todos.Add(task)
		err = todos.Store(todoFile)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(1)
		}
		events, _ := getEvents("")
		todos.Print(false, events)
	case *done > 0:
		err := todos.Complete(*done)
		err = todos.Store(todoFile)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(1)
		}
		events, _ := getEvents("")
		todos.Print(false, events)
	case *undone > 0:
		err := todos.Uncomplete(*undone)
		err = todos.Store(todoFile)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(1)
		}
		events, _ := getEvents("")
		todos.Print(false, events)
	case *del > 0:
		err := todos.Delete(*del)
		err = todos.Store(todoFile)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(1)
		}
		events, _ := getEvents("")
		todos.Print(false, events)
	case *show:
		events, _ := getEvents("")
		todos.Print(false, events)
	case *showall:
		events, _ := getEvents("")
		todos.Print(true, events)
	default:
		events, _ := getEvents("")
		todos.Print(false, events)
	}
}
func getInput(r io.Reader, args ...string) (string, error) {
	if len(args) > 0 {
		return strings.Join(args, " "), nil
	}
	scanner := bufio.NewScanner(r)
	scanner.Scan()
	if err := scanner.Err(); err != nil {
		return "", err
	}
	text := scanner.Text()
	if len(text) == 0 {
		return "", errors.New("empty input")
	}
	return text, nil
}

func getEvents(show string) (*calendar.Events, error) {
	ctx := context.Background()
	b, err := os.ReadFile("credentials.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}
	config, err := google.ConfigFromJSON(b, calendar.CalendarReadonlyScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	client := getClient(config)

	srv, err := calendar.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Fatalf("Unable to retrieve Calendar client: %v", err)
	}

	t := time.Now().Format(time.RFC3339)
	var events *calendar.Events
	fmt.Println(show)
	switch show {
	case "10":
		events, err = srv.Events.List("primary").ShowDeleted(false).
			SingleEvents(true).TimeMin(t).MaxResults(10).OrderBy("startTime").Do()
	case "100":
		events, err = srv.Events.List("primary").ShowDeleted(false).
			SingleEvents(true).TimeMin(t).MaxResults(100).OrderBy("startTime").Do()
	case "1000":
		events, err = srv.Events.List("primary").ShowDeleted(false).
			SingleEvents(true).TimeMin(t).MaxResults(1000).OrderBy("startTime").Do()
	case "tomorrow":
		t = time.Now().AddDate(0, 0, 1).Format(time.RFC3339)
		events, err = srv.Events.List("primary").ShowDeleted(false).
			SingleEvents(true).TimeMin(t[0:10] + "T00:00:00Z").TimeMax(t[0:10] + "T23:59:59Z").OrderBy("startTime").Do()
	case "week":
		te := time.Now().AddDate(0, 0, 7).Format(time.RFC3339)
		events, err = srv.Events.List("primary").ShowDeleted(false).
			SingleEvents(true).TimeMin(t[0:10] + "T00:00:00Z").TimeMax(te[0:10] + "T23:59:59Z").OrderBy("startTime").Do()
	case "month":
		te := time.Now().AddDate(0, 1, 0).Format(time.RFC3339)
		events, err = srv.Events.List("primary").ShowDeleted(false).
			SingleEvents(true).TimeMin(t[0:10] + "T00:00:00Z").TimeMax(te[0:10] + "T23:59:59Z").OrderBy("startTime").Do()
	case "year":
		t = time.Now().AddDate(0, 0, 1).Format(time.RFC3339)
		te := time.Now().AddDate(1, 0, 0).Format(time.RFC3339)
		events, err = srv.Events.List("primary").ShowDeleted(false).
			SingleEvents(true).TimeMin(t[0:10] + "T00:00:00Z").TimeMax(te[0:10] + "T23:59:59Z").OrderBy("startTime").Do()
	default:
		events, err = srv.Events.List("primary").ShowDeleted(false).
			SingleEvents(true).TimeMin(t[0:10] + "T00:00:00Z").TimeMax(t[0:10] + "T23:59:59Z").OrderBy("startTime").Do()
	}
	if err != nil {
		log.Fatalf("Unable to retrieve next ten of the user's events: %v", err)
	}
	fmt.Println(events.Items)
	return events, err
}

// Retrieve a token, saves the token, then returns the generated client.
func getClient(config *oauth2.Config) *http.Client {
	// The file token.json stores the user's access and refresh tokens, and is
	// created automatically when the authorization flow completes for the first
	// time.
	tokFile := "token.json"
	tok, err := tokenFromFile(tokFile)
	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken(tokFile, tok)
	}
	return config.Client(context.Background(), tok)
}

// Request a token from the web, then returns the retrieved token.
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		log.Fatalf("Unable to read authorization code: %v", err)
	}

	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web: %v", err)
	}
	return tok
}

// Retrieves a token from a local file.
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

// Saves a token to a file path.
func saveToken(path string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	json.NewEncoder(f).Encode(token)
}
