"use client";
import React, { useState, ReactNode, useCallback, useEffect } from "react";
import { AuthContext } from "./AuthContext";
import {
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

type AuthProviderProps = {
  children: ReactNode;
};

export const AuthProvider: React.FC<AuthProviderProps> = ({ children }) => {
  const [user, setUser] = useState<IUser | null>(null);
  const [owner, setOwner] = useState<IOwner | null>(null);
  const [accessToken, setAccessToken] = useState<string | null>(null);
  const [ownerAccessToken, setOwnerAccessToken] = useState<string | null>(null);
  const [loading, setLoading] = useState<boolean>(true);
  const router = useRouter();

  useEffect(() => {
    const storedAccessToken = localStorage.getItem("accessToken");
    const storedOwnerAccessToken = localStorage.getItem("ownerAccessToken");

    if (storedAccessToken) {
      setAccessToken(storedAccessToken);
      try {
        // eslint-disable-next-line @typescript-eslint/no-explicit-any
        const decodedToken: any = jwtDecode(storedAccessToken);
        setUser({
          id: decodedToken.id,
          email: decodedToken.email,
          username: "", // You may need to fetch the username from the backend
          name: "", // You may need to fetch the name from the backend
          // createdAt: "", // You may need to fetch the createdAt from the backend
          // updatedAt: "", // You may need to fetch the updatedAt from the backend
        });
      } catch (error) {
        console.error("Error decoding access token:", error);
      }
    }

    if (storedOwnerAccessToken) {
      setOwnerAccessToken(storedOwnerAccessToken);
      // Decode owner token and set owner state similarly
    }

    setLoading(false);
  }, []);
  const login = useCallback(
    async (data: ILoginRequestDTO): Promise<IAuthResponse> => {
      // Corrected
      setLoading(true);
      try {
        const response: IAuthResponse = await userLogin(data);
        setUser(response.user);
        setAccessToken(response.accessToken);
        localStorage.setItem("accessToken", response.accessToken);
        router.push("/dashboard/user"); // Redirect to user dashboard
        return response; // Return the response
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
      // Corrected
      setLoading(true);
      try {
        const response: IAuthResponse = await userRegister(data);
        setUser(response.user);
        setAccessToken(response.accessToken);
        localStorage.setItem("accessToken", response.accessToken);
        router.push("/dashboard/user"); // Redirect to user dashboard
        return response; // Return the response
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
      // Corrected
      setLoading(true);
      try {
        const response: IOwnerAuthResponse = await ownerLoginApi(data);
        setOwner(response.owner);
        setOwnerAccessToken(response.access_token);
        localStorage.setItem("ownerAccessToken", response.access_token);
        router.push("/dashboard/owner"); // Redirect to owner dashboard
        return response; // Return the response
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
      // Corrected
      setLoading(true);
      try {
        const response: IOwnerAuthResponse = await ownerRegisterApi(data);
        setOwner(response.owner);
        setOwnerAccessToken(response.access_token);
        localStorage.setItem("ownerAccessToken", response.access_token);
        router.push("/dashboard/owner"); // Redirect to owner dashboard
        return response; // Return the response
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
    localStorage.removeItem("accessToken");
    router.push("/login");
  }, [router]);

  const ownerLogout = useCallback(() => {
    setOwner(null);
    setOwnerAccessToken(null);
    localStorage.removeItem("ownerAccessToken");
    router.push("/owner-login");
  }, [router]);

  const isAuthenticated = () => {
    return !!(accessToken || ownerAccessToken);
  };

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
  };

  return (
    <AuthContext.Provider value={contextValue}>{children}</AuthContext.Provider>
  );
};
