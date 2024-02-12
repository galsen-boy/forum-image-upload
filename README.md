<p align="center">
  <a href="" rel="noopener">
 <img width=200px height=200px src="./static/images/logoForum.png" alt="Project logo"></a>
</p>

<h3 align="center">FORUM</h3>

<p align="center"> This project consists in creating a web forum that allows
    <br> 
</p>
<p>communication between users.</p>
<p>associating categories to posts.</p>
<p>liking and disliking posts and comments.</p>
<p>filtering posts.</p>

## 📝 Table of Contents

- [About](#about)
- [Usage](#usage)
- [Description](https://learn.zone01dakar.sn/git/root/public/src/branch/master/subjects/forum)
- [Authors](#authors)
- [Implementation details](#built_using)

## 🧐 About <a name = "about"></a>

SQLite

In order to store the data in your forum (like users, posts, comments, etc.) you will use the database library SQLite.

SQLite is a popular choice as an embedded database software for local/client storage in application software such as web browsers. It enables you to create a database as well as controlling it by using queries.

To structure your database and to achieve better performance we highly advise you to take a look at the entity relationship diagram and build one based on your own database.

    You must use at least one SELECT, one CREATE and one INSERT queries.

To know more about SQLite you can check the SQLite page.
Authentication

In this segment the client must be able to register as a new user on the forum, by inputting their credentials. You also have to create a login session to access the forum and be able to add posts and comments.

You should use cookies to allow each user to have only one opened session. Each of this sessions must contain an expiration date. It is up to you to decide how long the cookie stays "alive". The use of UUID is a Bonus task.

Instructions for user registration:

    Must ask for email
        When the email is already taken return an error response.
    Must ask for username
    Must ask for password
        The password must be encrypted when stored (this is a Bonus task)

The forum must be able to check if the email provided is present in the database and if all credentials are correct. It will check if the password is the same with the one provided and, if the password is not the same, it will return an error response.
Communication

In order for users to communicate between each other, they will have to be able to create posts and comments.

    Only registered users will be able to create posts and comments.
    When registered users are creating a post they can associate one or more categories to it.
        The implementation and choice of the categories is up to you.
    The posts and comments should be visible to all users (registered or not).
    Non-registered users will only be able to see posts and comments.

Likes and Dislikes

Only registered users will be able to like or dislike posts and comments.

The number of likes and dislikes should be visible by all users (registered or not).
Filter

You need to implement a filter mechanism, that will allow users to filter the displayed posts by :

    categories
    created posts
    liked posts

You can look at filtering by categories as subforums. A subforum is a section of an online forum dedicated to a specific topic.

Note that the last two are only available for registered users and must refer to the logged in user.
Docker

For the forum project you must use Docker. You can read about docker basics in the ascii-art-web-dockerize subject.

### Description

Instructions

    You must use SQLite.
    You must handle website errors, HTTP status.
    You must handle all sort of technical errors.
    The code must respect the good practices.
    It is recommended to have test files for unit testing.

Allowed packages

    All standard Go packages are allowed.
    sqlite3
    bcrypt
    UUID

This project will help you learn about:

    The basics of web :
        HTML
        HTTP
        Sessions and cookies
    Using and setting up Docker
        Containerizing an application
        Compatibility/Dependency
        Creating images
    SQL language
        Manipulation of databases
    The basics of encryption

## ✍️ Authors <a name = "authors"></a>

Author: ZONE 01 COPYRIGHT <br>
igueye <br>
nifaye<br>
ymadike <br>
daiba <br>
ndiba

## 🎈 Usage <a name="usage"></a>

To run this project you to do :

- copy that command to start the server

```
go run .
```

After that we gonna tape formuler and submit it to see haw does it work

## ⛏️ Implementation details <a name = "built_using"></a>

- [Golang](https://go.dev/) - for the back end we use golang
- [Html](https://developer.mozilla.org/fr/docs/Web/HTML) - for the front end
- [Css](https://developer.mozilla.org/fr/docs/Web/CSS) - for style
- [Docker](https://docs.docker.com/get-started/) - Docker is a platform for developing, shipping, and running applications inside lightweight, portable, and self-sufficient containers. Containers are a form of virtualization technology that packages an application and its dependencies together into a single, executable package, known as a container image. Docker provides tools and services to create, manage, and run these containers.
- [SQLite3](https://github.com/mattn/go-sqlite3) - for the Databse
- [bcrypt](https://pkg.go.dev/golang.org/x/crypto/bcrypt) - help to crypt the password for a user
- [UUID](https://github.com/gofrs/uuid) - Package uuid provides a pure Go implementation of Universally Unique Identifiers (UUID) variant as defined in RFC-4122

that all for this project
