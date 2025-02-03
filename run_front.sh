# Pull the latest changes from the repository
echo "Pulling latest changes from the repository..."
git pull

# Get the process ID of the running Streamlit application
PID_STREAMLIT=$(pgrep -f "streamlit run front/panel.py")

# If a process is found, kill it
if [ -n "$PID_STREAMLIT" ]; then
  echo "Killing existing Streamlit process with PID: $PID_STREAMLIT"
  kill $PID_STREAMLIT
  sleep 2  # Give it a moment to ensure the process is fully terminated
fi

# Start the Streamlit application
echo "Starting Streamlit application..."
nohup python -m streamlit run front/panel.py &