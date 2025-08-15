import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import { apiClient } from '@/lib/api-client'
import type {
  GenerateRequest,
  ValidateRequest,
  AmendRequest,
  DiffRequest
} from '@/types/api'

// Query Keys
export const queryKeys = {
  health: ['health'] as const,
  metrics: ['metrics'] as const,
  generate: (request: GenerateRequest) => ['generate', request] as const,
  validate: (request: ValidateRequest) => ['validate', request] as const,
  amend: (request: AmendRequest) => ['amend', request] as const,
  diff: (request: DiffRequest) => ['diff', request] as const,
}

// Health Check Query
export function useHealthQuery() {
  return useQuery({
    queryKey: queryKeys.health,
    queryFn: () => apiClient.health(),
    refetchInterval: 30000, // Refetch every 30 seconds
    retry: 3,
    staleTime: 10000, // Consider data fresh for 10 seconds
  })
}

// Metrics Query
export function useMetricsQuery() {
  return useQuery({
    queryKey: queryKeys.metrics,
    queryFn: () => apiClient.metrics(),
    refetchInterval: 60000, // Refetch every minute
    retry: 1,
    enabled: false, // Only fetch when explicitly requested
  })
}

// Generate Mutation
export function useGenerateMutation() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: apiClient.generate,
    onSuccess: (data, variables) => {
      // Cache the result for potential reuse
      queryClient.setQueryData(queryKeys.generate(variables), data)
    },
    onError: (error) => {
      console.error('Generate error:', error)
    },
  })
}

// Validate Mutation
export function useValidateMutation() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: apiClient.validate,
    onSuccess: (data, variables) => {
      queryClient.setQueryData(queryKeys.validate(variables), data)
    },
    onError: (error) => {
      console.error('Validate error:', error)
    },
  })
}

// Amend Mutation
export function useAmendMutation() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: apiClient.amend,
    onSuccess: (data, variables) => {
      queryClient.setQueryData(queryKeys.amend(variables), data)
      // Invalidate related queries
      queryClient.invalidateQueries({ queryKey: ['generate'] })
      queryClient.invalidateQueries({ queryKey: ['validate'] })
    },
    onError: (error) => {
      console.error('Amend error:', error)
    },
  })
}

// Diff Mutation
export function useDiffMutation() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: apiClient.diff,
    onSuccess: (data, variables) => {
      queryClient.setQueryData(queryKeys.diff(variables), data)
    },
    onError: (error) => {
      console.error('Diff error:', error)
    },
  })
}

// Combined hook for real-time validation
export function useRealtimeValidation(spec: string, format: 'yaml' | 'json' = 'yaml') {
  return useQuery({
    queryKey: queryKeys.validate({ spec, format }),
    queryFn: () => apiClient.validate({ spec, format }),
    enabled: spec.length > 0,
    refetchOnWindowFocus: false,
    staleTime: 1000, // Consider data fresh for 1 second
    retry: 1,
  })
}

// Hook for generating and validating in sequence
export function useGenerateAndValidate() {
  const generateMutation = useGenerateMutation()
  const validateMutation = useValidateMutation()

  const generateAndValidate = async (request: GenerateRequest) => {
    const generateResult = await generateMutation.mutateAsync(request)
    
    if (generateResult.spec) {
      const validateResult = await validateMutation.mutateAsync({
        spec: generateResult.spec,
        format: generateResult.format
      })
      
      return {
        generate: generateResult,
        validate: validateResult
      }
    }
    
    return { generate: generateResult, validate: null }
  }

  return {
    generateAndValidate,
    isLoading: generateMutation.isPending || validateMutation.isPending,
    error: generateMutation.error || validateMutation.error,
    reset: () => {
      generateMutation.reset()
      validateMutation.reset()
    }
  }
}