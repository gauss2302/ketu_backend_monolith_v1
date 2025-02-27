/* eslint-disable @typescript-eslint/no-explicit-any */
// _components/auth-form.tsx
"use client";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { useState, FormEvent } from "react";
import Link from "next/link";
import { Alert, AlertDescription, AlertTitle } from "@/components/ui/alert";
import { Eye, EyeOff } from "lucide-react";
import {
  Card,
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";

interface AuthFormProps {
  title: string;
  description: string;
  submitText: string;
  linkText: string;
  linkHref: string;
  linkMessage: string;
  onSubmit: (formData: any) => Promise<void>;
  initialFormData: any;
  error: string | null;
  loading: boolean;
  fields: {
    name: string;
    label: string;
    type: string;
    placeholder?: string;
  }[];
}

export function AuthForm({
  title,
  description,
  submitText,
  linkText,
  linkHref,
  linkMessage,
  onSubmit,
  initialFormData,
  error,
  loading,
  fields,
}: AuthFormProps) {
  const [formData, setFormData] = useState(initialFormData);
  const [showPassword, setShowPassword] = useState(false);

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setFormData({
      ...formData,
      [e.target.name]: e.target.value,
    });
  };

  const handleSubmit = async (e: FormEvent) => {
    e.preventDefault();
    await onSubmit(formData);
  };

  const togglePasswordVisibility = () => {
    setShowPassword(!showPassword);
  };

  return (
    <Card className="w-full max-w-md">
      <CardHeader>
        <CardTitle className="text-2xl font-semibold text-center">
          {title}
        </CardTitle>
        <CardDescription className="text-center">{description}</CardDescription>
      </CardHeader>
      <CardContent className="p-6 w-full">
        {error && (
          <Alert variant="destructive">
            <AlertTitle>Error</AlertTitle>
            <AlertDescription>{error}</AlertDescription>
          </Alert>
        )}
        <form onSubmit={handleSubmit}>
          {fields.map((field) => (
            <div className="mb-4" key={field.name}>
              <Label htmlFor={field.name}>{field.label}</Label>
              <Input
                type={
                  field.name === "password" && showPassword
                    ? "text"
                    : field.type
                }
                id={field.name}
                name={field.name}
                value={formData[field.name]}
                onChange={handleChange}
                required
                placeholder={field.placeholder}
              />
              {field.name === "password" && (
                <button
                  type="button"
                  className="absolute top-1/2 right-0 pr-3 flex items-center text-gray-400 transform -translate-y-1/2"
                  onClick={togglePasswordVisibility}
                >
                  {showPassword ? <EyeOff /> : <Eye />}
                </button>
              )}
            </div>
          ))}
          <Button type="submit" className="w-full" disabled={loading}>
            {loading ? "Loading..." : submitText}
          </Button>
        </form>
      </CardContent>
      <CardFooter className="text-center p-6">
        <p>
          {linkMessage}{" "}
          <Link href={linkHref} className="text-blue-500 hover:underline">
            {linkText}
          </Link>
        </p>
      </CardFooter>
    </Card>
  );
}
