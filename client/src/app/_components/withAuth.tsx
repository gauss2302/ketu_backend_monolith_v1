// app/_components/withAuth.tsx
"use client";
import { useAuth } from "./AuthContext";
import { useRouter, usePathname } from "next/navigation";
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
    const pathname = usePathname();
    const { userType } = options;

    useEffect(() => {
      if (!loading) { // Check the loading state first!
        const isLoggedIn = user && accessToken;
        const isOwnerLoggedIn = owner && ownerAccessToken;

        // Check if the user is already on the login or owner login page
        const isLoginPage = pathname === "/login";
        const isOwnerLoginPage = pathname === "/owner-login";

        if (!isLoggedIn && !isOwnerLoggedIn) {
          if (!isLoginPage && !isOwnerLoginPage) {
            router.push("/login"); // Redirect if not logged in and not on login pages
          }
        } else if (userType === "user" && !isLoggedIn) {
          router.push("/dashboard/owner");
        } else if (userType === "owner" && !isOwnerLoggedIn) {
          router.push("/dashboard/user");
        } else if ((isLoggedIn && isLoginPage) || (isOwnerLoggedIn && isOwnerLoginPage)){
          if (userType === "user" && isLoggedIn){
            router.push("/dashboard/user");
          } else if (userType === "owner" && isOwnerLoggedIn){
            router.push("/dashboard/owner");
          }
        }
      }
    }, [loading, user, owner, router, userType, accessToken, ownerAccessToken, pathname]);

    if (loading) {
      return <div>Loading...</div>;
    }
    if ((!user && userType === "user") || (!owner && userType === "owner")) {
      return null; // Or a loading indicator/placeholder
    }

    return <WrappedComponent {...props} />;
  };

  return Wrapper;
};

export default withAuth;