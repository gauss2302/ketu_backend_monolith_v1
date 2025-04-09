/* eslint-disable @typescript-eslint/no-explicit-any */
// src/app/_lib/authUtils.ts

import api from "@/app/_lib/api";
import { jwtDecode } from "jwt-decode";

export interface JwtPayload {
  exp: number;
  user_id?: number;
  id?: number;
  email: string;
  role?: string;
  [key: string]: any;
}

// Check if a token is expired
export const isTokenExpired = (token: string): boolean => {
  try {
    const decoded = jwtDecode<JwtPayload>(token);
    // Check if the expiration time is in the past
    // Subtract 60 seconds to refresh the token a bit before it actually expires
    return decoded.exp * 1000 < Date.now() - 60000;
  } catch (error) {
    console.error("Error decoding token:", error);
    return true; // Assume expired if we can't decode it
  }
};

// Refresh user token
export const refreshUserToken = async (token: string): Promise<string> => {
  try {
    const response = await api.post(
      "/auth/refresh",
      {},
      {
        headers: {
          Authorization: `Bearer ${token}`,
        },
      }
    );
    return response.data.accessToken;
  } catch (error) {
    console.error("Failed to refresh user token:", error);
    throw error;
  }
};

// Refresh owner token
export const refreshOwnerToken = async (token: string): Promise<string> => {
  try {
    const response = await api.post(
      "/owner/auth/refresh",
      {},
      {
        headers: {
          Authorization: `Bearer ${token}`,
        },
      }
    );
    return response.data.accessToken;
  } catch (error) {
    console.error("Failed to refresh owner token:", error);
    throw error;
  }
};

// Get a fresh token (refresh if needed)
export const getFreshToken = async (
  token: string | null,
  isOwner: boolean = false
): Promise<string | null> => {
  if (!token) return null;

  if (isTokenExpired(token)) {
    try {
      // Attempt to refresh the token
      const newToken = isOwner
        ? await refreshOwnerToken(token)
        : await refreshUserToken(token);

      // Store the new token
      if (isOwner) {
        localStorage.setItem("ownerAccessToken", newToken);
      } else {
        localStorage.setItem("accessToken", newToken);
      }

      return newToken;
      // eslint-disable-next-line @typescript-eslint/no-unused-vars
    } catch (error) {
      // If refresh fails, clear the token
      if (isOwner) {
        localStorage.removeItem("ownerAccessToken");
      } else {
        localStorage.removeItem("accessToken");
      }

      // Force a page reload to redirect to login
      if (typeof window !== "undefined") {
        window.location.href = isOwner ? "/owner-login" : "/login";
      }

      return null;
    }
  }

  return token;
};
