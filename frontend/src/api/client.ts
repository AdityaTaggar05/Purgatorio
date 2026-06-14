import type { ApiResponse } from "../types/api";

type HttpMethod = "GET" | "POST" | "PUT" | "DELETE";

interface ClientConfig {
  baseUrl: string;
  getToken: () => string | null;
  onTokenRefresh: () => Promise<string | null>;
  onAuthFailure: () => void;
}

export class ApiClient {
  private baseUrl: string;
  private getToken: ClientConfig["getToken"];
  private onTokenRefresh: ClientConfig["onTokenRefresh"];
  private onAuthFailure: ClientConfig["onAuthFailure"];

  constructor(config: ClientConfig) {
    this.baseUrl = config.baseUrl;
    this.getToken = config.getToken;
    this.onTokenRefresh = config.onTokenRefresh;
    this.onAuthFailure = config.onAuthFailure;
  }

  private async request<T>(
    path: string,
    method: HttpMethod,
    body?: unknown,
    retried = false
  ): Promise<ApiResponse<T>> {
    const headers: Record<string, string> = {};
    const token = this.getToken();
    if (token) {
      headers["Authorization"] = `Bearer ${token}`;
    }
    if (body !== undefined) {
      headers["Content-Type"] = "application/json";
    }

    const response = await fetch(`${this.baseUrl}${path}`, {
      method,
      headers,
      credentials: "include",
      body: body !== undefined ? JSON.stringify(body) : undefined,
    });

    if (response.status === 401 && !retried) {
      const freshToken = await this.onTokenRefresh();
      if (freshToken) {
        return this.request<T>(path, method, body, true);
      }
      this.onAuthFailure();
      return {
        success: false,
        data: null as unknown as T,
        error: { code: "Unauthorized", message: "Session expired" },
      };
    }

    if (response.status === 204) {
      return { success: true, data: null as unknown as T };
    }

    return response.json();
  }

  get<T>(path: string): Promise<ApiResponse<T>> {
    return this.request<T>(path, "GET");
  }

  post<T>(path: string, body?: unknown): Promise<ApiResponse<T>> {
    return this.request<T>(path, "POST", body);
  }

  put<T>(path: string, body?: unknown): Promise<ApiResponse<T>> {
    return this.request<T>(path, "PUT", body);
  }

  del<T>(path: string, body?: unknown): Promise<ApiResponse<T>> {
    return this.request<T>(path, "DELETE", body);
  }
}
