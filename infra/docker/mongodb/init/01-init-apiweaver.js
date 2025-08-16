// MongoDB initialization script for APIWeaver
// This script sets up the initial database structure and indexes

print('Starting APIWeaver MongoDB initialization...');

// Switch to the APIWeaver database
db = db.getSiblingDB('apiweaver');

// Create collections with validation schemas
print('Creating collections with validation schemas...');

// API Keys collection
db.createCollection('api_keys', {
  validator: {
    $jsonSchema: {
      bsonType: 'object',
      required: ['key_id', 'key_hash', 'name', 'created_at', 'is_active'],
      properties: {
        key_id: {
          bsonType: 'string',
          description: 'Unique identifier for the API key'
        },
        key_hash: {
          bsonType: 'string',
          description: 'Hashed API key value'
        },
        name: {
          bsonType: 'string',
          description: 'Human-readable name for the API key'
        },
        description: {
          bsonType: 'string',
          description: 'Optional description of the API key usage'
        },
        permissions: {
          bsonType: 'array',
          items: {
            bsonType: 'string'
          },
          description: 'Array of permissions granted to this API key'
        },
        rate_limit: {
          bsonType: 'object',
          properties: {
            requests_per_minute: {
              bsonType: 'int',
              minimum: 1
            },
            burst: {
              bsonType: 'int',
              minimum: 1
            }
          }
        },
        created_at: {
          bsonType: 'date',
          description: 'Timestamp when the API key was created'
        },
        updated_at: {
          bsonType: 'date',
          description: 'Timestamp when the API key was last updated'
        },
        last_used_at: {
          bsonType: 'date',
          description: 'Timestamp when the API key was last used'
        },
        expires_at: {
          bsonType: 'date',
          description: 'Optional expiration timestamp'
        },
        is_active: {
          bsonType: 'bool',
          description: 'Whether the API key is currently active'
        }
      }
    }
  }
});

// Usage statistics collection
db.createCollection('usage_stats', {
  validator: {
    $jsonSchema: {
      bsonType: 'object',
      required: ['api_key_id', 'endpoint', 'method', 'timestamp', 'status_code'],
      properties: {
        api_key_id: {
          bsonType: 'string',
          description: 'ID of the API key used'
        },
        endpoint: {
          bsonType: 'string',
          description: 'API endpoint accessed'
        },
        method: {
          bsonType: 'string',
          enum: ['GET', 'POST', 'PUT', 'DELETE', 'PATCH'],
          description: 'HTTP method used'
        },
        timestamp: {
          bsonType: 'date',
          description: 'When the request was made'
        },
        status_code: {
          bsonType: 'int',
          minimum: 100,
          maximum: 599,
          description: 'HTTP status code returned'
        },
        response_time_ms: {
          bsonType: 'int',
          minimum: 0,
          description: 'Response time in milliseconds'
        },
        request_size_bytes: {
          bsonType: 'long',
          minimum: 0,
          description: 'Size of the request in bytes'
        },
        response_size_bytes: {
          bsonType: 'long',
          minimum: 0,
          description: 'Size of the response in bytes'
        },
        user_agent: {
          bsonType: 'string',
          description: 'User agent of the client'
        },
        ip_address: {
          bsonType: 'string',
          description: 'Client IP address'
        }
      }
    }
  }
});

// Generated specs cache collection
db.createCollection('generated_specs', {
  validator: {
    $jsonSchema: {
      bsonType: 'object',
      required: ['spec_id', 'input_hash', 'output_spec', 'created_at'],
      properties: {
        spec_id: {
          bsonType: 'string',
          description: 'Unique identifier for the generated spec'
        },
        input_hash: {
          bsonType: 'string',
          description: 'Hash of the input markdown content'
        },
        input_content: {
          bsonType: 'string',
          description: 'Original markdown content'
        },
        output_spec: {
          bsonType: 'object',
          description: 'Generated OpenAPI specification'
        },
        output_format: {
          bsonType: 'string',
          enum: ['json', 'yaml'],
          description: 'Format of the output specification'
        },
        generation_options: {
          bsonType: 'object',
          description: 'Options used during generation'
        },
        created_at: {
          bsonType: 'date',
          description: 'When the spec was generated'
        },
        accessed_at: {
          bsonType: 'date',
          description: 'When the spec was last accessed'
        },
        access_count: {
          bsonType: 'int',
          minimum: 0,
          description: 'Number of times this spec has been accessed'
        },
        ttl: {
          bsonType: 'date',
          description: 'Time-to-live for cache expiration'
        }
      }
    }
  }
});

