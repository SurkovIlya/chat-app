# chat-app
Chat application using a web socket

## A local option to launch the application and work with it:
### Requirements
* Docker and Go
### Usage
Clone the repository with:
```bash
git clone github.com/SurkovIlya/chat-app
```
Copy the `env.example` file to a `.env` file.
```bash
cp .env.example .env
```
Update the postgres variables declared in the new `.env` to match your preference. 

Build and start the services with:
```bash
docker-compose up --build
```
### chat-app connect

<summary> <h4>{chat-app-host}/chat - connecting to the application on web socket</h4></summary>
  

