// app/_lib/api.ts
import axios, { AxiosInstance } from "axios";
import {
  ILoginRequestDTO,
  IRegisterRequestDTO,
  IOwnerLoginRequestDTO,
  IOwnerRegisterRequestDTO,
  IAuthResponse,
  IOwnerAuthResponse,
} from "../_interfaces/auth"; // Import *all* interfaces from _interfaces

const api: AxiosInstance = axios.create({
  baseURL: "http://localhost:8090/api/v1", // CORRECTED base URL
  headers: {
    "Content-Type": "application/json",
  },
  withCredentials: false,
});

// User Authentication
export const userLogin = async (
  data: ILoginRequestDTO
): Promise<IAuthResponse> => {
  try {
    const response = await api.post("/auth/login", data);
    return response.data;
  } catch (error: unknown) {
    console.error("User login error:", error);
    throw error;
  }
};

export const userRegister = async (
  data: IRegisterRequestDTO
): Promise<IAuthResponse> => {
  try {
    const response = await api.post("/auth/register", data);
    return response.data;
  } catch (error: unknown) {
    console.error("User registration error:", error);
    if (axios.isAxiosError(error) && error.response) {
      console.error("Server response:", error.response.data);
      throw error;
    }
    throw new Error("Registration failed: Network error");
  }
};

// Owner Authentication
export const ownerLogin = async (
  data: IOwnerLoginRequestDTO
): Promise<IOwnerAuthResponse> => {
  try {
    const response = await api.post("/owner/auth/login", data);
    return response.data;
  } catch (error: unknown) {
    console.error("Owner login error:", error);
    throw error;
  }
};

export const ownerRegister = async (
  data: IOwnerRegisterRequestDTO
): Promise<IOwnerAuthResponse> => {
  try {
    const response = await api.post("/owner/auth/register", data);
    return response.data;
  } catch (error: unknown) {
    console.error("Owner register error:", error);
    throw error;
  }
};

export default api;
