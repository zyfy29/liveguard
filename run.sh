# Pull the latest changes from the repository
echo "Pulling latest changes from the repository..."
git pull

# Get the process ID of the running Go application
PID_GO=$(pgrep -f "bearguard")

# If a process is found, kill it
if [ -n "$PID_GO" ]; then
  echo "Killing existing Go process with PID: $PID_GO"
  kill $PID_GO
  sleep 2  # Give it a moment to ensure the process is fully terminated
fi

# Build the Go application
echo "Building application..."
go build -o bearguard -a

# Start the Go application
echo "Starting Go application..."
nohup ./bearguard &
