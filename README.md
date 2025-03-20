# What is it
It is a rate limiter.

# For what
To limit the number of requests.

# Algorithms
Token Bucket, Leaky Bucket, Sliding Window, Fixed Window.

# Usage
```./stopper -host="0.0.0.0" -port=9090 -target="https://google.com" -limiter="token_bucket" -interval="1s" -limit=100```


# Help
```./stopper -h```