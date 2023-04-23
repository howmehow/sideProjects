Chatbot integrated in Slack app, using WIT.AI and Wolfram Search Query API. 

You can run it with go run main.go after adding your own environmental variables in .env file:

SLACK_BOT_TOKEN=""

SLACK_APP_TOKEN=""

WIT_AI_TOKEN=""

WOLFRAM_APP_ID=""


You will have to create app in https://api.slack.com/apps as well as  in https://developer.wolframalpha.com/portal/myapps/ and https://wit.ai/apps/ 

In slack app event subscription have to be created with Socket Mode enabled.
Oauth & Permissions have to be created with following scopes:
- app_mentions:read
- chat:write
- im:read
- im:write
- mpim:history
- mpim:read
- mpim:write
- users:read

You can find all required info in slack documentation. If for example you would like your bot to work in a channel you need to add following scope:
- channels:read
- channels:write
If there is any errors in the terminal, it might be with the scopes.