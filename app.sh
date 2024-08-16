COMMAND=$1

if [ "$COMMAND" = "start" ]; then
    echo "Starting React..."
    cd Client
    start "" npm run dev
    echo "Starting Server..."
    cd ../Server
    start "" CompileDaemon -command="./project-gin"
elif [ "$COMMAND" = "stop" ]; then
    echo "Stopping all running instance..."
    pkill -f "node"
elif [ "$COMMAND" = "install" ]; then
    echo "preparing project..."
    cd Client
    start "" npm install -y
    cd../Server
    start "" go run migrate/migrate.go
else
    echo "Invalid COMMAND. Please provide either 'start' or 'stop'."
fi