# Steak
A repository for recipes. This project was started mainly for two reasons. The first being that I was getting into cooking and needed a place to store my recipes and I wasn't about to start carrying around a notebook. The second was that I wanted to try web assemblies.

Website: https://alistairfink.com/steak

## Build
Backend - Build Docker Image (Actual Docker Repo is Private) 

Run the following from the project root.
```bash
$ docker build -t alistairfink/side_projects:steak_backend .
$ docker build -t alistairfink/side_projects:steak_db -f ./Backend/DB-Dockerfile .
$ docker-compose up -d
```
The db should now be available on port 5430 and the backend on 41690.

Frontend - Serve using any web server.

Example using serve from the root of the repo.
```bash
$ cd ./Frontend
$ GOOS=js GOARCH=wasm go build -o main.wasm -ldflags="-s -w"
$ serve -s . 
```

If using serve this will most likely bind to port 5000 and the front end should be available at localhost:5000. The "serverURL" in the frontend's main.go will need to be changed to point at localhost:41690. (Please don't spam my server). 

## What I Learned From This
- Web assembilies are kinda trash right now. If you're mostly doing dom manipulation like I was doing in this project them your code will very closely resemble jquery (or at least it did to me). This will most likely be fixed as more frameworks come out. 
- Web assembilies are kinda difficult to work with. Not only did I not know what I was doing but finding relevant resources was very difficult. Most of what I found were simple tutorials with very simple operations that you might do in the first 5 minutes of a project. Anything past that was a challenge to find information on. This will almost definitely be fixed with time as more people adopt the tech.
- Web assembilies are kinda big. The one compiled in this project is ~8mb. I have attempted to compress this but mostly being self taught in devops I have no idea how I would go about serving zipped files. This will be fixed as soon as I stop being retarded.

## Why is this Called Steak?
This is an app about food. As a result when I was trying to think of a good name for this all I could think about was food. Steak is my favourite food and I couldn't think of anything else so I thought why not?
