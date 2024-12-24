# RAG LLM using LangChainGO

This project implements a Retrieval-Augmented Generation (RAG) system using [LangChainGo](https://github.com/tmc/langchaingo). The architecture consists of: 
- MongoDB instance for storing user history
- PostgreSQL database with pgvector for document embeddings
- Ollama service for managing LLM interactions
  - **nomic-embed-text** for embedding
  - **sailor2:8b** as the LLM

## Prerequisites
- Docker
- NVIDIA GPU and drivers for GPU-accelerated Ollama (optional).

## Getting Started
1. **Clone the Repository**:
   ```bash
   git clone https://github.com/oniharnantyo/lanchaingo-rag-golang
   cd lanchaingo-rag-golang
   ```
2. **Copy backend config**:

   Run the Docker Compose setup:
   ```bash
   cd backend
   cp .env.example .env
   ```

3. **Start Services**:

   Run the Docker Compose setup:
   ```bash
   docker-compose up -d
   ```

4. **Stopping Services**:

   To stop the services, run:
   ```bash
   docker-compose down
   ```

## API Documentation

### Add Documents
Uploads a document to the backend for processing and embedding.

**Request**:
```bash
curl --request POST \
  --url http://localhost:8000/documents \
  --header 'content-type: multipart/form-data' \
  --form 'file=@/path/to/your/document.pdf'
```

### Query
Sends a query to the backend and retrieves relevant results from the indexed documents.

**Request**:
```bash
curl --request POST \
  --url http://localhost:8000/queries \
  --header 'content-type: application/json' \
  --data '{
  "query": "your question here"
}'
```

## License
This project is licensed under the MIT License. See the LICENSE file for details.