// Create indexes for performance
print('Creating indexes for better performance...');

// API Keys indexes
db.api_keys.createIndex({ 'key_id': 1 }, { unique: true });
db.api_keys.createIndex({ 'key_hash': 1 }, { unique: true });
db.api_keys.createIndex({ 'is_active': 1 });
db.api_keys.createIndex({ 'created_at': 1 });
db.api_keys.createIndex({ 'expires_at': 1 }, { sparse: true });

// Usage statistics indexes
db.usage_stats.createIndex({ 'api_key_id': 1, 'timestamp': -1 });
db.usage_stats.createIndex({ 'timestamp': -1 });
db.usage_stats.createIndex({ 'endpoint': 1, 'method': 1 });
db.usage_stats.createIndex({ 'status_code': 1 });

// Generated specs indexes
db.generated_specs.createIndex({ 'spec_id': 1 }, { unique: true });
db.generated_specs.createIndex({ 'input_hash': 1 });
db.generated_specs.createIndex({ 'created_at': -1 });
db.generated_specs.createIndex({ 'accessed_at': -1 });
db.generated_specs.createIndex({ 'ttl': 1 }, { expireAfterSeconds: 0 });

// Create default API key for development
if (db.api_keys.countDocuments() === 0) {
  print('Creating default development API key...');
  
  const defaultApiKey = {
    key_id: 'dev-key-001',
    key_hash: 'dev-hash-for-local-development-only',
    name: 'Development API Key',
    description: 'Default API key for local development - DO NOT USE IN PRODUCTION',
    permissions: ['generate', 'validate', 'amend', 'admin'],
    rate_limit: {
      requests_per_minute: 1000,
      burst: 50
    },
    created_at: new Date(),
    updated_at: new Date(),
    is_active: true
  };
  
  db.api_keys.insertOne(defaultApiKey);
  print('Default API key created successfully');
}

// Create test data for development
if (process.env.NODE_ENV === 'development') {
  print('Creating test data for development environment...');
  
  // Sample usage statistics
  const sampleStats = [
    {
      api_key_id: 'dev-key-001',
      endpoint: '/api/generate',
      method: 'POST',
      timestamp: new Date(Date.now() - 3600000), // 1 hour ago
      status_code: 200,
      response_time_ms: 1250,
      request_size_bytes: 2048,
      response_size_bytes: 8192,
      user_agent: 'APIWeaver-Web/1.0.0',
      ip_address: '127.0.0.1'
    },
    {
      api_key_id: 'dev-key-001',
      endpoint: '/api/validate',
      method: 'POST',
      timestamp: new Date(Date.now() - 1800000), // 30 minutes ago
      status_code: 200,
      response_time_ms: 800,
      request_size_bytes: 1024,
      response_size_bytes: 512,
      user_agent: 'APIWeaver-CLI/1.0.0',
      ip_address: '127.0.0.1'
    }
  ];
  
  db.usage_stats.insertMany(sampleStats);
  print('Sample usage statistics created');
}

// Create application configuration
db.createCollection('app_config');
db.app_config.createIndex({ 'key': 1 }, { unique: true });

// Insert default configuration
const defaultConfig = [
  {
    key: 'rate_limiting.default_requests_per_minute',
    value: 100,
    description: 'Default rate limit for API requests',
    updated_at: new Date()
  },
  {
    key: 'rate_limiting.default_burst',
    value: 10,
    description: 'Default burst limit for API requests',
    updated_at: new Date()
  },
  {
    key: 'cache.generated_specs_ttl_hours',
    value: 24,
    description: 'Time-to-live for cached generated specifications in hours',
    updated_at: new Date()
  },
  {
    key: 'security.api_key_expiry_days',
    value: 365,
    description: 'Default expiry period for API keys in days',
    updated_at: new Date()
  }
];

db.app_config.insertMany(defaultConfig);

print('APIWeaver MongoDB initialization completed successfully!');

// Print summary
print('=== Initialization Summary ===');
print('Collections created: api_keys, usage_stats, generated_specs, app_config');
print('Indexes created for optimal performance');
print('Default development API key created (key_id: dev-key-001)');
print('Default configuration values set');
print('Database is ready for APIWeaver application!');