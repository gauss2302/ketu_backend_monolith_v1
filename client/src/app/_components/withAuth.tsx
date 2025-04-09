/* eslint-disable @typescript-eslint/no-explicit-any */
// app/_components/withAuth.tsx
// In app/_components/withAuth.tsx

"use client";
import { useAuth } from "./AuthContext";
import { useRouter, usePathname } from "next/navigation";
import { useEffect, useState } from "react";

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
    const pathname = usePathname();
    const { userType } = options;
    const [shouldRender, setShouldRender] = useState(false);

    useEffect(() => {
      if (!loading) {
        const isLoggedIn = user && accessToken;
        const isOwnerLoggedIn = owner && ownerAccessToken;

        // Check if the user is already on the login or owner login page
        const isLoginPage = pathname === "/login";
        const isOwnerLoginPage = pathname === "/owner-login";

        if (!isLoggedIn && !isOwnerLoggedIn) {
          if (!isLoginPage && !isOwnerLoginPage) {
            router.replace("/login");
            return;
          }
        } else if (userType === "user" && !isLoggedIn) {
          router.replace("/dashboard/owner");
          return;
        } else if (userType === "owner" && !isOwnerLoggedIn) {
          router.replace("/dashboard/user");
          return;
        } else if (
          (isLoggedIn && isLoginPage) ||
          (isOwnerLoggedIn && isOwnerLoginPage)
        ) {
          if (userType === "user" && isLoggedIn) {
            router.replace("/dashboard/user");
            return;
          } else if (userType === "owner" && isOwnerLoggedIn) {
            router.replace("/dashboard/owner");
            return;
          }
        }

        // If we passed all the redirect conditions, we can render the component
        setShouldRender(true);
      }
    }, [
      loading,
      user,
      owner,
      router,
      userType,
      accessToken,
      ownerAccessToken,
      pathname,
    ]);

    if (loading) {
      return (
        <div className="flex items-center justify-center min-h-screen">
          <div className="animate-spin rounded-full h-12 w-12 border-t-2 border-b-2 border-primary"></div>
        </div>
      );
    }

    // Don't render the component until we're sure we don't need to redirect
    if (!shouldRender) {
      return null;
    }

    return <WrappedComponent {...props} />;
  };

  return Wrapper;
};

export default withAuth;
