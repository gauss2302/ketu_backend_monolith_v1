// app/_components/withAuth.tsx
"use client";
import { useAuth } from "./AuthContext";
import { useRouter } from "next/navigation";
import { useEffect, ReactNode } from "react";

type AuthOptions = {
  userType?: "user" | "owner";
};

const withAuth = (
  WrappedComponent: React.ComponentType<any>,
  options: AuthOptions = {}
) => {
  const Wrapper = (props: any) => {
    const { user, owner, loading, accessToken, ownerAccessToken } = useAuth();
    const router = useRouter();
    const { userType } = options;

    useEffect(() => {
      if (!loading) {
        const isLoggedIn = user && accessToken;
        const isOwnerLoggedIn = owner && ownerAccessToken;

        if (!isLoggedIn && !isOwnerLoggedIn) {
          router.push("/login"); // Redirect if no token
        } else if (userType === "user" && !isLoggedIn) {
          router.push("/dashboard/owner");
        } else if (userType === "owner" && !isOwnerLoggedIn) {
          router.push("/dashboard/user");
        }
      }
    }, [loading, user, owner, router, userType, accessToken, ownerAccessToken]);

    if (loading) {
      return <div>Loading...</div>;
    }
    if (!user && !owner) {
      return null; // Or a loading indicator/placeholder
    }

    return <WrappedComponent {...props} />;
  };

  return Wrapper;
};

export default withAuth;
