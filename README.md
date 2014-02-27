# Read Raptor

Read Raptor intelligently notifies your users about new content.

## Why would I use this on my app?

Ever noticed how iMessage doesn't message your phone if you've read it on your desktop? Read Raptor can make that happen for your app too.

Read Raptor keeps a register of articles in your app and tracks who has read them. You can tell her to send you webhooks at specified times which will let you know what each user hasn't read.

Here's 3 fun things she can help you with:

1. Show which articles are unread for each user
2. Email your users about new content only if they **haven't** seen it
3. Send out digest emails with grouped content that each user **hasn't** read

## Installation

### Local Setup

Make sure Postgres is installed.

    # Create a local development database
    psql -c "create database rr_development"

    # Install goose for running migrations
    go get bitbucket.org/liamstask/goose/cmd/goose

    # Copy sample env file
    cp .env.sample .env

    # Install forego (or use foreman)
    go get github.com/ddollar/forego

    # Run database migrations
    forego run goose up

    # Start up the server (or use something like gin to auto-reload)
    go get ./... && forego run


### Editing html/css

We're using `compass` which is a ruby gem to compile sass into css. Make sure ruby is installed then

    gem install compass sass susy
    compass watch

### Running the tests

Make sure the test database exists:

    psql -c "create database rr_test"

Run the tests:

    go test ./...


## Get Started

### Create an Account

    # Set some environment variables
    export RR_URL=http://localhost:5000

    curl -X POST $RR_URL/accounts \
         -d email=whatupdave@example.com

    # Response
    {
      "account": {
        "id":1,
        "created":"2014-02-24T18:36:10.653736897-08:00",
        "email":"whatupdave@example.com",
        "publicKey":"bf55ea13d87bc9b79e28974352cbc0fd0caca1d9",
        "privateKey":"7519a5328f4741d72a6895ff7a4ea9a446e36b17"
      }
    }

    # some more env vars
    export RR_PUBLIC_KEY=bf55ea13d87bc9b79e28974352cbc0fd0caca1d9
    export RR_PRIVATE_KEY=7519a5328f4741d72a6895ff7a4ea9a446e36b17


### Generating tracking urls

Tracking pixel urls are in the format: `/t/:username/:article_id/:user_id/:signature.gif`. The signature is a `sha1` hash of the `private_key` + args in the same order.

**Bash example**

    function _tracking_url {
      echo -n "$RR_URL/t/$RR_PUBLIC_KEY/$1/$2/`echo -n $RR_PRIVATE_KEY$RR_PUBLIC_KEY$1$2 | openssl dgst -sha1`.gif"
    }


**Ruby example**

    require 'digest/sha1'

    def tracking_url(article_id, user_id)
      sig = Digest::SHA1.hexdigest "#{ENV['RR_PRIVATE_KEY']}#{ENV['RR_PUBLIC_KEY']}#{article_id}#{user_id}"
      "#{ENV['RR_URL']}/t/#{ENV['RR_PUBLIC_KEY']}/#{article_id}/#{user_id}/#{sig}.gif"
    end


### Example 1: Show unread articles

Let's say your site gets a new post and you want to notify 3 users about it. First register the content item:

    curl -X POST $RR_URL/articles \
         -u $RR_PRIVATE_KEY: \
         -d '{
           "key": "post_1",
           "recipients": ["user_1", "user_2", "user_3"]
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

    curl -u $RR_PRIVATE_KEY: $RR_URL/articles/post_1

Response:

    {
      "article": {
        "key": "post_1",
        "pending": ["user_2", "user_3"]
      }
    }

Go ahead and display that content as unread.

### Example 2: Email users who haven't seen an article

You want to email users about a new article, but you don't want to email them if they've already seen it. Read Raptor let's you register callbacks that will notify you about content that users haven't seen.

Register some content:

    curl -X POST $RR_URL/articles \
         -u $RR_PRIVATE_KEY: \
         -d '{
           "key": "article_1",
           "recipients": ["user_1", "user_2", "user_3"],
           "via": [{
             "type": "webhook",
             "at": '"`date -v+1M +%s`"',
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
         -u $RR_PRIVATE_KEY: \
         -d '{
           "key": "article_1",
           "via": [{
             "type": "webhook",
             "at": '"`date -v+1M +%s`"',
             "recipients": ["immediate_user"],
             "url": "http://requestb.in/u3igzqu3"
           }, {
             "type": "webhook",
             "at": '"`date -v+24H +%s`"',
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

## Contributing

1. Sign up at Assembly ([https://assemblymade.com](https://assemblymade.com))
2. Create a Task for the work ([https://assemblymade.com/readraptor](https://assemblymade.com/readraptor))
3. Create your feature branch (`git checkout -b my-new-feature`)
4. Commit your changes, you can reference the task number (`git commit -am 'Add some feature for #123'`)
5. Push to the branch (`git push origin my-new-feature`)
6. Create new Pull Request
