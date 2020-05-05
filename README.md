## Boyfriend Management Console
This is a Silly Valentine's day app.  User (i.e. the partner you are trying to please) can add tasks with messages.
When a user asks for a task to be completed, server calls out to twilio api to send a text message to the boyfriend/girlfriend who is being managed by this webapp.

Is this simpler than just sending a text message? probably not, but it only took a couple of hours to make and was a nice novelty gift.


Magic build instructions for deployment 
```
GOOS=linux go build -ldflags="-s -w" ./
cp ashleyButtons ./dist
scp -r dist/ ec2-user@bf-management-console.treestack.io:~/
```

to run your service just start it
```
export TWILIO_ACCOUNT_ID=$ACCOUNT_ID
export TWILIO_AUTH_TOKEN=$AUTH_TOKEN
export BFMC_RECIEVER_NUMBER=$SOME_PHONE_NUMBER
export BFMC_FROM_NUMBER=$SOME_OTHER_NUMBER

./ashleyButtons
```
service runs by default on 8080.  One improvement to add would be to set the service port from an envvar or cli flag. 
