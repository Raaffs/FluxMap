#!/bin/bash

# Set the URL to send GET requests to
URL="http://localhost:4000"  # Change this to your target URL
# Set the number of requests
NUM_REQUESTS=1000  # Number of total requests to send
# Set the desired rate (e.g., 20 requests per second)
RATE=100

# Calculate the delay in seconds between requests
DELAY=$(bc -l <<< "1 / $RATE")

echo "Sending $NUM_REQUESTS GET requests to $URL at a rate of $RATE requests/second..."

# Function to send a single GET request
send_request() {
    curl -s -o /dev/null "$URL"
}

# Loop to send the requests
for ((i=0; i<NUM_REQUESTS; i++)); do
    send_request &  # Send request in background
    # If we hit the rate limit, wait
    if (( (i + 1) % RATE == 0 )); then
        sleep 1  # Sleep for 1 second
    fi
done

# Wait for all background processes to complete
wait

echo "Finished sending requests."
