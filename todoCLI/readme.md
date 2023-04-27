![image](https://user-images.githubusercontent.com/58289892/234943596-2d3f55d2-51ad-4520-a3e2-d956ebb034d0.png)



Created task manager for terminal with golang, it can be connected to google calendar and show upcoming events. There is many options to play around with. Main idea is to log in to your google calendar through oauth2 and get upcoming events. Then you can add todos and mark them as completed. You can also hide tasks from the list and show them later. You can move tasks to the top or bottom of the list. 

You can also show upcoming events from your calendar. You can also show upcoming events from your calendar for tomorrow, week, month, year or next 1000 events. Y

I'm not planning to publish the app, as it's more of an app made for learning/I might be the only user, but if you want to use it, you can clone the repo and build it yourself.

Before you build app you need file called credentials.json (to connect to google calendar) - follow instruction on this site to get this.

https://developers.google.com/calendar/api/quickstart/go?hl=en

after you get the file, put it next to main.go file and run 
```
go build main.go
```
then you will see the executable link in terminal to log into your google account. 
Link that is provided from google calendar looks like this 
```
http://localhost/?state=state-token&code=4/0AVHEtk5Pw1e4M5MynWG3Vj2asfdasdfasdfasdfSmfRnEhKxszIkDsLw_tGdyhyLygQ&scope=https://www.googleapis.com/auth/calendar.readonly
```
Cope what is between code= and &scope= and paste it into terminal. That's it. The app is go create the token.json which is your token to log into your google account (don't share it with anyone, as well as credentials.json).

Now you can happily use the app.

```
./main -add Walk the dog
./main -add Drink water
./main -show //this one is gonna show you the list of your tasks and events from the calendar 
```

```
Where it says int, you need to put index number of the task you want to perform action on.
Everything else works on it's own. 
    -add text
        add a new todo
    -cal10
        show ten next events from your calendar
    -cal100
        show hundred next events from your calendar
    -cal1000
        show next thousand events
    -calmonth
        show upcoming month
    -caltomorrow
        show events from your calendar for tomorrow
    -calweek
        show upcoming week
    -calyear
        show upcoming year
    -del int
        delete a todo
    -done int
        mark a todo as completed
    -help
        get the list of commands
    -hide int
        hide a task from the list
    -hideall
        hide all done tasks from the list
    -moveend int
        move a task to the bottom of the list
    -movestart int
        move a task to the top of the list
    -show
        show all todos
    -showall
        show all todos, even hidden
    -undone int
        mark a todo as not completed
    -unhide int
        unhide a task from the list
```

Caveats:
- Special signs kinda not working at all. Probably would have to move from json file to db to make it work. 

To make it faster I have done some aliases in my .zshrc file. I do recommend using some aliases, as it makes it much faster to use the app.
```
#Aliases
alias -g ma='/Users/michal/GolandProjects/taskManagerCLI/main -add'
alias -g md='/Users/michal/GolandProjects/taskManagerCLI/main -done'
alias -g m='/Users/michal/GolandProjects/taskManagerCLI/main'
```
Thanks for reading.

