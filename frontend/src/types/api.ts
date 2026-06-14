export interface ApiResponse<T> {
  success: boolean;
  data: T;
  message?: string;
  error?: ApiError;
}

export interface ApiError {
  code: string;
  message: string;
  details?: ValidationDetail[];
}

export interface ValidationDetail {
  field: string;
  message: string;
}
