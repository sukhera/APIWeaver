// API Request/Response Types
export interface GenerateRequest {
  markdown: string
  format?: 'yaml' | 'json'
  validate?: boolean
}

export interface GenerateResponse {
  spec: string
  format: 'yaml' | 'json'
  errors?: ValidationError[]
  warnings?: ValidationError[]
}

export interface ValidateRequest {
  spec: string
  format?: 'yaml' | 'json'
}

export interface ValidateResponse {
  valid: boolean
  errors?: ValidationError[]
  warnings?: ValidationError[]
  summary?: ValidationSummary
}

export interface AmendRequest {
  originalSpec: string
  changes: string
  format?: 'yaml' | 'json'
  dryRun?: boolean
}

export interface AmendResponse {
  amendedSpec: string
  diff?: DiffResult
  errors?: ValidationError[]
  warnings?: ValidationError[]
}

export interface DiffRequest {
  originalSpec: string
  modifiedSpec: string
  format?: 'yaml' | 'json'
}

export interface DiffResponse {
  diff: DiffResult
  summary: DiffSummary
}

// Validation Types
export interface ValidationError {
  line?: number
  column?: number
  path?: string
  message: string
  severity: 'error' | 'warning' | 'info'
  code?: string
}

export interface ValidationSummary {
  totalErrors: number
  totalWarnings: number
  endpoints: number
  schemas: number
  parameters: number
}

// Diff Types
export interface DiffResult {
  added: DiffLine[]
  removed: DiffLine[]
  modified: DiffLine[]
  unchanged: DiffLine[]
}

export interface DiffLine {
  lineNumber: number
  content: string
  path?: string
  type: 'added' | 'removed' | 'modified' | 'unchanged'
}

export interface DiffSummary {
  linesAdded: number
  linesRemoved: number
  linesModified: number
  endpointsAdded: number
  endpointsRemoved: number
  endpointsModified: number
}

// Health Check
export interface HealthResponse {
  status: 'healthy' | 'unhealthy'
  version: string
  uptime: number
  timestamp: string
}

// API Error Response
export interface ApiError {
  error: string
  message: string
  code?: string
  details?: Record<string, unknown>
}