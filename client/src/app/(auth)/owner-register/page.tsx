// app/(auth)/owner-register/page.tsx
"use client";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { useState, FormEvent } from "react";
import Link from "next/link";
import { IOwnerRegisterRequestDTO } from "@/app/_interfaces/auth";
import { Alert, AlertDescription, AlertTitle } from "@/components/ui/alert";
import { useAuth } from "@/app/_components/AuthContext";
import {
  Card,
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";

export default function OwnerRegisterPage() {
  const { ownerRegister, loading } = useAuth();
  const [formData, setFormData] = useState<IOwnerRegisterRequestDTO>({
    name: "",
    email: "",
    phone: "",
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
      await ownerRegister(formData);
    } catch (err: any) {
      setError(err.response?.data?.error || "An unknown error occurred");
    }
  };

  return (
    <div className="flex items-center justify-center min-h-screen w-full px-4">
      <Card className="w-full max-w-md shadow-lg">
        <CardHeader>
          <CardTitle className="text-xl sm:text-2xl font-semibold text-center">
            Owner Register
          </CardTitle>
          <CardDescription className="text-center">
            Create a new account to manage your restaurants
          </CardDescription>
        </CardHeader>
        <CardContent className="p-4 sm:p-6">
          {error && (
            <Alert variant="destructive" className="mb-4">
              <AlertTitle>Error</AlertTitle>
              <AlertDescription>{error}</AlertDescription>
            </Alert>
          )}
          <form onSubmit={handleSubmit} className="space-y-4">
            <div>
              <Label htmlFor="name" className="mb-1 block">
                Name
              </Label>
              <Input
                type="text"
                id="name"
                name="name"
                value={formData.name}
                onChange={handleChange}
                required
                placeholder="Enter your name"
              />
            </div>
            <div>
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
              />
            </div>
            <div>
              <Label htmlFor="phone" className="mb-1 block">
                Phone
              </Label>
              <Input
                type="tel"
                id="phone"
                name="phone"
                value={formData.phone}
                onChange={handleChange}
                required
                placeholder="Enter your phone number"
              />
            </div>
            <div>
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
              />
            </div>
            <Button
              type="submit"
              className="w-full"
              disabled={loading}
              size="lg"
            >
              {loading ? "Loading..." : "Register"}
            </Button>
          </form>
        </CardContent>
        <CardFooter className="text-center p-4 sm:p-6 flex justify-center">
          <p className="text-sm sm:text-base">
            Already have an account?{" "}
            <Link href="/owner-login" className="text-blue-500 hover:underline">
              Login
            </Link>
          </p>
        </CardFooter>
      </Card>
    </div>
  );
}
