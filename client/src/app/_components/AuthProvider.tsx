/* eslint-disable @typescript-eslint/no-explicit-any */
// In AuthProvider.tsx

"use client";
import React, { useState, ReactNode, useCallback, useEffect } from "react";
import { AuthContext } from "./AuthContext";
import api, {
  userLogin,
  userRegister,
  ownerLogin as ownerLoginApi,
  ownerRegister as ownerRegisterApi,
} from "../_lib/api";
import {
  IAuthResponse,
  IOwnerAuthResponse,
  IUser,
  IOwner,
  ILoginRequestDTO,
  IRegisterRequestDTO,
  IOwnerLoginRequestDTO,
  IOwnerRegisterRequestDTO,
} from "../_interfaces/auth";
import { useRouter } from "next/navigation";
import { jwtDecode } from "jwt-decode";
import Cookies from "js-cookie"; // We'll need to install this package

type AuthProviderProps = {
  children: ReactNode;
};

export const AuthProvider: React.FC<AuthProviderProps> = ({ children }) => {
  const [user, setUser] = useState<IUser | null>(null);
  const [owner, setOwner] = useState<IOwner | null>(null);
  const [accessToken, setAccessToken] = useState<string | null>(null);
  const [ownerAccessToken, setOwnerAccessToken] = useState<string | null>(null);
  const [loading, setLoading] = useState<boolean>(true);
  const [isClient, setIsClient] = useState<boolean>(false);
  const router = useRouter();

  // Set isClient to true on mount
  useEffect(() => {
    setIsClient(true);
  }, []);

  useEffect(() => {
    // Only run this effect on the client side
    if (!isClient) return;

    // Function to load authentication state
    const loadAuthState = () => {
      try {
        const storedAccessToken = Cookies.get("accessToken");
        const storedOwnerAccessToken = Cookies.get("ownerAccessToken");

        if (storedAccessToken) {
          setAccessToken(storedAccessToken);
          api.defaults.headers.common[
            "Authorization"
          ] = `Bearer ${storedAccessToken}`;
          try {
            const decodedToken: any = jwtDecode(storedAccessToken);
            setUser({
              id: decodedToken.id,
              email: decodedToken.email,
              username: decodedToken.username || "",
              name: decodedToken.name || "",
            });
          } catch (error) {
            console.error("Error decoding access token:", error);
            // Clear invalid token
            Cookies.remove("accessToken");
          }
        }

        if (storedOwnerAccessToken) {
          setOwnerAccessToken(storedOwnerAccessToken);
          try {
            const decodedToken: any = jwtDecode(storedOwnerAccessToken);
            setOwner({
              owner_id: decodedToken.id,
              email: decodedToken.email,
              name: decodedToken.name,
              phone: decodedToken.phone,
              created_at: decodedToken.created_at,
            });
          } catch (error) {
            console.error("Error decoding owner access token:", error);
            // Clear invalid token
            Cookies.remove("ownerAccessToken");
          }
        }
      } finally {
        setLoading(false);
      }
    };

    loadAuthState();

    // Add event listener to sync auth state across tabs
    window.addEventListener("storage", (event) => {
      if (event.key === "auth:logout") {
        setUser(null);
        setAccessToken(null);
        setOwner(null);
        setOwnerAccessToken(null);
      }
    });
  }, [isClient]);

  const login = useCallback(
    async (data: ILoginRequestDTO): Promise<IAuthResponse> => {
      setLoading(true);
      try {
        const response: IAuthResponse = await userLogin(data);
        setUser(response.user);
        setAccessToken(response.accessToken);

        // Store in cookies instead of localStorage
        Cookies.set("accessToken", response.accessToken, {
          expires: 7, // 7 days
          sameSite: "lax",
          path: "/",
        });

        api.defaults.headers.common[
          "Authorization"
        ] = `Bearer ${response.accessToken}`;
        console.log("Auth state updated, redirecting to dashboard");

        router.push("/dashboard/user");
        return response;
      } catch (error) {
        console.error("Login failed:", error);
        throw error;
      } finally {
        setLoading(false);
      }
    },
    [router]
  );

  const register = useCallback(
    async (data: IRegisterRequestDTO): Promise<IAuthResponse> => {
      setLoading(true);
      try {
        console.log("Registration attempt with data:", {
          ...data,
          password: "[REDACTED]",
        });
        const response: IAuthResponse = await userRegister(data);
        console.log("Registration successful, response:", response);

        // Make sure we have the user data and token
        if (!response.user || !response.accessToken) {
          console.error(
            "Registration response missing user or token:",
            response
          );
          throw new Error("Invalid registration response");
        }

        // Update state
        setUser(response.user);
        setAccessToken(response.accessToken);

        // Store token in cookies
        Cookies.set("accessToken", response.accessToken, {
          expires: 7,
          sameSite: "lax",
          path: "/",
        });

        console.log("Auth state updated, redirecting to dashboard");

        // Add a slight delay to ensure state is updated before navigation
        setTimeout(() => {
          router.push("/dashboard/user");
        }, 100);

        return response;
      } catch (error) {
        console.error("Registration failed:", error);
        throw error;
      } finally {
        setLoading(false);
      }
    },
    [router]
  );

  const ownerLogin = useCallback(
    async (data: IOwnerLoginRequestDTO): Promise<IOwnerAuthResponse> => {
      setLoading(true);
      try {
        const response: IOwnerAuthResponse = await ownerLoginApi(data);
        setOwner(response.owner);
        setOwnerAccessToken(response.access_token);

        // Store in cookies instead of localStorage
        Cookies.set("ownerAccessToken", response.access_token, {
          expires: 7,
          sameSite: "lax",
          path: "/",
        });

        api.defaults.headers.common[
          "Authorization"
        ] = `Bearer ${ownerAccessToken}`;

        router.push("/dashboard/owner");
        return response;
      } catch (error) {
        console.error("Owner login failed:", error);
        throw error;
      } finally {
        setLoading(false);
      }
    },
    [router]
  );

  const ownerRegister = useCallback(
    async (data: IOwnerRegisterRequestDTO): Promise<IOwnerAuthResponse> => {
      setLoading(true);
      try {
        const response: IOwnerAuthResponse = await ownerRegisterApi(data);
        setOwner(response.owner);
        setOwnerAccessToken(response.access_token);

        // Store in cookies instead of localStorage
        Cookies.set("ownerAccessToken", response.access_token, {
          expires: 7,
          sameSite: "lax",
          path: "/",
        });

        api.defaults.headers.common[
          "Authorization"
        ] = `Bearer ${ownerAccessToken}`;
        router.push("/dashboard/owner");
        return response;
      } catch (error) {
        console.error("Owner registration failed:", error);
        throw error;
      } finally {
        setLoading(false);
      }
    },
    [router]
  );

  const logout = useCallback(() => {
    setUser(null);
    setAccessToken(null);
    Cookies.remove("accessToken");

    // Broadcast logout event to other tabs
    localStorage.setItem("auth:logout", Date.now().toString());

    router.push("/login");
  }, [router]);

  const ownerLogout = useCallback(() => {
    setOwner(null);
    setOwnerAccessToken(null);
    Cookies.remove("ownerAccessToken");

    // Broadcast logout event to other tabs
    localStorage.setItem("auth:logout", Date.now().toString());

    router.push("/owner-login");
  }, [router]);

  const isAuthenticated = useCallback(() => {
    return !!(accessToken || ownerAccessToken);
  }, [accessToken, ownerAccessToken]);

  const refreshToken = useCallback(async (): Promise<void> => {
    if (!accessToken) return;

    try {
      // Implement token refresh logic here
      const response = await api.post("/auth/refresh");
      const newToken = response.data.accessToken;

      setAccessToken(newToken);
      Cookies.set("accessToken", newToken, {
        expires: 7,
        sameSite: "lax",
        path: "/",
      });
    } catch (error) {
      console.error("Failed to refresh token:", error);
      // If refresh fails, log the user out
      logout();
    }
  }, [accessToken, logout]);

  const refreshOwnerToken = useCallback(async (): Promise<void> => {
    if (!ownerAccessToken) return;

    try {
      // Implement owner token refresh logic here
      const response = await api.post("/owner/auth/refresh");
      const newToken = response.data.access_token;

      setOwnerAccessToken(newToken);
      Cookies.set("ownerAccessToken", newToken, {
        expires: 7,
        sameSite: "lax",
        path: "/",
      });
    } catch (error) {
      console.error("Failed to refresh owner token:", error);
      // If refresh fails, log the owner out
      ownerLogout();
    }
  }, [ownerAccessToken, ownerLogout]);

  const contextValue = {
    user,
    owner,
    accessToken,
    ownerAccessToken,
    login,
    register,
    ownerLogin,
    ownerRegister,
    logout,
    ownerLogout,
    loading,
    isAuthenticated,
    refreshToken, // Add this
    refreshOwnerToken,
  };

  return (
    <AuthContext.Provider value={contextValue}>
      {/* Render a loading indicator or the children based on client/server rendering */}
      {!isClient && loading ? (
        <div className="flex items-center justify-center min-h-screen">
          <div className="animate-spin rounded-full h-12 w-12 border-t-2 border-b-2 border-primary"></div>
        </div>
      ) : (
        children
      )}
    </AuthContext.Provider>
  );
};
