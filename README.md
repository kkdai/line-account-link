LINE account link sample code: A sample code for LINE account link
==============

 [![GoDoc](https://godoc.org/github.com/kkdai/line-account-link.svg?status.svg)](https://godoc.org/github.com/kkdai/line-account-link)  [![Build Status](https://travis-ci.org/kkdai/line-account-link.svg?branch=master)](https://travis-ci.org/kkdai/line-account-link)

[![goreportcard.com](https://goreportcard.com/badge/github.com/kkdai/line-account-link)](https://goreportcard.com/report/github.com/kkdai/line-account-link)

This is sample code to demostration LINE chatbot account link, refer to document https://developers.line.biz/en/docs/messaging-api/linking-accounts/

You could see each step mark in comment (from step 1 to 12)

![](https://developers.line.biz/media/messaging-api/linking-accounts/sequence-f0747c60.png)


1. The bot server calls the API that issues a link token from - the LINE user ID.
2. The LINE Platform returns the link token to the bot server.
3. The bot server calls the Messaging API to send a linking URL to the user.
4. The LINE Platform sends a linking URL to the user.
5. The user accesses the linking URL.
6. The web server displays the login screen.
7. The user enters his/her credentials.
8. The web server acquires the user ID from the provider's service and uses that to generate a nonce.
9. The web server redirects the user to the account-linking endpoint.
10. The user accesses the account-linking endpoint.
11. The LINE Platform sends an event (which includes the LINE user ID and nonce) via webhook to the bot server.
12. The bot server uses the nonce to acquire the user ID of the provider's service.


Installation and Usage
=============

- Create your LINE Developer account and create a developer trial account. 
- Deploy this project to heroku
    - You need to fill following configuration during you create project.
        - LINECORP_PLATFORM_CHANNEL_CHANNELSECRET
        - LINECORP_PLATFORM_CHANNEL_CHANNELTOKEN
        - LINECORP_PLATFORM_CHANNEL_SERVERURL
- 



Deploy on Heroku
=============
[![Deploy](https://www.herokucdn.com/deploy/button.svg)](https://heroku.com/deploy)


