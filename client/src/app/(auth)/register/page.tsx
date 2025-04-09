// In app/(auth)/register/page.tsx

"use client";
import { useState } from "react";
import { useAuth } from "@/app/_components/AuthContext";
import { IRegisterRequestDTO } from "@/app/_interfaces/auth";
import { WavyBackground } from "@/components/ui/wavy-background";
import { AuthForm } from "@/app/_components/auth-form";
import { useRouter } from "next/navigation";

export default function RegisterPage() {
  const { register, loading } = useAuth();
  const [error, setError] = useState<string | null>(null);
  const router = useRouter();

  const handleSubmit = async (formData: IRegisterRequestDTO) => {
    setError(null);
    try {
      console.log("Submitting registration form:", {
        ...formData,
        password: "[REDACTED]",
      });
      const result = await register(formData);
      console.log("Registration successful:", result);

      // Add manual navigation as a backup
      router.push("/dashboard/user");
      // eslint-disable-next-line @typescript-eslint/no-explicit-any
    } catch (err: any) {
      console.error("Registration error:", err);
      setError(err.response?.data?.error || "An unknown error occurred");
    }
  };

  return (
    <WavyBackground
      colors={["#6366f1", "#8b5cf6", "#a855f7", "#d946ef", "#ec4899"]}
      blur={5}
      waveWidth={30}
      waveOpacity={0.7}
    >
      <div className="flex items-center justify-center min-h-screen">
        <AuthForm
          title="User Register"
          description="Create your account to get started."
          submitText="Register"
          linkText="Login"
          linkHref="/login"
          linkMessage="Already have an account?"
          onSubmitAction={handleSubmit}
          initialFormData={{ username: "", name: "", email: "", password: "" }}
          error={error}
          loading={loading}
          fields={[
            {
              name: "username",
              label: "Username",
              type: "text",
              placeholder: "Enter your username",
            },
            {
              name: "name",
              label: "Name",
              type: "text",
              placeholder: "Enter your name",
            },
            {
              name: "email",
              label: "Email",
              type: "email",
              placeholder: "Enter your email",
            },
            {
              name: "password",
              label: "Password",
              type: "password",
              placeholder: "Enter your password",
            },
          ]}
        />
      </div>
    </WavyBackground>
  );
}
