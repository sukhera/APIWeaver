import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest'
import { apiClient, ApiError } from '../api-client'
import type { GenerateRequest, ValidateRequest } from '@/types/api'

describe('ApiClient', () => {
  beforeEach(() => {
    vi.mocked(global.fetch).mockClear()
  })

  afterEach(() => {
    vi.resetAllMocks()
  })

  describe('generate', () => {
    it('makes a successful generate request', async () => {
      const mockResponse = {
        spec: 'openapi: 3.1.0\ninfo:\n  title: Test API\n  version: 1.0.0',
        format: 'yaml',
        errors: [],
        warnings: [],
      }

      vi.mocked(global.fetch).mockResolvedValueOnce({
        ok: true,
        headers: new Headers({ 'content-type': 'application/json' }),
        json: async () => mockResponse,
      })

      const request: GenerateRequest = {
        markdown: '# Test API\n\n## GET /users',
        format: 'yaml',
        validate: true,
      }

      const result = await apiClient.generate(request)

      expect(global.fetch).toHaveBeenCalledWith('/api/generate', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(request),
      })

      expect(result).toEqual(mockResponse)
    })




  })

  describe('validate', () => {
    it('makes a successful validate request', async () => {
      const mockResponse = {
        valid: true,
        errors: [],
        warnings: [],
        summary: {
          totalErrors: 0,
          totalWarnings: 0,
          endpoints: 1,
          schemas: 0,
          parameters: 0,
        },
      }

      vi.mocked(global.fetch).mockResolvedValueOnce({
        ok: true,
        headers: new Headers({ 'content-type': 'application/json' }),
        json: async () => mockResponse,
      })

      const request: ValidateRequest = {
        spec: 'openapi: 3.1.0\ninfo:\n  title: Test API\n  version: 1.0.0',
        format: 'yaml',
      }

      const result = await apiClient.validate(request)

      expect(global.fetch).toHaveBeenCalledWith('/api/validate', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(request),
      })

      expect(result).toEqual(mockResponse)
    })

    it('handles validation errors', async () => {
      const mockResponse = {
        valid: false,
        errors: [
          {
            line: 5,
            column: 10,
            path: '/info/title',
            message: 'Title is required',
            severity: 'error' as const,
            code: 'MISSING_TITLE',
          },
        ],
        warnings: [],
      }

      vi.mocked(global.fetch).mockResolvedValueOnce({
        ok: true,
        headers: new Headers({ 'content-type': 'application/json' }),
        json: async () => mockResponse,
      })

      const request: ValidateRequest = {
        spec: 'openapi: 3.1.0\ninfo:\n  version: 1.0.0',
        format: 'yaml',
      }

      const result = await apiClient.validate(request)
      expect(result.valid).toBe(false)
      expect(result.errors).toHaveLength(1)
      expect(result.errors![0].message).toBe('Title is required')
    })
  })

  describe('health', () => {
    it('makes a successful health check request', async () => {
      const mockResponse = {
        status: 'healthy',
        version: '1.0.0',
        uptime: 3600,
        timestamp: '2023-01-01T00:00:00Z',
      }

      vi.mocked(global.fetch).mockResolvedValueOnce({
        ok: true,
        headers: new Headers({ 'content-type': 'application/json' }),
        json: async () => mockResponse,
      })

      const result = await apiClient.health()

      expect(global.fetch).toHaveBeenCalledWith('/api/health', {
        method: 'GET',
        headers: {
          'Content-Type': 'application/json',
        },
      })

      expect(result).toEqual(mockResponse)
    })
  })

  describe('error handling', () => {
    it('handles non-JSON error responses', async () => {
      vi.mocked(global.fetch).mockResolvedValueOnce({
        ok: false,
        status: 500,
        statusText: 'Internal Server Error',
        headers: new Headers({ 'content-type': 'text/plain' }),
        json: async () => {
          throw new Error('Not JSON')
        },
      })

      const request: GenerateRequest = {
        markdown: '# Test API',
        format: 'yaml',
      }

      await expect(apiClient.generate(request)).rejects.toThrow(ApiError)
    })

    it('handles empty responses', async () => {
      vi.mocked(global.fetch).mockResolvedValueOnce({
        ok: true,
        headers: new Headers({ 'content-type': 'text/plain' }),
        text: async () => 'Success',
        json: async () => {
          throw new Error('Not JSON')
        },
      })

      const result = await apiClient.health()
      expect(result).toBe('Success')
    })
  })

  describe('ApiError class', () => {
    it('creates error with message and status', () => {
      const error = new ApiError('Test error', 400)
      
      expect(error.message).toBe('Test error')
      expect(error.status).toBe(400)
      expect(error.name).toBe('ApiError')
      expect(error).toBeInstanceOf(Error)
    })

    it('has default status of 0', () => {
      const error = new ApiError('Test error')
      
      expect(error.status).toBe(0)
    })
  })
})