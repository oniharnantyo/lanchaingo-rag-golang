#!/bin/bash

# Start Ollama in the background.
/bin/ollama serve &
# Record Process ID.
pid=$!

# Pause for Ollama to start.
sleep 5

echo "🔴 Retrieve nomic-embed-text model..."
ollama pull nomic-embed-text
echo "🟢 Done!"

echo "🔴 Retrieve ollama run sailor2:8b model..."
ollama pull sailor2:8b
echo "🟢 ollama run sailor2:8b Done!"

# Wait for Ollama process to finish.
wait $pid