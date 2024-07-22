"use client";
import { useAuth } from "@/components/PocketBaseAuthProvider";
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
import { useRouter } from "next/navigation";
import PocketBase, { ClientResponseError } from "pocketbase";
import { useState } from "react";

export default function LoginForm() {
  const pb = new PocketBase("http://localhost:8090");
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [confirmPassword, setConfirmPassword] = useState("");
  const [error, setError] = useState("");
  const router = useRouter();

  const handleLogin = async (event: { preventDefault: () => void }) => {
    event.preventDefault();
    try {
      const authData = await pb.collection("users").create({
        email: email,
        password: password,
        passwordConfirm: confirmPassword,
      });
      console.log("Registered successfully:", authData);
      router.push("/login");
    } catch (err) {
      if (err instanceof ClientResponseError) {
        console.error("Failed to create account:", err.response);
        setError(
          "Failed to create account:\n" +
            Object.keys(err.response.data)
              .map((key) => err.response.data[key].message)
              .join("\n")
        );
      }
    }
  };

  return (
    <section className="flex flex-col items-center justify-center min-h-screen">
      <Card className="w-full max-w-sm">
        <CardHeader>
          <CardTitle className="text-2xl">Sign Up</CardTitle>
          <CardDescription>
            Enter your email below to create your account.
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
            <div className="grid gap-2">
              <Label htmlFor="password">Confirm Password</Label>
              <Input
                id="password"
                type="password"
                value={confirmPassword}
                onChange={(e) => setConfirmPassword(e.target.value)}
                required
              />
            </div>
            {error && (
              <div className="text-red-500">
                {error.split("\n").map((line, index) => (
                  <p key={index}>{line}</p>
                ))}
              </div>
            )}
          </CardContent>
          <CardFooter>
            <Button type="submit" className="w-full">
              Sign up
            </Button>
          </CardFooter>
        </form>
      </Card>
    </section>
  );
}
