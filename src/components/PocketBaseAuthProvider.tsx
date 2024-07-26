"use client";
import { createContext, useContext, useEffect, useState } from "react";
import PocketBase from "pocketbase";
import { usePathname, useRouter } from "next/navigation";

type AuthContextType = {
  isAuthenticated: boolean | null;
};

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export const AuthProvider: React.FC<{ children: React.ReactNode }> = ({
  children,
}) => {
  const [isAuthenticated, setIsAuthenticated] = useState<boolean | null>(null);
  const [loading, setLoading] = useState(true);
  const router = useRouter();
  const pathname = usePathname();
  const pb = new PocketBase(process.env.NEXT_PUBLIC_POCKETBASE_URL);

  useEffect(() => {
    const checkAuth = async () => {
      if (pb.authStore.isValid) {
        setIsAuthenticated(true);
      } else {
        setIsAuthenticated(false);
      }
      setLoading(false); // Set loading to false after checking authentication
    };

    checkAuth();
  }, [pb.authStore.isValid, pathname, router]);

  if (loading) {
    return <div>loading...</div>; // Display a loading state until authentication is checked
  }

  return (
    <AuthContext.Provider value={{ isAuthenticated }}>
      {children}
    </AuthContext.Provider>
  );
};

export const useAuth = () => {
  const context = useContext(AuthContext);
  if (context === undefined) {
    throw new Error("useAuth must be used within an AuthProvider");
  }
  return context;
};
