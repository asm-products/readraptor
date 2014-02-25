# Read Raptor

Read Raptor intelligently notifies your users about new content.

## Why would I use this on my app?

Ever noticed how iMessage doesn't message your phone if you've read it on your desktop? Read Raptor can make that happen for your app too.

Read Raptor keeps a register of articles in your app and tracks who has read them. You can tell her to send you webhooks at specified times which will let you know what each user hasn't read.

Here's 3 fun things she can help you with:

1. Show which articles are unread for each user
2. Email your users about new content only if they **haven't** seen it
3. Send out digest emails with grouped content that each user **hasn't** read

## Get Started

### Create an Account

    # Set some environment variables
    export RR_URL=http://localhost:5000

    curl -X POST $RR_URL/accounts \
         -d username=whatupdave

    # Response
    {
      "account": {
        "username": "whatupdave",
        "api_key": "api_3c12d9556813"
      }
    }

    # some more env vars
    export RR_API_KEY=api_3c12d9556813
    export RR_USERNAME=whatupdave


### Generating tracking urls

Tracking pixel urls are in the format: `/t/:username/:article_id/:user_id/:signature.gif`. The signature is a `sha1` hash of the `api_key` + args in the same order.

**Bash example**

    function _tracking_url {
      echo -n "$RR_URL/t/$RR_USERNAME/$1/$2/`echo -n $RR_API_KEY$RR_USERNAME$1$2 | openssl dgst -sha1`.gif"
    }


**Ruby example**

    require 'digest/sha1'

    def tracking_url(article_id, user_id)
      sig = Digest::SHA1.hexdigest "#{ENV['RR_API_KEY']}#{ENV['RR_USERNAME']}#{article_id}#{user_id}"
      "#{ENV['RR_URL']}/t/#{ENV['RR_USERNAME']}/#{article_id}/#{user_id}/#{sig}.gif"
    end


### Example 1: Show unread articles

Let's say your site gets a new post and you want to notify 3 users about it. First register the content item:

    curl -X POST $RR_URL/articles \
         -u $RR_API_KEY: \
         -d '{
           "key": "post_1",
           "pending": ["user_1", "user_2", "user_3"]
         }'

Response:

    {
      "article": {
        "key": "post_1",
        "pending": ["user_1", "user_2", "user_3"]
      }
    }


Mark first user as seen by requesting the tracking pixel. :

    curl -I `_tracking_url post_1 user_1`

Now get the list of users who have not seen the content:

    curl -u $RR_API_KEY: $RR_URL/articles/post_1

Response:

    {
      "article": {
        "key": "post_1",
        "pending": ["user_2", "user_3"]
      }
    }

Go ahead and display that content as unread.

### Example 2: Email users who haven't seen an article

You want to email users about a new article, but you don't want to email them if they've already seen it. Read Raptor let's you register callbacks that will notify you about content that users haven't seen. The delay argument accepts strings such as "2h45m" or "1m", see http://golang.org/pkg/time/#ParseDuration for more.

Register some content:

    curl -X POST $RR_URL/articles \
         -u $RR_API_KEY: \
         -d '{
           "key": "article_1",
           "pending": ["user_1", "user_2", "user_3"],
           "callbacks": [{
             "delay": "60s",
             "url": "http://requestb.in/u3igzqu3"
           }]
         }'

Response:

    {
      "article": {
        "key": "article_1",
        "pending": ["user_1", "user_2", "user_3"]
      }
    }

One of the users sees the content:

    curl -I `_tracking_url article_1 user_1`

You'll receive a callback per user that hasn't seen some content:

    # First callback
    {
      "callback" {
        "user": "user_1",
        "pending": ["article_1"]
      }
    }

    # Second callback
    {
      "callback" {
        "user": "user_3",
        "pending": ["article_1"]
      }
    }

### Example 3: Digest emails

What if some users want immediate emails and others want daily digests? No problemo! You can register a callback for 1 minute, then 24 hours and specify the recipients in each callback.

    curl -X POST $RR_URL/articles \
         -u $RR_API_KEY: \
         -d '{
           "key": "article_1",
           "callbacks": [{
             "delay": "1m",
             "recipients": ["immediate_user"],
             "url": "http://requestb.in/u3igzqu3"
           }, {
             "delay": "24h",
             "recipients": ["digest_user"],
             "url": "http://requestb.in/u3igzqu3"
           }]
         }'

Response:

    {
      "article": {
        "key": "article_1",
        "pending": ["immediate_user", "digest_user"]
      }
    }

Read Raptor will send you both webhooks unless the users see that piece of content at some point along the way.

## Local Setup

Make sure Postgres is installed. Create a local development database:

    psql -c "create database rr_development"

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


### Running the tests

Make sure the test database exists:

    psql -c "create database rr_test"

Run the tests:

    go test ./...

## Contributing

1. Sign up at Assembly ([https://assemblymade.com](https://assemblymade.com))
2. Create a Task for the work ([https://assemblymade.com/readraptor](https://assemblymade.com/readraptor))
3. Create your feature branch (`git checkout -b my-new-feature`)
4. Commit your changes, you can reference the task number (`git commit -am 'Add some feature for #123'`)
5. Push to the branch (`git push origin my-new-feature`)
6. Create new Pull Request