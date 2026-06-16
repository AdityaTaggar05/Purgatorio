import type { ApiResponse } from "../../types/api";
import type { LoginRequest, LoginResponse, RefreshResponse, RegisterRequest, RegisterResponse } from "./types";

const API_URL = "/auth";

export async function login(payload: LoginRequest): Promise<ApiResponse<LoginResponse>> {
  const response = await fetch(
    `${API_URL}/login`,
    {
      method: "POST",
      credentials: "include",
      headers: {
        "Content-Type":
          "application/json",
      },
      body: JSON.stringify(payload),
    }
  );

  return response.json();
}

export async function register(payload: RegisterRequest): Promise<ApiResponse<RegisterResponse>> {
  const response = await fetch(
    `${API_URL}/register`,
    {
      method: "POST",
      credentials: "include",
      headers: {
        "Content-Type": "application/json"
      },
      body: JSON.stringify(payload)
    }
  )

  return response.json();
}

export async function logout(): Promise<boolean> {
  const response = await fetch(
    `${API_URL}/logout`,
    {
      method: "POST",
      credentials: "include"
    }
  )

  return response.ok
}

export async function refresh(): Promise<ApiResponse<RefreshResponse>> {
  const response = await fetch(
    `${API_URL}/refresh`,
    {
      method: "POST",
      credentials: "include"
    }
  )

  return response.json()
}
