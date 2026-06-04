import type { ApiResponse } from "../../types/api";
import type { LoginRequest, LoginResponse, RegisterRequest, RegisterResponse } from "./types";

const API_URL =
  "http://localhost:8080/auth";

export async function login(payload: LoginRequest): Promise<ApiResponse<LoginResponse>> {
  const response = await fetch(
    `${API_URL}/login`,
    {
      method: "POST",
      headers: {
        "Content-Type":
          "application/json",
      },
      body: JSON.stringify({
        "email": payload.email,
        "password": payload.password
      }),
    }
  );

  return response.json();
}

export async function register(payload: RegisterRequest): Promise<ApiResponse<RegisterResponse>> {
  const response = await fetch(
    `${API_URL}/register`,
    {
      method: "POST",
      headers: {
        "Content-Type": "application/json"
      },
      body: JSON.stringify({
        "email": payload.email,
        "username": payload.username,
        "password": payload.password
      })
    }
  )

  return response.json();
}
