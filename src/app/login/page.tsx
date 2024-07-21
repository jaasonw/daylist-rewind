"use client";
import { useAuth } from "@/components/PocketBaseAuthProvidor";
import { Button } from "@/components/ui/button";
import {
  Card,
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
// import { useRouter } from "next/router";
import { useRouter } from "next/navigation";
import PocketBase from "pocketbase";
import { useState } from "react";

export default function LoginForm() {
  const pb = new PocketBase("http://localhost:8090");
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [error, setError] = useState("");
  const router = useRouter();
  const { isAuthenticated } = useAuth();

  // loading
  if (isAuthenticated === null) return null;
  if (!isAuthenticated) return <></>;

  const handleLogin = async (event: { preventDefault: () => void }) => {
    event.preventDefault();
    try {
      const authData = await pb
        .collection("users")
        .authWithPassword(email, password);
      console.log("Logged in successfully:", authData);
      router.push("/dashboard");
      // Handle successful login, e.g., redirect or update state
    } catch (err) {
      console.error("Failed to log in:", err);
      setError("Login failed. Please check your email and password.");
    }
  };

  return (
    <section className="flex flex-col items-center justify-center min-h-screen">
      <Card className="w-full max-w-sm">
        <CardHeader>
          <CardTitle className="text-2xl">Login</CardTitle>
          <CardDescription>
            Enter your email below to login to your account.
          </CardDescription>
        </CardHeader>
        <form onSubmit={handleLogin}>
          <CardContent className="grid gap-4">
            <div className="grid gap-2">
              <Label htmlFor="email">Email</Label>
              <Input
                id="email"
                type="email"
                placeholder="m@example.com"
                value={email}
                onChange={(e) => setEmail(e.target.value)}
                required
              />
            </div>
            <div className="grid gap-2">
              <Label htmlFor="password">Password</Label>
              <Input
                id="password"
                type="password"
                value={password}
                onChange={(e) => setPassword(e.target.value)}
                required
              />
            </div>
            {error && <p className="text-red-500">{error}</p>}
          </CardContent>
          <CardFooter>
            <Button type="submit" className="w-full">
              Sign in
            </Button>
          </CardFooter>
        </form>
      </Card>
    </section>
  );
}
