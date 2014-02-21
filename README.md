# Read Raptor

**API for read receipts**

## Setup

Make sure Postgres is installed. Create a local development database:

    psql -c "create database lw_development"

Install goose for running migrations:

    go get bitbucket.org/liamstask/goose/cmd/goose

Copy sample env file:

    cp .env.sample .env

Install forego (or use foreman):

    go get github.com/ddollar/forego

Run migrations:

    forego run goose up

Start up the server (or use something like gin to auto-reload):

    forego run

### Create an Account

    curl -X POST $LW_URL/accounts \
         -d username=whatupdave

    # Response
    {
      "account": {
        "username": "whatupdave",
        "api_key": "api_3c12d9556813"
      }
    }

### Generating tracking urls

Generate the url in the format: `/t/:username/:content_item_id/:user_id/:signature.gif`. The signature is a `sha1` hash of the `api_key` + args in the same order.

    # bash example
    export LW_API_KEY=api_3c12d9556813
    export LW_USERNAME=whatupdave
    export LW_URL=http://localhost:5000

    function _tracking_url {
      echo -n "$LW_URL/t/$LW_USERNAME/$1/$2/`echo -n $LW_API_KEY$LW_USERNAME$1$2 | openssl dgst -sha1`.gif"
    }


### Basic Scenario

Let's say your site gets a new post and you want to notify 3 users. First register the content item:

    curl -X POST $LW_URL/content_items \
         -u $LW_API_KEY: \
         -d '{
           "key": "post_1",
           "unseen": ["user_1", "user_2", "user_3"]
         }'

Response:

    {
      "content_item": {
        "key": "post_1",
        "unseen": ["user_1", "user_2", "user_3"]
      }
    }


Mark first user as seen by requesting the tracking pixel. :

    curl -I `_tracking_url post_1 user_1`

Now get the list of users who have not seen the content:

    curl -u $LW_API_KEY: $LW_URL/content_items/post_1

Response:

    {
      "content_item": {
        "key": "post_1",
        "unseen": ["user_2", "user_3"]
      }
    }

### Content updates

Want to know if a thread is unread or how many comments are new? Sure thang!

Register your thread:

    curl -X POST $LW_URL/content_items \
         -u $LW_API_KEY: \
         -d '{
           "key": "thread_1",
           "unseen": ["user_1", "user_2"]
         }'

Response:

    {
      "content_item": {
        "key": "thread_1",
        "unseen": ["user_1", "user_2"]
      }
    }

User 1 sees the thread:

    curl -I `_tracking_url thread_1 user_1`

Now register an update:

    curl -X POST $LW_URL/content_items \
         -u $LW_API_KEY: \
         -d '{
           "parent": "thread_1",
           "key": "comment_1",
           "unseen": ["user_1", "user_2"]
         }'

Request the content item:

    curl -u $LW_API_KEY: $LW_URL/content_items/thread_1

Response:

    {
      "content_item": {
        "key": "thread_1",
        "unseen": ["user_2"],
        "children": [{
          "key": "comment_1",
          "unseen": ["user_1", "user_2"]
        }]
      }
    }

Now you know that `user_1` has 1 new comment and `user_2` hasn't seen the thread at all.

### User callbacks

Now, let's say you want to email users about content updates. Read Raptor let's you wait a minute and then only email the users that didn't see the content.

Register some content:

    curl -X POST $LW_URL/content \
         -u $LW_API_KEY: \
         -d '{
           "key": "content_1",
           "unseen": ["user_1", "user_2", "user_3"],
           "callbacks": [{
             "seconds": 60,
             "url": "http://requestb.in/u3igzqu3"
           }]
         }'

Response:

    {
      "content_item": {
        "key": "content_1",
        "unseen": ["user_1", "user_2", "user_3"]
      }
    }

One of the users sees the content:

    curl -I `_tracking_url content_1 user_1`

You'll receive a callback per user that hasn't seen some content:

    # First callback
    {
      "callback" {
        "user": "user_1",
        "unseen": ["content_1"]
      }
    }

    # Second callback
    {
      "callback" {
        "user": "user_3",
        "unseen": ["content_1"]
      }
    }

### Example: Push Notification + Daily Digest

What if you want to send a push notification and also include the update in a daily digest? No problemo! You can register a callback for 1 minute, then 24 hours.

    curl -X POST $LW_URL/content \
         -u $LW_API_KEY: \
         -d '{
           "key": "content_1",
           "unseen": ["user_1"],
           "callbacks": [{
             "seconds": 60,
             "url": "http://requestb.in/u3igzqu3"
           }, {
             "hours": 24,
             "url": "http://requestb.in/u3igzqu3"
           }]
         }'

Response:

    {
      "content_item": {
        "key": "content_1",
        "unseen": ["user_1"]
      }
    }

Read Raptor will send you both notifications unless the user sees that piece of content at some point along the way.

### Running the tests

Make sure the test database exists:

    psql -c "create database lw_test"

Run the tests:

    go test

## Contributing

1. Sign up at Assembly ([https://assemblymade.com](https://assemblymade.com))
2. Create a Task for the work ([https://assemblymade.com/readraptor](https://assemblymade.com/readraptor))
3. Create your feature branch (`git checkout -b my-new-feature`)
4. Commit your changes, you can reference the task number (`git commit -am 'Add some feature for #123'`)
5. Push to the branch (`git push origin my-new-feature`)
6. Create new Pull Request