## Installation
In order to run the project, please follow the following steps:
1. Clone the Repo
2. Run `go get` or `go mod dowload` or `go mod tidy`
## Routes
|  API Path   | Method |                                  What it does                                   |
| :---------: | :----: | :-----------------------------------------------------------------------------: |
| /api/scanAV |  POST  | scans provided paths and send result to s3, and a email to hardcoded recipients |
## Enviroment Variables
Add this varaibles to your .env file
- `AWS_KEY=""`
- `AWS_SECREY_KEY=""`
## Run
Make sure that your system is using go 1.19.0,
Run `curl -X POST http://localhost:8088/api/scanAv \
-H "Content-Type: application/json" \
-d '{"Paths": ["/tmp/", "/var/tmp/", "/root/apps/dev/svg-ui/"]}'`
## Deployment
In order to monitor and check logs will be using `PM2` to run `Go` binary
1. Check for port 8088, proceed to kill if is already been used, with `npx kill-port 8088`
2. Run `make run` to create new server file
3. Run `pm2 start ./server --name "svg_scan_av"`
4. Enjoy
