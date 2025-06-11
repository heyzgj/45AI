// Basic entities
export interface User {
  id: number;
  wechat_openid: string;
  nickname: string;
  avatar_url: string;
  credits: number;
  created_at: string;
  updated_at: string;
}

export interface Template {
  id: number;
  name: string;
  description: string;
  preview_image_url: string;
  credit_cost: number;
  is_active: boolean;
  created_at: string;
}

export interface Transaction {
  id: number;
  user_id: number;
  type: 'purchase' | 'generation';
  amount: number;
  description: string;
  external_payment_id?: string;
  related_template_id?: number;
  created_at: string;
}

// Request/Response types for Authentication
export interface LoginRequest {
  code: string;
}

export interface LoginResponse {
  user: User;
  token: string;
  refreshToken?: string;
}

// Request/Response types for Templates
export interface GetTemplatesRequest {
  // No parameters needed for getting all templates
}

export interface GetTemplatesResponse {
  templates?: Template[];
  list?: Template[];
  total?: number;
  hasMore?: boolean;
}

// Request/Response types for Image Generation
export interface GenerateImageRequest {
  template_id: number;
  image?: File | string;
  user_photo_url?: string;
}

export interface GenerateImageResponse {
  job_id?: string;
  task_id?: string;
  status: 'pending' | 'processing' | 'completed' | 'failed' | 'succeeded';
  images?: string[];
  error?: string;
  estimated_time?: number;
}

export interface GetGenerationStatusRequest {
  job_id?: string;
  task_id?: string;
}

export interface GetGenerationStatusResponse {
  job_id?: string;
  task_id?: string;
  status: 'pending' | 'processing' | 'completed' | 'failed' | 'succeeded';
  images?: string[];
  result_url?: string;
  error?: string;
  error_message?: string;
  progress?: number;
}

// Request/Response types for Transactions and Credits
export interface GetUserCreditsResponse {
  credits: number;
}

export interface PurchaseCreditsRequest {
  product_id: string;
  platform: 'wechat' | 'apple';
  receipt: string;
}

export interface PurchaseCreditsResponse {
  success: boolean;
  credits_added: number;
  new_balance: number;
  transaction_id: number;
}

export interface GetTransactionsRequest {
  limit?: number;
  offset?: number;
}

export interface GetTransactionsResponse {
  transactions: Transaction[];
  total: number;
}

// Generic HTTP request types
export interface RequestOptions {
  url: string;
  method?: 'GET' | 'POST' | 'PUT' | 'DELETE';
  data?: any;
  header?: Record<string, string>;
  timeout?: number;
  showLoading?: boolean;
  loadingText?: string;
  retries?: number;
}

export interface ResponseData<T = any> {
  code: number;
  message: string;
  data: T;
} 