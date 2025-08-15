import type {
  GenerateRequest,
  GenerateResponse,
  ValidateRequest,
  ValidateResponse,
  AmendRequest,
  AmendResponse,
  DiffRequest,
  DiffResponse,
  HealthResponse,
  ApiError as ApiErrorType
} from '@/types/api'

class ApiClient {
  private baseURL: string

  constructor(baseURL: string = '/api') {
    this.baseURL = baseURL
  }

  private async request<T>(
    endpoint: string,
    options: RequestInit = {}
  ): Promise<T> {
    const url = `${this.baseURL}${endpoint}`
    
    const defaultHeaders = {
      'Content-Type': 'application/json',
    }

    const config: RequestInit = {
      ...options,
      headers: {
        ...defaultHeaders,
        ...options.headers,
      },
    }

    try {
      const response = await fetch(url, config)
      
      if (!response.ok) {
        let errorData: ApiErrorType
        try {
          errorData = await response.json()
        } catch {
          errorData = {
            error: 'Network Error',
            message: `HTTP ${response.status}: ${response.statusText}`,
          }
        }
        throw new ApiError(errorData.message, response.status)
      }

      // Handle empty responses
      const contentType = response.headers.get('content-type')
      if (contentType && contentType.includes('application/json')) {
        return await response.json()
      } else {
        // Return response text for non-JSON responses
        const text = await response.text()
        return text as unknown as T
      }
    } catch (error) {
      if (error instanceof ApiError) {
        throw error
      }
      
      // Network or other errors
      throw new ApiError(
        error instanceof Error ? error.message : 'Unknown error occurred',
        0
      )
    }
  }

  // Generate OpenAPI spec from Markdown
  async generate(request: GenerateRequest): Promise<GenerateResponse> {
    return this.request<GenerateResponse>('/generate', {
      method: 'POST',
      body: JSON.stringify(request),
    })
  }

  // Validate OpenAPI spec
  async validate(request: ValidateRequest): Promise<ValidateResponse> {
    return this.request<ValidateResponse>('/validate', {
      method: 'POST',
      body: JSON.stringify(request),
    })
  }

  // Amend existing OpenAPI spec
  async amend(request: AmendRequest): Promise<AmendResponse> {
    return this.request<AmendResponse>('/amend', {
      method: 'POST',
      body: JSON.stringify(request),
    })
  }

  // Get diff between two specs
  async diff(request: DiffRequest): Promise<DiffResponse> {
    return this.request<DiffResponse>('/diff', {
      method: 'POST',
      body: JSON.stringify(request),
    })
  }

  // Health check
  async health(): Promise<HealthResponse> {
    return this.request<HealthResponse>('/health', {
      method: 'GET',
    })
  }

  // Get service metrics (if available)
  async metrics(): Promise<Record<string, unknown>> {
    return this.request<Record<string, unknown>>('/metrics', {
      method: 'GET',
    })
  }
}

// Custom error class for API errors
class ApiError extends Error {
  constructor(
    message: string,
    public status: number = 0
  ) {
    super(message)
    this.name = 'ApiError'
  }
}

// Export singleton instance
export const apiClient = new ApiClient()
export { ApiClient, ApiError }