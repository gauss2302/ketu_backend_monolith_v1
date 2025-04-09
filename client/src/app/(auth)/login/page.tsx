// In app/(auth)/login/page.tsx

"use client";
import { useState } from "react";
import { useAuth } from "@/app/_components/AuthContext";
import { ILoginRequestDTO } from "@/app/_interfaces/auth";
import { WavyBackground } from "@/components/ui/wavy-background";
import { AuthForm } from "@/app/_components/auth-form";

export default function LoginPage() {
  const { login, loading } = useAuth();
  const [error, setError] = useState<string | null>(null);

  const handleSubmit = async (formData: ILoginRequestDTO) => {
    setError(null);
    try {
      await login(formData);
      // eslint-disable-next-line @typescript-eslint/no-explicit-any
    } catch (err: any) {
      setError(err.response?.data?.error || "An error occurred");
    }
  };

  return (
    <WavyBackground
      colors={["#6366f1", "#8b5cf6", "#a855f7", "#d946ef", "#ec4899"]}
      blur={5}
      waveWidth={30}
      waveOpacity={0.7}
    >
      <div className="flex items-center justify-center min-h-screen relative z-10">
        <AuthForm
          title="User Login"
          description="Enter your credentials to access your account."
          submitText="Login"
          linkText="Register"
          linkHref="/register"
          linkMessage="Don't have an account?"
          onSubmitAction={handleSubmit} // Changed from onSubmit to onSubmitAction
          initialFormData={{ email: "", password: "" }}
          error={error}
          loading={loading}
          fields={[
            { name: "email", label: "Email", type: "email" },
            { name: "password", label: "Password", type: "password" },
          ]}
        />
      </div>
    </WavyBackground>
  );
}
