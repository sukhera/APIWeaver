export const DESIGN_TOKENS = {
  colors: {
    primary: 'hsl(214, 95%, 67%)', // #3b82f6
    secondary: 'hsl(152, 81%, 39%)', // #10b981
    accent: 'hsl(43, 96%, 50%)', // #f59e0b
    error: 'hsl(0, 84%, 60%)', // #ef4444
    warning: 'hsl(43, 96%, 50%)', // #f59e0b
    success: 'hsl(152, 81%, 39%)', // #10b981
  },
  spacing: {
    xs: '0.25rem', // 4px
    sm: '0.5rem',  // 8px
    md: '1rem',    // 16px
    lg: '1.5rem',  // 24px
    xl: '2rem',    // 32px
    '2xl': '3rem', // 48px
  },
  breakpoints: {
    sm: '640px',
    md: '768px',
    lg: '1024px',
    xl: '1280px',
    '2xl': '1536px',
  },
  container: {
    workspace: '1400px',
    content: '1200px',
    prose: '800px',
    sidebar: '320px',
    panel: '480px',
  }
} as const

export const API_ENDPOINTS = {
  generate: '/api/generate',
  amend: '/api/amend',
  validate: '/api/validate',
  health: '/api/health',
  metrics: '/api/metrics',
} as const

export const LOCAL_STORAGE_KEYS = {
  theme: 'apiweaver-theme',
  recentFiles: 'apiweaver-recent-files',
  editorContent: 'apiweaver-editor-content',
  settings: 'apiweaver-settings',
} as const

export const VALIDATION_MESSAGES = {
  required: 'This field is required',
  invalidMarkdown: 'Invalid Markdown format',
  invalidJSON: 'Invalid JSON format',
  invalidYAML: 'Invalid YAML format',
  fileTooBig: 'File size too large',
  unsupportedFileType: 'Unsupported file type',
} as const

export const FILE_TYPES = {
  markdown: ['.md', '.markdown'],
  yaml: ['.yaml', '.yml'],
  json: ['.json'],
  openapi: ['.yaml', '.yml', '.json'],
} as const

export const MAX_FILE_SIZE = 10 * 1024 * 1024 // 10MB

export const SUPPORTED_EXTENSIONS = [
  ...FILE_TYPES.markdown,
  ...FILE_TYPES.yaml,
  ...FILE_TYPES.json,
] as const

export const DEFAULT_MARKDOWN_CONTENT = `---
title: My API
version: 1.0.0
servers:
  - url: https://api.example.com
---

# My API Documentation

## GET /users
Get all users from the system.

**Parameters:**

| Name | Type | Required | Description |
|------|------|----------|-------------|
| limit | integer | No | Maximum number of users to return |
| offset | integer | No | Number of users to skip |

**Responses:**

| Status | Description | Schema |
|--------|-------------|--------|
| 200 | Success | User array |
| 400 | Bad request | Error object |

**Response Schema:**

\`\`\`json
{
  "users": [
    {
      "id": "uuid",
      "name": "string",
      "email": "string",
      "created_at": "datetime"
    }
  ],
  "total": "integer"
}
\`\`\`

## POST /users
Create a new user.

**Request Body:**

\`\`\`json
{
  "name": "string",
  "email": "string"
}
\`\`\`

**Responses:**

| Status | Description | Schema |
|--------|-------------|--------|
| 201 | User created | User object |
| 400 | Invalid input | Error object |
| 409 | Email exists | Error object |
`

export const EXAMPLE_OPENAPI_SPEC = {
  openapi: '3.1.0',
  info: {
    title: 'My API',
    version: '1.0.0',
  },
  servers: [
    {
      url: 'https://api.example.com',
    },
  ],
  paths: {
    '/users': {
      get: {
        summary: 'Get all users from the system',
        parameters: [
          {
            name: 'limit',
            in: 'query',
            required: false,
            schema: {
              type: 'integer',
            },
            description: 'Maximum number of users to return',
          },
          {
            name: 'offset',
            in: 'query',
            required: false,
            schema: {
              type: 'integer',
            },
            description: 'Number of users to skip',
          },
        ],
        responses: {
          '200': {
            description: 'Success',
            content: {
              'application/json': {
                schema: {
                  type: 'object',
                  properties: {
                    users: {
                      type: 'array',
                      items: {
                        $ref: '#/components/schemas/User',
                      },
                    },
                    total: {
                      type: 'integer',
                    },
                  },
                },
              },
            },
          },
          '400': {
            description: 'Bad request',
            content: {
              'application/json': {
                schema: {
                  $ref: '#/components/schemas/Error',
                },
              },
            },
          },
        },
      },
      post: {
        summary: 'Create a new user',
        requestBody: {
          content: {
            'application/json': {
              schema: {
                type: 'object',
                properties: {
                  name: {
                    type: 'string',
                  },
                  email: {
                    type: 'string',
                  },
                },
                required: ['name', 'email'],
              },
            },
          },
        },
        responses: {
          '201': {
            description: 'User created',
            content: {
              'application/json': {
                schema: {
                  $ref: '#/components/schemas/User',
                },
              },
            },
          },
          '400': {
            description: 'Invalid input',
            content: {
              'application/json': {
                schema: {
                  $ref: '#/components/schemas/Error',
                },
              },
            },
          },
          '409': {
            description: 'Email exists',
            content: {
              'application/json': {
                schema: {
                  $ref: '#/components/schemas/Error',
                },
              },
            },
          },
        },
      },
    },
  },
  components: {
    schemas: {
      User: {
        type: 'object',
        properties: {
          id: {
            type: 'string',
            format: 'uuid',
          },
          name: {
            type: 'string',
          },
          email: {
            type: 'string',
            format: 'email',
          },
          created_at: {
            type: 'string',
            format: 'date-time',
          },
        },
        required: ['id', 'name', 'email', 'created_at'],
      },
      Error: {
        type: 'object',
        properties: {
          error: {
            type: 'string',
          },
          message: {
            type: 'string',
          },
          code: {
            type: 'integer',
          },
        },
        required: ['error', 'message', 'code'],
      },
    },
  },
}