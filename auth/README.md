# Yume
For storing users' namespace and fetching/insert vectors into them.

## Setup
- Setup pinecone instance (name: "atari", size: 768, metric: "cosine", spec is aws us-east-1)
- `encore secret set --dev pineconeKey` - Set the pinecone key for the project (dev environment)
