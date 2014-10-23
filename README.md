# Parfait backend

## Running

`AWS_SECRET_ACCESS_KEY=foo AWS_ACCESS_KEY_ID=bar CID_DEVELOPER_KEY=baz CID_DEVELOPER_SECRET=qux ./parfait`

Now you should be able to curl against it

`curl --verbose -X POST localhost:3000/users/<user_id>/email_accounts/<label>`


## Start scanning a new account (implemented)

This assumes the account has already been connected to Context.IO.  See http://context.io/docs/lite/users/email_accounts for the field definitions of `:id` and `:label`

`POST /users/:id/email_accounts/:label`

* 200 - Account added and being scanned
* 404 - Account not found for the CIO account
* 500 - Shit be crazy


## Get unread messages

Returns up to 25 results

`GET /users/:id/email_accounts/:label/unread_messages`

* 200 - Found data (see JSON schema)
* 400 - Account isn't being scanned
* 500 - Sit be crazy
 

```javascript
{
   "id": "account id",
   "label": "account label",
   "unread_messages": [
      {
         "message_id": "<message-id@example.com>", 
         "subject": "hello there!",
         "sender": {
            "address": "sender@example.com",
            "total_count": 100,
            "unread_count": 20,
            "answered_count": 4,
            "flagged_count": 1,
            "draft_count": 0
         },
      }
   ]
}
```

## Get sender data (Maybe, if there's time?)

`GET /users/:id/email_accounts/:label/senders`

* 200 - Found data (see JSON schema)
* 404 - Account isn't being scanned
* 500 - Shit be crazy

```javascript
{
   "id": "account id",
   "label": "account label",
   "senders": [
      {
         "address": "sender@example.com",
         "read_count": 100,
         "unread_count": 20,
         "forwarded_count": 4,
         "flagged_count": 1
      }
   ]
}
```
