// app/_components/AuthProvider.tsx
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

type AuthProviderProps = {
  children: ReactNode;
};

export const AuthProvider: React.FC<AuthProviderProps> = ({ children }) => {
  const [user, setUser] = useState<IUser | null>(null);
  const [owner, setOwner] = useState<IOwner | null>(null);
  const [accessToken, setAccessToken] = useState<string | null>(null);
  const [ownerAccessToken, setOwnerAccessToken] = useState<string | null>(null);
  const [loading, setLoading] = useState<boolean>(false);
  const router = useRouter();

  // Load token from local storage on initial load
  useEffect(() => {
    const storedAccessToken = localStorage.getItem("accessToken");
    const storedOwnerAccessToken = localStorage.getItem("ownerAccessToken");

    if (storedAccessToken) {
      setAccessToken(storedAccessToken);
      // You might also fetch user details here if needed.
    }
    if (storedOwnerAccessToken) {
      setOwnerAccessToken(storedOwnerAccessToken);
    }
  }, []);

  const login = useCallback(
    async (data: ILoginRequestDTO) => {
      setLoading(true);
      try {
        const response: IAuthResponse = await userLogin(data);
        setUser(response.user);
        setAccessToken(response.accessToken);
        localStorage.setItem("accessToken", response.accessToken);
        router.push("/dashboard/user"); // Redirect to user dashboard
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
    async (data: IRegisterRequestDTO) => {
      setLoading(true);
      try {
        const response: IAuthResponse = await userRegister(data);
        setUser(response.user);
        setAccessToken(response.accessToken);
        localStorage.setItem("accessToken", response.accessToken);
        router.push("/dashboard/user"); // Redirect to user dashboard
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
    async (data: IOwnerLoginRequestDTO) => {
      setLoading(true);
      try {
        const response: IOwnerAuthResponse = await ownerLoginApi(data);
        setOwner(response.owner);
        setOwnerAccessToken(response.access_token);
        localStorage.setItem("ownerAccessToken", response.access_token);
        router.push("/dashboard/owner"); // Redirect to owner dashboard
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
    async (data: IOwnerRegisterRequestDTO) => {
      setLoading(true);
      try {
        const response: IOwnerAuthResponse = await ownerRegisterApi(data);
        setOwner(response.owner);
        setOwnerAccessToken(response.access_token);
        localStorage.setItem("ownerAccessToken", response.access_token);
        router.push("/dashboard/owner"); // Redirect to owner dashboard
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
  };

  return (
    <AuthContext.Provider value={contextValue}>{children}</AuthContext.Provider>
  );
};
