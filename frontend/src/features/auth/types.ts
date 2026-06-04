export interface LoginRequest {
  email: string;
  password: string;
}

export interface LoginResponse {
  user: User;
  tokens: JwtTokens;
}

export interface RegisterRequest {
  email: string;
  username: string;
  password: string;
}

export interface RegisterResponse {
  user: User;
  tokens: JwtTokens;
}

export interface User {
  id: string;
  username: string;
  xp: number;
  level: number;
  terrace_level: number;
}

export interface JwtTokens {
  access_token: string;
  refresh_token: string;
}
