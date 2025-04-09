// app/(auth)/owner-login/page.tsx
"use client";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { useState, FormEvent } from "react";
import Link from "next/link";
import { Alert, AlertDescription, AlertTitle } from "@/components/ui/alert";
import { useAuth } from "@/app/_components/AuthContext";
import { IOwnerLoginRequestDTO } from "@/app/_interfaces/auth";
import {
  Card,
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";

export default function OwnerLoginPage() {
  const { ownerLogin, loading } = useAuth();
  const [formData, setFormData] = useState<IOwnerLoginRequestDTO>({
    email: "",
    password: "",
  });
  const [error, setError] = useState<string | null>(null);

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setFormData({
      ...formData,
      [e.target.name]: e.target.value,
    });
  };

  const handleSubmit = async (e: FormEvent) => {
    e.preventDefault();
    setError(null);
    try {
      await ownerLogin(formData);
    } catch (err: any) {
      setError(err.response?.data?.error || "An unknown error occurred");
    }
  };

  return (
    <div className="flex items-center justify-center min-h-screen w-full px-4">
      <Card className="w-full max-w-md shadow-lg">
        <CardHeader>
          <CardTitle className="text-xl sm:text-2xl font-semibold text-center">
            Owner Login
          </CardTitle>
          <CardDescription className="text-center">
            Login to manage your restaurants
          </CardDescription>
        </CardHeader>
        <CardContent className="p-4 sm:p-6">
          {error && (
            <Alert variant="destructive" className="mb-4">
              <AlertTitle>Error</AlertTitle>
              <AlertDescription>{error}</AlertDescription>
            </Alert>
          )}
          <form onSubmit={handleSubmit}>
            <div className="mb-4">
              <Label htmlFor="email" className="mb-1 block">
                Email
              </Label>
              <Input
                type="email"
                id="email"
                name="email"
                value={formData.email}
                onChange={handleChange}
                required
                placeholder="Enter your email"
                className="w-full"
              />
            </div>
            <div className="mb-6">
              <Label htmlFor="password" className="mb-1 block">
                Password
              </Label>
              <Input
                type="password"
                id="password"
                name="password"
                value={formData.password}
                onChange={handleChange}
                required
                placeholder="Enter your password"
                className="w-full"
              />
            </div>
            <Button
              type="submit"
              className="w-full"
              disabled={loading}
              size="lg"
            >
              {loading ? "Loading..." : "Login"}
            </Button>
          </form>
        </CardContent>
        <CardFooter className="text-center p-4 sm:p-6 flex justify-center">
          <p className="text-sm sm:text-base">
            Don't have an account?{" "}
            <Link
              href="/owner-register"
              className="text-blue-500 hover:underline"
            >
              Register
            </Link>
          </p>
        </CardFooter>
      </Card>
    </div>
  );
}
