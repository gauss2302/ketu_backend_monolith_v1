"use client";
import React, { createContext, useContext } from "react";
import { IUser, IOwner } from "../_interfaces/auth"; // Import user and owner interfaces
import {
  ILoginRequestDTO,
  IRegisterRequestDTO,
  IOwnerLoginRequestDTO,
  IOwnerRegisterRequestDTO,
} from "../_interfaces/auth";

interface AuthContextType {
  user: IUser | null;
  owner: IOwner | null;
  accessToken: string | null;
  ownerAccessToken: string | null;
  loading: boolean;
  login: (data: any) => Promise<void>;
  register: (data: any) => Promise<void>;
  ownerLogin: (data: any) => Promise<void>;
  ownerRegister: (data: any) => Promise<void>;
  logout: () => void;
  ownerLogout: () => void;
}

export const AuthContext = createContext<AuthContextType | null>(null);

export const useAuth = () => {
  const context = useContext(AuthContext);
  if (!context) {
    throw new Error("useAuth must be used within an AuthProvider");
  }
  return context;
};
