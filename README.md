# X Feeder Bot for Bluesky

This open-source Go application automatically forwards posts from a Twitter/X account to Bluesky. It periodically checks for new tweets and posts them to your Bluesky account, including images and videos. This uses the twitter scraper from [Imperatrona https://github.com/imperatrona/twitter-scraper](https://github.com/imperatrona/twitter-scraper) to find tweets for a given account in the environment variables.

## Features

- Automatic tweet fetching from a specified X account
- Posts text content to Bluesky
- Uploads and embeds images and videos
- Configurable posting schedule via cron
- Duplicate post detection to avoid reposting
- Docker support for easy deployment to Dokploy, Coolify, and more.
- Testing

## Prerequisites

- Go 1.18 or later
- Twitter/X auth token and CSRF token (from Twitter Developer Dashboard)
- Bluesky account

## Environment Variables

Create a `.env` file in the root directory with the following variables from the `.env.example` file:

```
# Bluesky Credentials
HOST=https://bsky.social
HANDLE=your-bsky-handle.bsky.social
PASSWORD=your-bsky-password

# Twitter/X Credentials
AUTH_TOKEN=your-twitter-auth-token
CSRF_TOKEN=your-twitter-csrf-token
TWITTER_ACCOUNT=twitter-account-to-follow

# Optional Configuration
JOB_SPEC=* * * * *
```

## Installation

1. Clone the repository:

```bash
git clone https://github.com/petermazzocco/go-x-feeder-bot.git
cd go-x-feeder-bot
```

2. Install dependencies:

```bash
go mod download
```

## Running Locally

### Option 1: Using Go directly

```bash
go run main.go
```

### Option 2: Using Air (Hot Reload)

If you have [Air](https://github.com/cosmtrek/air) installed:

```bash
air
```

## Deployment

### Docker

Build and run using Docker:

```bash
docker build -t yourname/go-x-feeder-bot:latest .
docker run --env-file .env x-feeder-bot
```

### Dokploy

How to deploy with GitHub:
1. Push your code to a Git repository
2. Connect your repository to Dokploy
3. Create a new app and select your repository
4. Configure environment variables in Dokploy dashboard
5. Deploy the application

How to deploy with Docker:
1. Build and push the Docker image:
   ```bash
   docker buildx build --platform linux/amd64 --push -t yourname/go-x-feeder-bot:latest .
   ```
2. On Dokploy, trigger a build on your service

## Getting Twitter Auth and CSRF Tokens

To get the required Twitter/X tokens:

1. Log in to Twitter in your web browser
2. Open Developer Tools (F12)
3. Go to the Application tab > Storage > Cookies
4. Find the `auth_token` cookie for the twitter.com domain
5. Find the `ct0` cookie, which is your CSRF token

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is open source and available under the [MIT License](LICENSE).
